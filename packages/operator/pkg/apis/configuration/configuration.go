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

package configuration

import (
	common_conf "github.com/odahu/odahu-flow/packages/operator/pkg/config/common"
	"github.com/spf13/viper"
)

// For now it is very limited configuration.
// The main future to expose external links to Odahuflow resources.
// But it will represent the full odahuflow services configuration in the future.
// TODO: support all odahuflow configuration
type Configuration struct {
	// Common secretion of configuration
	CommonConfiguration CommonConfiguration `json:"common"`
	// Configuration describe training process
	TrainingConfiguration TrainingConfiguration `json:"training"`
}

type CommonConfiguration struct {
	// The collection of external urls, for example: metrics, edge, service catalog and so on
	ExternalURLs []ExternalUrl `json:"externalUrls"`
}

type ExternalUrl struct {
	// Human readable name
	Name string `json:"name"`
	// Link to a resource
	URL string `json:"url"`
	// Optional link to an image which represents a type of the resource, for example the logo of Grafana
	ImageURL string `json:"imageUrl"`
}

type TrainingConfiguration struct {
	MetricURL string `json:"metricUrl"`
}

func ExportExternalUrlsFromConfig() []ExternalUrl {
	// TODO: move it to a different file
	// TODO: manually mapping is a bad approach
	externalURLs := make([]ExternalUrl, 0)
	configExternalURLs := viper.Get(common_conf.ExternalURLs).([]interface{})
	for _, externalURL := range configExternalURLs {
		externalURL := externalURL.(map[interface{}]interface{})

		name := externalURL["name"].(string)
		url := externalURL["url"].(string)

		var imageURL string
		if externalURL["image_url"] != nil {
			imageURL = externalURL["image_url"].(string)
		}
		externalURLs = append(externalURLs, ExternalUrl{
			Name:     name,
			URL:      url,
			ImageURL: imageURL,
		})
	}

	return externalURLs
}
