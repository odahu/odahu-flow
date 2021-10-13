package controller

import (
	"database/sql"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/connection"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	"github.com/odahu/odahu-flow/packages/operator/pkg/controller/adapters/v1/batch"
	"github.com/odahu/odahu-flow/packages/operator/pkg/controller/adapters/v1/deployment"
	"github.com/odahu/odahu-flow/packages/operator/pkg/controller/adapters/v1/packaging"
	"github.com/odahu/odahu-flow/packages/operator/pkg/controller/adapters/v1/route"
	"github.com/odahu/odahu-flow/packages/operator/pkg/controller/adapters/v1/training"
	odahu_errors "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	batch_kube_client "github.com/odahu/odahu-flow/packages/operator/pkg/kubeclient/batchclient"
	deploy_kube_client "github.com/odahu/odahu-flow/packages/operator/pkg/kubeclient/deploymentclient"
	pack_kube_client "github.com/odahu/odahu-flow/packages/operator/pkg/kubeclient/packagingclient"
	train_kube_client "github.com/odahu/odahu-flow/packages/operator/pkg/kubeclient/trainingclient"
	batch_repo "github.com/odahu/odahu-flow/packages/operator/pkg/repository/batch/postgres"
	deploy_repo "github.com/odahu/odahu-flow/packages/operator/pkg/repository/deployment/postgres"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/outbox"
	pack_repo "github.com/odahu/odahu-flow/packages/operator/pkg/repository/packaging/postgres"
	route_repo "github.com/odahu/odahu-flow/packages/operator/pkg/repository/route/postgres"
	train_repo "github.com/odahu/odahu-flow/packages/operator/pkg/repository/training/postgres"
	batch_service "github.com/odahu/odahu-flow/packages/operator/pkg/service/batch"
	dep_service "github.com/odahu/odahu-flow/packages/operator/pkg/service/deployment"
	pack_service "github.com/odahu/odahu-flow/packages/operator/pkg/service/packaging"
	route_service "github.com/odahu/odahu-flow/packages/operator/pkg/service/route"
	train_service "github.com/odahu/odahu-flow/packages/operator/pkg/service/training"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// Controller does not create InferenceJobs from API Server (only fetches them)
// Thus connectionGetter is not used. We build dummy connectionGetter to init InferenceJobService
type dummyConnGetter struct {
}

func (dc *dummyConnGetter) GetConnection(id string, encrypted bool) (*connection.Connection, error) {
	return nil, odahu_errors.NotFoundError{Entity: id}
}

func SetupRunners(runMgr *WorkersManager, kubeMgr manager.Manager, db *sql.DB, cfg config.Config) {

	kClient := kubeMgr.GetClient()
	kConfig := kubeMgr.GetConfig()

	if cfg.Training.Enabled {
		trainService := train_service.NewService(train_repo.TrainingRepo{DB: db})
		trainKubeClient := train_kube_client.NewClient(
			cfg.Training.Namespace,
			cfg.Training.TrainingIntegrationNamespace,
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
		packService := pack_service.NewService(pack_repo.PackagingRepo{DB: db})
		packKubeClient := pack_kube_client.NewClient(
			cfg.Packaging.Namespace,
			cfg.Packaging.PackagingIntegrationNamespace,
			kClient,
			kConfig,
		)

		packWorker := NewGenericWorker(
			"packaging", cfg.Common.LaunchPeriod,
			packaging.NewAdapter(packService, packKubeClient, kubeMgr),
		)
		runMgr.AddRunnable(&packWorker)
	}

	if cfg.Deployment.Enabled {
		depService := dep_service.NewService(deploy_repo.DeploymentRepo{DB: db}, route_repo.RouteRepo{DB: db},
			outbox.EventPublisher{DB: db})
		deployKubeClient := deploy_kube_client.NewClient(cfg.Deployment.Namespace, kClient)

		deployWorker := NewGenericWorker(
			"deployment", cfg.Common.LaunchPeriod,
			deployment.NewAdapter(depService, deployKubeClient, kubeMgr),
		)
		runMgr.AddRunnable(&deployWorker)

		routeService := route_service.NewService(route_repo.RouteRepo{DB: db}, outbox.EventPublisher{DB: db})

		routeWorker := NewGenericWorker(
			"route", cfg.Common.LaunchPeriod,
			route.NewAdapter(routeService, deployKubeClient, kubeMgr),
		)
		runMgr.AddRunnable(&routeWorker)

	}

	if cfg.Batch.Enabled {

		connService := dummyConnGetter{}

		batchJobService := batch_service.NewJobService(
			batch_repo.BIJRepo{DB: db}, batch_repo.BISRepo{DB: db}, &connService)
		batchServiceService := batch_service.NewInferenceServiceService(batch_repo.BISRepo{DB: db})
		batchKubeClient := batch_kube_client.NewClient(kClient, cfg.Batch.Namespace, kConfig)

		batchWorker := NewGenericWorker(
			"batch", cfg.Common.LaunchPeriod,
			batch.NewAdapter(kubeMgr, batchKubeClient, batchJobService, batchServiceService),
		)
		runMgr.AddRunnable(&batchWorker)
	}

}
