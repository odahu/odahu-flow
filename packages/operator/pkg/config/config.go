//
//    Copyright 2020 EPAM Systems
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
//

package config

import (
	"fmt"
	"github.com/prometheus/common/log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

const (
	defaultConfigPathForDev = "packages/operator"
)

var (
	CfgFile string
	logC    = logf.Log.WithName("config")
)

type AuthConfig struct {
	// ODAHU API URL
	APIURL string `json:"apiUrl"`
	// It is a mock for the future. Currently, it is always empty.
	APIToken string `json:"apiToken"`
	// OpenID client_id credential for service account
	ClientID string `json:"clientId"`
	// OpenID client_secret credential for service account
	ClientSecret string `json:"clientSecret"`
	// OpenID token url
	OAuthOIDCTokenEndpoint string `json:"oauthOidcTokenEndpoint"`
}

type Config struct {
	API            APIConfig             `json:"api"`
	Common         CommonConfig          `json:"common"`
	Users          UserConfig            `json:"users"`
	Connection     ConnectionConfig      `json:"connection"`
	Deployment     ModelDeploymentConfig `json:"deployment"`
	ServiceCatalog ServiceCatalog        `json:"serviceCatalog"`
	Trainer        TrainerConfig         `json:"trainer"`
	Packager       PackagerConfig        `json:"packager"`
	Training       ModelTrainingConfig   `json:"training"`
	Packaging      ModelPackagingConfig  `json:"packaging"`
	Operator       OperatorConfig        `json:"operator"`
	Batch          BatchConfig           `json:"batch"`
}

func LoadConfig() (*Config, error) {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if CfgFile != "" {
		viper.SetConfigFile(CfgFile)
	} else {
		viper.AddConfigPath(defaultConfigPathForDev)
	}

	if err := viper.ReadInConfig(); err != nil {
		logC.Info(fmt.Sprintf("Error during reading of the odahuflow config: %s", err.Error()))
	}

	config := &Config{
		API:            NewDefaultAPIConfig(),
		Common:         NewDefaultCommonConfig(),
		Users:          NewDefaultUserConfig(),
		Connection:     NewDefaultConnectionConfig(),
		Deployment:     NewDefaultModelDeploymentConfig(),
		ServiceCatalog: NewDefaultServiceCatalogConfig(),
		Trainer:        NewDefaultTrainerConfig(),
		Packager:       NewDefaultPackagerConfig(),
		Training:       NewDefaultModelTrainingConfig(),
		Packaging:      NewDefaultModelPackagingConfig(),
		Operator:       NewDefaultOperatorConfig(),
	}

	err := viper.Unmarshal(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// Remove sensitive fields from config
func (c Config) CleanupSensitiveFields() Config {
	// We receive config by value, so it is a copy and we can modify it.
	c.Packager.Auth = AuthConfig{}
	c.Trainer.Auth = AuthConfig{}
	c.Operator.Auth = AuthConfig{}

	return c
}

func MustLoadConfig() *Config {
	config, err := LoadConfig()
	if err != nil {
		log.Error(err, "Can not load the application config")

		panic(err)
	}

	return config
}

func InitBasicParams(cmd *cobra.Command) {
	setUpLogger()

	cmd.PersistentFlags().StringVar(&CfgFile, "config", "", "config file")
}

func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func setUpLogger() {
	logf.SetLogger(logf.ZapLoggerTo(os.Stdout, false))
}

func NewDefaultConfig() *Config {
	return &Config{
		API:            NewDefaultAPIConfig(),
		Common:         NewDefaultCommonConfig(),
		Users:          NewDefaultUserConfig(),
		Connection:     NewDefaultConnectionConfig(),
		Deployment:     NewDefaultModelDeploymentConfig(),
		ServiceCatalog: NewDefaultServiceCatalogConfig(),
		Trainer:        NewDefaultTrainerConfig(),
		Packager:       NewDefaultPackagerConfig(),
		Training:       NewDefaultModelTrainingConfig(),
		Packaging:      NewDefaultModelPackagingConfig(),
		Operator:       NewDefaultOperatorConfig(),
	}
}
