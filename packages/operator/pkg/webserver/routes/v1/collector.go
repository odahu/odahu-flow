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
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	connection_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/connection"
	k8s_connection_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/connection/kubernetes"
	vault_connection_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/connection/vault"
	deployment_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/deployment/kubernetes"
	packaging_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/packaging"
	k8s_packaging_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/packaging/kubernetes"
	postgres_packaging_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/packaging/postgres"
	training_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/training"
	k8s_training_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/training/kubernetes"
	postgres_training_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/training/postgres"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils"
	"github.com/odahu/odahu-flow/packages/operator/pkg/webserver/routes"
	"github.com/odahu/odahu-flow/packages/operator/pkg/webserver/routes/v1/configuration"
	"github.com/odahu/odahu-flow/packages/operator/pkg/webserver/routes/v1/connection"
	"github.com/odahu/odahu-flow/packages/operator/pkg/webserver/routes/v1/deployment"
	"github.com/odahu/odahu-flow/packages/operator/pkg/webserver/routes/v1/packaging"
	"github.com/odahu/odahu-flow/packages/operator/pkg/webserver/routes/v1/training"
	userinfo "github.com/odahu/odahu-flow/packages/operator/pkg/webserver/routes/v1/user"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	OdahuflowV1ApiVersion = "api/v1"
)

func SetupV1Routes(
	routeGroup *gin.RouterGroup,
	k8sClient client.Client,
	k8sConfig *rest.Config,
	odahuConfig config.Config) (err error) {
	var connRepository connection_repository.Repository
	var tiRepository training_repository.ToolchainRepository
	var piRepository packaging_repository.PackagingIntegrationRepository

	// Setup the connection repository
	switch odahuConfig.Connection.RepositoryType {
	case config.RepositoryKubernetesType:
		connRepository = k8s_connection_repository.NewRepository(
			odahuConfig.Connection.Namespace,
			k8sClient,
		)
	case config.RepositoryVaultType:
		connRepository, err = vault_connection_repository.NewRepositoryFromConfig(odahuConfig.Connection.Vault)
		if err != nil {
			return err
		}
	default:
		return errors.New("unexpect connection repository type")
	}

	depRepository := deployment_repository.NewRepository(odahuConfig.Deployment.Namespace, k8sClient)
	packRepository := k8s_packaging_repository.NewRepository(
		odahuConfig.Packaging.Namespace,
		odahuConfig.Packaging.PackagingIntegrationNamespace,
		k8sClient,
		k8sConfig,
	)
	trainRepository := k8s_training_repository.NewRepository(
		odahuConfig.Training.Namespace,
		odahuConfig.Training.ToolchainIntegrationNamespace,
		k8sClient,
		k8sConfig,
	)

	// Setup the training toolchain repository
	switch odahuConfig.Training.ToolchainIntegrationRepositoryType {
	case config.RepositoryKubernetesType:
		tiRepository = k8s_training_repository.NewRepository(
			odahuConfig.Training.Namespace,
			odahuConfig.Training.ToolchainIntegrationNamespace,
			k8sClient,
			k8sConfig,
		)
	case config.RepositoryPostgresType:
		db, err := sql.Open("postgres", odahuConfig.Common.DatabaseConnectionString)
		if err != nil {
			return err
		}
		tiRepository = postgres_training_repository.ToolchainRepository{DB: db}
	default:
		return errors.New("unexpect toolchain repository type")
	}

	// Setup the packaging integration repository
	switch odahuConfig.Packaging.PackagingIntegrationRepositoryType {
	case config.RepositoryKubernetesType:
		piRepository = k8s_packaging_repository.NewRepository(
			odahuConfig.Packaging.Namespace,
			odahuConfig.Packaging.PackagingIntegrationNamespace,
			k8sClient,
			k8sConfig,
		)
	case config.RepositoryPostgresType:
		db, err := sql.Open("postgres", odahuConfig.Common.DatabaseConnectionString)
		if err != nil {
			return err
		}
		piRepository = postgres_packaging_repository.PackagingIntegrationRepository{DB: db}
	default:
		return errors.New("unexpect toolchain repository type")
	}

	connection.ConfigureRoutes(routeGroup, connRepository, utils.EvaluatePublicKey, odahuConfig.Connection)
	deployment.ConfigureRoutes(
		routeGroup, depRepository, odahuConfig.Deployment, odahuConfig.Common.ResourceGPUName,
	)

	packagingRouteGroup := routeGroup.Group("", routes.DisableAPIMiddleware(odahuConfig.Packaging.Enabled))
	packaging.ConfigureRoutes(packagingRouteGroup, packRepository, piRepository, connRepository, odahuConfig.Packaging, odahuConfig.Common.ResourceGPUName)
	packaging.ConfigurePiRoutes(packagingRouteGroup, piRepository)

	trainingRouteGroup := routeGroup.Group("", routes.DisableAPIMiddleware(odahuConfig.Training.Enabled))
	training.ConfigureRoutes(trainingRouteGroup, trainRepository, tiRepository, connRepository, odahuConfig.Training, odahuConfig.Common.ResourceGPUName)
	training.ConfigureToolchainRoutes(
		trainingRouteGroup, tiRepository,
	)

	configuration.ConfigureRoutes(routeGroup, odahuConfig)
	userinfo.ConfigureRoutes(routeGroup, odahuConfig.Users.Claims)

	return err
}
