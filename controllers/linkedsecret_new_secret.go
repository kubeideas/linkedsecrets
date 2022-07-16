package controllers

import (
	"context"
	"fmt"
	securityv1 "kubeideas/linkedsecrets/api/v1"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// NewSecret create new kubernetes secret, fetch data from cloud secret manager and add cronjob to get it synchronized autimatically
func (r *LinkedSecretReconciler) NewSecret(ctx context.Context, linkedsecret *securityv1.LinkedSecret) error {

	log := log.FromContext(ctx)

	log.V(1).Info("BEGIN ADD NEW SECRET", "secret", fmt.Sprintf("%s/%s", linkedsecret.Namespace, linkedsecret.Spec.SecretName))

	// Remove cronjob
	//if err := r.removeCronJob(ctx, linkedsecret); err != nil {
	//	return err
	//}

	// create kubernetes secret  with cloud secret
	secret, err := r.GetCloudSecret(ctx, linkedsecret)

	if err != nil {
		r.Recorder.Event(linkedsecret, "Warning", "FailSynching", err.Error())
		linkedsecret.Status.CurrentSecretStatus = STATUSNOTSYNCHED
		linkedsecret.Status.NextScheduleExecution = nil
		linkedsecret.Status.CronJobStatus = JOBNOTSCHEDULED
		linkedsecret.Status.CurrentSchedule = linkedsecret.Spec.Schedule
		linkedsecret.Status.ObservedGeneration = linkedsecret.Generation
		linkedsecret.Status.CurrentSecret = linkedsecret.Spec.SecretName
		if err := r.Status().Update(ctx, linkedsecret); err != nil {
			return err
		}
		return err
	}

	// if false secret will be deleted with along with linkedsecret.
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

	log.V(1).Info("Secret created", "secret", fmt.Sprintf("%s/%s", secret.Namespace, secret.Name))

	// set status
	linkedsecret.Status.CurrentSecretStatus = STATUSSYNCHED
	linkedsecret.Status.LastScheduleExecution = &metav1.Time{Time: time.Now()}
	linkedsecret.Status.CurrentSchedule = linkedsecret.Spec.Schedule
	linkedsecret.Status.ObservedGeneration = linkedsecret.Generation
	linkedsecret.Status.CurrentSecret = linkedsecret.Spec.SecretName

	// set status for linkedsecret without synchronization
	if linkedsecret.Spec.Suspended {
		linkedsecret.Status.CronJobStatus = JOBSUSPENDED
		linkedsecret.Status.NextScheduleExecution = nil

		//update linkedsecret
		if err := r.Status().Update(ctx, linkedsecret); err != nil {
			return err
		}
		r.Recorder.Event(linkedsecret, "Warning", "Sync suspended", linkedsecret.Name)
	}

	// create secret cronjob if a schedule was defined
	if linkedsecret.Spec.Schedule != "" && !linkedsecret.Spec.Suspended {
		if err := r.addCronjob(ctx, linkedsecret); err != nil {
			return err
		}
	}

	// Record secret created
	r.Recorder.Event(linkedsecret, "Normal", "Created", fmt.Sprintf("Secret %s/%s", secret.Namespace, secret.Name))

	log.V(1).Info("END ADD NEW SECRET", "secret", fmt.Sprintf("%s/%s", linkedsecret.Namespace, linkedsecret.Spec.SecretName))

	return nil
}
