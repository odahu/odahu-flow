package main

import (
	migrator_package "github.com/odahu/odahu-flow/packages/operator/pkg/database/migrators/postgres"
	"github.com/spf13/cobra"
	"log"
)

var (
	databaseConnection string
)

var rootCmd = &cobra.Command{
	Use:   "odahu-migrate",
	Short: "odahu-migrate is tool that migrate odahu database forward",
	Run: func(cmd *cobra.Command, args []string) {
		migr := migrator_package.Migrator{databaseConnection}
		err := migr.MigrateToLatest()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {

	rootCmd.PersistentFlags().StringVar(&databaseConnection,
		"database", "", "connection string to database (required)",
	)
	err := rootCmd.MarkPersistentFlagRequired("database")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
