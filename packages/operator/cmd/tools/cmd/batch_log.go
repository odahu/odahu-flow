package cmd

import (
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/predict_v2"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
)

// Model log commands

type ModelOutputLogger interface {
	Log(requestID string, request predict_v2.InferenceResponse) error
}

var logCommand = &cobra.Command{
	Use:  "log",
	Short: "Catch model input or output from json files to fluentd service",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			os.Exit(0)
		}
	},
}


type ModelInputLogger interface {
	Log(requestID string, request predict_v2.InferenceRequest) error
}

var logModelInputCommand = &cobra.Command{
	Use:                        "input",
	Short: "log model input to feedback storage",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			return
		}
		inputLocation := args[0]
		zap.S().Info(inputLocation)
	},
}
var logModelOutputCommand = &cobra.Command{
	Use:                        "output",
	Short: "log model output to feedback storage",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

	},
}