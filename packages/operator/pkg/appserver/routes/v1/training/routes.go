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

package training

import (
	"github.com/gin-gonic/gin"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	conn_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/connection"
	mt_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/training"
)

// Services that should be injected to training webserver
type Services struct {
	// Repository of toolchain integrations
	ToolchainStorage mt_repository.ToolchainRepository
	// Repository of model trainings
	Storage mt_repository.Repository
	// Training Compute Service
	TrainService mt_repository.Service
	// Repository of connections
	ConnStorage conn_repository.Repository
}

func ConfigureRoutes(
	routeGroup *gin.RouterGroup,
	config config.ModelTrainingConfig,
	gpuResourceName string, services Services) {

	mtController := ModelTrainingController{
		trainStorage: services.Storage,
		trainService: services.TrainService,
		validator: NewMtValidator(
			services.ToolchainStorage,
			services.ConnStorage,
			config.DefaultResources,
			config.OutputConnectionID,
			gpuResourceName,
		),
	}

	routeGroup.GET(GetModelTrainingURL, mtController.getMT)
	routeGroup.GET(GetAllModelTrainingURL, mtController.getAllMTs)
	routeGroup.GET(GetModelTrainingLogsURL, mtController.getModelTrainingLog)
	routeGroup.POST(CreateModelTrainingURL, mtController.createMT)
	routeGroup.PUT(UpdateModelTrainingURL, mtController.updateMT)
	routeGroup.PUT(SaveModelTrainingResultURL, mtController.saveMTResult)
	routeGroup.DELETE(DeleteModelTrainingURL, mtController.deleteMT)

}
