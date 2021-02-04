package feedback

import (
	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/spf13/viper"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	maxRetryToDeliver = 5
	maxRetryWait      = 1000
)

var logger = log.Log.WithName("config")

func NewDataLogger() (DataLogging, error) {
	host := viper.GetString(CfgFluentdHost)
	port := viper.GetInt(CfgFluentdPort)

	logger.Info("Connecting to FluentD", "host", host, "port", port)

	return fluent.New(fluent.Config{
		FluentPort:   port,
		FluentHost:   host,
		MaxRetry:     maxRetryToDeliver,
		Async:        true,
		MaxRetryWait: maxRetryWait,
	})
}
