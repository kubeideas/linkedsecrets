package controllers

import (
	"context"
	securityv1 "linkedsecrets/api/v1"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

//GetGCPSecret return secret data from Google Secret Manager
func GetGCPSecret(linkedsecret *securityv1.LinkedSecret) ([]byte, error) {

	// get options informed in linkedsecret spec
	project := linkedsecret.Spec.ProviderOptions["project"]
	name := linkedsecret.Spec.ProviderOptions["secret"]
	version := linkedsecret.Spec.ProviderOptions["version"]

	// Secret name with its path
	secretPath := "projects/" + project + "/secrets/" + name + "/versions/" + version

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	defer client.Close()

	if err != nil {
		return nil, err
	}

	// Build a request
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: secretPath,
	}

	// Access secret
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return nil, err
	}

	//return error if payload data is empty
	if len(result.Payload.Data) == 0 {
		return nil, &EmptySecret{name, "is empty"}

	}

	return result.Payload.Data, nil
}
