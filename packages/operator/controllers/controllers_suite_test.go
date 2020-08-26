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
	"github.com/emicklei/go-restful/log"
	"github.com/odahu/odahu-flow/packages/operator/pkg/testhelpers/testenvs"
	"os"
	"path/filepath"
	ctrl "sigs.k8s.io/controller-runtime"
	"sync"
	"testing"

	. "github.com/onsi/gomega"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var (
	cfg        *rest.Config
)

func testMainWrapper(m *testing.M) int {

	var err error

	var closeKube func() error

	_, cfg, closeKube, _, err = testenvs.SetupTestKube(
		filepath.Join("..", "config", "crds"),
		filepath.Join("..", "hack", "tests", "thirdparty_crds"),
	)

	defer func() {
		if err := closeKube(); err != nil {
			log.Print("Error during release test Kube Environment resources")
		}
	}()
	if err != nil {
		return -1
	}

	return m.Run()
}

func TestMain(m *testing.M) {

	os.Exit(testMainWrapper(m))
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
