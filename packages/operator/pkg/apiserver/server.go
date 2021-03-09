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

package apiserver

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/routes"
	v1Routes "github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/routes/v1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	"github.com/odahu/odahu-flow/packages/operator/pkg/testhelpers/testenvs"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils"
	"github.com/rakyll/statik/fs"
	"go.uber.org/multierr"
	"net/http"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	logutils "github.com/odahu/odahu-flow/packages/operator/pkg/utils/log"
)

const (
	URLPrefix = "/api/v1"
)

var log = logf.Log.WithName("api-server")

type Server struct {
	webServer   *http.Server
	kubeManager manager.Manager
}

type KubeTestServer struct {
	apiServer *Server
	closeTestKubeEnv func() error
}

func (k KubeTestServer) Run(errorCh chan<- error) {
	k.apiServer.Run(errorCh)
}

func (k KubeTestServer) Close(ctx context.Context) (err error) {
	err = k.apiServer.webServer.Close()
	err = multierr.Append(err, k.closeTestKubeEnv())
	return err
}

type ServerI interface {
	Run(errorCh chan<- error)
	Close(ctx context.Context) error
}


func setupServer(cfg *config.Config, kubeMgr manager.Manager) (*Server, error) {
	db, err := sql.Open("postgres", cfg.Common.DatabaseConnectionString)
	if err != nil {
		return nil, err
	}

	staticFS, err := fs.New()
	if err != nil {
		return nil, err
	}

	gin.SetMode(gin.ReleaseMode)

	server := gin.New()
	server.Use(gin.Recovery())

	logutils.SetupLogBindingMiddleware(server)
	routes.SetUpHealthCheck(server)
	routes.SetUpSwagger(server.Group(""), staticFS)
	routes.SetUpPrometheus(server)
	v1Group := server.Group(URLPrefix)
	err = v1Routes.SetupV1Routes(v1Group, kubeMgr, db, *cfg)

	return &Server{
		webServer: &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.API.Port),
			Handler: server,
		},
		kubeManager:    kubeMgr,
	}, err
}

func NewAPIServer(cfg *config.Config) (ServerI, error) {
	// LocalBackendType means setup test Kubernetes Control Plane
	if cfg.API.Backend.Type == config.LocalBackendType {
		log.Info("Setting local Kubernetes Control Plane for API Server")
		_, _, closeTestKube, kubeMgr, err := testenvs.SetupTestKube(cfg.API.Backend.Local.LocalBackendCRDPath)
		if err != nil {
			return nil, err
		}
		server, err := setupServer(cfg, kubeMgr)
		if err != nil {
			return nil, err
		}
		return KubeTestServer{
			apiServer:        server,
			closeTestKubeEnv: closeTestKube,
		}, nil
	}

	kubeMgr, err := utils.NewManager(ctrl.Options{
		MetricsBindAddress: "0",
		LeaderElection: false,  // We use common leader election worker for kube-manager and DB syncer goroutines
	})
	if err != nil {
		return nil, err
	}
	return setupServer(cfg, kubeMgr)
}

// @title API Gateway
// @version 1.0
// @description This is an API Gateway webServer.
// @termsOfService http://swagger.io/terms/

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
func (s *Server) Run(errorCh chan<- error) {
	
	go func() {
		if cerr := s.webServer.ListenAndServe(); cerr != nil && cerr != http.ErrServerClosed {
			log.Error(cerr, "error in web webServer")
			errorCh <- cerr
		} else {
			errorCh <- nil
		}
	}()
}

// Shutdown gracefully
func (s *Server) Close(ctx context.Context) error {
	if cerr := s.webServer.Shutdown(ctx); cerr != nil {
		log.Error(cerr, "Unable to shutdown web server")
		return cerr
	}
	log.Info("Web server is stopped")
	return nil
}
