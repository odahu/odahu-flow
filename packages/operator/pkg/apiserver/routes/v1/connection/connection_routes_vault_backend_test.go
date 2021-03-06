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

package connection_test

import (
	conn_vault_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/connection/vault"
	odahuflow_vault_utils "github.com/odahu/odahu-flow/packages/operator/pkg/repository/util/vault"
	conn_service "github.com/odahu/odahu-flow/packages/operator/pkg/service/connection"
	"github.com/stretchr/testify/suite"
	"net"
	"testing"
)

const (
	testSecretMountPath   = "test-path"
	testSecretMountEngine = "kv"
	testDecryptToken      = "test-decrypt-token" //nolint
)

type ConnectionRouteVaultBackendSuite struct {
	ConnectionRouteGenericSuite
	vaultServer net.Listener
}

func (s *ConnectionRouteVaultBackendSuite) SetupSuite() {
	vaultClient, vaultServer := odahuflow_vault_utils.CreateTestVault(
		s.T(),
		testSecretMountPath,
		testSecretMountEngine,
	)

	connRepo := conn_vault_repository.NewRepository(vaultClient, testSecretMountPath)
	s.connService = conn_service.NewService(connRepo)
	s.vaultServer = vaultServer
	s.connDecryptToken = testDecryptToken

	s.ConnectionRouteGenericSuite.SetupSuite()
}

func (s *ConnectionRouteVaultBackendSuite) TearDownSuite() {
	err := s.vaultServer.Close()
	if err != nil {
		s.T().Fatal("Cannot shutdown test vault server")
	}
}

func TestConnectionRouteVaultBackend(t *testing.T) {
	suite.Run(t, new(ConnectionRouteVaultBackendSuite))
}
