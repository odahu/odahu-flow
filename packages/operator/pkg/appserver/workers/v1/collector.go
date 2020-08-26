package v1

import (
	"database/sql"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	"github.com/odahu/odahu-flow/packages/operator/pkg/services"
	"github.com/odahu/odahu-flow/packages/operator/pkg/appserver/workers/v1/deployment"
	"github.com/odahu/odahu-flow/packages/operator/pkg/appserver/workers/v1/packaging"
	"github.com/odahu/odahu-flow/packages/operator/pkg/appserver/workers/v1/training"
	"github.com/odahu/odahu-flow/packages/operator/pkg/workersmanager"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func SetupRunners(runMgr *workersmanager.WorkersManager, kubeMgr manager.Manager, db *sql.DB, cfg config.Config) {

	kClient := kubeMgr.GetClient()
	kConfig := kubeMgr.GetConfig()

	if cfg.Training.Enabled {
		trainStorage := services.InitTrainStorage(cfg, db)
		trainService := services.InitTrainService(cfg, kClient, kConfig)

		trainWorker := NewGenericWorker(
			kubeMgr, "training", cfg.Common.LaunchPeriod,
			training.NewSyncer(trainStorage, trainService),
		)
		runMgr.AddRunnable(&trainWorker)
	}

	if cfg.Packaging.Enabled {
		packStorage := services.InitPackStorage(cfg, db)
		packService := services.InitPackService(cfg, kClient, kConfig)

		packWorker := NewGenericWorker(
			kubeMgr, "packaging", cfg.Common.LaunchPeriod,
			packaging.NewSyncer(packStorage, packService),
		)
		runMgr.AddRunnable(&packWorker)
	}

	if cfg.Deployment.Enabled {
		deployStorage := services.InitDeployStorage(cfg, db)
		deployService := services.InitDeployService(cfg, kClient)

		deployWorker := NewGenericWorker(
			kubeMgr, "deployment", cfg.Common.LaunchPeriod,
			deployment.NewSyncer(deployStorage, deployService),
		)
		runMgr.AddRunnable(&deployWorker)
	}

}
