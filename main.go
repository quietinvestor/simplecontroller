package main

import (
	"github.com/quietinvestor/simplecontroller/controllers"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/klog/v2/textlogger"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
)

func main() {
	loggerConfig := textlogger.NewConfig()
	logger := textlogger.NewLogger(loggerConfig).WithName("simplecontroller")
	ctrl.SetLogger(logger)

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Cache: cache.Options{
			DefaultLabelSelector: labels.SelectorFromSet(map[string]string{controllers.WatchKey: controllers.WatchValue}),
		},
	})
	if err != nil {
		panic(err)
	}

	if err := (&controllers.PodLabelReconciler{Client: mgr.GetClient()}).SetupWithManager(mgr); err != nil {
		panic(err)
	}

	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		panic(err)
	}
}
