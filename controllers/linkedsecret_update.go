package controllers

import (
	"context"
	"fmt"
	securityv1 "linkedsecrets/api/v1"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// UpdateLinkedSecret apply any change made on linkedsecret and synchronize  kubernetes secret
func (r *LinkedSecretReconciler) UpdateLinkedSecret(ctx context.Context, linkedsecret *securityv1.LinkedSecret) error {

	log := r.Log.WithValues("linkedsecret", fmt.Sprintf("%s/%s", linkedsecret.Namespace, linkedsecret.Name))

	// Remove cronjob if schedule or keepSecretOnDelete were changed or synchronization was suspended
	//if linkedsecret.Status.CurrentSchedule != linkedsecret.Spec.Schedule ||
	//	linkedsecret.Spec.Suspended ||
	//	linkedsecret.Status.KeepSecretOnDelete != linkedsecret.Spec.KeepSecretOnDelete {
	//	if err := r.RemoveCronJob(ctx, linkedsecret); err != nil {
	//		return err
	//	}
	//}

	// Remove cronjob
	if err := r.RemoveCronJob(ctx, linkedsecret); err != nil {
		return err
	}

	// update secret with provider data
	secret, err := r.GetProviderSecret(linkedsecret)

	if err != nil {
		r.Recorder.Event(linkedsecret, "Warning", "FailSynching", err.Error())
		linkedsecret.Status.CurrentSecretStatus = STATUSNOTSYNCHED
		if err := r.Status().Update(ctx, linkedsecret); err != nil {
			return err
		}
		return err
	}
	log.V(1).Info("Synchronize secret data on update", "secret", fmt.Sprintf("Secret %s/%s", secret.Namespace, secret.Name))

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

	// Deployment rollout update
	if !isEqual {
		r.rolloutUpdate(ctx, linkedsecret)
	}

	// Suspend cronjob
	if linkedsecret.Spec.Suspended {
		//set cronjob suspended
		linkedsecret.Status.CronJobStatus = JOBSUSPENDED
		r.Recorder.Event(linkedsecret, "Warning", "Cronjob suspended", linkedsecret.Name)
	}

	// Add cronjob
	if !linkedsecret.Spec.Suspended {
		if err := r.AddCronjob(ctx, linkedsecret); err != nil {
			return err
		}
	}

	// update linkedsecret status
	linkedsecret.Status.CurrentSecretStatus = STATUSSYNCHED
	linkedsecret.Status.CreatedSecret = secret.Name
	linkedsecret.Status.CreatedSecretNamespace = secret.Namespace
	linkedsecret.Status.CurrentProvider = linkedsecret.Spec.Provider
	linkedsecret.Status.CurrentProviderOptions = linkedsecret.Spec.ProviderOptions
	linkedsecret.Status.CurrentSchedule = linkedsecret.Spec.Schedule
	linkedsecret.Status.KeepSecretOnDelete = linkedsecret.Spec.KeepSecretOnDelete

	if err := r.Status().Update(ctx, linkedsecret); err != nil {
		return err
	}

	//debug info
	log.V(1).Info("Update linkedsecret", "CurrentSecretStatus", linkedsecret.Status.CurrentSecretStatus)
	log.V(1).Info("Update linkedsecret", "CreatedSecret", linkedsecret.Status.CreatedSecret)
	log.V(1).Info("Update linkedsecret", "CreatedSecretNamespace", linkedsecret.Status.CreatedSecretNamespace)
	log.V(1).Info("Update linkedsecret", "CurrentProvider", linkedsecret.Status.CurrentProvider)
	log.V(1).Info("Update linkedsecret", "CurrentProviderOptions", linkedsecret.Status.CurrentProviderOptions)
	log.V(1).Info("Update linkedsecret", "KeepSecretOnDelete", linkedsecret.Status.KeepSecretOnDelete)
	log.V(1).Info("Update linkedsecret", "CurrentSchedule", linkedsecret.Status.CurrentSchedule)
	log.V(1).Info("Update linkedsecret", "CronJobStatus", linkedsecret.Status.CronJobStatus)
	log.V(1).Info("Update linkedsecret", "Cronjob map", len(r.Cronjob))

	// Record linkedsecret updated
	r.Recorder.Event(linkedsecret, "Normal", "Updated", linkedsecret.Name)

	return nil
}
