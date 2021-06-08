package controllers

import (
	"context"
	"fmt"
	securityv1 "linkedsecrets/api/v1"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

//GetGCPSecret return secret data from Google Secret Manager
func (r *LinkedSecretReconciler) GetGCPSecret(linkedsecret *securityv1.LinkedSecret) ([]byte, error) {

	log := r.Log.WithValues("linkedsecret", fmt.Sprintf("%s/%s", linkedsecret.Namespace, linkedsecret.Name))

	// get provider options informed in linkedsecret spec
	project := linkedsecret.Spec.ProviderOptions["project"]
	name := linkedsecret.Spec.ProviderOptions["secret"]
	version := linkedsecret.Spec.ProviderOptions["version"]

	// Secret name with its path
	secretPath := "projects/" + project + "/secrets/" + name + "/versions/" + version

	// Create the client.
	ctx := context.Background()
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
