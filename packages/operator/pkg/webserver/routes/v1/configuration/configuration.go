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
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/configuration"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config/training"
	"github.com/odahu/odahu-flow/packages/operator/pkg/webserver/routes"
	"github.com/spf13/viper"
)

const (
	GetConfigurationURL    = "/configuration"
	UpdateConfigurationURL = "/configuration"
)

func ConfigureRoutes(routeGroup *gin.RouterGroup) {
	routeGroup.GET(GetConfigurationURL, getConfiguration)
	routeGroup.PUT(UpdateConfigurationURL, updateConfiguration)
}

// @Summary Get the Odahuflow service configuration
// @Description Get the Odahuflow service configuration
// @Tags Configuration
// @Accept  json
// @Produce  json
// @Success 200 {object} configuration.Configuration
// @Router /api/v1/configuration [get]
func getConfiguration(c *gin.Context) {
	c.JSON(http.StatusOK, &configuration.Configuration{
		CommonConfiguration: configuration.CommonConfiguration{
			ExternalURLs: configuration.ExportExternalUrlsFromConfig(),
		},
		TrainingConfiguration: configuration.TrainingConfiguration{
			MetricURL: viper.GetString(training.MetricURL),
		},
	})
}

// @Summary Update a Odahuflow service configuration
// @Description Update a Configuration
// @Tags Configuration
// @Param configuration body configuration.Configuration true "Create a Configuration"
// @Accept  json
// @Produce  json
// @Success 200 {object} routes.HTTPResult
// @Router /api/v1/configuration [put]
func updateConfiguration(c *gin.Context) {
	// TODO: find the best way to implement it

	c.JSON(http.StatusOK, routes.HTTPResult{Message: "This is stub for now"})
}
