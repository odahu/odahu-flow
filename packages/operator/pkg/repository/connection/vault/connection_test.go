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

package vault_test

import (
	"net"
	"testing"

	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/connection"
	odahuflow_errors "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	conn_repo "github.com/odahu/odahu-flow/packages/operator/pkg/repository/connection"
	conn_vault_repo "github.com/odahu/odahu-flow/packages/operator/pkg/repository/connection/vault"
	odahuflow_vault_utils "github.com/odahu/odahu-flow/packages/operator/pkg/repository/util/vault"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"
)

const (
	testSecretMountPath   = "test-path"
	testDecryptToken      = "test-decrypt-token" //nolint
	testSecretMountEngine = "kv"
	notFoundConnID        = "not-found"
)

type VaultConnRepoSuite struct {
	suite.Suite
	g           *GomegaWithT
	connRepo    conn_repo.Repository
	vaultServer net.Listener
}

func (s *VaultConnRepoSuite) SetupSuite() {
	vaultClient, vaultServer := odahuflow_vault_utils.CreateTestVault(
		s.T(),
		testSecretMountPath,
		testSecretMountEngine,
	)

	s.connRepo = conn_vault_repo.NewRepository(vaultClient, testSecretMountPath)
	s.vaultServer = vaultServer
}

func (s *VaultConnRepoSuite) TearDownSuite() {
	err := s.vaultServer.Close()
	if err != nil {
		s.T().Fatal("Cannot shutdown test vault server")
	}
}

func (s *VaultConnRepoSuite) SetupTest() {
	s.g = NewWithT(s.T())
}

func TestEncDecSuite(t *testing.T) {
	suite.Run(t, new(VaultConnRepoSuite))
}

func (s *VaultConnRepoSuite) TestGetEmptyList() {
	connections, err := s.connRepo.GetConnectionList()

	s.g.Expect(err).ShouldNot(HaveOccurred())
	s.g.Expect(connections).Should(BeEmpty())
}

func (s *VaultConnRepoSuite) TestGetNotFound() {
	_, err := s.connRepo.GetConnection(notFoundConnID)

	s.g.Expect(err).Should(And(
		HaveOccurred(),
		MatchError(odahuflow_errors.NotFoundError{Entity: notFoundConnID}),
	))
}

func (s *VaultConnRepoSuite) TestConnectionRepository() {
	const (
		connID   = "test-conn-id"
		connType = connection.GITType
	)

	created := &connection.Connection{
		ID: connID,
		Spec: v1alpha1.ConnectionSpec{
			Type: connType,
		},
	}

	err := s.connRepo.SaveConnection(created)
	s.g.Expect(err).NotTo(HaveOccurred())

	fetched, err := s.connRepo.GetConnection(connID)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.g.Expect(fetched.ID).To(Equal(created.ID))
	s.g.Expect(fetched.Spec).To(Equal(created.DeleteSensitiveData().Spec))

	newConnType := connection.GcsType
	updated := &connection.Connection{
		ID: connID,
		Spec: v1alpha1.ConnectionSpec{
			Type: newConnType,
		},
	}

	err = s.connRepo.UpdateConnection(updated)
	s.g.Expect(err).NotTo(HaveOccurred())

	fetched, err = s.connRepo.GetConnection(connID)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.g.Expect(fetched.ID).To(Equal(updated.ID))
	s.g.Expect(fetched.Spec).To(Equal(updated.DeleteSensitiveData().Spec))
	s.g.Expect(fetched.Spec.Type).To(Equal(newConnType))

	s.g.Expect(s.connRepo.DeleteConnection(connID)).NotTo(HaveOccurred())
	_, err = s.connRepo.GetConnection(connID)
	s.g.Expect(err).To(HaveOccurred())
}
