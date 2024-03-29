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

// UpdateLinkedSecret apply any change made on linkedsecret and synchronize  kubernetes secret
func (r *LinkedSecretReconciler) UpdateLinkedSecret(ctx context.Context, linkedsecret *securityv1.LinkedSecret) error {

	log := log.FromContext(ctx)

	log.V(1).Info("BEGIN RECONCILING - UPDATE", "linkedsecret", fmt.Sprintf("%s/%s", linkedsecret.Namespace, linkedsecret.Name))

	// Remove cronjob
	if err := r.removeCronJob(ctx, linkedsecret); err != nil {
		return err
	}

	// Get cloud secret data
	secret, err := r.GetCloudSecret(ctx, linkedsecret)

	if err != nil {
		r.Recorder.Event(linkedsecret, "Warning", "FailSynchSecret", err.Error())
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

	log.V(1).Info("Secret updated", "secret", fmt.Sprintf("Secret %s/%s", secret.Namespace, secret.Name))

	// Deployment rollout update
	if !isEqual {
		r.rolloutUpdateDeployment(ctx, linkedsecret)
	}

	// update linkedsecret status
	linkedsecret.Status.CurrentSecretStatus = STATUSSYNCHED
	linkedsecret.Status.CurrentSchedule = linkedsecret.Spec.Schedule
	linkedsecret.Status.ObservedGeneration = linkedsecret.Generation
	linkedsecret.Status.CurrentSecret = linkedsecret.Spec.SecretName
	linkedsecret.Status.LastScheduleExecution = &metav1.Time{Time: time.Now()}

	// Suspend cronjob
	if linkedsecret.Spec.Suspended {
		//set cronjob suspended
		linkedsecret.Status.CronJobStatus = JOBSUSPENDED
		linkedsecret.Status.NextScheduleExecution = nil

		if err := r.Status().Update(ctx, linkedsecret); err != nil {
			return err
		}
		r.Recorder.Event(linkedsecret, "Warning", "Cronjob suspended", linkedsecret.Name)
	}

	// Add cronjob
	if linkedsecret.Spec.Schedule != "" && !linkedsecret.Spec.Suspended {
		if err := r.addCronjob(ctx, linkedsecret); err != nil {
			return err
		}
	}

	//debug info
	log.V(1).Info("Linkedsecret updated", "ObservedGeneration", linkedsecret.Status.ObservedGeneration)

	// Record linkedsecret updated
	r.Recorder.Event(linkedsecret, "Normal", "Updated", linkedsecret.Name)

	log.V(1).Info("END RECONCILING - UPDATE", "linkedsecret", fmt.Sprintf("%s/%s", linkedsecret.Namespace, linkedsecret.Name))

	return nil
}
