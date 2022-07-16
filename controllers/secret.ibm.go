package controllers

import (
	"context"
	"encoding/base64"
	"fmt"
	securityv1 "kubeideas/linkedsecrets/api/v1"
	"os"

	"github.com/IBM/go-sdk-core/core"
	sm "github.com/IBM/secrets-manager-go-sdk/secretsmanagerv1"
	"github.com/go-openapi/strfmt"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// Credentials will be provided by environment variables:

func (r *LinkedSecretReconciler) GetIBMSecret(ctx context.Context, linkedsecret *securityv1.LinkedSecret) ([]byte, error) {

	log := log.FromContext(ctx)

	// check required provider options

	if _, ok := linkedsecret.Spec.ProviderOptions["secretManagerInstanceId"]; !ok {
		return nil, &InvalidIBMSecretManagerID{}
	}

	//secretManagerInstanceId must be UUID formated
	if !strfmt.IsUUID(linkedsecret.Spec.ProviderOptions["secretManagerInstanceId"]) {
		return nil, &InvalidUUIDFormat{
			fieldName: "secretManagerInstanceId",
		}
	}

	if _, ok := linkedsecret.Spec.ProviderOptions["region"]; !ok {
		return nil, &InvalidIBMRegion{}
	}

	if _, ok := linkedsecret.Spec.ProviderOptions["secretId"]; !ok {
		return nil, &InvalidSecretOption{}
	}

	//SecretId must be UUID formated
	if !strfmt.IsUUID(linkedsecret.Spec.ProviderOptions["secretId"]) {
		return nil, &InvalidUUIDFormat{
			fieldName: "secretId",
		}
	}

	// get provider options informed in linkedsecret spec
	secretManagerInstanceId := linkedsecret.Spec.ProviderOptions["secretManagerInstanceId"]
	region := linkedsecret.Spec.ProviderOptions["region"]
	secretId := linkedsecret.Spec.ProviderOptions["secretId"]

	//Service API key is used to access IBM Secrets
	serviceApiKey, defVar := os.LookupEnv("IBM_SERVICE_API_KEY")

	if !defVar {
		log.V(1).Info("[ IBM_SERVICE_API_KEY ] environment variable not defined.")
		return nil, &InvalidIBMServiceApiKey{}
	}

	// Create manager
	secretsManager, err := sm.NewSecretsManagerV1(&sm.SecretsManagerV1Options{
		URL: fmt.Sprintf("https://%s.%s.secrets-manager.appdomain.cloud", secretManagerInstanceId, region),
		Authenticator: &core.IamAuthenticator{
			ApiKey: serviceApiKey,
		},
	})

	if err != nil {
		log.V(1).Info("IBM Error creating secret manager", "Error", err)
		return nil, err
	}

	// Get secret
	getSecretRes, _, err := secretsManager.GetSecret(&sm.GetSecretOptions{
		SecretType: core.StringPtr(sm.GetSecretOptionsSecretTypeArbitraryConst),
		ID:         core.StringPtr(secretId),
	})

	if err != nil {

		log.V(1).Info("IBM Error getting secret", "Error", err)
		return nil, err
	}

	secret := getSecretRes.Resources[0].(*sm.SecretResource)

	secretData := secret.SecretData.(map[string]interface{})

	var secretPayload []byte

	// Check for base64 encode label
	// If payload base64 decoding returns any error, it will be considered not encoded.
	if decPayload, err := base64.StdEncoding.DecodeString(secretData["payload"].(string)); err == nil {
		secretPayload = decPayload
	} else {
		log.V(1).Info("Secret returned from IBM may be not encoded", "Error", err)
		log.V(1).Info("IBM secret will be assigned directly", "notDecoded", "secretData")
		secretPayload = []byte(secretData["payload"].(string))
	}

	log.V(1).Info("IBM return secret", "secretId", secretId)
	return secretPayload, nil
}
