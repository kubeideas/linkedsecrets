package controllers

import (
	"context"
	securityv1 "kubeideas/linkedsecrets/api/v1"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

//GetGCPSecret return secret data from Google Secret Manager
func (r *LinkedSecretReconciler) GetGCPSecret(ctx context.Context, linkedsecret *securityv1.LinkedSecret) ([]byte, error) {

	log := log.FromContext(ctx)

	// check required provider options
	if _, ok := linkedsecret.Spec.ProviderOptions["project"]; !ok {
		return nil, &InvalidGoogleCloudProject{}
	}

	if _, ok := linkedsecret.Spec.ProviderOptions["secret"]; !ok {
		return nil, &InvalidSecretOption{}
	}

	// get provider options informed in linkedsecret spec
	project := linkedsecret.Spec.ProviderOptions["project"]
	name := linkedsecret.Spec.ProviderOptions["secret"]
	// set default "latest" if providerOption version was not specified
	version := "latest"

	// get version if defined
	if _, ok := linkedsecret.Spec.ProviderOptions["version"]; ok {
		version = linkedsecret.Spec.ProviderOptions["version"]
	}

	// Secret name with its path
	secretPath := "projects/" + project + "/secrets/" + name + "/versions/" + version

	// Create the client.
	client, err := secretmanager.NewClient(ctx)

	if err != nil {
		log.V(1).Info("GCP Error creating client", name, err)
		return nil, err
	}
	// defer close connection
	defer client.Close()

	// Build a request
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: secretPath,
	}

	// Access secret
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		log.V(1).Info("GCP Error getting secret", name, err)
		return nil, err
	}

	//return error if payload data is empty
	if len(result.Payload.Data) == 0 {
		log.V(1).Info("GCP Error empty secret", name, err)
		return nil, &EmptySecret{name, "is empty"}

	}
	log.V(1).Info("GCP return secret", "secret", name)
	return result.Payload.Data, nil
}
