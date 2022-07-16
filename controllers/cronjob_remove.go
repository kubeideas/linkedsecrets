package controllers

import (
	"context"
	securityv1 "kubeideas/linkedsecrets/api/v1"

	"sigs.k8s.io/controller-runtime/pkg/log"
)

// RemoveCronJob remove cron job
func (r *LinkedSecretReconciler) removeCronJob(ctx context.Context, linkedsecret *securityv1.LinkedSecret) error {

	log := log.FromContext(ctx)

	// remove cron only if it exists on map
	if _, ok := r.Cronjob[linkedsecret.UID]; ok {
		r.Cronjob[linkedsecret.UID].Remove(linkedsecret.Status.CronJobID)
		r.Cronjob[linkedsecret.UID].Stop()
		delete(r.Cronjob, linkedsecret.UID)
		//debug info
		log.V(1).Info("Remove cronjob", "ID", linkedsecret.Status.CronJobID)
	} else {
		// debug info
		log.V(1).Info("Remove cronjob", "schedule", "no schedule to be removed")
		return nil
	}

	// update cronjob status
	linkedsecret.Status.CronJobStatus = JOBNOTSCHEDULED
	linkedsecret.Status.CronJobID = -1

	if err := r.Status().Update(ctx, linkedsecret); err != nil {
		return err
	}

	log.V(1).Info("Cronjob Removed", "CronJobStatus", linkedsecret.Status.CronJobStatus)

	//record deletion event
	r.Recorder.Event(linkedsecret, "Normal", "RemovedCronJob", linkedsecret.Status.CurrentSchedule)
	return nil
}
