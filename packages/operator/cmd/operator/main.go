//
//    Copyright 2019 EPAM Systems
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
//

package main

import (
	"fmt"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	"os"

	"github.com/odahu/odahu-flow/packages/operator/pkg/apis"
	"github.com/odahu/odahu-flow/packages/operator/pkg/controller"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	k8s_config "sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
)

var log = logf.Log.WithName("operator")

var mainCmd = &cobra.Command{
	Use:   "operator",
	Short: "odahu-flow operator",
	Run:   startOperator,
}

func init() {
	config.InitBasicParams(mainCmd)

	const (
		enableTrainControllerName  = "enable-training-controller"
		enablePackControllerName   = "enable-packaging-controller"
		enableDepControllerName    = "enable-deployment-controller"
		trainingEnabledConfigKey   = "training.enabled"
		packagingEnabledConfigKey  = "packaging.enabled"
		deploymentEnabledConfigKey = "deployment.enabled"
	)

	mainCmd.PersistentFlags().Bool(enableTrainControllerName, false, "config file")
	mainCmd.PersistentFlags().Bool(enablePackControllerName, false, "config file")
	mainCmd.PersistentFlags().Bool(enableDepControllerName, false, "config file")

	if err := viper.BindPFlag(
		trainingEnabledConfigKey, mainCmd.PersistentFlags().Lookup(enableTrainControllerName),
	); err != nil {
		panic(err)
	}

	if err := viper.BindPFlag(
		packagingEnabledConfigKey, mainCmd.PersistentFlags().Lookup(enablePackControllerName),
	); err != nil {
		panic(err)
	}

	if err := viper.BindPFlag(
		deploymentEnabledConfigKey, mainCmd.PersistentFlags().Lookup(enableDepControllerName),
	); err != nil {
		panic(err)
	}
}

func main() {
	if err := mainCmd.Execute(); err != nil {
		log.Error(err, "Shutdown operator")
		os.Exit(1)
	}
}

func startOperator(cmd *cobra.Command, args []string) {
	odahuConfig := config.MustLoadConfig()

	logf.SetLogger(logf.ZapLogger(true))
	log := logf.Log.WithName("entrypoint")

	log.Info("setting up client for manager")

	cfg, err := k8s_config.GetConfig()
	if err != nil {
		log.Error(err, "unable to set up client config")
		os.Exit(1)
	}

	log.Info("setting up manager")
	mgr, err := manager.New(
		cfg,
		manager.Options{
			MetricsBindAddress: fmt.Sprintf(":%d", odahuConfig.Operator.MonitoringPort),
		},
	)

	if err != nil {
		log.Error(err, "unable to set up overall controller manager")
		os.Exit(1)
	}

	log.Info("Registering Components.")

	log.Info("Setting up odahu-flow scheme")
	if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
		log.Error(err, "unable add odahu-flow APIs to scheme")
		os.Exit(1)
	}

	log.Info("Setting up controller")
	if err := controller.AddToManager(mgr, odahuConfig); err != nil {
		log.Error(err, "unable to register controllers to the manager")
		os.Exit(1)
	}

	log.Info("Starting the operator.")
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		log.Error(err, "unable to run the manager")
		os.Exit(1)
	}
}
