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

package controllers_test

import (
	stdlog "log"
	"os"
	"path/filepath"
	ctrl "sigs.k8s.io/controller-runtime"
	"sync"
	"testing"

	odahuv1alpha1 "github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	istioschema "github.com/aspenmesh/istio-client-go/pkg/client/clientset/versioned/scheme"
	tektonschema "github.com/tektoncd/pipeline/pkg/client/clientset/versioned/scheme"
	knservingv1 "knative.dev/serving/pkg/apis/serving/v1"
)

var cfg *rest.Config

func TestMain(m *testing.M) {
	t := &envtest.Environment{
		CRDDirectoryPaths: []string{
			filepath.Join("..", "config", "crds"),
			filepath.Join("..", "hack", "tests", "thirdparty_crds"),
		},
	}

	err := odahuv1alpha1.AddToScheme(scheme.Scheme)
	if err != nil {
		panic(err)
	}
	istioschema.AddToScheme(scheme.Scheme)

	if err := knservingv1.AddToScheme(scheme.Scheme); err != nil {
		panic(err)
	}
	if err := tektonschema.AddToScheme(scheme.Scheme); err != nil {
		panic(err)
	}

	if cfg, err = t.Start(); err != nil {
		stdlog.Fatal(err)
	}

	code := m.Run()

	if err = t.Stop(); err != nil {
		panic(err)
	}

	os.Exit(code)
}

// SetupTestReconcile returns a reconcile.Reconcile implementation that delegates to inner and
// writes the request to requests after Reconcile is finished.
func SetupTestReconcile(inner reconcile.Reconciler) (reconcile.Reconciler, chan reconcile.Request) {
	requests := make(chan reconcile.Request)
	fn := reconcile.Func(func(req reconcile.Request) (reconcile.Result, error) {
		result, err := inner.Reconcile(req)
		requests <- req
		return result, err
	})
	return fn, requests
}

// StartTestManager adds recFn
func StartTestManager(mgr manager.Manager, g *GomegaWithT) (chan struct{}, *sync.WaitGroup) {
	stop := make(chan struct{})
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		g.Expect(mgr.Start(stop)).NotTo(HaveOccurred())
	}()
	return stop, wg
}

type ReconcilerWithSetup interface {
	reconcile.Reconciler
	SetupBuilder(mgr ctrl.Manager) *ctrl.Builder
}

type ReconcilerWrapper struct {
	embeddedReconciler ReconcilerWithSetup
	requests           chan reconcile.Request
}

func NewReconcilerWrapper(embedded ReconcilerWithSetup, requests chan reconcile.Request) *ReconcilerWrapper {
	return &ReconcilerWrapper{
		embeddedReconciler: embedded,
		requests:           requests,
	}
}

func (rw *ReconcilerWrapper) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	result, err := rw.embeddedReconciler.Reconcile(req)
	rw.requests <- req
	return result, err
}

func (rw *ReconcilerWrapper) SetupWithManager(mgr ctrl.Manager) error {
	return rw.embeddedReconciler.SetupBuilder(mgr).Complete(rw)
}
