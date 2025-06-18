package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/quietinvestor/simplecontroller/internal/setup"

	"github.com/spf13/cobra"
	"k8s.io/klog/v2/textlogger"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func main() {
	loggerConfig := textlogger.NewConfig()
	loggerConfig.AddFlags(flag.CommandLine)

	var namespace string

	rootCmd := &cobra.Command{
		Use:   "simplecontroller",
		Short: "Simple controller for labeling Pods",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := config.GetConfigOrDie()

			mgr, err := setup.Setup(cfg, namespace, *loggerConfig)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to initialize manager: %v\n", err)
				os.Exit(1)
			}

			if err := mgr.Start(cmd.Context()); err != nil {
				fmt.Fprintf(os.Stderr, "controller runtime manager exited: %v\n", err)
				os.Exit(1)
			}
		},
	}

	rootCmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)
	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "Namespace to watch. If empty, watch all.")

	flag.Parse()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
