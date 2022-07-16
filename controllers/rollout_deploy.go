package controllers

import (
	"context"
	"fmt"
	securityv1 "kubeideas/linkedsecrets/api/v1"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// Apply deployment rollout update in order to update environment variables with new secret data
func (r *LinkedSecretReconciler) rolloutUpdateDeployment(ctx context.Context, linkedsecret *securityv1.LinkedSecret) error {

	log := log.FromContext(ctx)

	if linkedsecret.Spec.RolloutRestartDeploy == "" {
		log.V(1).Info("Rollout update ignored", "deployment", "not defined")
		return &DeploymentNotDefined{}
	}

	// Get deploy to update
	var lnsDeployment appsv1.Deployment

	deployName := client.ObjectKey{Namespace: linkedsecret.Namespace, Name: linkedsecret.Spec.RolloutRestartDeploy}

	if err := r.Get(ctx, deployName, &lnsDeployment); err != nil {
		log.V(1).Info("Rollout update ignored", "Error", &DeploymentNotFound{name: linkedsecret.Spec.RolloutRestartDeploy})
		return &DeploymentNotFound{name: linkedsecret.Spec.RolloutRestartDeploy}
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
		return err
	}

	log.V(1).Info("Restart annotation added", "Deployment", fmt.Sprintf("%s/%s", lnsDeployment.Namespace, lnsDeployment.Name))

	return nil

}
