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

package v1

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	di "github.com/odahu/odahu-flow/packages/operator/pkg/services"
	conn_service "github.com/odahu/odahu-flow/packages/operator/pkg/service/connection"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils"
	"github.com/odahu/odahu-flow/packages/operator/pkg/appserver/routes"
	"github.com/odahu/odahu-flow/packages/operator/pkg/appserver/routes/v1/configuration"
	"github.com/odahu/odahu-flow/packages/operator/pkg/appserver/routes/v1/connection"
	"github.com/odahu/odahu-flow/packages/operator/pkg/appserver/routes/v1/deployment"
	"github.com/odahu/odahu-flow/packages/operator/pkg/appserver/routes/v1/packaging"
	"github.com/odahu/odahu-flow/packages/operator/pkg/appserver/routes/v1/training"
	userinfo "github.com/odahu/odahu-flow/packages/operator/pkg/appserver/routes/v1/user"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func SetupV1Routes(routeGroup *gin.RouterGroup, kubeMgr manager.Manager, db *sql.DB, cfg config.Config) (err error) {

	k8sClient := kubeMgr.GetClient()
	k8sConfig := kubeMgr.GetConfig()

	connStorage, err := di.InitConnStorage(cfg, k8sClient)
	if err != nil {
		return err
	}
	connService := conn_service.NewService(connStorage)
	toolchainStorage := di.InitToolchainStorage(cfg, db)
	piStorage := di.InitPackagingIntStorage(cfg, db)
	trainStorage := di.InitTrainStorage(cfg, db)
	packStorage := di.InitPackStorage(cfg, db)
	deployStorage := di.InitDeployStorage(cfg, db)

	trainService := di.InitTrainService(cfg, k8sClient, k8sConfig)
	packService := di.InitPackService(cfg, k8sClient, k8sConfig)
	deployService := di.InitDeployService(cfg, k8sClient)


	connection.ConfigureRoutes(routeGroup, connService, utils.EvaluatePublicKey, cfg.Connection)

	deployment.ConfigureRoutes(routeGroup, deployStorage, deployService, cfg.Deployment, cfg.Common.ResourceGPUName)
	packagingRouteGroup := routeGroup.Group("", routes.DisableAPIMiddleware(cfg.Packaging.Enabled))
	packaging.ConfigureRoutes(
		packagingRouteGroup, packService, packStorage,
		piStorage, connStorage, cfg.Packaging, cfg.Common.ResourceGPUName,
	)
	packaging.ConfigurePiRoutes(packagingRouteGroup, piStorage)

	trainingRouteGroup := routeGroup.Group("", routes.DisableAPIMiddleware(cfg.Training.Enabled))
	training.ConfigureRoutes(
		trainingRouteGroup, cfg.Training, cfg.Common.ResourceGPUName, training.Services{
			ToolchainStorage: toolchainStorage,
			TrainService:     trainService,
			Storage:          trainStorage,
			ConnStorage:      connStorage,
		})

	training.ConfigureToolchainRoutes(
		trainingRouteGroup, toolchainStorage,
	)

	configuration.ConfigureRoutes(routeGroup, cfg)
	userinfo.ConfigureRoutes(routeGroup, cfg.Users.Claims)

	return err
}
