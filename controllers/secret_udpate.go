package controllers

import (
	"context"
	"fmt"
	securityv1 "linkedsecrets/api/v1"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// UpdateSecret synchronize kubernetes secret
func (r *LinkedSecretReconciler) UpdateSecret(ctx context.Context, linkedsecret *securityv1.LinkedSecret) error {

	log := r.Log.WithValues("linkedsecret", fmt.Sprintf("%s/%s", linkedsecret.Namespace, linkedsecret.Name))

	secret, err := r.GetProviderSecret(linkedsecret)

	if err != nil {
		r.Recorder.Event(linkedsecret, "Warning", "FailSynching", err.Error())
		linkedsecret.Status.CurrentSecretStatus = STATUSNOTSYNCHED
		if err := r.Status().Update(ctx, linkedsecret); err != nil {
			return err
		}
		return err
	}

	// Set the controller reference so that we know which object owns this.
	// Secret will be deleted when Linkedsecret is deleted.
	if !linkedsecret.Spec.KeepSecretOnDelete {
		if err := ctrl.SetControllerReference(linkedsecret, &secret, r.Scheme); err != nil {
			return err
		}
	}

	// check secret data changes
	isEqual := r.checkDiff(ctx, linkedsecret, secret)

	// update existent secret
	updateOpts := &client.UpdateOptions{FieldManager: linkedsecret.Name}
	if err := r.Update(ctx, &secret, updateOpts); err != nil {
		r.Recorder.Event(linkedsecret, "Warning", "FailUpdating", err.Error())
		linkedsecret.Status.CurrentSecretStatus = STATUSNOTSYNCHED
		if err := r.Status().Update(ctx, linkedsecret); err != nil {
			return err
		}
		return err
	}

	log.V(1).Info("Synchronized secret data on schedule", "secret", fmt.Sprintf("%s/%s", secret.Namespace, secret.Name))

	// Deployment rollout update
	if !isEqual {
		r.rolloutUpdate(ctx, linkedsecret)
	}

	return nil
}
