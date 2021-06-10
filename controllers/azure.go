package controllers

import (
	"context"
	"fmt"
	securityv1 "linkedsecrets/api/v1"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/keyvault/keyvault"
	kvauth "github.com/Azure/azure-sdk-for-go/services/keyvault/auth"
)

// Credentials will be provided by environment variables:
// AZURE_TENANT_ID, AZURE_CLIENT_ID and AZURE_CLIENT_SECRET
// GetAzureSecret return secret data from AWS Secret Manage
func (r *LinkedSecretReconciler) GetAzureSecret(linkedsecret *securityv1.LinkedSecret) ([]byte, error) {

	log := r.Log.WithValues("linkedsecret", fmt.Sprintf("%s/%s", linkedsecret.Namespace, linkedsecret.Name))

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
