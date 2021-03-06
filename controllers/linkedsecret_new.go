package controllers

import (
	"context"
	"fmt"
	securityv1 "linkedsecrets/api/v1"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NewLinkedSecret create new kubernetes secret, fetch data from cloud secret manager and add cronjob to get it synchronized autimatically
func (r *LinkedSecretReconciler) NewLinkedSecret(ctx context.Context, linkedsecret *securityv1.LinkedSecret) error {

	log := r.Log.WithValues("linkedsecret", fmt.Sprintf("%s/%s", linkedsecret.Namespace, linkedsecret.Name))

	// remove schedule job if secret name was changed
	if linkedsecret.Status.CronJobID > 0 {
		if err := r.RemoveCronJob(ctx, linkedsecret); err != nil {
			return err
		}
	}

	// update provider status
	linkedsecret.Status.CurrentProvider = linkedsecret.Spec.Provider
	linkedsecret.Status.CurrentProviderOptions = linkedsecret.Spec.ProviderOptions

	// create secret with providerdata
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

	createOptions := &client.CreateOptions{FieldManager: linkedsecret.Name}

	if err = r.Create(ctx, &secret, createOptions); err != nil {
		r.Recorder.Event(linkedsecret, "Warning", "FailCreating", err.Error())
		linkedsecret.Status.CurrentSecretStatus = STATUSNOTSYNCHED

		if err := r.Status().Update(ctx, linkedsecret); err != nil {
			return err
		}
		return err
	}

	log.V(1).Info("Synchronize secret data on creation", "secret", fmt.Sprintf("%s/%s", secret.Namespace, secret.Name))

	// set status for linkedsecret without synchronization
	if linkedsecret.Spec.Suspended {
		linkedsecret.Status.CronJobStatus = JOBSUSPENDED
	}

	// update linkedsecret status
	linkedsecret.Status.CurrentSecretStatus = STATUSSYNCHED
	linkedsecret.Status.CreatedSecret = secret.Name
	linkedsecret.Status.CreatedSecretNamespace = secret.Namespace
	linkedsecret.Status.LastScheduleExecution = &metav1.Time{Time: time.Now()}
	linkedsecret.Status.KeepSecretOnDelete = linkedsecret.Spec.KeepSecretOnDelete
	linkedsecret.Status.CurrentSchedule = linkedsecret.Spec.Schedule

	if err := r.Status().Update(ctx, linkedsecret); err != nil {
		return err
	}

	// Record secret created
	r.Recorder.Event(linkedsecret, "Normal", "Created", fmt.Sprintf("Secret %s/%s", secret.Namespace, secret.Name))

	// create secret cronjob if a schedule was defined
	if linkedsecret.Spec.Schedule != "" && !linkedsecret.Spec.Suspended {
		if err := r.AddCronjob(ctx, linkedsecret); err != nil {
			return err
		}
	}

	return nil
}
