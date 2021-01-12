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

package servicecatalog

import (
	"github.com/gin-gonic/gin"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	"github.com/odahu/odahu-flow/packages/operator/pkg/servicecatalog/servicecatalogroutes"
	_ "github.com/odahu/odahu-flow/packages/operator/pkg/static" //nolint
	"github.com/rakyll/statik/fs"
	"net/http"
)


// @title Service Catalog
// @version 1.0
// @description Service catalog serves information about deployed models
// @termsOfService http://swagger.io/terms/

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
func SetUPMainServer(
	mrc *ModelRouteCatalog,
	config config.ServiceCatalog,
) (*http.Server, error) {
	staticFS, err := fs.New()
	if err != nil {
		return nil, err
	}

	router := gin.Default()
	rootRouteGroup := router.Group(config.BaseURL)

	servicecatalogroutes.SetUpCatalogSwagger(rootRouteGroup, staticFS, mrc.ProcessSwaggerJSON)
	servicecatalogroutes.SetUpSwagger(rootRouteGroup, staticFS)
	servicecatalogroutes.SetupDeployedModelRoute(rootRouteGroup, mrc.GetDeployedModel)
	servicecatalogroutes.SetUpHealthCheck(router)


	server := &http.Server{
		Addr:    ":5000",
		Handler: router,
	}

	return server, nil
}
