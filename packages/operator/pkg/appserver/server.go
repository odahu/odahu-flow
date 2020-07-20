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

package appserver

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/odahu/odahu-flow/packages/operator/pkg/appserver/routes"
	v1Routes "github.com/odahu/odahu-flow/packages/operator/pkg/appserver/routes/v1"
	v1Runners "github.com/odahu/odahu-flow/packages/operator/pkg/appserver/workers/v1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils"
	"github.com/odahu/odahu-flow/packages/operator/pkg/workersmanager"
	"github.com/rakyll/statik/fs"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"net/http"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sync"
	"time"
)

const (
	URLPrefix = "/api/v1"
	leaderElectionLockName = "odahu-flow-api-le-lock"

	defaultLeaseDuration = 15 * time.Second
	defaultRenewDeadline = 10 * time.Second
	defaultRetryPeriod   = 2 * time.Second
)

var log = logf.Log.WithName("api-server")

type Server struct {
	webServer   *http.Server
	runManager  *workersmanager.WorkersManager
	kubeManager manager.Manager

	kubeMgrStop    chan struct{}
	kubeMgrStopped chan struct{}

	leNamespace  string
}

func NewAPIServer(config *config.Config) (*Server, error) {

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

	staticFS, err := fs.New()
	if err != nil {
		return nil, err
	}

	gin.SetMode(gin.ReleaseMode)

	server := gin.Default()

	routes.SetUpHealthCheck(server)
	routes.SetUpSwagger(server.Group(""), staticFS)
	routes.SetUpPrometheus(server)
	v1Group := server.Group(URLPrefix)
	err = v1Routes.SetupV1Routes(v1Group, kubeMgr, db, *config)

	runManager := workersmanager.NewWorkersManager()

	// Workers can be disabled
	if !config.API.DisableWorkers {
		v1Runners.SetupRunners(runManager, kubeMgr, db, *config)
	} else {
		log.Info("Sync workers are disabled. " +
			"External services will not be synced with storage by current process")
	}

	return &Server{
		webServer: &http.Server{
			Addr:    fmt.Sprintf(":%d", config.API.Port),
			Handler: server,
		},
		runManager:     runManager,
		kubeManager:    kubeMgr,
		kubeMgrStopped: make(chan struct{}),
		kubeMgrStop:    make(chan struct{}),
		leNamespace:    "odahu-flow",
	}, err
}

// @title API Gateway
// @version 1.0
// @description This is an API Gateway webServer.
// @termsOfService http://swagger.io/terms/

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
func (s *Server) Run(errorCh chan<- error) error {


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

	leaderElector, err := leaderelection.NewLeaderElector(leaderelection.LeaderElectionConfig{
		Lock:          lock,
		LeaseDuration: defaultLeaseDuration,
		RenewDeadline: defaultRenewDeadline,
		RetryPeriod:   defaultRetryPeriod,
		Callbacks:     leaderelection.LeaderCallbacks{
			OnStartedLeading: func(_ context.Context) {
				go func() {
					if cerr := s.runManager.Run(); cerr != nil {
						log.Error(cerr, "error in runners manager")
						errorCh <- cerr
					} else {
						errorCh <- nil
					}
				}()
				go func() {
					if cerr := s.kubeManager.Start(s.kubeMgrStop); cerr != nil {
						log.Error(cerr, "error in kubernetes manager")
						errorCh <- cerr
					} else {
						errorCh <- nil
					}
					close(s.kubeMgrStopped) // We need a way to notify Shutdown function about completed stop
				}()
			},
			OnStoppedLeading: func() {
				errorCh <- fmt.Errorf("leader election lost")
			},
		},
	})
	if err != nil {
		return err
	}

	go leaderElector.Run(context.Background())

	go func() {
		if cerr := s.webServer.ListenAndServe(); cerr != nil && cerr != http.ErrServerClosed {
			log.Error(cerr, "error in web webServer")
			errorCh <- cerr
		} else {
			errorCh <- nil
		}
	}()


	return nil


}

// Shutdown gracefully
func (s *Server) Close(ctx context.Context) error {
	wg := sync.WaitGroup{}
	wg.Add(3)
	var err error
	go func() {
		defer wg.Done()
		if cerr := s.webServer.Shutdown(ctx); cerr != nil {
			log.Error(cerr, "Unable to shutdown web server")
			err = cerr
		} else {
			log.Info("Web server is stopped")
		}
	}()

	go func() {
		defer wg.Done()
		if cerr := s.runManager.Shutdown(ctx); cerr != nil {
			log.Error(cerr, "Unable to shutdown runners manager")
			err = cerr
		} else {
			log.Info("Runners manager is stopped")
		}
	}()

	go func() {
		defer wg.Done()
		close(s.kubeMgrStop)
		select {
		case <-s.kubeMgrStopped:
			log.Info("Kube manager is stopped")
		case <-ctx.Done():
			log.Error(ctx.Err(), "Unable to shutdown kube manager")
			err = ctx.Err()
		}

	}()

	wg.Wait()
	return err
}
