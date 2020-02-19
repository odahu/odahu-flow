/*
 * Copyright 2019 EPAM Systems
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package common

import (
	"github.com/spf13/viper"
)

const (
	// HTTP port with "/metric" endpoint
	MonitoringPort = "operator.monitoring_port"
	// API URL
	APIURL = "operator.api_url"
	// It is a mock for the future. Currently, it is always empty.
	APIToken = "operator.api_token"
	// OpenID client_id credential for service account
	ClientID = "operator.client_id"
	// OpenID client_secret credential for service account
	ClientSecret = "operator.client_secret"
	// OpenID token url
	OAuthOIDCTokenEndpoint = "operator.oauth_oidc_token_endpoint" // #nosec
)

func init() {
	viper.SetDefault(MonitoringPort, 7777)
	viper.SetDefault(APIURL, "http://localhost:5000")
	viper.SetDefault(APIToken, "")
	viper.SetDefault(ClientID, "")
	viper.SetDefault(ClientSecret, "")
	viper.SetDefault(OAuthOIDCTokenEndpoint, "")
}
