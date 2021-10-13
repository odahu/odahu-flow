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
	job_routes "github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/routes/v1/batch/job"
	service_routes "github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/routes/v1/batch/service"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/routes/v1/configuration"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/routes/v1/connection"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/routes/v1/deployment"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/routes/v1/packaging"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/routes/v1/training"
	userinfo "github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/routes/v1/user"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	pack_kube_client "github.com/odahu/odahu-flow/packages/operator/pkg/kubeclient/packagingclient"
	train_kube_client "github.com/odahu/odahu-flow/packages/operator/pkg/kubeclient/trainingclient"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/connection/kubernetes"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/connection/memory"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/connection/vault"
	batch_service "github.com/odahu/odahu-flow/packages/operator/pkg/service/batch"
	conn_service "github.com/odahu/odahu-flow/packages/operator/pkg/service/connection"
	md_service "github.com/odahu/odahu-flow/packages/operator/pkg/service/deployment"
	mp_service "github.com/odahu/odahu-flow/packages/operator/pkg/service/packaging"
	"github.com/odahu/odahu-flow/packages/operator/pkg/service/packaging_integration"
	mr_service "github.com/odahu/odahu-flow/packages/operator/pkg/service/route"
	mt_service "github.com/odahu/odahu-flow/packages/operator/pkg/service/training"
	"github.com/odahu/odahu-flow/packages/operator/pkg/service/training_integration"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	batch_repo "github.com/odahu/odahu-flow/packages/operator/pkg/repository/batch/postgres"
	conn_repo_type "github.com/odahu/odahu-flow/packages/operator/pkg/repository/connection"
	deploy_repo "github.com/odahu/odahu-flow/packages/operator/pkg/repository/deployment/postgres"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/outbox"
	pack_repo "github.com/odahu/odahu-flow/packages/operator/pkg/repository/packaging/postgres"
	route_repo "github.com/odahu/odahu-flow/packages/operator/pkg/repository/route/postgres"
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
	case config.RepositoryMemoryType:
		connRepository = memory.NewRepository()
	default:
		return errors.New("unexpect connection repository type")
	}

	trainingIntegrationRepo := train_repo.TrainingIntegrationRepo{DB: db}
	trainingIntegrationService := training_integration.NewService(trainingIntegrationRepo)

	trainRepo := train_repo.TrainingRepo{DB: db}
	piRepo := pack_repo.PackagingIntegrationRepository{DB: db}
	piService := packaging_integration.NewService(&piRepo)
	packRepo := pack_repo.PackagingRepo{DB: db}
	deployRepo := deploy_repo.DeploymentRepo{DB: db}
	routeRepo := route_repo.RouteRepo{DB: db}
	batchServiceRepo := batch_repo.BISRepo{DB: db}
	batchJobRepo := batch_repo.BIJRepo{DB: db}

	trainKubeClient := train_kube_client.NewClient(
		cfg.Training.Namespace,
		cfg.Training.TrainingIntegrationNamespace,
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
	depService := md_service.NewService(deployRepo, routeRepo, outbox.EventPublisher{DB: db})
	mrService := mr_service.NewService(routeRepo, outbox.EventPublisher{DB: db})
	batchServiceService := batch_service.NewInferenceServiceService(batchServiceRepo)
	batchJobService := batch_service.NewJobService(batchJobRepo, batchServiceRepo, connService)

	connection.ConfigureRoutes(routeGroup, connService, utils.EvaluatePublicKey, cfg.Connection)

	mdEventGetter := outbox.DeploymentEventGetter{DB: db}
	mrEventGetter := outbox.RouteEventGetter{DB: db}

	deployment.ConfigureRoutes(routeGroup, depService, mdEventGetter, mrService, mrEventGetter,
		cfg.Deployment, cfg.Common.ResourceGPUName)
	packagingRouteGroup := routeGroup.Group("", routes.DisableAPIMiddleware(cfg.Packaging.Enabled))
	packaging.ConfigureRoutes(
		packagingRouteGroup, packKubeClient, packService,
		piService, connRepository, cfg.Packaging, cfg.Common.ResourceGPUName,
	)
	packaging.ConfigurePiRoutes(packagingRouteGroup, piService)

	trainingRouteGroup := routeGroup.Group("", routes.DisableAPIMiddleware(cfg.Training.Enabled))

	training.ConfigureRoutes(
		trainingRouteGroup,
		cfg.Training,
		cfg.Common.ResourceGPUName,
		trainService, trainingIntegrationService, connRepository, trainKubeClient)

	training.ConfigureTrainingIntegrationRoutes(
		trainingRouteGroup, trainingIntegrationService,
	)

	configuration.ConfigureRoutes(routeGroup, cfg)
	userinfo.ConfigureRoutes(routeGroup, cfg.Users.Claims)

	batchRouteGroup := routeGroup.Group("", routes.DisableAPIMiddleware(cfg.Batch.Enabled))
	service_routes.SetupRoutes(batchRouteGroup, batchServiceService)
	batchJobRouteGroup := routeGroup.Group("", routes.DisableAPIMiddleware(cfg.Batch.Enabled))
	job_routes.SetupRoutes(batchJobRouteGroup, batchJobService)

	return err
}
