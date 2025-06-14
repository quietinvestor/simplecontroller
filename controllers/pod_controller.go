package controllers

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	WatchKey   = "simplecontroller.io/watch"
	WatchValue = "true"

	ColorKey   = "simplecontroller.io/color"
	ColorValue = "blue"
)

type PodLabelReconciler struct {
	client.Client
}

func (r *PodLabelReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var pod corev1.Pod

	if err := r.Get(ctx, req.NamespacedName, &pod); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if pod.Labels[WatchKey] != WatchValue {
		return ctrl.Result{}, nil
	}

	if pod.Labels[ColorKey] == ColorValue {
		return ctrl.Result{}, nil
	}

	updated := pod.DeepCopy()
	updated.Labels[ColorKey] = ColorValue

	return ctrl.Result{}, r.Patch(ctx, updated, client.MergeFrom(&pod))
}

func (r *PodLabelReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).For(&corev1.Pod{}).Complete(r)
}
