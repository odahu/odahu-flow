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
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	"github.com/odahu/odahu-flow/packages/operator/pkg/packager"
	connection_http_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/connection/http"
	packaging_http_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/packaging/http"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("packager-main")

const (
	mpFileCLIParam   = "mp-file"
	mpIDCLIParam     = "mp-id"
	apiURLCLIParam   = "api-url"
	MPFile           = "packager.mpFile"
	APIURL           = "packager.auth.apiUrl"
	ModelPackagingID = "packager.modelPackagingId"
)

var mainCmd = &cobra.Command{
	Use:              "packager",
	Short:            "odahu-flow packager cli",
	TraverseChildren: true,
}

var packagerSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Prepare environment for a packager",
	Run: func(cmd *cobra.Command, args []string) {
		if err := newPackagerWithHTTPRepositories(config.MustLoadConfig().Packager).SetupPackager(); err != nil {
			log.Error(err, "Packaging setup failed")
			os.Exit(1)
		}
	},
}

var saveCmd = &cobra.Command{
	Use:   "result",
	Short: "Save a packer result",
	Run: func(cmd *cobra.Command, args []string) {
		if err := newPackagerWithHTTPRepositories(config.MustLoadConfig().Packager).SaveResult(); err != nil {
			log.Error(err, "Result saving failed")
			os.Exit(1)
		}
	},
}

func init() {
	config.InitBasicParams(mainCmd)

	mainCmd.PersistentFlags().String(mpFileCLIParam, "mp.json", "File with model packaging content")
	config.PanicIfError(viper.BindPFlag(MPFile, mainCmd.PersistentFlags().Lookup(mpFileCLIParam)))

	mainCmd.PersistentFlags().String(mpIDCLIParam, "", "ID of the model packaging")
	config.PanicIfError(viper.BindPFlag(ModelPackagingID, mainCmd.PersistentFlags().Lookup(mpIDCLIParam)))

	mainCmd.PersistentFlags().String(apiURLCLIParam, "", "API URL")
	config.PanicIfError(viper.BindPFlag(APIURL, mainCmd.PersistentFlags().Lookup(apiURLCLIParam)))

	mainCmd.AddCommand(packagerSetupCmd, saveCmd)
}

func newPackagerWithHTTPRepositories(config config.PackagerConfig) *packager.Packager {
	packRepo := packaging_http_repository.NewRepository(
		config.Auth.APIURL,
		config.Auth.APIToken,
		config.Auth.ClientID,
		config.Auth.ClientSecret,
		config.Auth.OAuthOIDCTokenEndpoint,
	)
	connRepo := connection_http_repository.NewRepository(
		config.Auth.APIURL,
		config.Auth.APIToken,
		config.Auth.ClientID,
		config.Auth.ClientSecret,
		config.Auth.OAuthOIDCTokenEndpoint,
	)

	return packager.NewPackager(packRepo, connRepo, config)
}

func main() {
	if err := mainCmd.Execute(); err != nil {
		log.Error(err, "packager CLI command failed")

		os.Exit(1)
	}
}
