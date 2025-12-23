/*
Copyright 2025.

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

package controller

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	monitoringv1 "github.com/lswcomp/BI_WS2025/api/v1"
)

// ObservationReconciler reconciles a Observation object
type ObservationReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// Definitions to manage status conditions
const (
	// typeAvailableObservation represents the status of the Observation reconciliation
	typeAvailableObservation = "Available"
	// typeProgressingObservation represents the status used when the Observation is being reconciled
	typeProgressingObservation = "Progressing"
	// typeDegradedObservation represents the status used when the Observation has encountered an error
	typeDegradedObservation = "Degraded"
)

// +kubebuilder:rbac:groups=monitoring.bi-ws2025.de,resources=observations,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=monitoring.bi-ws2025.de,resources=observations/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=monitoring.bi-ws2025.de,resources=observations/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=persistentVolumeClaims,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Observation object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.22.1/pkg/reconcile
func (r *ObservationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	var observation monitoringv1.Observation
	if err := r.Get(ctx, req.NamespacedName, &observation); err != nil {
		if apierrors.IsNotFound(err) {
			// If the custom resource is not found then it usually means that it was deleted or not created
			// In this way, we will stop the reconciliation
			log.Info("Observation resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get Observation")
		return ctrl.Result{}, err
	}

	// Initialize status conditions if not yet present
	if len(observation.Status.Conditions) == 0 {
		meta.SetStatusCondition(&observation.Status.Conditions, metav1.Condition{
			Type:    typeProgressingObservation,
			Status:  metav1.ConditionUnknown,
			Reason:  "Reconciling",
			Message: "Starting reconciliation",
		})
		if err := r.Status().Update(ctx, &observation); err != nil {
			log.Error(err, "Failed to update Observation status")
			return ctrl.Result{}, err
		}

		// Re-fetch the CronJob after updating the status
		if err := r.Get(ctx, req.NamespacedName, &observation); err != nil {
			log.Error(err, "Failed to re-fetch Observation")
			return ctrl.Result{}, err
		}
	}

	var deployments appsv1.DeploymentList
	if err := r.List(ctx, &deployments, client.InNamespace(req.Namespace), client.MatchingFields{jobOwnerKey: req.Name}); err != nil {
		log.Error(err, "unable to list Deployments")
		// Update status condition to reflect the error
		meta.SetStatusCondition(&observation.Status.Conditions, metav1.Condition{
			Type:    typeDegradedObservation,
			Status:  metav1.ConditionTrue,
			Reason:  "ReconciliationError",
			Message: fmt.Sprintf("Failed to list Deployments: %v", err),
		})
		if statusErr := r.Status().Update(ctx, &observation); statusErr != nil {
			log.Error(statusErr, "Failed to update Observation status")
		}
		return ctrl.Result{}, err
	}

	var services corev1.ServiceList
	if err := r.List(ctx, &services, client.InNamespace(req.Namespace), client.MatchingFields{jobOwnerKey: req.Name}); err != nil {
		log.Error(err, "unable to list Services")
		// Update status condition to reflect the error
		meta.SetStatusCondition(&observation.Status.Conditions, metav1.Condition{
			Type:    typeDegradedObservation,
			Status:  metav1.ConditionTrue,
			Reason:  "ReconciliationError",
			Message: fmt.Sprintf("Failed to list Services: %v", err),
		})
		if statusErr := r.Status().Update(ctx, &observation); statusErr != nil {
			log.Error(statusErr, "Failed to update Observation status")
		}
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

var (
	jobOwnerKey = ".metadata.controller"
	apiGVStr    = monitoringv1.GroupVersion.String()
)

// SetupWithManager sets up the controller with the Manager.
func (r *ObservationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&monitoringv1.Observation{}).
		Named("observation").
		Complete(r)
}
