package setup

import (
	"fmt"
	"net/http"

	"github.com/quietinvestor/simplecontroller/internal/controller"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/klog/v2/textlogger"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
)

// SetupManager creates a new controller-runtime Manager,
// configures namespace and label-based caching, adds health checks,
// and registers the PodLabelReconciler.
func SetupManager(namespace string, loggerConfig textlogger.Config) (ctrl.Manager, error) {
	logger := textlogger.NewLogger(&loggerConfig).WithName("simplecontroller")
	ctrl.SetLogger(logger)

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Cache: cache.Options{
			DefaultNamespaces: map[string]cache.Config{
				namespace: {},
			},
			DefaultLabelSelector: labels.SelectorFromSet(map[string]string{
				controller.WatchKey: controller.WatchValue,
			}),
		},
		HealthProbeBindAddress: ":8081",
	})
	if err != nil {
		logger.Error(err, "unable to create manager")
		return nil, fmt.Errorf("unable to create manager: %w", err)
	}

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		logger.Error(err, "unable to set up liveness check")
		return nil, fmt.Errorf("unable to set up liveness check: %w", err)
	}

	if err := mgr.AddReadyzCheck("readyz",
		func(req *http.Request) error {
			if !mgr.GetCache().WaitForCacheSync(req.Context()) {
				return fmt.Errorf("cache not synced")
			}
			return nil
		},
	); err != nil {
		logger.Error(err, "unable to set up readiness check")
		return nil, fmt.Errorf("unable to set up readiness check: %w", err)
	}

	if err := (&controller.PodLabelReconciler{
		Client: mgr.GetClient(),
	}).SetupWithManager(mgr); err != nil {
		logger.Error(err, "unable to set up PodLabelReconciler")
		return nil, fmt.Errorf("unable to set up PodLabelReconciler: %w", err)
	}

	return mgr, nil
}
