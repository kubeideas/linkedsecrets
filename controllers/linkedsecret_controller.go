/*
Copyright 2021.

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
	"fmt"
	"reflect"

	securityv1 "linkedsecrets/api/v1"

	"github.com/go-logr/logr"
	"github.com/robfig/cron/v3"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// LinkedSecretReconciler reconciles a LinkedSecret object
type LinkedSecretReconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
	Cronjob  map[types.UID]*cron.Cron
}

// +kubebuilder:rbac:groups=security.kubeideas.io,resources=linkedsecrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=security.kubeideas.io,resources=linkedsecrets/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=security.kubeideas.io,resources=linkedsecrets/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;update;patch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete

// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *LinkedSecretReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	//ctx := context.Background()
	log := r.Log.WithValues("linkedsecret", req.NamespacedName)

	var linkedsecret securityv1.LinkedSecret

	if err := r.Get(ctx, req.NamespacedName, &linkedsecret); err != nil {
		log.Info("linkedsecret does not exists", "linkedsecret", err)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Add finalizer if necessary
	if err := r.AddFinalizer(&linkedsecret); err != nil {
		return ctrl.Result{}, err
	}

	// if linkedsecret not being deleted add finalizer
	// if linkedsecret is being delete cronjob entry must be deleted first
	if stopReconcile, err := r.CheckObjectDeletion(ctx, &linkedsecret); stopReconcile && err != nil {
		return ctrl.Result{}, err
	} else if stopReconcile && err == nil {
		return ctrl.Result{}, nil
	}

	// create new secret with provider data if does not exists
	var secret corev1.Secret
	secretName := client.ObjectKey{Namespace: linkedsecret.Namespace, Name: linkedsecret.Spec.SecretName}
	if err := r.Get(ctx, secretName, &secret); err != nil {
		log.V(1).Info("BEGIN RECONCILING - NEW SECRET", "secret", fmt.Sprintf("%s/%s", linkedsecret.Namespace, linkedsecret.Spec.SecretName))
		// create new kubernetes secret but keep old one
		if err := r.NewLinkedSecret(ctx, &linkedsecret); err != nil {
			return ctrl.Result{}, err
		}
		log.V(1).Info("END RECONCILING - NEW SECRET", "secret", fmt.Sprintf("%s/%s", linkedsecret.Namespace, linkedsecret.Spec.SecretName))

		return ctrl.Result{}, nil
	}

	// Update linkedsecret if spec was changed
	if linkedsecret.Status.CurrentSecretStatus != STATUSSYNCHED ||
		linkedsecret.Status.CurrentProvider != linkedsecret.Spec.Provider ||
		!reflect.DeepEqual(linkedsecret.Status.CurrentProviderOptions, linkedsecret.Spec.ProviderOptions) ||
		linkedsecret.Status.CurrentSchedule != linkedsecret.Spec.Schedule ||
		linkedsecret.Status.KeepSecretOnDelete != linkedsecret.Spec.KeepSecretOnDelete ||
		(linkedsecret.Spec.Suspended && linkedsecret.Status.CronJobStatus != JOBSUSPENDED) ||
		(!linkedsecret.Spec.Suspended && linkedsecret.Status.CronJobStatus == JOBSUSPENDED) {

		log.V(1).Info("BEGIN RECONCILING - UPDATE", "linkedsecret", fmt.Sprintf("%s/%s", linkedsecret.Namespace, linkedsecret.Name))
		if err := r.UpdateLinkedSecret(ctx, &linkedsecret); err != nil {
			return ctrl.Result{}, err
		}

		log.V(1).Info("END RECONCILING - UPDATE", "linkedsecret", fmt.Sprintf("%s/%s", linkedsecret.Namespace, linkedsecret.Name))

		return ctrl.Result{}, nil

	}

	// Recreate schedule if controller was restarted
	//if !linkedsecret.Spec.Suspended && linkedsecret.Status.CronJobStatus != JOBFAILPARSESCHEDULE {
	if _, ok := r.Cronjob[linkedsecret.UID]; !ok {
		log.V(1).Info("BEGIN RECONCILING - RESTART MANAGER", "schedule", "recreated")

		// skip linkedsecrets without schedule
		if linkedsecret.Spec.Schedule != "" {
			// Add cronjob
			r.AddCronjob(ctx, &linkedsecret)
		}

		log.V(1).Info("END RECONCILING - RESTART MANAGER", "schedule", "recreated")
		return ctrl.Result{}, nil
	}
	//}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LinkedSecretReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&securityv1.LinkedSecret{}).
		Complete(r)
}
