package controllers

import (
	"context"
	"fmt"
	securityv1 "kubeideas/linkedsecrets/api/v1"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Apply deployment rollout update in order to update environment variables with new secret data
func (r *LinkedSecretReconciler) rolloutUpdate(ctx context.Context, linkedsecret *securityv1.LinkedSecret) {

	log := r.Log.WithValues("linkedsecret", fmt.Sprintf("%s/%s", linkedsecret.Namespace, linkedsecret.Name))

	if linkedsecret.Spec.Deployment == "" {
		log.V(1).Info("Rollout update ignored", "deployment", "not defined")

	}

	// Get deploy to update
	var lnsDeployment appsv1.Deployment

	deployName := client.ObjectKey{Namespace: linkedsecret.Namespace, Name: linkedsecret.Spec.Deployment}

	if err := r.Get(ctx, deployName, &lnsDeployment); err != nil {
		log.V(1).Info("Rollout update ignored", "Deployment", fmt.Sprintf("%s/%s not found", linkedsecret.Namespace, linkedsecret.Spec.Deployment))
	}

	// Create restart annotation
	// Restart will happen in 5 seconds from now
	annotations := map[string]string{"kubectl.kubernetes.io/restartedAt": time.Now().Add(time.Second * time.Duration(5)).Format(time.RFC3339)}

	// Add annotation to restart deployment
	lnsDeployment.Spec.Template.ObjectMeta.Annotations = annotations

	// update deployment
	updateOpts := &client.UpdateOptions{}
	if err := r.Update(ctx, &lnsDeployment, updateOpts); err != nil {
		log.Info("Error Adding annotations to deployment", "Error", err)
	}

	log.V(1).Info("Rollout update succeed", "Deployment", fmt.Sprintf("%s/%s", lnsDeployment.Namespace, lnsDeployment.Name))

}
