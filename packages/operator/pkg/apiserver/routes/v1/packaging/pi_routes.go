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

package packaging

import (
	"github.com/gin-gonic/gin"
)

func ConfigurePiRoutes(routeGroup *gin.RouterGroup, piService packagingIntegrationService) {

	piController := &PackagingIntegrationController{
		service:   piService,
		validator: NewPiValidator(),
	}

	routeGroup.GET(GetPackagingIntegrationURL, piController.getPackagingIntegration)
	routeGroup.GET(GetAllPackagingIntegrationURL, piController.getAllPackagingIntegrations)
	routeGroup.POST(CreatePackagingIntegrationURL, piController.createPackagingIntegration)
	routeGroup.PUT(UpdatePackagingIntegrationURL, piController.updatePackagingIntegration)
	routeGroup.DELETE(DeletePackagingIntegrationURL, piController.deletePackagingIntegration)
}
