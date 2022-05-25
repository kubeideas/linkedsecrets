package controllers

import (
	"context"
	"fmt"
	securityv1 "kubeideas/linkedsecrets/api/v1"
	"reflect"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Check if current secret data is different from Cloud secret data
func (r *LinkedSecretReconciler) checkDiff(ctx context.Context, linkedsecret *securityv1.LinkedSecret, newSecret corev1.Secret) bool {

	log := r.Log.WithValues("linkedsecret", fmt.Sprintf("%s/%s", linkedsecret.Namespace, linkedsecret.Name))

	//Get current secret
	var currSecret corev1.Secret
	secretName := client.ObjectKey{Namespace: linkedsecret.Namespace, Name: linkedsecret.Spec.SecretName}
	if err := r.Get(ctx, secretName, &currSecret); err != nil {
		log.V(1).Info("checkSecretDiff - Secret not found", "secret", fmt.Sprintf("%s/%s", linkedsecret.Namespace, linkedsecret.Spec.SecretName))
		return false
	}

	// Compare both secret data
	isEqual := reflect.DeepEqual(currSecret.Data, newSecret.Data)

	log.V(1).Info("checkSecretDiff succeed", "IsEqual", isEqual)

	return isEqual
}
