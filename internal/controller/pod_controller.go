package controller

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Label keys and values used to determine whether to patch a Pod.
const (
	WatchKey   = "simplecontroller.io/watch"
	WatchValue = "true"

	ColorKey   = "simplecontroller.io/color"
	ColorValue = "blue"
)

// PodLabelReconciler holds the client used to reconcile Pods based on label criteria.
type PodLabelReconciler struct {
	client.Client
}

// Reconcile adds the color label to a Pod if it has the watch label.
// It only patches the Pod if the label is missing or has an incorrect value.
func (r *PodLabelReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := ctrl.LoggerFrom(ctx)

	var pod corev1.Pod

	if err := r.Get(ctx, req.NamespacedName, &pod); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		logger.Error(err, "failed to get pod")
		return ctrl.Result{}, err
	}

	if pod.Labels[WatchKey] != WatchValue {
		return ctrl.Result{}, nil
	}

	if v, ok := pod.Labels[ColorKey]; ok && v == ColorValue {
		return ctrl.Result{}, nil
	}

	updated := pod.DeepCopy()
	updated.Labels[ColorKey] = ColorValue
	logger.V(2).Info("successfully labeled pod")

	if err := r.Patch(ctx, updated, client.MergeFrom(&pod)); err != nil {
		logger.Error(err, "failed to patch pod")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager registers the PodLabelReconciler with the manager
// and configures it to watch Pod resources.
func (r *PodLabelReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).For(&corev1.Pod{}).Complete(r)
}
