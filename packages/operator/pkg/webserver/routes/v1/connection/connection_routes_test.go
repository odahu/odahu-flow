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

package connection_test

import (
	"bytes"
	"encoding/json"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/connection"
	odahuflow_errors "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	conn_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/connection"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils"
	"github.com/odahu/odahu-flow/packages/operator/pkg/webserver/routes"
	conn_route "github.com/odahu/odahu-flow/packages/operator/pkg/webserver/routes/v1/connection"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"
)

const (
	connID        = "testconn"
	connID1       = "testconn1"
	connID2       = "testconn2"
	connReference = "refs/heads/master"
	connURI       = "git@github.com:odahu/odahu-flow.git"
	creds         = "bG9sCg=="
)

func stubKeyEvaluator(uri string) (string, error) {
	return "stub-key", nil
}

type ConnectionRouteGenericSuite struct {
	suite.Suite
	g                *GomegaWithT
	server           *gin.Engine
	routeGroup       *gin.RouterGroup
	connRepository   conn_repository.Repository
	connDecryptToken string
}

func (s *ConnectionRouteGenericSuite) SetupSuite() {
	utils.SetupLogger()
}

func (s *ConnectionRouteGenericSuite) TearDownTest() {
	for _, connID := range []string{connID, connID1, connID2} {
		if err := s.connRepository.DeleteConnection(connID); err != nil && !odahuflow_errors.IsNotFoundError(err) {
			// If a connection is not found then it was not created during a test case
			s.T().Fatalf("Unexpected error: %v", err)
		}
	}
}

func (s *ConnectionRouteGenericSuite) SetupTest() {
	s.g = NewGomegaWithT(s.T())

	s.registerHTTPHandlers(config.ConnectionConfig{
		Enabled: true,
	})
}

func (s *ConnectionRouteGenericSuite) registerHTTPHandlers(connectionConfig config.ConnectionConfig) {
	s.server = gin.Default()
	s.routeGroup = s.server.Group("")
	conn_route.ConfigureRoutes(s.routeGroup, s.connRepository, stubKeyEvaluator, connectionConfig)
}

func (s *ConnectionRouteGenericSuite) newMultipleConnStubs() []*connection.Connection {
	conn1 := newConnStub()
	conn1.ID = connID1
	_, err := s.connRepository.CreateConnection(conn1)
	s.g.Expect(err).NotTo(HaveOccurred())

	conn2 := newConnStub()
	conn2.ID = connID2
	_, err = s.connRepository.CreateConnection(conn2)
	s.g.Expect(err).NotTo(HaveOccurred())

	return []*connection.Connection{conn1, conn2}
}

func newConnStub() *connection.Connection {
	return &connection.Connection{
		ID: connID,
		Spec: v1alpha1.ConnectionSpec{
			Type:     connection.DockerType,
			URI:      connURI,
			Username: "username",
			Password: "password",
		},
	}
}

func (s *ConnectionRouteGenericSuite) TestGetConnection() {
	conn := newConnStub()
	_, err := s.connRepository.CreateConnection(conn)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodGet,
		strings.Replace(conn_route.GetConnectionURL, ":id", connID, -1),
		nil)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result connection.Connection
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(result.ID).Should(Equal(conn.ID))
	s.g.Expect(result.Spec).Should(Equal(conn.DeleteSensitiveData().Spec))
}

func (s *ConnectionRouteGenericSuite) TestGetConnectionNotFound() {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodGet,
		strings.Replace(conn_route.GetConnectionURL, ":id", "not-found", -1),
		nil,
	)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result routes.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusNotFound))
	s.g.Expect(result.Message).Should(ContainSubstring("not found"))
}

func (s *ConnectionRouteGenericSuite) TestGetAllConnections() {
	conn := newConnStub()
	_, err := s.connRepository.CreateConnection(conn)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodGet,
		conn_route.GetAllConnectionURL,
		nil,
	)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result []connection.Connection
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(result).Should(HaveLen(1))
	s.g.Expect(result[0].ID).Should(Equal(conn.ID))
	s.g.Expect(result[0].Spec).Should(Equal(conn.DeleteSensitiveData().Spec))
}

func (s *ConnectionRouteGenericSuite) TestGetAllEmptyConnections() {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodGet,
		conn_route.GetAllConnectionURL,
		nil,
	)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result []connection.Connection
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(result).Should(HaveLen(0))
}

