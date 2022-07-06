package controllers

import (
	"context"
	"errors"
	securityv1 "kubeideas/linkedsecrets/api/v1"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

//GetCloudSecret get Cloud Secret data
func (r *LinkedSecretReconciler) GetCloudSecret(ctx context.Context, linkedsecret *securityv1.LinkedSecret) (corev1.Secret, error) {

	log := log.FromContext(ctx)

	// Default secret type
	var secretType corev1.SecretType = "Opaque"

	var err error
	data := []byte{}
	secret := corev1.Secret{}
	var secretMap map[string][]byte

	// Fake test
	fakeData := func(linkedsecret *securityv1.LinkedSecret) ([]byte, error) {

		secretJSON := `{"user": "user01","password": "pwd12345"}`
		secretPlain := []byte("user=user02  \n password = pwd78910")

		secretDockerPlain := []byte("docker-server = https://mydockersever.example \n\n\ndocker-username = tiger \n\ndocker-password = pass113 \n\n docker-email = tiger@acme.com")
		secretDockerJSON := `{"docker-username": "tiger", "docker-password": "pass113", "docker-email": "tiger@acme.com"}`

		if linkedsecret.Spec.ProviderOptions["secretId"] == "ff000000-f000-f000-f000-ffffffffff00" {
			return nil, errors.New("secret not found")
		}

		if linkedsecret.Spec.ProviderOptions["secret"] == "secret-not-found" {
			return nil, errors.New("secret not found")
		}

		if linkedsecret.Spec.ProviderOptions["secret"] == "docker-secret-json" {
			return []byte(secretDockerJSON), nil
		}

		if linkedsecret.Spec.ProviderOptions["secret"] == "docker-secret-plain" {
			return secretDockerPlain, nil
		}

		if linkedsecret.Spec.ProviderSecretFormat == "JSON" {
			return []byte(secretJSON), nil
		} else {
			return secretPlain, nil
		}

	}

	//retrieve Cloud secret data
	switch linkedsecret.Spec.Provider {
	case GOOGLE:
		//data, err = r.GetGCPSecret(ctx, linkedsecret)
		data, err = fakeData(linkedsecret)
	case AWS:
		//data, err = r.GetAWSSecret(ctx, linkedsecret)
		data, err = fakeData(linkedsecret)
	case AZURE:
		//data, err = r.GetAzureSecret(ctx, linkedsecret)
		data, err = fakeData(linkedsecret)
	case IBM:
		//data, err = r.GetIBMSecret(ctx, linkedsecret)
		data, err = fakeData(linkedsecret)
	}

	//return error retrieving Cloud secret data
	if err != nil {
		return corev1.Secret{}, err
	}

	// create key/value map based on data format
	if linkedsecret.Spec.ProviderSecretFormat == JSONFORMAT {
		secretMap, err = parseJSON(data)
	} else {
		secretMap, err = parsePlainData(data)
	}

	if err != nil {
		return secret, err
	}

	// infer docker secret
	if dockerSecret, ok := inferDockerConfig(secretMap); ok {
		secretMap = dockerSecret
		secretType = "kubernetes.io/dockerconfigjson"
		log.V(1).Info("Docker Config inferred", "secret type", secretType)
	}

	// create new secret object and add data
	secret = corev1.Secret{
		TypeMeta: v1.TypeMeta{
			Kind:       "Secret",
			APIVersion: corev1.SchemeGroupVersion.String(),
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      linkedsecret.Spec.SecretName,
			Namespace: linkedsecret.Namespace,
		},
		Data: secretMap,
		Type: secretType,
	}

	return secret, nil

}
