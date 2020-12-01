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
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/routes"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/routes/v1/configuration"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/routes/v1/connection"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/routes/v1/deployment"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/routes/v1/packaging"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/routes/v1/training"
	userinfo "github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/routes/v1/user"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/connection/kubernetes"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/connection/vault"
	conn_service "github.com/odahu/odahu-flow/packages/operator/pkg/service/connection"
	md_service "github.com/odahu/odahu-flow/packages/operator/pkg/service/deployment"
	mr_service "github.com/odahu/odahu-flow/packages/operator/pkg/service/route"
	mt_service "github.com/odahu/odahu-flow/packages/operator/pkg/service/training"
	mp_service "github.com/odahu/odahu-flow/packages/operator/pkg/service/packaging"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	pack_kube_client "github.com/odahu/odahu-flow/packages/operator/pkg/kubeclient/packagingclient"
	train_kube_client "github.com/odahu/odahu-flow/packages/operator/pkg/kubeclient/trainingclient"

	conn_repo_type "github.com/odahu/odahu-flow/packages/operator/pkg/repository/connection"
	deploy_repo "github.com/odahu/odahu-flow/packages/operator/pkg/repository/deployment/postgres"
	route_repo "github.com/odahu/odahu-flow/packages/operator/pkg/repository/route/postgres"
	pack_repo "github.com/odahu/odahu-flow/packages/operator/pkg/repository/packaging/postgres"
	train_repo "github.com/odahu/odahu-flow/packages/operator/pkg/repository/training/postgres"
)

func SetupV1Routes(routeGroup *gin.RouterGroup, kubeMgr manager.Manager, db *sql.DB, cfg config.Config) (err error) {

	k8sClient := kubeMgr.GetClient()
	k8sConfig := kubeMgr.GetConfig()

	var connRepository conn_repo_type.Repository
	switch cfg.Connection.RepositoryType {
	case config.RepositoryKubernetesType:
		connRepository = kubernetes.NewRepository(
			cfg.Connection.Namespace,
			k8sClient,
		)
	case config.RepositoryVaultType:
		var err error
		connRepository, err = vault.NewRepositoryFromConfig(cfg.Connection.Vault)
		if err != nil {
			return err
		}
	default:
		return errors.New("unexpect connection repository type")
	}


	toolchainRepo := train_repo.ToolchainRepo{DB: db}
	trainRepo := train_repo.TrainingRepo{DB: db}
	piRepo := pack_repo.PackagingIntegrationRepository{DB: db}
	packRepo := pack_repo.PackagingRepo{DB: db}
	deployRepo := deploy_repo.DeploymentRepo{DB: db}
	routeRepo := route_repo.RouteRepo{DB: db}

	trainKubeClient := train_kube_client.NewClient(
		cfg.Training.Namespace,
		cfg.Training.ToolchainIntegrationNamespace,
		k8sClient,
		k8sConfig,
	)
	packKubeClient := pack_kube_client.NewClient(
		cfg.Packaging.Namespace,
		cfg.Packaging.PackagingIntegrationNamespace,
		k8sClient,
		k8sConfig,
	)

	connService := conn_service.NewService(connRepository)
	trainService := mt_service.NewService(trainRepo)
	packService := mp_service.NewService(packRepo)
	depService := md_service.NewService(deployRepo, routeRepo)
	mrService := mr_service.NewService(routeRepo)


	connection.ConfigureRoutes(routeGroup, connService, utils.EvaluatePublicKey, cfg.Connection)

	deployment.ConfigureRoutes(routeGroup, depService, mrService, nil, cfg.Deployment, cfg.Common.ResourceGPUName)
	packagingRouteGroup := routeGroup.Group("", routes.DisableAPIMiddleware(cfg.Packaging.Enabled))
	packaging.ConfigureRoutes(
		packagingRouteGroup, packKubeClient, packService,
		piRepo, connRepository, cfg.Packaging, cfg.Common.ResourceGPUName,
	)
	packaging.ConfigurePiRoutes(packagingRouteGroup, piRepo)

	trainingRouteGroup := routeGroup.Group("", routes.DisableAPIMiddleware(cfg.Training.Enabled))

	training.ConfigureRoutes(
		trainingRouteGroup,
		cfg.Training,
		cfg.Common.ResourceGPUName,
		trainService, toolchainRepo, connRepository, trainKubeClient)

	training.ConfigureToolchainRoutes(
		trainingRouteGroup, toolchainRepo,
	)

	configuration.ConfigureRoutes(routeGroup, cfg)
	userinfo.ConfigureRoutes(routeGroup, cfg.Users.Claims)

	return err
}
