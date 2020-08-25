/*
 * Copyright 2020 EPAM Systems
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package connection_test

import (
	"encoding/base64"
	"errors"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/connection"
	odahu_errors "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	conn_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/connection"
	conn_service "github.com/odahu/odahu-flow/packages/operator/pkg/service/connection"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

const connID = "some-id"
const notBase64 = "not base64"

func TestConnectionserviceSuite(t *testing.T) {
	suite.Run(t, new(ConnectionServiceTestSuite))
}

type ConnectionServiceTestSuite struct {
	suite.Suite
	mockRepo          conn_service.RepositoryMock
	connectionService conn_service.Service
}

func (s *ConnectionServiceTestSuite) SetupTest() {
	s.mockRepo = conn_service.RepositoryMock{}
	s.connectionService = conn_service.NewService(&s.mockRepo)
}

// Tests that Service retrieves a connection from repository by ID and masks secrets when
// "encrypted" flag set to true. At the same time public key must be base64-encoded.
func (s *ConnectionServiceTestSuite) TestGetConnection_Encrypted() {
	connectionFromRepo := stubConnection()
	expectedPublicKey := base64.StdEncoding.EncodeToString([]byte(connectionFromRepo.Spec.PublicKey))

	s.mockRepo.On("GetConnection", connID).Return(&connectionFromRepo, nil)

	conn, err := s.connectionService.GetConnection(connID, true)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), connID, conn.ID)

	assert.Equal(s.T(), connection.DecryptedDataMask, conn.Spec.Password)
	assert.Equal(s.T(), connection.DecryptedDataMask, conn.Spec.KeySecret)
	assert.Equal(s.T(), connection.DecryptedDataMask, conn.Spec.KeyID)

	assert.Equal(s.T(), expectedPublicKey, conn.Spec.PublicKey)
}

// Tests getting decrypted connection from Service
// Expect: secrets and public key are base64-encoded
func (s *ConnectionServiceTestSuite) TestGetConnection_Decrypted() {
	connectionFromRepo := stubConnection()
	expectedPassword := base64.StdEncoding.EncodeToString([]byte(connectionFromRepo.Spec.Password))
	expectedKeySecret := base64.StdEncoding.EncodeToString([]byte(connectionFromRepo.Spec.KeySecret))
	expectedKeyID := base64.StdEncoding.EncodeToString([]byte(connectionFromRepo.Spec.KeyID))
	expectedPublicKey := base64.StdEncoding.EncodeToString([]byte(connectionFromRepo.Spec.PublicKey))

	s.mockRepo.On("GetConnection", connID).Return(&connectionFromRepo, nil)

	conn, err := s.connectionService.GetConnection(connID, false)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), connID, conn.ID)

	assert.Equal(s.T(), expectedPassword, conn.Spec.Password)
	assert.Equal(s.T(), expectedKeySecret, conn.Spec.KeySecret)
	assert.Equal(s.T(), expectedKeyID, conn.Spec.KeyID)
	assert.Equal(s.T(), expectedPublicKey, conn.Spec.PublicKey)
}

func (s *ConnectionServiceTestSuite) TestGetConnection_ErrorFromRepo() {
	errorFromRepo := errors.New("my error")
	s.mockRepo.On("GetConnection", connID).Return(nil, errorFromRepo)

	conn, err := s.connectionService.GetConnection(connID, false)
	assert.Nil(s.T(), conn)
	assert.Equal(s.T(), errorFromRepo, err)
}

func (s *ConnectionServiceTestSuite) TestGetConnectionList() {
	options := []conn_repository.ListOption{
		conn_repository.Page(0),
		conn_repository.Size(100),
	}
	connectionsFromRepo := []connection.Connection{stubConnection(), stubConnection()}
	expectedConnections := append([]connection.Connection{}, connectionsFromRepo...)

	s.mockRepo.On("GetConnectionList", options).Return(connectionsFromRepo, nil)

	actualConnections, err := s.connectionService.GetConnectionList(options...)
	assert.Nil(s.T(), err)

	for i, conn := range actualConnections {
		assert.Equal(s.T(), expectedConnections[i].ID, conn.ID)
		assert.Equal(s.T(), connection.DecryptedDataMask, conn.Spec.Password)
		assert.Equal(s.T(), connection.DecryptedDataMask, conn.Spec.KeySecret)
		assert.Equal(s.T(), connection.DecryptedDataMask, conn.Spec.KeyID)
		expectedPublicKey := base64.StdEncoding.EncodeToString([]byte(expectedConnections[i].Spec.PublicKey))
		assert.Equal(s.T(), expectedPublicKey, conn.Spec.PublicKey)
	}
}

func (s *ConnectionServiceTestSuite) TestGetConnectionList_ErrorFromRepo() {
	options := []conn_repository.ListOption{
		conn_repository.Page(0),
		conn_repository.Size(100),
	}
	s.mockRepo.On("GetConnectionList", options).Return([]connection.Connection{}, errors.New("some"))

	actualConnections, err := s.connectionService.GetConnectionList(options...)
	assert.Empty(s.T(), actualConnections)
	assert.Error(s.T(), err)
}

// Tests that Service just proxies the call to repo
func (s *ConnectionServiceTestSuite) TestDeleteConnection() {
	errorFromRepo := errors.New("some error")
	s.mockRepo.On("DeleteConnection", connID).Return(errorFromRepo)
	err := s.connectionService.DeleteConnection(connID)
	assert.Equal(s.T(), errorFromRepo, err)
}

func (s *ConnectionServiceTestSuite) TestUpdateConnection() {
	originalConnection := stubConnection()
	connectionForService := originalConnection

	connectionForService.EncodeBase64Fields()

	s.mockRepo.
		On("UpdateConnection", mock.AnythingOfType("*connection.Connection")).
		Return(nil)

	updatedConnection, err := s.connectionService.UpdateConnection(connectionForService)

	assert.Nil(s.T(), err)

	// Check that connection was base64-decoded before sending to repo
	assert.Equal(s.T(), originalConnection, s.mockRepo.UpdatedConnection)

	// Check that secrets are masked in result
	assert.Equal(s.T(), connection.DecryptedDataMask, updatedConnection.Spec.KeyID)
	assert.Equal(s.T(), connection.DecryptedDataMask, updatedConnection.Spec.KeySecret)
	assert.Equal(s.T(), connection.DecryptedDataMask, updatedConnection.Spec.Password)
	// Check that public key is base64-encoded
	expectedPublicKey := base64.StdEncoding.EncodeToString([]byte(originalConnection.Spec.PublicKey))
	assert.Equal(s.T(), expectedPublicKey, updatedConnection.Spec.PublicKey)
}

func (s *ConnectionServiceTestSuite) TestUpdateConnection_DecodingFail() {
	connectionForService := stubConnection()
	connectionForService.Spec.Password = notBase64
	connectionForService.Spec.KeySecret = notBase64
	connectionForService.Spec.PublicKey = notBase64
	connectionForService.Spec.KeyID = notBase64

	updatedConnection, err := s.connectionService.UpdateConnection(connectionForService)

	// Checks that service returns an InvalidEntityError if base64 decoding fails
	assert.Nil(s.T(), updatedConnection)
	assert.Error(s.T(), err)
	assert.IsType(s.T(), odahu_errors.InvalidEntityError{}, err)
}

func (s *ConnectionServiceTestSuite) TestUpdateConnection_ErrorFromRepo() {
	connectionForService := stubConnection()
	connectionForService.EncodeBase64Fields()

	s.mockRepo.
		On("UpdateConnection", mock.AnythingOfType("*connection.Connection")).
		Return(errors.New("error from repo"))

	updatedConnection, err := s.connectionService.UpdateConnection(connectionForService)

	assert.Nil(s.T(), updatedConnection)
	assert.Error(s.T(), err)
}

func (s *ConnectionServiceTestSuite) TestCreateConnection() {
	originalConnection := stubConnection()
	connectionForService := originalConnection

	connectionForService.EncodeBase64Fields()

	s.mockRepo.
		On("CreateConnection", mock.AnythingOfType("*connection.Connection")).
		Return(nil)

	updatedConnection, err := s.connectionService.CreateConnection(connectionForService)

	assert.Nil(s.T(), err)

	// Check that connection was base64-decoded before sending to repo
	assert.Equal(s.T(), originalConnection, s.mockRepo.CreatedConnection)

	// Check that secrets are masked in result
	assert.Equal(s.T(), connection.DecryptedDataMask, updatedConnection.Spec.KeyID)
	assert.Equal(s.T(), connection.DecryptedDataMask, updatedConnection.Spec.KeySecret)
	assert.Equal(s.T(), connection.DecryptedDataMask, updatedConnection.Spec.Password)
	// Check that public key is base64-encoded
	expectedPublicKey := base64.StdEncoding.EncodeToString([]byte(originalConnection.Spec.PublicKey))
	assert.Equal(s.T(), expectedPublicKey, updatedConnection.Spec.PublicKey)
}

func (s *ConnectionServiceTestSuite) TestCreateConnection_DecodingFail() {
	connectionForService := stubConnection()
	connectionForService.Spec.Password = notBase64
	connectionForService.Spec.KeySecret = notBase64
	connectionForService.Spec.PublicKey = notBase64
	connectionForService.Spec.KeyID = notBase64

	createdConnection, err := s.connectionService.CreateConnection(connectionForService)

	// Checks that service returns an InvalidEntityError if base64 decoding fails
	assert.Nil(s.T(), createdConnection)
	assert.Error(s.T(), err)
	assert.IsType(s.T(), odahu_errors.InvalidEntityError{}, err)
}

func (s *ConnectionServiceTestSuite) TestCreateConnection_ErrorFromRepo() {
	connectionForService := stubConnection()
	connectionForService.EncodeBase64Fields()

	s.mockRepo.
		On("CreateConnection", mock.AnythingOfType("*connection.Connection")).
		Return(errors.New("error from repo"))

	createdConnection, err := s.connectionService.CreateConnection(connectionForService)

	assert.Nil(s.T(), createdConnection)
	assert.Error(s.T(), err)
}

func stubConnection() connection.Connection {
	return connection.Connection{
		ID: connID,
		Spec: v1alpha1.ConnectionSpec{
			Password:  "123456",
			KeySecret: "secret",
			KeyID:     "secretID",
			PublicKey: "public",
		},
	}
}
