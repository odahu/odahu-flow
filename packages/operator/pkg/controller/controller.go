package controller

import (
	"context"
	"database/sql"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sync"
	"time"
)

const (
	leaderElectionLockName = "odahu-flow-sys-controller-le-lock"

	defaultLeaseDuration = 15 * time.Second
	defaultRenewDeadline = 10 * time.Second
	defaultRetryPeriod   = 2 * time.Second
)

var log = logf.Log.WithName("odahu-system-controller")

type Controller struct {
	runManager  *WorkersManager
	kubeManager manager.Manager
	leNamespace  string
}

func NewController(config *config.Config) (*Controller, error) {

	kubeMgr, err := utils.NewManager(ctrl.Options{
		MetricsBindAddress: "0",
		LeaderElection: false,  // We use common leader election worker for kube-manager and DB syncer goroutines
	})
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("postgres", config.Common.DatabaseConnectionString)
	if err != nil {
		return nil, err
	}

	runManager := NewWorkersManager()

	SetupRunners(runManager, kubeMgr, db, *config)

	return &Controller{
		runManager:  runManager,
		kubeManager: kubeMgr,
		leNamespace: "odahu-flow",
	}, err
}

func (s *Controller) Run(ctx context.Context) error {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Construct client for leader election
	client, err := kubernetes.NewForConfig(rest.AddUserAgent(s.kubeManager.GetConfig(), "leader-election"))
	if err != nil {
		return err
	}
	// Leader id, needs to be unique
	id, err := os.Hostname()
	if err != nil {
		return err
	}
	id = id + "_" + string(uuid.NewUUID())

	lock, err := resourcelock.New(resourcelock.ConfigMapsResourceLock,
		s.leNamespace,
		leaderElectionLockName,
		client.CoreV1(),
		client.CoordinationV1(),
		resourcelock.ResourceLockConfig{
			Identity: id,
		})

	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	leaderElector, err := leaderelection.NewLeaderElector(leaderelection.LeaderElectionConfig{
		Lock:          lock,
		LeaseDuration: defaultLeaseDuration,
		RenewDeadline: defaultRenewDeadline,
		RetryPeriod:   defaultRetryPeriod,
		Callbacks:     leaderelection.LeaderCallbacks{
			OnStartedLeading: func(ctx context.Context) {
				kubeMgrStop := make(chan struct{})

				wg.Add(1)
				go func() {
					defer wg.Done()
					if cerr := s.runManager.Run(ctx); cerr != nil {
						log.Error(cerr, "Run manager is finished with error")
						cancel()
					} else {
						log.Info("Run Manager is finished")
					}
				}()

				go func() {
					<-ctx.Done()
					log.Info("Context canceled. Send signal to stop kube manager")
					close(kubeMgrStop)

				}()
				wg.Add(1)
				go func() {
					defer wg.Done()
					if cerr := s.kubeManager.Start(kubeMgrStop); cerr != nil {
						log.Error(cerr, "Kube manager is finished with error")
						cancel()
					} else {
						log.Info("Kube manager is finished")
					}
				}()
			},
			OnStoppedLeading: func() {
				log.Info("Stop leading...")
				log.Info("Wait until leader election runners will be stopped")
				wg.Wait()
				log.Info("Leader election runners were stopped successfully")
			},
		},
	})
	if err != nil {
		return err
	}

	leaderElector.Run(ctx)
	return nil
}
