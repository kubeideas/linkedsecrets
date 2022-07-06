package controllers

import (
	"context"
	securityv1 "kubeideas/linkedsecrets/api/v1"

	"sigs.k8s.io/controller-runtime/pkg/log"
)

//CheckObjectDeletion check if linkedsecret is being delete and remove cronjob before it happens
func (r *LinkedSecretReconciler) CheckObjectDeletion(ctx context.Context, linkedsecret *securityv1.LinkedSecret) (stopReconcile bool, err error) {

	log := log.FromContext(ctx)

	// timestamp != 0 object is under deletion
	if !linkedsecret.ObjectMeta.DeletionTimestamp.IsZero() {

		log.V(1).Info("BEGIN DELETE LINKEDSECRET")

		if containString(linkedsecret.ObjectMeta.Finalizers, LINKEDSECRETFINALIZER) {

			// Remove Secret cronjob before deletion if there is a valid job id
			if err := r.removeCronJob(ctx, linkedsecret); err != nil {
				return true, err
			}

			// remove finalizer from the list and update it.
			linkedsecret.ObjectMeta.Finalizers = removeString(linkedsecret.ObjectMeta.Finalizers, LINKEDSECRETFINALIZER)
			if err := r.Update(context.Background(), linkedsecret); err != nil {
				log.V(1).Error(err, "Remove finalizer", "string", LINKEDSECRETFINALIZER)
				return true, err
			}
		}
		log.V(1).Info("END DELETE LINKEDSECRET")
		// Stop reconciliation as the item is being deleted
		return true, nil
	}

	return false, nil
}
