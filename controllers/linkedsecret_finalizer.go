package controllers

import (
	"context"
	securityv1 "kubeideas/linkedsecrets/api/v1"

	"sigs.k8s.io/controller-runtime/pkg/log"
)

// AddFinalizer add finalized to linkedsecret if it does not have one
func (r *LinkedSecretReconciler) AddFinalizer(ctx context.Context, linkedsecret *securityv1.LinkedSecret) (err error) {
	log := log.FromContext(ctx)

	//timestamp == 0 object is not being delete
	if linkedsecret.ObjectMeta.DeletionTimestamp.IsZero() {

		if !containString(linkedsecret.ObjectMeta.Finalizers, LINKEDSECRETFINALIZER) {

			log.V(1).Info("BEGIN ADD FINALIZER")

			linkedsecret.ObjectMeta.Finalizers = append(linkedsecret.ObjectMeta.Finalizers, LINKEDSECRETFINALIZER)
			if err := r.Update(context.Background(), linkedsecret); err != nil {
				return err
			}
			log.V(1).Info("Append Finalizer", "linkedsecret-finalizer", LINKEDSECRETFINALIZER)

			log.V(1).Info("END ADD FINALIZER")
			return nil
		}

	}
	return nil
}
