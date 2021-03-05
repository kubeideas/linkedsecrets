package controllers

import (
	"context"
	"fmt"
	securityv1 "linkedsecrets/api/v1"
	"time"

	"github.com/robfig/cron/v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CronJobParser parse schedule format
// valid formats : "@every 5m", "*/5 * * * * *"
func (r *LinkedSecretReconciler) CronJobParser(ctx context.Context, linkedsecret *securityv1.LinkedSecret) error {

	log := r.Log.WithValues("linkedsecret", fmt.Sprintf("%s/%s", linkedsecret.Namespace, linkedsecret.Name))

	parser := cron.NewParser(
		cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)

	if _, err := parser.Parse(linkedsecret.Spec.Schedule); err != nil {
		r.Recorder.Event(linkedsecret, "Warning", "FailParsing", fmt.Sprintf("%s", err))
		// update cronjob status
		linkedsecret.Status.CronJobStatus = JOBFAILPARSESCHEDULE
		linkedsecret.Status.CronJobID = -1
		linkedsecret.Status.CurrentSchedule = linkedsecret.Spec.Schedule
		if err := r.Status().Update(ctx, linkedsecret); err != nil {
			return err
		}
		//debug info
		log.V(1).Info("Parse schedule", "schedule", linkedsecret.Status.CurrentSchedule)
		log.V(1).Info("Parse schedule", "jobId", linkedsecret.Status.CronJobID)
		log.V(1).Info("Parse schedule", "jobStatus", linkedsecret.Status.CronJobStatus)

		return err
	}

	return nil
}

// CronUpdateJob is the job to be executed to get secret synchronized
func (r *LinkedSecretReconciler) CronUpdateJob(ctx context.Context, linkedsecret *securityv1.LinkedSecret) error {

	//Fetch cloud provider data and update secret
	if err := r.UpdateSecret(ctx, linkedsecret); err != nil {
		//set job status to Scheduled
		linkedsecret.Status.CronJobStatus = JOBSCHEDULED
		if err := r.Status().Update(ctx, linkedsecret); err != nil {
			return err
		}
		r.Recorder.Event(linkedsecret, "Warning", "FailSchedule", fmt.Sprintf("%s", err))
		return err
	}

	//set job status to Scheduled
	linkedsecret.Status.CronJobStatus = JOBSCHEDULED
	linkedsecret.Status.LastScheduleExecution = &metav1.Time{Time: time.Now()}

	if err := r.Status().Update(ctx, linkedsecret); err != nil {
		return err
	}

	return nil
}

// AddCronjob add new execution of CronUpdateJob
func (r *LinkedSecretReconciler) AddCronjob(ctx context.Context, linkedsecret *securityv1.LinkedSecret) error {

	log := r.Log.WithValues("linkedsecret", fmt.Sprintf("%s/%s", linkedsecret.Namespace, linkedsecret.Name))

	//parse schedule
	if err := r.CronJobParser(ctx, linkedsecret); err != nil {
		return err
	}

	// Add parse to support seconds and start cron
	parser := cron.NewParser(
		cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)
	cronjob := cron.New(cron.WithParser(parser))
	r.Cronjob[linkedsecret.UID] = cronjob
	r.Cronjob[linkedsecret.UID].Start()

	ID, err := r.Cronjob[linkedsecret.UID].AddFunc(linkedsecret.Spec.Schedule, func() {
		//fmt.Printf("\ntick [ %s ]\n", linkedsecret.Spec.Schedule)
		r.CronUpdateJob(ctx, linkedsecret)
	})

	if err != nil {
		return err
	}

	//set status for cronjob
	linkedsecret.Status.CronJobID = ID
	linkedsecret.Status.CronJobStatus = JOBSCHEDULED
	linkedsecret.Status.CurrentSchedule = linkedsecret.Spec.Schedule
	if err := r.Status().Update(ctx, linkedsecret); err != nil {
		return err
	}

	//debug info
	log.V(1).Info("Add cronjob", "schedule", linkedsecret.Status.CurrentSchedule)
	log.V(1).Info("Add cronjob", "jobId", linkedsecret.Status.CronJobID)
	log.V(1).Info("Add cronjob", "jobStatus", linkedsecret.Status.CronJobStatus)

	r.Recorder.Event(linkedsecret, "Normal", "CreatedSchedule", fmt.Sprintf("Secret %s/%s", linkedsecret.Status.CreatedSecretNamespace, linkedsecret.Status.CreatedSecret))

	return nil
}

// RemoveCronJob remove an execution of CronUpdateJob
func (r *LinkedSecretReconciler) RemoveCronJob(ctx context.Context, linkedsecret *securityv1.LinkedSecret) error {

	log := r.Log.WithValues("linkedsecret", fmt.Sprintf("%s/%s", linkedsecret.Namespace, linkedsecret.Name))

	// remove cron only if it exists on map
	if _, ok := r.Cronjob[linkedsecret.UID]; ok {
		r.Cronjob[linkedsecret.UID].Remove(linkedsecret.Status.CronJobID)
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

	// debug info
	log.V(1).Info("Remove cronjob", "schedule", linkedsecret.Status.CurrentSchedule)
	log.V(1).Info("Remove cronjob", "jobId", linkedsecret.Status.CronJobID)
	log.V(1).Info("Remove cronjob", "jobStatus", linkedsecret.Status.CronJobStatus)
	log.V(1).Info("Remove cronjob", "cronEntries", len(r.Cronjob[linkedsecret.UID].Entries()))

	//record deletion event
	r.Recorder.Event(linkedsecret, "Normal", "RemovedSchedule", fmt.Sprintf(" %s/%s", linkedsecret.Status.CreatedSecretNamespace, linkedsecret.Status.CreatedSecret))
	return nil
}