func (s *ConnectionRouteGenericSuite) TestGetAllConnectionsByType() {
	connGit := &connection.Connection{
		ID: connID1,
		Spec: v1alpha1.ConnectionSpec{
			Type:      connection.GITType,
			URI:       connURI,
			Reference: connReference,
			KeySecret: creds,
		},
	}
	_, err := s.connRepository.CreateConnection(connGit)
	s.g.Expect(err).NotTo(HaveOccurred())

	connDocker := &connection.Connection{
		ID: connID2,
		Spec: v1alpha1.ConnectionSpec{
			Type:      connection.DockerType,
			URI:       connURI,
			Reference: connReference,
			KeySecret: creds,
		},
	}
	_, err = s.connRepository.CreateConnection(connDocker)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, conn_route.GetAllConnectionURL, nil)
	s.g.Expect(err).NotTo(HaveOccurred())

	query := req.URL.Query()
	query.Set("type", string(connection.GITType))
	req.URL.RawQuery = query.Encode()

	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result []connection.Connection
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(result).Should(HaveLen(1))
	s.g.Expect(result[0].ID).Should(Equal(connGit.ID))
	s.g.Expect(result[0].Spec).Should(Equal(connGit.DeleteSensitiveData().Spec))
}

func (s *ConnectionRouteGenericSuite) TestGetAllConnectionsMultipleFiltersByType() {
	connGit := &connection.Connection{
		ID: connID1,
		Spec: v1alpha1.ConnectionSpec{
			Type:      connection.GITType,
			URI:       connURI,
			Reference: connReference,
			KeySecret: creds,
		},
	}
	_, err := s.connRepository.CreateConnection(connGit)
	s.g.Expect(err).NotTo(HaveOccurred())

	connDocker := &connection.Connection{
		ID: connID2,
		Spec: v1alpha1.ConnectionSpec{
			Type:      connection.DockerType,
			URI:       connURI,
			Reference: connReference,
			KeySecret: creds,
		},
	}
	_, err = s.connRepository.CreateConnection(connDocker)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, conn_route.GetAllConnectionURL, nil)
	s.g.Expect(err).NotTo(HaveOccurred())

	query := req.URL.Query()
	query.Set("type", string(connection.GITType))
	query.Add("type", string(connection.DockerType))
	req.URL.RawQuery = query.Encode()

	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result []connection.Connection
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(result).Should(HaveLen(2))
}

func (s *ConnectionRouteGenericSuite) TestGetAllConnectionsPaging() {
	s.newMultipleConnStubs()

	connNames := map[string]interface{}{connID1: nil, connID2: nil}

	// Return first page
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, conn_route.GetAllConnectionURL, nil)
	s.g.Expect(err).NotTo(HaveOccurred())

	query := req.URL.Query()
	query.Set("size", "1")
	query.Set("page", "0")
	req.URL.RawQuery = query.Encode()

	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result []connection.Connection
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(result).Should(HaveLen(1))
	delete(connNames, result[0].ID)

	// Return second page
	w = httptest.NewRecorder()
	req, err = http.NewRequest(http.MethodGet, conn_route.GetAllConnectionURL, nil)
	s.g.Expect(err).NotTo(HaveOccurred())

	query = req.URL.Query()
	query.Set("size", "1")
	query.Set("page", "1")
	req.URL.RawQuery = query.Encode()

	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(result).Should(HaveLen(1))
	delete(connNames, result[0].ID)

	// Return third empty page
	w = httptest.NewRecorder()
	req, err = http.NewRequest(http.MethodGet, conn_route.GetAllConnectionURL, nil)
	s.g.Expect(err).NotTo(HaveOccurred())

	query = req.URL.Query()
	query.Set("size", "1")
	query.Set("page", "2")
	req.URL.RawQuery = query.Encode()

	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(result).Should(HaveLen(0))
	s.g.Expect(result).Should(BeEmpty())
}

func (s *ConnectionRouteGenericSuite) TestCreateConnection() {
	connEntity := newConnStub()

	connEntityBody, err := json.Marshal(connEntity)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, conn_route.CreateConnectionURL, bytes.NewReader(connEntityBody))
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var connResponse connection.Connection
	err = json.Unmarshal(w.Body.Bytes(), &connResponse)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(http.StatusCreated).Should(Equal(w.Code))
	s.g.Expect(connResponse.ID).Should(Equal(connEntity.ID))
	s.g.Expect(connResponse.Spec).To(Equal(connEntity.DeleteSensitiveData().Spec))

	conn, err := s.connRepository.GetConnection(connID)
	s.g.Expect(err).ShouldNot(HaveOccurred())
	s.g.Expect(conn.ID).To(Equal(connEntity.ID))
	s.g.Expect(conn.Spec).To(Equal(connEntity.DeleteSensitiveData().Spec))
}

