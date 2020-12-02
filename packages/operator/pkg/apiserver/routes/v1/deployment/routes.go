//
//    Copyright 2019 EPAM Systems
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

package deployment

import (
	"github.com/gin-gonic/gin"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/routes"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	md_service "github.com/odahu/odahu-flow/packages/operator/pkg/service/deployment"
	mr_service "github.com/odahu/odahu-flow/packages/operator/pkg/service/route"
)

func ConfigureRoutes(routeGroup *gin.RouterGroup, mdService md_service.Service, mrService mr_service.Service,
	mrEventsReader RoutesEventGetter, deploymentConfig config.ModelDeploymentConfig, gpuResourceName string, ) {

	mdController := ModelDeploymentController{
		mdService:   mdService,
		mdValidator: NewModelDeploymentValidator(deploymentConfig, gpuResourceName),
	}
	routeGroup = routeGroup.Group("", routes.DisableAPIMiddleware(deploymentConfig.Enabled))

	routeGroup.GET(GetModelDeploymentURL, mdController.getMD)
	routeGroup.GET(GetAllModelDeploymentURL, mdController.getAllMDs)
	routeGroup.POST(CreateModelDeploymentURL, mdController.createMD)
	routeGroup.PUT(UpdateModelDeploymentURL, mdController.updateMD)
	routeGroup.DELETE(DeleteModelDeploymentURL, mdController.deleteMD)
	routeGroup.GET(GetModelDeploymentDefaultRouteURL, mdController.getDefaultRoute)

	mrController := ModelRouteController{
		service: mrService,
		validator:        NewMrValidator(mdService),
		eventsReader: mrEventsReader,
	}
	routeGroup.GET(GetModelRouteURL, mrController.getMR)
	routeGroup.GET(GetAllModelRouteURL, mrController.getAllMRs)
	routeGroup.POST(CreateModelRouteURL, mrController.createMR)
	routeGroup.PUT(UpdateModelRouteURL, mrController.updateMR)
	routeGroup.DELETE(DeleteModelRouteURL, mrController.deleteMR)
	routeGroup.GET(EventsModelRouteURL, mrController.getRouteEvents)
}
