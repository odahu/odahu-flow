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
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils"
	"github.com/rakyll/statik/fs"
	"net/http"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

const (
	URLPrefix = "/api/v1"
)

var log = logf.Log.WithName("api-server")

type Server struct {
	webServer   *http.Server
	kubeManager manager.Manager
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

	return &Server{
		webServer: &http.Server{
			Addr:    fmt.Sprintf(":%d", config.API.Port),
			Handler: server,
		},
		kubeManager:    kubeMgr,
	}, err
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
	return


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
