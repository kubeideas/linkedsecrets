package controllers

import (
	"context"
	securityv1 "kubeideas/linkedsecrets/api/v1"

	"github.com/robfig/cron/v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// AddCronjob add new execution of CronUpdateJob
func (r *LinkedSecretReconciler) addCronjob(ctx context.Context, linkedsecret *securityv1.LinkedSecret) error {

	log := log.FromContext(ctx)

	//parse schedule
	if err := r.cronJobParser(ctx, linkedsecret); err != nil {
		return err
	}

	// Add parse to support seconds and start cron if necessary
	if _, ok := r.Cronjob[linkedsecret.UID]; !ok {
		cronOptions := cron.WithParser(cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor))
		cronjob := cron.New(cronOptions)
		r.Cronjob[linkedsecret.UID] = cronjob
	}

	ID, err := r.Cronjob[linkedsecret.UID].AddFunc(linkedsecret.Spec.Schedule, func() {
		r.cronUpdateSecretJob(ctx, linkedsecret)
	})

	if err != nil {
		return err
	}

	// start cron
	r.Cronjob[linkedsecret.UID].Start()

	//set status
	linkedsecret.Status.CronJobID = ID
	linkedsecret.Status.CronJobStatus = JOBSCHEDULED
	linkedsecret.Status.NextScheduleExecution = &metav1.Time{Time: r.Cronjob[linkedsecret.UID].Entry(linkedsecret.Status.CronJobID).Next}

	// update linkedsecret status
	if err := r.Status().Update(ctx, linkedsecret); err != nil {
		return err
	}

	//debug info
	log.V(1).Info("Cronjob added", "jobId", linkedsecret.Status.CronJobID)
	log.V(1).Info("Cronjob added", "jobStatus", linkedsecret.Status.CronJobStatus)

	//r.Recorder.Event(linkedsecret, "Normal", "CreatedSchedule", fmt.Sprintf("Secret %s/%s", linkedsecret.Namespace, linkedsecret.Spec.SecretName))
	r.Recorder.Event(linkedsecret, "Normal", "CreatedCronJob", linkedsecret.Status.CurrentSchedule)

	return nil
}
