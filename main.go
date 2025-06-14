package main

import (
	"flag"
	"os"

	"github.com/quietinvestor/simplecontroller/controllers"
	"github.com/spf13/cobra"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/klog/v2/textlogger"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
)

func main() {
	loggerConfig := textlogger.NewConfig()
	loggerConfig.AddFlags(flag.CommandLine)

	var namespace string

	rootCmd := &cobra.Command{
		Use:   "simplecontroller",
		Short: "Simple controller for labeling Pods",
		Run: func(cmd *cobra.Command, args []string) {
			logger := textlogger.NewLogger(loggerConfig).WithName("simplecontroller")
			ctrl.SetLogger(logger)

			mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
				Cache: cache.Options{
					DefaultNamespaces: map[string]cache.Config{
						namespace: {},
					},
					DefaultLabelSelector: labels.SelectorFromSet(map[string]string{
						controllers.WatchKey: controllers.WatchValue,
					}),
				},
			})
			if err != nil {
				logger.Error(err, "unable to create manager")
				os.Exit(1)
			}

			if err := (&controllers.PodLabelReconciler{Client: mgr.GetClient()}).SetupWithManager(mgr); err != nil {
				logger.Error(err, "unable to set up PodLabelReconciler")
				os.Exit(1)
			}

			if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
				logger.Error(err, "problem running manager")
				os.Exit(1)
			}
		},
	}

	rootCmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)
	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "namespace to watch (empty: watch all)")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
