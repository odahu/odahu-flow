package controller

import (
	"database/sql"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	deploy_kube_client "github.com/odahu/odahu-flow/packages/operator/pkg/kubeclient/deploymentclient"
	pack_kube_client "github.com/odahu/odahu-flow/packages/operator/pkg/kubeclient/packagingclient"
	train_kube_client "github.com/odahu/odahu-flow/packages/operator/pkg/kubeclient/trainingclient"
	deploy_repo "github.com/odahu/odahu-flow/packages/operator/pkg/repository/deployment/postgres"
	pack_repo "github.com/odahu/odahu-flow/packages/operator/pkg/repository/packaging/postgres"
	train_repo "github.com/odahu/odahu-flow/packages/operator/pkg/repository/training/postgres"
	train_service "github.com/odahu/odahu-flow/packages/operator/pkg/service/training"
	"github.com/odahu/odahu-flow/packages/operator/pkg/controller/adapters/v1/deployment"
	"github.com/odahu/odahu-flow/packages/operator/pkg/controller/adapters/v1/packaging"
	"github.com/odahu/odahu-flow/packages/operator/pkg/controller/adapters/v1/training"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func SetupRunners(runMgr *WorkersManager, kubeMgr manager.Manager, db *sql.DB, cfg config.Config) {

	kClient := kubeMgr.GetClient()
	kConfig := kubeMgr.GetConfig()

	if cfg.Training.Enabled {
		trainService := train_service.NewService(train_repo.TrainingRepo{DB: db}, db)
		trainKubeClient := train_kube_client.NewClient(
			cfg.Training.Namespace,
			cfg.Training.ToolchainIntegrationNamespace,
			kClient,
			kConfig,
		)

		trainWorker := NewGenericWorker(
			"training", cfg.Common.LaunchPeriod,
			training.NewAdapter(trainService, trainKubeClient, kubeMgr),
		)
		runMgr.AddRunnable(&trainWorker)
	}

	if cfg.Packaging.Enabled {
		packRepo := pack_repo.PackagingRepo{DB: db}
		packKubeClient := pack_kube_client.NewClient(
			cfg.Packaging.Namespace,
			cfg.Packaging.PackagingIntegrationNamespace,
			kClient,
			kConfig,
		)

		packWorker := NewGenericWorker(
			"packaging", cfg.Common.LaunchPeriod,
			packaging.NewAdapter(packRepo, packKubeClient, kubeMgr),
		)
		runMgr.AddRunnable(&packWorker)
	}

	if cfg.Deployment.Enabled {
		deployRepo := deploy_repo.DeploymentRepo{DB: db}
		deployKubeClient := deploy_kube_client.NewClient(cfg.Deployment.Namespace, kClient)

		deployWorker := NewGenericWorker(
			"deployment", cfg.Common.LaunchPeriod,
			deployment.NewAdapter(deployRepo, deployKubeClient, kubeMgr),
		)
		runMgr.AddRunnable(&deployWorker)
	}

}
