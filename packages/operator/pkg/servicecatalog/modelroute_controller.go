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

package servicecatalog

import (
	"context"
	"fmt"
	odahuflowv1alpha1 "github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	"github.com/odahu/odahu-flow/packages/operator/pkg/servicecatalog/catalog"
	"io/ioutil"
	"net/http"
	"net/url"
	ctrl "sigs.k8s.io/controller-runtime"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var (
	log                 = logf.Log.WithName("service-catalog-controller")
	defaultRequeueDelay = 1 * time.Second
)

const (
	ServiceGatalogMaxConcurrentReconciles = 10
	defaultUpdatePeriod                   = 20 * time.Second
	defaultModelRequestTimeout            = 10 * time.Second
)

func NewModelRouteReconciler(
	mgr manager.Manager,
	mrc *catalog.ModelRouteCatalog,
	deploymentConfig config.ModelDeploymentConfig,
) *ModelRouteReconciler {
	rmr := ModelRouteReconciler{
		Client:           mgr.GetClient(),
		scheme:           mgr.GetScheme(),
		mrc:              mrc,
		ticker:           time.NewTicker(defaultUpdatePeriod),
		deploymentConfig: deploymentConfig,
	}

	go func() {
		rmr.StartUpdate()
	}()

	return &rmr
}

type ModelRouteReconciler struct {
	client.Client
	scheme           *runtime.Scheme
	mrc              *catalog.ModelRouteCatalog
	ticker           *time.Ticker
	deploymentConfig config.ModelDeploymentConfig
}

func (r *ModelRouteReconciler) SetupBuilder(mgr ctrl.Manager) *ctrl.Builder {
	return ctrl.NewControllerManagedBy(mgr).For(&odahuflowv1alpha1.ModelRoute{})
}

func (r *ModelRouteReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return r.SetupBuilder(mgr).Complete(r)
}

func (r *ModelRouteReconciler) StartUpdate() {
	k8sRouteList := &odahuflowv1alpha1.ModelRouteList{}

	for range r.ticker.C {
		err := r.List(
			context.TODO(),
			k8sRouteList,
			&client.ListOptions{
				Namespace: r.deploymentConfig.Namespace,
			},
		)

		if err != nil {
			log.Error(err, "Can not get list of model routes")
		}

		log.Info("Found alive model routes", "model routes", k8sRouteList)

		r.mrc.UpdateModelRouteCatalog(k8sRouteList)
	}
}

func (r *ModelRouteReconciler) generateModelRequest(mr *odahuflowv1alpha1.ModelRoute) (*http.Request, error) {
	modelURL := &url.URL{
		Scheme: "http",
		Host: fmt.Sprintf(
			"%s.%s",
			r.deploymentConfig.Istio.ServiceName,
			r.deploymentConfig.Istio.Namespace,
		),
		Path: mr.Spec.URLPrefix,
	}

	edgeHostURL := r.deploymentConfig.Edge.Host
	parsedExternalEdgeURL, err := url.Parse(edgeHostURL)
	if err != nil {
		log.Error(err, "Can not parse the edge host url", "edge host", edgeHostURL)

		return nil, err
	}

	return &http.Request{
		Method: http.MethodGet,
		URL:    modelURL,
		Host:   parsedExternalEdgeURL.Host,
	}, nil
}

func (r *ModelRouteReconciler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	modelRouteCR := &odahuflowv1alpha1.ModelRoute{}
	err := r.Get(context.TODO(), request.NamespacedName, modelRouteCR)
	if err != nil {
		if errors.IsNotFound(err) {
			r.mrc.DeleteModelRoute(request.NamespacedName.Name)

			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}

	if modelRouteCR.Status.State != odahuflowv1alpha1.ModelRouteStateReady {
		log.Info("Model is not ready", "mr id", modelRouteCR.Name)
		return reconcile.Result{RequeueAfter: defaultRequeueDelay}, nil
	}

	modelRequest, err := r.generateModelRequest(modelRouteCR)
	if err != nil {
		log.Error(err, "Can not generate model request", "model route id", modelRouteCR.Name)

		return reconcile.Result{RequeueAfter: defaultRequeueDelay}, nil
	}

	httpClient := http.Client{
		Timeout: defaultModelRequestTimeout,
	}
	response, err := httpClient.Do(modelRequest)
	if err != nil {
		log.Error(
			err, "Can not get swagger response for model",
			"mr id", modelRouteCR.Name,
		)
		return reconcile.Result{RequeueAfter: defaultRequeueDelay}, nil
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)

	log.Info("Get response from model", "model route id", modelRouteCR.Name, "content", string(contents))

	if err != nil {
		log.Error(err, "Can not get swagger response for model", "mr id", modelRouteCR.Name)
		return reconcile.Result{RequeueAfter: defaultRequeueDelay}, nil
	}

	err = r.mrc.AddModelRoute(modelRouteCR, contents)
	if err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
