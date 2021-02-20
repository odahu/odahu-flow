package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(batchCommand)
	batchCommand.AddCommand(validateInputCommand)
	batchCommand.AddCommand(validateOutputCommand)
}

var batchCommand = &cobra.Command{
	Use:  "batch",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var validateInputCommand = &cobra.Command{
	Use:  "validate-input",
	Short: "validate input for user batch inference container",
	Run: func(cmd *cobra.Command, args []string) {

	},
}


var validateOutputCommand = &cobra.Command{
	Use:  "validate-output",
	Short: "validate output of user batch inference container",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
