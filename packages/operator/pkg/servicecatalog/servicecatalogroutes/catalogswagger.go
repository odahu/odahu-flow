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

package servicecatalogroutes

import (
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/swagger"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/odahu/odahu-flow/packages/operator/docs" //nolint
)

// SetUpCatalogSwagger serves swagger documentation that combine all deployed models on a single page
func SetUpCatalogSwagger(rg *gin.RouterGroup, apiStaticFS http.FileSystem, reader swagger.DefinitionReader) {
	rg.GET("/catalog/*any", swagger.Handler(apiStaticFS, reader))
}
