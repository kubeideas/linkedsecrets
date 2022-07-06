package controllers

import (
	"context"
	"fmt"
	securityv1 "kubeideas/linkedsecrets/api/v1"

	"github.com/robfig/cron/v3"

	"sigs.k8s.io/controller-runtime/pkg/log"
)

// CronJobParser parse schedule format
// valid formats : "@every 5m", "*/5 * * * * *"
func (r *LinkedSecretReconciler) cronJobParser(ctx context.Context, linkedsecret *securityv1.LinkedSecret) error {

	log := log.FromContext(ctx)

	parser := cron.NewParser(
		cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)

	if _, err := parser.Parse(linkedsecret.Spec.Schedule); err != nil {
		r.Recorder.Event(linkedsecret, "Warning", "FailParsing", fmt.Sprintf("%s", err))
		// update cronjob status
		linkedsecret.Status.CronJobStatus = JOBFAILPARSESCHEDULE
		linkedsecret.Status.CronJobID = -1
		if err := r.Status().Update(ctx, linkedsecret); err != nil {
			log.V(1).Info("Cronjob update status", "error", err)
			return err
		}

		//debug info
		log.V(1).Info("Cronjob parse schedule", "schedule", linkedsecret.Status.CurrentSchedule)
		log.V(1).Info("Cronjob parse schedule", "jobId", linkedsecret.Status.CronJobID)
		log.V(1).Info("Cronjob parse schedule", "jobStatus", linkedsecret.Status.CronJobStatus)

		return err
	}

	return nil
}
