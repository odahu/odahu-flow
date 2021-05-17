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
	mt_kube_client "github.com/odahu/odahu-flow/packages/operator/pkg/kubeclient/trainingclient"
	conn_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/connection"
	"github.com/odahu/odahu-flow/packages/operator/pkg/service/toolchain"
	mt_service "github.com/odahu/odahu-flow/packages/operator/pkg/service/training"
)

func ConfigureRoutes(
	routeGroup *gin.RouterGroup,
	config config.ModelTrainingConfig,
	gpuResourceName string,
	trainService mt_service.Service,
	toolchainService toolchain.Service,
	connRepo conn_repository.Repository,
	trainKubeClient mt_kube_client.Client) {

	mtController := ModelTrainingController{
		trainService: trainService,
		kubeClient:   trainKubeClient,
		validator: NewMtValidator(
			toolchainService,
			connRepo,
			config,
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
