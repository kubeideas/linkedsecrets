/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	"github.com/robfig/cron/v3"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	securityv1 "kubeideas/linkedsecrets/api/v1"
)

// LinkedSecretReconciler reconciles a LinkedSecret object
type LinkedSecretReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
	Cronjob  map[types.UID]*cron.Cron
}

//+kubebuilder:rbac:groups=security.kubeideas.io,resources=linkedsecrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=security.kubeideas.io,resources=linkedsecrets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=security.kubeideas.io,resources=linkedsecrets/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;update;patch
//+kubebuilder:rbac:groups="",resources=events,verbs=create;patch
//+kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete

// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.2/pkg/reconcile
func (r *LinkedSecretReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	log := log.FromContext(ctx)

	var linkedsecret securityv1.LinkedSecret

	if err := r.Get(ctx, req.NamespacedName, &linkedsecret); err != nil {
		log.Info("linkedsecret does not exists", "linkedsecret", err)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// if linkedsecret is not being deleted add finalizer
	if err := r.AddFinalizer(ctx, &linkedsecret); err != nil {
		return ctrl.Result{}, err
	}

	// if linkedsecret is being deleted cronjob entry must be deleted first
	if stopReconcile, err := r.CheckObjectDeletion(ctx, &linkedsecret); stopReconcile && err != nil {
		return ctrl.Result{}, err
	} else if stopReconcile && err == nil {
		return ctrl.Result{}, nil
	}

	// Get kubernetes secret
	var secret corev1.Secret
	secretName := client.ObjectKey{Namespace: linkedsecret.Namespace, Name: linkedsecret.Spec.SecretName}

	// create new secret if new linkedsecret is being created or if secretName was changed.
	if err := r.Get(ctx, secretName, &secret); err != nil {

		// create new kubernetes secret with cloud secret
		if err := r.NewSecret(ctx, &linkedsecret); err != nil {
			return ctrl.Result{}, err
		}

		return ctrl.Result{}, nil
	}

	// update secret if any spec field was changed
	if linkedsecret.Generation != linkedsecret.Status.ObservedGeneration {

		if err := r.UpdateLinkedSecret(ctx, &linkedsecret); err != nil {
			return ctrl.Result{}, err
		}

	}

	// recreate job if manager were restarted
	if _, ok := r.Cronjob[linkedsecret.UID]; !ok &&
		!linkedsecret.Spec.Suspended &&
		linkedsecret.Spec.Schedule != "" {

		if err := r.UpdateLinkedSecret(ctx, &linkedsecret); err != nil {
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}

	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LinkedSecretReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&securityv1.LinkedSecret{}).
		Complete(r)
}
