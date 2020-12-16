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
	"context"
	"fmt"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apiserver"
	migrator_package "github.com/odahu/odahu-flow/packages/operator/pkg/database/migrators/postgres"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"path/filepath"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"syscall"
)

var log = logf.Log.WithName("api-main")

const (
	BackendTypeParam    = "backend-type"
	CRDPathParam        = "crd-path"
	BackendType         = "api.backend.type"
	LocalBackendCRDPath = "api.backend.local.crdPath"
)

var mainCmd = &cobra.Command{
	Use:   "api",
	Short: "odahu-flow API server",
	Run: func(cmd *cobra.Command, args []string) {

		cfg := config.MustLoadConfig()

		migr, err := migrator_package.NewMigrator(cfg.Common.DatabaseConnectionString)
		if err != nil {
			log.Error(err, "Unable to create migrator")
			os.Exit(1)
		}
		err = migr.MigrateToLatest()
		if err != nil {
			log.Error(err, "Unable to migrate")
			os.Exit(1)
		}

		// Run API Server
		apiServer, err := apiserver.NewAPIServer(cfg)
		if err != nil {
			log.Error(err, "unable set up api server")
			os.Exit(1)
		}
		errCh := make(chan error)
		apiServer.Run(errCh)


		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-quit:
			log.Info("SIGINT was received")
		case err := <-errCh:
			log.Error(err, "Error during execution one of components")
		}

		log.Info("Graceful shutdown of api server components is started")
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Common.GracefulTimeout)
		defer cancel()
		if closeErr := apiServer.Close(ctx); closeErr != nil {
			log.Error(closeErr, "Graceful shutdown of api server components is failed. Exit immediately")
			os.Exit(1)
		}
		log.Info("Graceful shutdown of api server components is successful")
	},
}

func init() {
	config.InitBasicParams(mainCmd)

	mainCmd.PersistentFlags().String(BackendTypeParam, config.ConfigBackendType, "Backend type")
	config.PanicIfError(viper.BindPFlag(BackendType, mainCmd.PersistentFlags().Lookup(BackendTypeParam)))

	mainCmd.PersistentFlags().String(
		CRDPathParam, filepath.Join("config", "crds"), "Path to a directory with Odahu-flow CRDs",
	)
	config.PanicIfError(viper.BindPFlag(LocalBackendCRDPath, mainCmd.PersistentFlags().Lookup(CRDPathParam)))

}

func main() {
	if err := mainCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
