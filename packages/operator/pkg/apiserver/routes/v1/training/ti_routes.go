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

package training

import (
	"github.com/gin-gonic/gin"
	"github.com/odahu/odahu-flow/packages/operator/pkg/service/toolchain"
)

func ConfigureToolchainRoutes(
	routeGroup *gin.RouterGroup,
	tiService toolchain.Service,
) {

	tiController := &ToolchainIntegrationController{
		service:   tiService,
		validator: NewTiValidator(),
	}

	routeGroup.GET(GetToolchainIntegrationURL, tiController.getToolchainIntegration)
	routeGroup.GET(GetAllToolchainIntegrationURL, tiController.getAllToolchainIntegrations)
	routeGroup.POST(CreateToolchainIntegrationURL, tiController.createToolchainIntegration)
	routeGroup.PUT(UpdateToolchainIntegrationURL, tiController.updateToolchainIntegration)
	routeGroup.DELETE(DeleteToolchainIntegrationURL, tiController.deleteToolchainIntegration)
}
