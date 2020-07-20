package services

import (
	"database/sql"
	"errors"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/connection"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/connection/kubernetes"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/connection/vault"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/deployment"
	post_depl "github.com/odahu/odahu-flow/packages/operator/pkg/repository/deployment/postgres"
	kube_depl "github.com/odahu/odahu-flow/packages/operator/pkg/repository/deployment/kubernetes"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/packaging"
	kube_pack "github.com/odahu/odahu-flow/packages/operator/pkg/repository/packaging/kubernetes"
	post_pack "github.com/odahu/odahu-flow/packages/operator/pkg/repository/packaging/postgres"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/training"
	kube_train "github.com/odahu/odahu-flow/packages/operator/pkg/repository/training/kubernetes"
	post_train "github.com/odahu/odahu-flow/packages/operator/pkg/repository/training/postgres"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	connRepository connection.Repository
	trainStorage   training.Repository
	packStorage    packaging.Repository
	deployStorage  deployment.Repository

	toolchainStorage training.ToolchainRepository
	piStorage        packaging.PackagingIntegrationRepository

	trainService  training.Service
	packService   packaging.Service
	deployService deployment.Repository
)

func InitConnStorage(cfg config.Config, k8sClient client.Client) (connection.Repository, error) {

	if connRepository != nil {
		return connRepository, nil
	}

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
			return nil, err
		}
	default:
		return nil, errors.New("unexpect connection repository type")
	}

	return connRepository, nil
}

func InitToolchainStorage(_ config.Config, db *sql.DB) training.ToolchainRepository {

	if toolchainStorage != nil {
		return toolchainStorage
	}
	toolchainStorage = post_train.ToolchainRepository{DB: db}
	return toolchainStorage
}

func InitPackagingIntStorage(_ config.Config, db *sql.DB) packaging.PackagingIntegrationRepository {

	if piStorage != nil {
		return piStorage
	}
	piStorage = post_pack.PackagingIntegrationRepository{DB: db}
	return piStorage
}

func InitTrainStorage(_ config.Config, db *sql.DB) training.Repository {
	if trainStorage != nil {
		return trainStorage
	}
	trainStorage = post_train.TrainingPostgresRepo{
		DB: db,
	}
	return trainStorage
}

func InitPackStorage(_ config.Config, db *sql.DB) packaging.Repository {
	if packStorage != nil {
		return packStorage
	}
	packStorage = post_pack.PackagingPostgresRepo{DB: db}
	return packStorage
}

func InitDeployStorage(_ config.Config, db *sql.DB) deployment.Repository {
	if deployStorage != nil {
		return deployStorage
	}
	deployStorage = post_depl.DeploymentPostgresRepo{DB: db}
	return deployStorage
}


// Services that provide training, packaging and deployment

func InitTrainService(cfg config.Config, k8sClient client.Client, k8sConfig *rest.Config) training.Service {

	if trainService != nil {
		return trainService
	}
	trainService = kube_train.NewRepository(cfg.Training.Namespace,
		cfg.Training.ToolchainIntegrationNamespace,
		k8sClient,
		k8sConfig,
	)
	return trainService
}

func InitPackService(cfg config.Config, k8sClient client.Client,
	k8sConfig *rest.Config) packaging.Service {
	if packService != nil {
		return packService
	}
	packService = kube_pack.NewRepository(cfg.Packaging.Namespace,
		cfg.Packaging.PackagingIntegrationNamespace,
		k8sClient,
		k8sConfig,
	)
	return packService
}

func InitDeployService(cfg config.Config, k8sClient client.Client) deployment.Repository {
	if deployService != nil {
		return deployService
	}
	deployService = kube_depl.NewRepository(cfg.Deployment.Namespace, k8sClient)
	return deployService

}

