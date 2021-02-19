package cmd

import (
	"go.uber.org/zap"
	"os"

	"github.com/spf13/cobra"
	connAPI "github.com/odahu/odahu-flow/packages/operator/pkg/apiclient/connection"
	connTools "github.com/odahu/odahu-flow/packages/operator/pkg/tools/connections"
)

const (
	configureRCloneCommandUsage = `List of ODAHU Connection IDs that must exist in ODAHU Cluster 
Each connection ID can be suffixed by :<RClone Storage Name>. 
If connection ID is not suffixed then RClone storage name is generated as "odahu-<Connection ID>".
For example next flags:
--conn model-output:output --conn input-data 
will lead to adding two storages to Rclone config with names accordingly:
[output]
...
[odahu-input-data]
...`
)

var (
	fileValues    []string
	clusterValues []string
)

func init() {
	rootCmd.AddCommand(authCommand)
	authCommand.AddCommand(configureRCloneCommand)
	configureRCloneCommand.Flags().StringArrayVar(&clusterValues, "conn", []string{}, configureRCloneCommandUsage)
}

var authCommand = &cobra.Command{
	Use:  "auth",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var configureRCloneCommand = &cobra.Command{
	Use:  "configure-rclone",
	Long: `Generate rclone configuration based on ODAHU object storage connection`,
	Run: func(cmd *cobra.Command, args []string) {
		client := connAPI.NewClient(cfg.Auth.APIURL, "",
			cfg.Auth.ClientID, cfg.Auth.ClientSecret, cfg.Auth.OAuthOIDCTokenEndpoint)
		err := connTools.GenerateRClone(clusterValues, fileValues, client)
		if err != nil {
			zap.S().Errorw("There are errors during generating RClone storages", zap.Error(err))
			os.Exit(1)
		}
		zap.S().Info("Completed")

	},
}
