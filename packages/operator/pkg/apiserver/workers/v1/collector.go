package v1

import (
	"database/sql"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	deploy_kube_client "github.com/odahu/odahu-flow/packages/operator/pkg/kubeclient/deploymentclient"
	pack_kube_client "github.com/odahu/odahu-flow/packages/operator/pkg/kubeclient/packagingclient"
	train_kube_client "github.com/odahu/odahu-flow/packages/operator/pkg/kubeclient/trainingclient"
	deploy_repo "github.com/odahu/odahu-flow/packages/operator/pkg/repository/deployment/postgres"
	pack_repo "github.com/odahu/odahu-flow/packages/operator/pkg/repository/packaging/postgres"
	train_repo "github.com/odahu/odahu-flow/packages/operator/pkg/repository/training/postgres"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/workers/v1/deployment"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/workers/v1/packaging"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/workers/v1/training"
	"github.com/odahu/odahu-flow/packages/operator/pkg/workersmanager"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func SetupRunners(runMgr *workersmanager.WorkersManager, kubeMgr manager.Manager, db *sql.DB, cfg config.Config) {

	kClient := kubeMgr.GetClient()
	kConfig := kubeMgr.GetConfig()

	if cfg.Training.Enabled {
		trainRepo := train_repo.TrainingPostgresRepo{DB: db}
		trainKubeClient := train_kube_client.NewClient(
			cfg.Training.Namespace,
			cfg.Training.ToolchainIntegrationNamespace,
			kClient,
			kConfig,
		)

		trainWorker := NewGenericWorker(
			kubeMgr, "training", cfg.Common.LaunchPeriod,
			training.NewSyncer(trainRepo, trainKubeClient),
		)
		runMgr.AddRunnable(&trainWorker)
	}

	if cfg.Packaging.Enabled {
		packRepo := pack_repo.PackagingPostgresRepo{DB: db}
		packKubeClient := pack_kube_client.NewClient(
			cfg.Packaging.Namespace,
			cfg.Packaging.PackagingIntegrationNamespace,
			kClient,
			kConfig,
		)

		packWorker := NewGenericWorker(
			kubeMgr, "packaging", cfg.Common.LaunchPeriod,
			packaging.NewSyncer(packRepo, packKubeClient),
		)
		runMgr.AddRunnable(&packWorker)
	}

	if cfg.Deployment.Enabled {
		deployRepo := deploy_repo.DeploymentPostgresRepo{DB: db}
		deployKubeClient := deploy_kube_client.NewClient(cfg.Deployment.Namespace, kClient)

		deployWorker := NewGenericWorker(
			kubeMgr, "deployment", cfg.Common.LaunchPeriod,
			deployment.NewSyncer(deployRepo, deployKubeClient),
		)
		runMgr.AddRunnable(&deployWorker)
	}

}
