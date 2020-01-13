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
	api_config "github.com/odahu/odahu-flow/packages/operator/pkg/config/api"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	_ "github.com/odahu/odahu-flow/packages/operator/pkg/config/connection"
	_ "github.com/odahu/odahu-flow/packages/operator/pkg/config/deployment"
	_ "github.com/odahu/odahu-flow/packages/operator/pkg/config/packaging"
	_ "github.com/odahu/odahu-flow/packages/operator/pkg/config/training"
	"github.com/odahu/odahu-flow/packages/operator/pkg/webserver"
	"github.com/spf13/cobra"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("api")

const (
	BackendTypeParam = "backend-type"
	CRDPathParam     = "crd-path"
)

var mainCmd = &cobra.Command{
	Use:   "api",
	Short: "odahu-flow API server",
	Run: func(cmd *cobra.Command, args []string) {
		apiServer, err := webserver.NewAPIServer()
		if err != nil {
			log.Error(err, "Can't set up api server")
			os.Exit(1)
		}

		go func() {
			if err := apiServer.Run(); err != nil && err != http.ErrServerClosed {
				log.Error(err, "Closing of api server")
				os.Exit(1)
			}
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		log.Info("Shutdown Server ...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := apiServer.Close(ctx); err != nil {
			log.Error(err, "Server shutdowns")
			os.Exit(1)
		}
	},
}

func init() {
	config.InitBasicParams(mainCmd)

	mainCmd.PersistentFlags().String(BackendTypeParam, "", "Backend type")
	config.PanicIfError(viper.BindPFlag(api_config.BackendType, mainCmd.PersistentFlags().Lookup(BackendTypeParam)))

	mainCmd.PersistentFlags().String(CRDPathParam, "", "Path to a directory with Odahu-flow CRDs")
	config.PanicIfError(viper.BindPFlag(api_config.LocalBackendCRDPath, mainCmd.PersistentFlags().Lookup(CRDPathParam)))
}

func main() {
	if err := mainCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
