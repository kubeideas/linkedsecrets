package controllers

import (
	"context"
	"fmt"
	securityv1 "linkedsecrets/api/v1"
)

//LINKEDSECRETFINALIZER identify linkendsecret to be intercept before delete
const LINKEDSECRETFINALIZER = "cronjob.finalizers.linkedsecrets.kubeidea.io"

// AddFinalizer add finalized to linkedsecret if it does not have one
func (r *LinkedSecretReconciler) AddFinalizer(linkedsecret *securityv1.LinkedSecret) error {

	log := r.Log.WithValues("linkedsecret", fmt.Sprintf("%s/%s", linkedsecret.Namespace, linkedsecret.Name))

	//timestamp == 0 object is not being delete
	if linkedsecret.ObjectMeta.DeletionTimestamp.IsZero() {

		if !containsString(linkedsecret.ObjectMeta.Finalizers, LINKEDSECRETFINALIZER) {

			log.V(1).Info("BEGIN RECONCILING - FINILIZER LINKEDSECRET")

			linkedsecret.ObjectMeta.Finalizers = append(linkedsecret.ObjectMeta.Finalizers, LINKEDSECRETFINALIZER)
			if err := r.Update(context.Background(), linkedsecret); err != nil {
				return err
			}
			log.V(1).Info("Append Finalizer", "linkedsecret-finalizer", LINKEDSECRETFINALIZER)

			log.V(1).Info("END RECONCILING - FINILIZER LINKEDSECRET")
		}

	}
	return nil
}

//CheckObjectDeletion check if linkedsecret is being delete and remove cronjob before it happens
func (r *LinkedSecretReconciler) CheckObjectDeletion(ctx context.Context, linkedsecret *securityv1.LinkedSecret) (stopReconcile bool, err error) {

	log := r.Log.WithValues("linkedsecret", fmt.Sprintf("%s/%s", linkedsecret.Namespace, linkedsecret.Name))

	// timestamp != 0 object is under deletion
	if !linkedsecret.ObjectMeta.DeletionTimestamp.IsZero() {

		log.V(1).Info("BEGIN RECONCILING - DELETE LINKEDSECRET")

		if containsString(linkedsecret.ObjectMeta.Finalizers, LINKEDSECRETFINALIZER) {

			// Remove Secret cronjob before deletion if there is a valid job id
			if linkedsecret.Status.CronJobID > 0 {
				if err := r.RemoveCronJob(ctx, linkedsecret); err != nil {
					return true, err
				}
			}

			// remove finalizer from the list and update it.
			linkedsecret.ObjectMeta.Finalizers = removeString(linkedsecret.ObjectMeta.Finalizers, LINKEDSECRETFINALIZER)
			if err := r.Update(context.Background(), linkedsecret); err != nil {
				return true, err
			}
			log.V(1).Info("Remove Finalizer", "linkedsecret-finalizer", LINKEDSECRETFINALIZER)
		}
		log.V(1).Info("END RECONCILING - DELETE LINKEDSECRET")
		// Stop reconciliation as the item is being deleted
		return true, nil
	}

	return false, nil
}

// Helper functions to check and remove string from a slice of strings.
func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

// func removeString(slice []string, s string) (result []string) {
// 	for _, item := range slice {
// 		if item == s {
// 			continue
// 		}
// 		result = append(result, item)
// 	}
// 	return
// }

func removeString(slice []string, s string) []string {
	for i, item := range slice {
		if item == s {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}
