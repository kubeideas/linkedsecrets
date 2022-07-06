package controllers

import (
	"context"
	securityv1 "kubeideas/linkedsecrets/api/v1"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/keyvault/keyvault"
	kvauth "github.com/Azure/azure-sdk-for-go/services/keyvault/auth"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// Credentials will be provided by environment variables:
// AZURE_TENANT_ID, AZURE_CLIENT_ID and AZURE_CLIENT_SECRET
// GetAzureSecret return secret data from AWS Secret Manage
func (r *LinkedSecretReconciler) GetAzureSecret(ctx context.Context, linkedsecret *securityv1.LinkedSecret) ([]byte, error) {

	log := log.FromContext(ctx)

	// check required provider options
	if _, ok := linkedsecret.Spec.ProviderOptions["keyvault"]; !ok {
		return nil, &InvalidAzureKeyvault{}
	}

	if _, ok := linkedsecret.Spec.ProviderOptions["secret"]; !ok {
		return nil, &InvalidSecretOption{}
	}

	// get provider options informed in linkedsecret spec
	kvName := linkedsecret.Spec.ProviderOptions["keyvault"]
	name := linkedsecret.Spec.ProviderOptions["secret"]

	// set default "" if providerOption version was not specified
	// "" means latest secret version
	version := ""

	// get version if defined
	if _, ok := linkedsecret.Spec.ProviderOptions["version"]; ok {
		version = linkedsecret.Spec.ProviderOptions["version"]
	}

	// create authorizer
	authorizer, err := kvauth.NewAuthorizerFromEnvironment()
	if err != nil {
		log.V(1).Info("Azure Error creating vault authorizer", kvName, err)
		return nil, err
	}

	//create session
	basicClient := keyvault.New()
	basicClient.RetryAttempts = 1
	basicClient.RetryDuration = 1 * time.Second
	basicClient.Authorizer = authorizer

	// get secret
	secret, err := basicClient.GetSecret(context.Background(), "https://"+kvName+".vault.azure.net", name, version)
	if err != nil {
		log.V(1).Info("Azure Error getting secret", name, err)
		return nil, err
	}

	log.V(1).Info("Azure return secret", "secret", name)
	return []byte(*secret.Value), nil
}
