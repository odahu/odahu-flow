package cmd

import (
	predict_v2_tools "github.com/odahu/odahu-flow/packages/operator/pkg/tools/predict_v2"
	"github.com/spf13/cobra"
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
	cpuprofile string
	model string
)

func init() {
	rootCmd.AddCommand(batchCommand)
	batchCommand.AddCommand(validateCommand)
	batchCommand.AddCommand(logCommand)


	validateCommand.AddCommand(validateInputCommand)
	validateCommand.AddCommand(validateOutputCommand)
	validateCommand.AddCommand(validateModelCommand)
	validateCommand.PersistentFlags().StringVarP(
		&source, "source", "s", ".", inputCommandUsage,
	)
	validateCommand.PersistentFlags().StringVarP(
		&destination, "destination", "d", ".", outputCommandUsage,
	)

	logCommand.AddCommand(logModelInputCommand)
	logCommand.AddCommand(logModelOutputCommand)
	logCommand.PersistentFlags().StringVarP(
		&model, "model", "m", ".", "directory with ML model files",
	)
	_ = logCommand.MarkFlagRequired("model")
}

var batchCommand = &cobra.Command{
	Use:  "batch",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			os.Exit(0)
		}
	},
}

var validateCommand = &cobra.Command{
	Use:  "validate",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			os.Exit(0)
		}
	},
}


var validateInputCommand = &cobra.Command{
	Use:  "validate-input",
	Short: "validate input for user batch inference container",
	Run: func(cmd *cobra.Command, args []string) {

		_, err := predict_v2_tools.ValidateDir(source, predict_v2_tools.ValidateInput, &destination)
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

		_, err := predict_v2_tools.ValidateDir(source, predict_v2_tools.ValidateOutput, &destination)
		if err != nil {
			zap.S().Errorw("There are errors during validation", zap.Error(err))
			os.Exit(1)
		}
	},
}

var validateModelCommand = &cobra.Command{
	Use:                        "validate-model",
	Short: "validate model for batch inference",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

