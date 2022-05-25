package controllers

// Use this code snippet in your app.
// If you need more information about configurations or implementing the sample code, visit the AWS docs:
// https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/setting-up.html

import (
	"encoding/base64"
	"fmt"
	securityv1 "kubeideas/linkedsecrets/api/v1"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

// Credentials will be provided by environment variables:
// AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY
// GetAWSSecret return secret data from AWS Secret Manager
func (r *LinkedSecretReconciler) GetAWSSecret(linkedsecret *securityv1.LinkedSecret) ([]byte, error) {

	log := r.Log.WithValues("linkedsecret", fmt.Sprintf("%s/%s", linkedsecret.Namespace, linkedsecret.Name))

	// check required provider options
	if _, ok := linkedsecret.Spec.ProviderOptions["region"]; !ok {
		return nil, &InvalidAWSRegion{}
	}

	if _, ok := linkedsecret.Spec.ProviderOptions["secret"]; !ok {
		return nil, &InvalidSecretOption{}
	}

	// get provider options informed in linkedsecret spec
	name := linkedsecret.Spec.ProviderOptions["secret"]
	region := linkedsecret.Spec.ProviderOptions["region"]

	// set default "AWSCURRENT" if providerOption version was not specified
	version := "AWSCURRENT"

	// get version if defined
	if _, ok := linkedsecret.Spec.ProviderOptions["version"]; ok {
		version = linkedsecret.Spec.ProviderOptions["version"]
	}

	// new cloud session
	sess, err := session.NewSession(aws.NewConfig().WithRegion(region))
	if err != nil {
		fmt.Println("Error creating session ", err)
		return nil, err
	}

	//Create a Secrets Manager client with informed region
	svc := secretsmanager.New(sess)

	//prepare input with informed name and version
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(name),
		VersionStage: aws.String(version), // VersionStage defaults to AWSCURRENT if unspecified
	}

	// In this sample we only handle the specific exceptions for the 'GetSecretValue' API.
	// See https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html

	// retrieve secret
	result, err := svc.GetSecretValue(input)
	if err != nil {
		log.V(1).Info("AWS Error getting secret", name, err)
		return nil, err
	}

	//var secretString, decodedBinarySecret string
	var data []byte
	if result.SecretString != nil {
		data = []byte(*result.SecretString)
	} else {
		// if SecretBinary base64 decoding returns any error, it will be considered not encoded.
		if decSecretBinary, err := base64.StdEncoding.DecodeString(string(result.SecretBinary)); err == nil {
			data = decSecretBinary
		} else {
			log.V(1).Info("Secret returned from AWS may be not encoded", "Error", err)
			log.V(1).Info("AWS secret will be assigned directly", "notDecoded", "SecretBinary")
			data = result.SecretBinary
		}
	}

	log.V(1).Info("AWS return secret", "secret", name)
	return data, nil
}
