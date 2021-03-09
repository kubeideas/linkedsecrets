package controllers

// Use this code snippet in your app.
// If you need more information about configurations or implementing the sample code, visit the AWS docs:
// https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/setting-up.html

import (
	"encoding/base64"
	"fmt"
	securityv1 "linkedsecrets/api/v1"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

// GetAWSSecret return secret data from AWS Secret Manage
func (r *LinkedSecretReconciler) GetAWSSecret(linkedsecret *securityv1.LinkedSecret) ([]byte, error) {

	log := r.Log.WithValues("linkedsecret", fmt.Sprintf("%s/%s", linkedsecret.Namespace, linkedsecret.Name))

	// get provider options informed in linkedsecret spec
	name := linkedsecret.Spec.ProviderOptions["secret"]
	region := linkedsecret.Spec.ProviderOptions["region"]
	version := linkedsecret.Spec.ProviderOptions["version"]

	//Create a Secrets Manager client with informed region
	svc := secretsmanager.New(session.New(), aws.NewConfig().WithRegion(region))

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
		//secretString = *result.SecretString
		data = []byte(*result.SecretString)
	} else {
		//decode binary secret
		data = make([]byte, base64.StdEncoding.DecodedLen(len(result.SecretBinary)))
		_, err := base64.StdEncoding.Decode(data, result.SecretBinary)
		if err != nil {
			log.V(1).Info("AWS Error decoding secret", name, err)
			return data, err

		}
		//decodedBinarySecret = string(decodedBinarySecretBytes[:len])
		//fmt.Println("secretString = ", secretString)
		//fmt.Println("decodedBinarySecret = ", decodedBinarySecret)
	}

	log.V(1).Info("AWS return secret", "secret", name)
	return data, nil
}
