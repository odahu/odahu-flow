
package main

import (
	"context"
	"fmt"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	"github.com/odahu/odahu-flow/packages/operator/pkg/controller"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"syscall"
	"time"
)

var log = logf.Log.WithName("controller-main")

var mainCmd = &cobra.Command{
	Use:   "api",
	Short: "odahu-flow controller",
	Long: "Odahu flow controller is responsible to reconcile processers in runtime (e.g. Kubernetes) " +
		  "according to desired state in app storage",
	Run: func(cmd *cobra.Command, args []string) {

		cfg := config.MustLoadConfig()
		controller, err := controller.NewController(cfg)
		if err != nil {
			log.Error(err, "unable set up controller")
			os.Exit(1)
		}

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		ctx, cancel := context.WithCancel(context.Background())

		go func() {
			<-quit
			log.Info("SIGINT was received. Try shutdown gracefully")
			cancel()

			t := time.NewTimer(cfg.Common.GracefulTimeout)
			defer t.Stop()
			<-t.C
			log.Info("Unable to shutdown gracefully")
			os.Exit(1)
		}()

		if err := controller.Run(ctx); err != nil {
			log.Error(err, "Error in controller")
			os.Exit(1)
		}

		log.Info("Graceful shutdown is successful")

	},
}

func init() {
	config.InitBasicParams(mainCmd)
}

func main() {
	if err := mainCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
