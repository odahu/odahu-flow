package cmd

import (
	"github.com/spf13/cobra"
	predict_v2_tools "github.com/odahu/odahu-flow/packages/operator/pkg/tools/predict_v2"
	"go.uber.org/zap"
	"os"
)

const inputCommandUsage = `
Path to directory that contains set of files with .json extensions in "Kubeflow Predict V2 Inference Request" format
`
const outputCommandUsage = `
Path to directory where json files that passed validation will be dumped
`

var (
	source      string
	destination string
)

func init() {
	rootCmd.AddCommand(batchCommand)
	batchCommand.AddCommand(validateInputCommand)
	batchCommand.AddCommand(validateOutputCommand)
	validateInputCommand.Flags().StringVarP(
		&source, "source", "s", ".", inputCommandUsage,
	)
	validateOutputCommand.Flags().StringVarP(
		&destination, "destination", "d", ".", outputCommandUsage,
	)
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


		_, err := predict_v2_tools.ValidateDir(source, &destination)
		if err != nil {
			zap.S().Errorw("There are errors during validation", zap.Error(err))
			os.Exit(1)
		}
	},
}


var validateOutputCommand = &cobra.Command{
	Use:  "validate-output",
	Short: "validate output of user batch inference container",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