// CreatedAt and UpdatedAt field should automatically be updated after create request
func (s *ConnectionRouteGenericSuite) TestCreateConnectionModifiable() {
	newResource := newConnStub()

	newResourceBody, err := json.Marshal(newResource)
	s.g.Expect(err).NotTo(HaveOccurred())

	reqTime := routes.GetTimeNowTruncatedToSeconds()
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, conn_route.CreateConnectionURL, bytes.NewReader(newResourceBody))
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var resp connection.Connection
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusCreated))
	s.g.Expect(resp.Status.CreatedAt).NotTo(BeNil())
	createdAtWasNotUpdated := reqTime.Before(resp.Status.CreatedAt) || reqTime.Equal(resp.Status.CreatedAt)
	s.g.Expect(createdAtWasNotUpdated).Should(Equal(true))
	s.g.Expect(resp.Status.UpdatedAt).NotTo(BeNil())
	updatedAtWasUpdated := reqTime.Before(resp.Status.CreatedAt) || reqTime.Equal(resp.Status.CreatedAt)
	s.g.Expect(updatedAtWasUpdated).Should(Equal(true))
}

func (s *ConnectionRouteGenericSuite) TestCreateDuplicateConnection() {
	conn := newConnStub()

	_, err := s.connRepository.CreateConnection(conn)
	s.g.Expect(err).NotTo(HaveOccurred())

	connEntityBody, err := json.Marshal(conn)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, conn_route.CreateConnectionURL, bytes.NewReader(connEntityBody))
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result routes.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusConflict))
	s.g.Expect(result.Message).Should(ContainSubstring("already exists"))
}

func (s *ConnectionRouteGenericSuite) TestValidateCreateConnection() {
	conn := newConnStub()
	conn.Spec.Type = "not-found-type"

	connEntityBody, err := json.Marshal(conn)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, conn_route.CreateConnectionURL, bytes.NewReader(connEntityBody))
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result routes.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusBadRequest))
	s.g.Expect(result.Message).Should(ContainSubstring("unknown type: not-found-type"))
}

func (s *ConnectionRouteGenericSuite) TestUpdateConnection() {
	conn := newConnStub()
	_, err := s.connRepository.CreateConnection(conn)
	s.g.Expect(err).NotTo(HaveOccurred())

	connEntity := newConnStub()
	connEntity.Spec.URI = "new-uri"

	connEntityBody, err := json.Marshal(connEntity)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPut, conn_route.UpdateConnectionURL, bytes.NewReader(connEntityBody))
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var connResponse connection.Connection
	err = json.Unmarshal(w.Body.Bytes(), &connResponse)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(connResponse.ID).Should(Equal(connEntity.ID))
	s.g.Expect(connResponse.Spec).Should(Equal(connEntity.DeleteSensitiveData().Spec))

	conn, err = s.connRepository.GetConnection(connID)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.g.Expect(conn.Spec).To(Equal(connEntity.DeleteSensitiveData().Spec))
}

func (s *ConnectionRouteGenericSuite) TestUpdateConnectionNotFound() {
	connEntity := newConnStub()

	connEntityBody, err := json.Marshal(connEntity)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPut, conn_route.UpdateConnectionURL, bytes.NewReader(connEntityBody))
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result routes.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusNotFound))
	s.g.Expect(result.Message).Should(ContainSubstring("not found"))
}

func (s *ConnectionRouteGenericSuite) TestValidateUpdateConnection() {
	conn := newConnStub()
	conn.Spec.Type = "not-found-type"
	_, err := s.connRepository.CreateConnection(conn)
	s.g.Expect(err).NotTo(HaveOccurred())

	connEntityBody, err := json.Marshal(conn)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, conn_route.CreateConnectionURL, bytes.NewReader(connEntityBody))
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result routes.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusBadRequest))
	s.g.Expect(result.Message).Should(ContainSubstring("unknown type: not-found-type"))
}

func (s *ConnectionRouteGenericSuite) TestDeleteConnection() {
	conn := newConnStub()
	_, err := s.connRepository.CreateConnection(conn)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodDelete,
		strings.Replace(conn_route.DeleteConnectionURL, ":id", connID, -1),
		nil,
	)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result routes.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(result.Message).Should(ContainSubstring("was deleted"))

	connList, err := s.connRepository.GetConnectionList()
	s.g.Expect(err).NotTo(HaveOccurred())
	s.g.Expect(connList).To(HaveLen(0))
}

func (s *ConnectionRouteGenericSuite) TestDeleteConnectionNotFound() {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodDelete,
		strings.Replace(conn_route.DeleteConnectionURL, ":id", "not-found", -1),
		nil,
	)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result routes.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusNotFound))
	s.g.Expect(result.Message).Should(ContainSubstring("not found"))
}

