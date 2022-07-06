package controllers

import (
	"context"
	"fmt"
	securityv1 "kubeideas/linkedsecrets/api/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// CronUpdateJob is the job to be executed to get secret synchronized
func (r *LinkedSecretReconciler) cronUpdateSecretJob(ctx context.Context, linkedsecret *securityv1.LinkedSecret) error {

	log := log.FromContext(ctx)

	// Get cloud secret data
	secret, err := r.GetCloudSecret(ctx, linkedsecret)
	if err != nil {
		r.Recorder.Event(linkedsecret, "Warning", "FailSynchSecret", err.Error())
		linkedsecret.Status.CurrentSecretStatus = STATUSNOTSYNCHED
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

	log.V(1).Info("Secret updated by Cronjob", "secret", fmt.Sprintf("Secret %s/%s", secret.Namespace, secret.Name))

	// Deployment rollout update
	if !isEqual {
		r.rolloutUpdateDeployment(ctx, linkedsecret)
	}

	//set job status to Scheduled
	linkedsecret.Status.CronJobStatus = JOBSCHEDULED
	linkedsecret.Status.LastScheduleExecution = &metav1.Time{Time: r.Cronjob[linkedsecret.UID].Entry(linkedsecret.Status.CronJobID).Prev}
	linkedsecret.Status.NextScheduleExecution = &metav1.Time{Time: r.Cronjob[linkedsecret.UID].Entry(linkedsecret.Status.CronJobID).Next}

	if err := r.Status().Update(ctx, linkedsecret); err != nil {
		return err
	}

	return nil
}
