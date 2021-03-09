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

	// always set the controller reference so that we know which object owns this.
	if err := ctrl.SetControllerReference(linkedsecret, &secret, r.Scheme); err != nil {
		return err
	}

	// update existent secret
	updateOpts := &client.UpdateOptions{}
	if err := r.Update(ctx, &secret, updateOpts); err != nil {
		r.Recorder.Event(linkedsecret, "Warning", "FailUpdating", err.Error())
		linkedsecret.Status.CurrentSecretStatus = STATUSNOTSYNCHED
		if err := r.Status().Update(ctx, linkedsecret); err != nil {
			return err
		}
		return err
	}

	log.V(1).Info("Synchronized secret data", "secret", fmt.Sprintf("%s/%s", secret.Namespace, secret.Name))

	return nil
}
