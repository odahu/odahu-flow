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

package utils

import (
	"fmt"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	"k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp" //nolint
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	k8s_config "sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var logM = logf.Log.WithName("k8s-manager")

type ManagerCloser interface {
	manager.Manager
	// Closing of a manager. The manager will be unusable after the execution of this function.
	Close() error
}

// Cleanup the kubernetes environment if it is provided.
type managerWrapper struct {
	k8sEnvironment *envtest.Environment
	manager.Manager
}

func (m *managerWrapper) Close() error {
	if m.k8sEnvironment == nil {
		return nil
	}

	if err := m.k8sEnvironment.Stop(); err != nil {
		logM.Error(err, "Error during closing of local kubernetes environment")

		return err
	}

	return nil
}

func NewClient(cache cache.Cache, config *rest.Config, options client.Options) (client.Client, error) {
	c, err := client.New(config, options)
	if err != nil {
		return nil, err
	}

	// TODO: enable caching for k8s entities
	return &client.DelegatingClient{
		Reader:       c,
		Writer:       c,
		StatusClient: c,
	}, nil
}

func newLocalManager(localConfig config.APILocalBackendConfig) (ManagerCloser, error) {
	var cfg *rest.Config

	k8sEnvironment := &envtest.Environment{
		CRDDirectoryPaths: []string{localConfig.LocalBackendCRDPath},
	}

	err := v1alpha1.AddToScheme(scheme.Scheme)
	if err != nil {
		logM.Error(err, "Cannot setup the odahuflow schema")
		return nil, err
	}

	cfg, err = k8sEnvironment.Start()
	if err != nil {
		logM.Error(err, "Cannot setup the test k8s api")
		return nil, err
	}

	mgr, err := manager.New(cfg, manager.Options{NewClient: NewClient})
	if err != nil {
		logM.Error(err, "Cannot setup the test k8s manager")
		return nil, err
	}

	return &managerWrapper{k8sEnvironment: k8sEnvironment, Manager: mgr}, nil
}

func NewManager(backendConfig config.APIBackendConfig) (ManagerCloser, error) {
	switch backendConfig.Type {
	case config.ConfigBackendType:
		return newConfigManager()
	case config.LocalBackendType:
		return newLocalManager(backendConfig.Local)
	default:
		return nil, fmt.Errorf("unexpected backend type: %s", backendConfig.Type)
	}
}

func newConfigManager() (ManagerCloser, error) {
	cfg, err := k8s_config.GetConfig()
	if err != nil {
		logM.Error(err, "K8s config creation")
		return nil, err
	}

	mgr, err := manager.New(cfg, manager.Options{NewClient: NewClient})
	if err != nil {
		logM.Error(err, "Manager creation")
		return nil, err
	}

	if err := v1alpha1.AddToScheme(mgr.GetScheme()); err != nil {
		logM.Error(err, "Update schema")
		return nil, err
	}

	return &managerWrapper{Manager: mgr}, nil
}
