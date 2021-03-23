/*
 *
 *     Copyright 2021 EPAM Systems
 *
 *     Licensed under the Apache License, Version 2.0 (the "License");
 *     you may not use this file except in compliance with the License.
 *     You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 *     Unless required by applicable law or agreed to in writing, software
 *     distributed under the License is distributed on an "AS IS" BASIS,
 *     WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *     See the License for the specific language governing permissions and
 *     limitations under the License.
 */

package cmd

import (
	connAPI "github.com/odahu/odahu-flow/packages/operator/pkg/apiclient/connection"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/model_registry/objectstorage"
	"github.com/spf13/cobra"
	"fmt"
	"os"
)

// flag values
var (
	connName   string
	remotePath string
	localPath  string
)

func init() {
	rootCmd.AddCommand(registryCommand)
	registryCommand.AddCommand(objectStorageRegistryCommand)
	objectStorageRegistryCommand.PersistentFlags().StringVar(
		&connName, "conn", "", "connection: s3, gcs, azureblob types are supported",
	)
	_ = objectStorageRegistryCommand.MarkPersistentFlagRequired("conn")
	objectStorageRegistryCommand.PersistentFlags().StringVar(
		&remotePath, "remotePath", "", "remote path that will be appended to conn.URI",
	)

	objectStorageRegistryCommand.AddCommand(syncModelCommand)
	syncModelCommand.PersistentFlags().StringVar(
		&localPath, "localPath", "", "Where sync model locally",
	)

	objectStorageRegistryCommand.AddCommand(modelInfoCommand)
}

var registryCommand = &cobra.Command{
	Use:  "registry",
	Short: "Operations with model registries",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			os.Exit(0)
		}
	},
}
var objectStorageRegistryCommand = &cobra.Command{
	Use:  "object-storage",
	Short: "Operations with models that are stored in object storage",
	Long: "object-storage registry is not true model registry but just a way " +
		"to sync models from object storage with minimum ODAHU conventions about files",
	Example: "odahu-tools registry object-storage sync --conn model-output --path wine-1.0.zip",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			os.Exit(0)
		}
	},
}

var syncModelCommand = &cobra.Command{
	Use:  "sync",
	Short: "Sync from object storage to local folder",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := connAPI.NewClient(cfg.Auth.APIURL, "",
			cfg.Auth.ClientID, cfg.Auth.ClientSecret, cfg.Auth.OAuthOIDCTokenEndpoint)
		registry := objectstorage.NewModelRegistry(client)

		if _, err := registry.SyncModel(connName, remotePath, localPath); err != nil {
			return err
		}
		fmt.Printf("Model files were successfully synced to %s\n", localPath)
		return nil
	},
}


var modelInfoCommand = &cobra.Command{
	Use:  "info",
	Short: "Show model name and version",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := connAPI.NewClient(cfg.Auth.APIURL, "",
			cfg.Auth.ClientID, cfg.Auth.ClientSecret, cfg.Auth.OAuthOIDCTokenEndpoint)
		registry := objectstorage.NewModelRegistry(client)

		name, ver, err := registry.Meta(connName, remotePath)
		if err != nil {
			return err
		}
		fmt.Printf("Model name: %s\nModel version: 	%s\n", name, ver)
		return nil
	},
}


