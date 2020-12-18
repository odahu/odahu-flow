/*
 * Copyright 2019 EPAM Systems
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"fmt"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apiclient/event"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/util/http"
	"github.com/odahu/odahu-flow/packages/operator/pkg/servicecatalog"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"log"
	nethttp "net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)


func initReflector(cfg config.ServiceCatalog, logger *zap.SugaredLogger,
	catalog servicecatalog.Catalog) (servicecatalog.Reflector, error) {
	aCfg := cfg.Auth
	httpClient := http.NewBaseAPIClient(
		aCfg.APIURL, aCfg.APIToken, aCfg.ClientID, aCfg.ClientSecret, aCfg.OAuthOIDCTokenEndpoint, "api/v1",
	)

	edgeURL, err := url.Parse(cfg.EdgeURL)
	if err != nil {
		return servicecatalog.Reflector{}, err
	}

	return servicecatalog.NewReflector(logger, servicecatalog.UpdateHandler{
		Discoverers: []servicecatalog.ModelServerDiscoverer{
			servicecatalog.OdahuMLServerDiscoverer{
				EdgeURL:    *edgeURL,
				EdgeHost:   cfg.EdgeHost,
				HTTPClient: &httpClient,
			},
		},
		Catalog:     catalog,
	}, servicecatalog.RouteEventFetcher{
		APIClient: event.ModelRouteEventClient{
			HTTPClient: &httpClient,
			Log:        logger,
		},
	},
	servicecatalog.ReflectorOpts{
		WorkersCount: cfg.WorkersCount,
		FetchTimeout: time.Duration(cfg.FetchTimeout) * time.Second,
	}), nil

}


var mainCmd = &cobra.Command{
	Use:   "service-catalog",
	Short: "Odahu-flow service catalog server",
	Run: func(cmd *cobra.Command, args []string) {
		odahuConfig := config.MustLoadConfig()

		// Initialize

		var err error
		var logger *zap.Logger
		if odahuConfig.ServiceCatalog.Debug {
			logger, err = zap.NewDevelopment()
		} else {
			logger, err = zap.NewProduction(zap.WithCaller(false))
		}
		if err != nil {
			log.Fatal("Unable to initialize logger")
		}

		sLogger := logger.Sugar().With("OdahuVersion", odahuConfig.Common.Version)

		if odahuConfig.ServiceCatalog.Debug {
			sLogger.Info("Debug mode is enabled")
		}

		routeCatalog := servicecatalog.NewModelRouteCatalog(sLogger)


		// Run reflector (keep state of service catalog up to date with ODAHU API Server)

		reflector, err := initReflector(odahuConfig.ServiceCatalog, sLogger, routeCatalog)
		if err != nil {
			sLogger.Fatalf("Unable set up service-catalog reflector. Error %v", err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		wg := sync.WaitGroup{}
		wg.Add(2)
		go func() {
			defer wg.Done()
			sLogger.Info("Starting the reflector.")
			if err := reflector.Run(ctx); err != nil {
				sLogger.Errorw("Reflector was stopped with errors", zap.Error(err))
			} else {
				sLogger.Info("Reflector was stopped gracefully")
			}
			cancel()
		}()

		// Run webserver. API for getting information about deployed models, swaggers, metadata, etc

		mainServer, err := servicecatalog.SetUPMainServer(routeCatalog, odahuConfig.ServiceCatalog)
		if err != nil {
			sLogger.Fatalf("Unable set up service-catalog server. Error %v", err)
		}
		go func() {
			defer wg.Done()
			if err := mainServer.ListenAndServe(); err != nil && err != nethttp.ErrServerClosed {
				sLogger.Errorw("Web server was stopped with errors", zap.Error(err))
			} else {
				sLogger.Info("Reflector was stopped gracefully")
			}
			cancel()
		}()
		go func() {  // Monitors cancel signal and calls .Shutdown for webserver
			<-ctx.Done()
			if err := mainServer.Shutdown(context.TODO()); err != nil {
				sLogger.Errorw("Error during server shutdown", zap.Error(err))
			}
		}()


		// Wait signal or error in some goroutine

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		select {
		case sig := <-sigs:  // Signal to stop
			sLogger.Info("Getting signal. Stop", "signal", sig.String())
			cancel()
		case <-ctx.Done():  // Error in goroutine
			sLogger.Info("Application trying to stop itself because of error")
		}

		sLogger.Info("Try to stop gracefully")
		go func() {
			t := time.NewTimer(time.Second * 5)
			<-t.C
			sLogger.Warn("Timeout to stop gracefully. Exit process")
			os.Exit(1)
		}()
		wg.Wait()
		sLogger.Info("Application was gracefully stopped")

	},
}

func init() {
	config.InitBasicParams(mainCmd)
}

func main() {
	if err := mainCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