func (s *ConnectionRouteGenericSuite) TestGetDecryptedConnection() {
	conn := newConnStub()
	_, err := s.connRepository.CreateConnection(conn)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodGet,
		strings.Replace(conn_route.GetDecryptedConnectionURL, ":id", connID, -1),
		nil,
	)
	s.g.Expect(err).NotTo(HaveOccurred())

	query := req.URL.Query()
	query.Set(conn_route.ConnDecryptTokenQueryParam, s.connDecryptToken)
	req.URL.RawQuery = query.Encode()

	s.server.ServeHTTP(w, req)

	var result connection.Connection
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(result.ID).Should(Equal(conn.ID))
	s.g.Expect(result.Spec).Should(Equal(conn.Spec))
}

func (s *ConnectionRouteGenericSuite) TestGetDecryptedConnectionNotFound() {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodGet,
		strings.Replace(conn_route.GetDecryptedConnectionURL, ":id", "not-found", -1),
		nil,
	)
	s.g.Expect(err).NotTo(HaveOccurred())

	query := req.URL.Query()
	query.Set(conn_route.ConnDecryptTokenQueryParam, s.connDecryptToken)
	req.URL.RawQuery = query.Encode()

	s.server.ServeHTTP(w, req)

	var result routes.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusNotFound))
	s.g.Expect(result.Message).Should(ContainSubstring("not found"))
}

func (s *ConnectionRouteGenericSuite) TestDisabledAPIGetConnection() {
	s.registerHTTPHandlers(config.ConnectionConfig{
		Enabled: false,
	})

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodGet,
		strings.Replace(conn_route.GetConnectionURL, ":id", "12345", -1),
		nil,
	)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result routes.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusNotFound))
	s.g.Expect(result.Message).Should(ContainSubstring("not found"))
}

func (s *ConnectionRouteGenericSuite) TestDisabledAPIGetAllConnections() {
	s.registerHTTPHandlers(config.ConnectionConfig{
		Enabled: false,
	})

	conn := newConnStub()
	_, err := s.connRepository.CreateConnection(conn)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodGet,
		conn_route.GetAllConnectionURL,
		nil,
	)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result []connection.Connection
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(result).Should(HaveLen(1))
	s.g.Expect(result[0].ID).Should(Equal(conn.ID))
	s.g.Expect(result[0].Spec).Should(Equal(conn.DeleteSensitiveData().Spec))
}

func (s *ConnectionRouteGenericSuite) TestDisabledAPIGetDecryptedConnection() {
	s.registerHTTPHandlers(config.ConnectionConfig{
		Enabled: false,
	})

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodGet,
		strings.Replace(conn_route.GetDecryptedConnectionURL, ":id", "12345", -1),
		nil,
	)
	s.g.Expect(err).NotTo(HaveOccurred())

	query := req.URL.Query()
	query.Set(conn_route.ConnDecryptTokenQueryParam, s.connDecryptToken)
	req.URL.RawQuery = query.Encode()

	s.server.ServeHTTP(w, req)

	var result routes.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusNotFound))
	s.g.Expect(result.Message).Should(ContainSubstring("not found"))
}

func (s *ConnectionRouteGenericSuite) TestDisabledAPICreateConnection() {
	s.registerHTTPHandlers(config.ConnectionConfig{
		Enabled: false,
	})

	connEntity := newConnStub()

	connEntityBody, err := json.Marshal(connEntity)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, conn_route.CreateConnectionURL, bytes.NewReader(connEntityBody))
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result routes.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusBadRequest))
	s.g.Expect(result.Message).Should(ContainSubstring(routes.DisabledAPIErrorMessage))

	_, err = s.connRepository.GetConnection(connID)
	s.g.Expect(err).Should(Equal(odahuflow_errors.NotFoundError{Entity: connID}))
}

func (s *ConnectionRouteGenericSuite) TestDisabledAPIUpdateConnection() {
	s.registerHTTPHandlers(config.ConnectionConfig{
		Enabled: false,
	})

	connEntity := newConnStub()

	connEntityBody, err := json.Marshal(connEntity)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPut, conn_route.UpdateConnectionURL, bytes.NewReader(connEntityBody))
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result routes.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusBadRequest))
	s.g.Expect(result.Message).Should(ContainSubstring(routes.DisabledAPIErrorMessage))

	_, err = s.connRepository.GetConnection(connID)
	s.g.Expect(err).Should(Equal(odahuflow_errors.NotFoundError{Entity: connID}))
}

func (s *ConnectionRouteGenericSuite) TestDisabledAPIDeleteConnection() {
	s.registerHTTPHandlers(config.ConnectionConfig{
		Enabled: false,
	})
	w := httptest.NewRecorder()

	req, err := http.NewRequest(
		http.MethodDelete,
		strings.Replace(conn_route.DeleteConnectionURL, ":id", "12345", -1),
		nil,
	)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result routes.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusBadRequest))
	s.g.Expect(result.Message).Should(ContainSubstring(routes.DisabledAPIErrorMessage))
}
