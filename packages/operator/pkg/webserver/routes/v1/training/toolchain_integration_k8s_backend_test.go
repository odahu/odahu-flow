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

package training_test

import (
	conn_k8s_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/connection/kubernetes"
	mt_k8s_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/training/kubernetes"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils"
	"github.com/stretchr/testify/suite"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"testing"
)

type TIK8SRouteSuite struct {
	TIGenericRouteSuite
}

func (s *TIK8SRouteSuite) SetupSuite() {
	mgr, err := manager.New(cfg, manager.Options{NewClient: utils.NewClient})
	if err != nil {
		panic(err)
	}

	s.mtRepository = mt_k8s_repository.NewRepository(testNamespace, testNamespace, mgr.GetClient(), nil)
	s.connRepository = conn_k8s_repository.NewRepository(testNamespace, mgr.GetClient())
}

func TestToolchainIntegrationK8SRouteSuite(t *testing.T) {
	suite.Run(t, new(TIK8SRouteSuite))
}
