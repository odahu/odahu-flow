package cmd

import (
	"fmt"
	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/predict_v2"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/feedback"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/url"
	feedback_utils "odahu-commons/feedback"
	"os"
	"strconv"
)

// Model log commands

const (
	defaultRequestTag = "request_response"
	defaultFluentdHost = "localhost"
	defaultFluentdPort = 24224
	defaultResponseTag = "response_body"
	maxRetryToDeliver = 100
	maxRetryWait      = 1000
)

var (
	requestTag string
	responseTag string
	model string
	version string
	requestID string
)

func init () {
	batchCommand.AddCommand(logCommand)
	logCommand.AddCommand(logModelInputCommand)
	logCommand.AddCommand(logModelOutputCommand)
	logCommand.PersistentFlags().StringVarP(
		&model, "model", "m", "", "ML Model name",
	)
	_ = logCommand.MarkPersistentFlagRequired("model")

	logCommand.PersistentFlags().StringVar(
		&version, "version",  "", "ML Model version",
	)
	_ = logCommand.MarkPersistentFlagRequired("version")


	logCommand.PersistentFlags().StringVarP(
		&requestID, "request-id", "r", "",
		"request id for which this request/response data is logged",
	)
	_ = logCommand.MarkFlagRequired("request-id")

	logCommand.PersistentFlags().StringVar(
		&apiURL, "fluentd", "", "fluentd base URL (schema://host:port)",
	)
	_ = viper.BindPFlag("feedback.fluentd.baseurl", logCommand.PersistentFlags().Lookup("fluentd"))

	logModelInputCommand.Flags().StringVar(&requestTag, "tag", defaultRequestTag, "tag model request")
	logModelOutputCommand.Flags().StringVar(&responseTag, "tag", defaultResponseTag, "tag model response")

}

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

func initFluentd() (*fluent.Fluent, error) {

	var host string
	var port int
	rawBaseURL := cfg.Feedback.Fluentd.BaseURL
	if rawBaseURL == "" {
		host = defaultFluentdHost
		port = defaultFluentdPort
	} else {
		baseURL, err := url.Parse(rawBaseURL)
		if err != nil {
			return nil, fmt.Errorf("unable to parse fluend base url: %s", rawBaseURL)
		}
		host = baseURL.Hostname()
		portString := baseURL.Port()
		if portString != "" {
			port, err = strconv.Atoi(portString)
			if err != nil {
				return nil, fmt.Errorf("fluentd port must be integer %s", portString)
			}
		} else {
			port = defaultFluentdPort
		}
	}

	zap.S().Infof("Connecting to fluentd using host %s and port %d", host, port)
	f, err := fluent.New(fluent.Config{
		FluentPort:   port,
		FluentHost:   host,
		MaxRetry:     maxRetryToDeliver,
		Async:        true,
		MaxRetryWait: maxRetryWait,
	})
	return f, err


}

func getRequestWrapper(modelName string, modelVersion string) func(content string)interface{}{
	return func(content string) interface{} {
		return feedback_utils.RequestResponse{
			RequestID:           requestID,
			RequestContent:      content,
			ModelVersion:        modelVersion,
			ModelName:           modelName,
		}
	}
}

func getResponseWrapper(modelName string, modelVersion string) func(content string)interface{}{
	return func(content string) interface{} {
		return feedback_utils.ResponseBody{
			RequestID:       requestID,
			ModelVersion:    modelVersion,
			ModelName:       modelName,
			ResponseContent: content,
		}
	}
}

var logModelInputCommand = &cobra.Command{
	Use:   "input",
	Short: "log model input to feedback storage",
	Args:  cobra.ExactArgs(1),
	Example: "odahu-tools batch log input <path-to-folder-with-json-files>",
	RunE: func(cmd *cobra.Command, args []string) error {
		logEngine, err := initFluentd()
		if err != nil {
			return err
		}
		defer func() {
			if err := logEngine.Close(); err != nil {
				zap.S().Errorw("Error closing fluentd", zap.Error(err))
			} else {
				zap.S().Info("Fluentd logs are flushed")
			}
		}()

		dataLogger := feedback.NewLogger(logEngine)

		modelName, modelVer := model, version

		wrap := getRequestWrapper(modelName, modelVer)

		for _, source := range args {
			zap.S().Infof("Handle %s directory", source)
			if err := dataLogger.LogDir(source, requestTag, wrap); err != nil {
				zap.S().Errorw("Error during logging model input", zap.Error(err))
				return err
			}
		}
		return nil
	},
}

var logModelOutputCommand = &cobra.Command{
	Use:                        "output",
	Short: "log model output to feedback storage",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		logEngine, err := initFluentd()
		if err != nil {
			return err
		}
		defer func() {
			if err := logEngine.Close(); err != nil {
				zap.S().Errorw("Error closing fluentd", zap.Error(err))
			} else {
				zap.S().Info("Fluentd logs are flushed")
			}
		}()

		dataLogger := feedback.NewLogger(logEngine)

		modelName, modelVer := model, version

		wrap := getResponseWrapper(modelName, modelVer)

		for _, source := range args {
			zap.S().Infof("Handle %s directory", source)
			if err := dataLogger.LogDir(source, responseTag, wrap); err != nil {
				zap.S().Errorw("Error during logging model output", zap.Error(err))
				return err
			}
		}
		return nil
	},
}