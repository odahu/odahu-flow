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
	conn_http_storage "github.com/odahu/odahu-flow/packages/operator/pkg/repository/connection/http"
	train_http_storage "github.com/odahu/odahu-flow/packages/operator/pkg/repository/training/http"
	conn_service "github.com/odahu/odahu-flow/packages/operator/pkg/service/connection"
	"github.com/odahu/odahu-flow/packages/operator/pkg/trainer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("trainer-main")

const (
	mtFileCLIParam             = "mt-file"
	mtIDCLIParam               = "mt-id"
	apiURLCLIParam             = "api-url"
	outputTrainingDirCLIParam  = "output-dir"
	MTFileConfigKey            = "trainer.mtFile"
	OutputTrainingDirConfigKey = "trainer.outputDir"
	APIURLConfigKey            = "trainer.auth.apiUrl"
	ModelTrainingIDConfigKey   = "trainer.modelTrainingId"
)

var mainCmd = &cobra.Command{
	Use:              "trainer",
	Short:            "odahu-flow trainer cli",
	TraverseChildren: true,
}

var trainerSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Prepare environment for a trainer",
	Run: func(cmd *cobra.Command, args []string) {
		if err := newTrainerWithHTTPRepositories(config.MustLoadConfig().Trainer).Setup(); err != nil {
			log.Error(err, "Training setup failed")
			os.Exit(1)
		}
	},
}

var saveCmd = &cobra.Command{
	Use:   "result",
	Short: "Save a trainer result",
	Run: func(cmd *cobra.Command, args []string) {
		if err := newTrainerWithHTTPRepositories(config.MustLoadConfig().Trainer).SaveResult(); err != nil {
			log.Error(err, "Result saving failed")
			os.Exit(1)
		}
	},
}

func init() {
	currentDir, err := os.Getwd()
	if err != nil {
		// Impossible situation
		panic(err)
	}

	config.InitBasicParams(mainCmd)

	mainCmd.PersistentFlags().String(mtFileCLIParam, "mt.json", "File with model training content")
	config.PanicIfError(viper.BindPFlag(MTFileConfigKey, mainCmd.PersistentFlags().Lookup(mtFileCLIParam)))

	mainCmd.PersistentFlags().String(mtIDCLIParam, "", "ID of the model training")
	config.PanicIfError(viper.BindPFlag(ModelTrainingIDConfigKey, mainCmd.PersistentFlags().Lookup(mtIDCLIParam)))

	mainCmd.PersistentFlags().String(apiURLCLIParam, "", "API URL")
	config.PanicIfError(viper.BindPFlag(APIURLConfigKey, mainCmd.PersistentFlags().Lookup(apiURLCLIParam)))

	mainCmd.PersistentFlags().String(
		outputTrainingDirCLIParam, currentDir,
		"The path to the dir when a user trainer will save their result",
	)
	config.PanicIfError(viper.BindPFlag(
		OutputTrainingDirConfigKey, mainCmd.PersistentFlags().Lookup(outputTrainingDirCLIParam),
	))

	mainCmd.AddCommand(trainerSetupCmd, saveCmd)
}

func newTrainerWithHTTPRepositories(config config.TrainerConfig) *trainer.ModelTrainer {
	log.Info(fmt.Sprintf("OAuthOIDCTokenEndpoint: %s", viper.GetString(config.Auth.OAuthOIDCTokenEndpoint)))
	trainRepo := train_http_storage.NewRepository(
		config.Auth.APIURL,
		config.Auth.APIToken,
		config.Auth.ClientID,
		config.Auth.ClientSecret,
		config.Auth.OAuthOIDCTokenEndpoint,
	)
	connRepo := conn_http_storage.NewRepository(
		config.Auth.APIURL,
		config.Auth.APIToken,
		config.Auth.ClientID,
		config.Auth.ClientSecret,
		config.Auth.OAuthOIDCTokenEndpoint,
	)
	connService := conn_service.NewService(connRepo)

	return trainer.NewModelTrainer(trainRepo, connService, config)
}

func main() {
	if err := mainCmd.Execute(); err != nil {
		log.Error(err, "trainer CLI command failed")

		os.Exit(1)
	}
}
