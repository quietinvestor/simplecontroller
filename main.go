package main

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/klog/v2/textlogger"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	watchKey   = "simplecontroller.io/watch"
	watchValue = "true"

	colorKey   = "simplecontroller.io/color"
	colorValue = "blue"
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

	if pod.Labels[watchKey] != watchValue {
		return ctrl.Result{}, nil
	}

	if pod.Labels[colorKey] == colorValue {
		return ctrl.Result{}, nil
	}

	updated := pod.DeepCopy()
	updated.Labels[colorKey] = colorValue

	return ctrl.Result{}, r.Patch(ctx, updated, client.MergeFrom(&pod))
}

func (r *PodLabelReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).For(&corev1.Pod{}).Complete(r)
}

func main() {
	loggerConfig := textlogger.NewConfig()
	logger := textlogger.NewLogger(loggerConfig).WithName("simplecontroller")
	ctrl.SetLogger(logger)

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Cache: cache.Options{
			DefaultLabelSelector: labels.SelectorFromSet(map[string]string{watchKey: watchValue}),
		},
	})
	if err != nil {
		panic(err)
	}

	if err := (&PodLabelReconciler{Client: mgr.GetClient()}).SetupWithManager(mgr); err != nil {
		panic(err)
	}

	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		panic(err)
	}
}
