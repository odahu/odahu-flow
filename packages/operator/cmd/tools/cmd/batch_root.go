package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

func init() {
	rootCmd.AddCommand(batchCommand)
}


var batchCommand = &cobra.Command{
	Use:  "batch",
	Short: "Support tools to prepare environment for batch inference",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			os.Exit(0)
		}
	},
}