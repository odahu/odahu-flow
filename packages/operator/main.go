/*


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

package main

import (
	conn_api_client "github.com/odahu/odahu-flow/packages/operator/pkg/apiclient/connection"
	"os"

	istioschema "github.com/aspenmesh/istio-client-go/pkg/client/clientset/versioned/scheme"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tektonschema "github.com/tektoncd/pipeline/pkg/client/clientset/versioned/scheme"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	knservingv1 "knative.dev/serving/pkg/apis/serving/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/odahu/odahu-flow/packages/operator/pkg/config"

	odahuflowv1alpha1 "github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	train_api_client "github.com/odahu/odahu-flow/packages/operator/pkg/apiclient/training"
	mp_api_client "github.com/odahu/odahu-flow/packages/operator/pkg/apiclient/packaging"
	"github.com/odahu/odahu-flow/packages/operator/controllers"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

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

	_ = clientgoscheme.AddToScheme(scheme)

	istioschema.AddToScheme(scheme)
	_ = tektonschema.AddToScheme(scheme)
	_ = knservingv1.AddToScheme(scheme)

	_ = odahuflowv1alpha1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

func main() {
	if err := mainCmd.Execute(); err != nil {
		setupLog.Error(err, "Shutdown operator")
		os.Exit(1)
	}
}

var mainCmd = &cobra.Command{
	Use:   "operator",
	Short: "odahu-flow operator",
	Run:   start,
}

func start(cmd *cobra.Command, args []string) {

	odahuConfig := config.MustLoadConfig()

	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: ":8080",
		Port:               9443,
		LeaderElection:     false,
		LeaderElectionID:   "23d32765.odahu.org",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	// Setup the training toolchain repository
	authCfg := odahuConfig.Operator.Auth

	connAPI := conn_api_client.NewClient(
		authCfg.APIURL,
		authCfg.APIToken,
		authCfg.ClientID,
		authCfg.ClientSecret,
		authCfg.OAuthOIDCTokenEndpoint,
	)

	if odahuConfig.Training.Enabled {
		trainAPIClient := train_api_client.NewClient(
			authCfg.APIURL, authCfg.APIToken, authCfg.ClientID, authCfg.ClientSecret, authCfg.OAuthOIDCTokenEndpoint,
		)

		if err = controllers.NewModelTrainingReconciler(
			mgr, *odahuConfig, trainAPIClient,
		).SetupWithManager(mgr); err != nil {

			setupLog.Error(err, "unable to create controller", "controller", "ModelTraining")
			os.Exit(1)

		}
	}

	if odahuConfig.Packaging.Enabled {

		packAPIClient := mp_api_client.NewClient(
			authCfg.APIURL, authCfg.APIToken, authCfg.ClientID, authCfg.ClientSecret, authCfg.OAuthOIDCTokenEndpoint,
		)

		if err = controllers.NewModelPackagingReconciler(
			mgr, *odahuConfig, packAPIClient,
		).SetupWithManager(mgr); err != nil {

			setupLog.Error(err, "unable to create controller", "controller", "ModelPackaging")
			os.Exit(1)
		}
	}

	if odahuConfig.Deployment.Enabled {
		if err = controllers.NewModelDeploymentReconciler(mgr, *odahuConfig).SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "ModelDeployment")
			os.Exit(1)
		}
		if err = controllers.NewModelRouteReconciler(mgr, *odahuConfig).SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "ModelRoute")
			os.Exit(1)
		}
	}

	if odahuConfig.Batch.Enabled {

		batchOpts := controllers.BatchInferenceJobReconcilerOptions{
			Mgr:               mgr,
			ConnGetter:        connAPI,
			Cfg:               *odahuConfig,
		}

		if err = controllers.NewBatchInferenceJobReconciler(batchOpts).SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "Batch")
			os.Exit(1)
		}
	}

	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
