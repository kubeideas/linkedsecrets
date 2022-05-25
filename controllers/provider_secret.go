package controllers

import (
	"fmt"
	securityv1 "kubeideas/linkedsecrets/api/v1"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//GetProviderSecret access provider and return filled secret object
func (r *LinkedSecretReconciler) GetProviderSecret(linkedsecret *securityv1.LinkedSecret) (corev1.Secret, error) {

	log := r.Log.WithValues("linkedsecret", fmt.Sprintf("%s/%s", linkedsecret.Namespace, linkedsecret.Name))

	// Default secret type
	var secretType corev1.SecretType = "Opaque"

	var err error
	data := []byte{}
	secret := corev1.Secret{}
	var secretMap map[string][]byte

	//retrieve Cloud secret data
	switch linkedsecret.Spec.Provider {
	case GOOGLE:
		data, err = r.GetGCPSecret(linkedsecret)
	case AWS:
		data, err = r.GetAWSSecret(linkedsecret)
	case AZURE:
		data, err = r.GetAzureSecret(linkedsecret)
	case IBM:
		data, err = r.GetIBMSecret(linkedsecret)
	}

	//return error retrieving Cloud secret data
	if err != nil {
		return corev1.Secret{}, err
	}

	// create key/value map based on choosen format
	if linkedsecret.Spec.ProviderDataFormat == JSONFORMAT {
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
