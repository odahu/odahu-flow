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

package packaging_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/connection"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/packaging"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/routes"
	pack_route "github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/routes/v1/packaging"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	odahuErrors "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	"github.com/odahu/odahu-flow/packages/operator/pkg/service/packaging_integration"
	httputil "github.com/odahu/odahu-flow/packages/operator/pkg/utils/httputil"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	piID           = "ti-test"
	piEntrypoint   = "test-entrypoint"
	piDefaultImage = "test:image"
	piPrivileged   = false
	testNamespace  = "default"
	testMpID1      = "test-model-1"
	testMpID2      = "test-model-2"
	mpImage        = "docker-rest"
)

var (
	mpArtifactName = "mock-artifact-name-id"
	piArguments    = packaging.JsonSchema{
		Properties: []packaging.Property{
			{
				Name: "argument-1",
				Parameters: []packaging.Parameter{
					{
						Name:  "minimum",
						Value: float64(5),
					},
					{
						Name:  "type",
						Value: "number",
					},
				},
			},
		},
		Required: []string{"argument-1"},
	}
	piTargets = []v1alpha1.TargetSchema{
		{
			Name: "target-1",
			ConnectionTypes: []string{
				string(connection.S3Type),
				string(connection.GcsType),
				string(connection.AzureBlobType),
			},
			Required: false,
		},
		{
			Name: "target-2",
			ConnectionTypes: []string{
				string(connection.DockerType),
			},
			Required: true,
		},
	}
)

type PackagingIntegrationRouteSuite struct {
	suite.Suite
	g             *GomegaWithT
	server        *gin.Engine
	piServiceMock *packaging_integration.MockService
}

func TestPIGenericRouteSuite(t *testing.T) {
	suite.Run(t, new(PackagingIntegrationRouteSuite))
}

func (s *PackagingIntegrationRouteSuite) SetupTest() {
	s.g = NewGomegaWithT(s.T())
	s.piServiceMock = &packaging_integration.MockService{}
	s.registerHandlers(config.NewDefaultModelPackagingConfig())
}

func (s *PackagingIntegrationRouteSuite) registerHandlers(packagingConfig config.ModelPackagingConfig) {
	s.server = gin.Default()
	v1Group := s.server.Group("")
	packGroup := v1Group.Group("", routes.DisableAPIMiddleware(packagingConfig.Enabled))
	pack_route.ConfigurePiRoutes(packGroup, s.piServiceMock)
}

func newPackagingIntegration() *packaging.PackagingIntegration {
	return &packaging.PackagingIntegration{
		ID: piID,
		Spec: packaging.PackagingIntegrationSpec{
			Entrypoint:   piEntrypoint,
			DefaultImage: piDefaultImage,
			Privileged:   piPrivileged,
			Schema: packaging.Schema{
				Targets:   piTargets,
				Arguments: piArguments,
			},
		},
	}
}

func (s *PackagingIntegrationRouteSuite) TestGetPackagingIntegration() {
	pi := newPackagingIntegration()
	s.piServiceMock.On("GetPackagingIntegration", pi.ID).Return(pi, nil)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/packaging/integration/%s", piID), nil)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result packaging.PackagingIntegration
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(result.Spec).Should(Equal(pi.Spec))
}

func (s *PackagingIntegrationRouteSuite) TestGetPackagingIntegrationNotFound() {
	notFoundID := "not-found"
	s.piServiceMock.
		On("GetPackagingIntegration", notFoundID).
		Return(nil, odahuErrors.NotFoundError{})

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/packaging/integration/%s", notFoundID), nil)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusNotFound))
	s.g.Expect(result.Message).Should(ContainSubstring("not found"))
}

func (s *PackagingIntegrationRouteSuite) TestCreatePackagingIntegration() {
	piEntity := newPackagingIntegration()

	s.piServiceMock.On("CreatePackagingIntegration", piEntity).Return(nil)

	piEntityBody, err := json.Marshal(piEntity)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/packaging/integration", bytes.NewReader(piEntityBody))
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var piResponse packaging.PackagingIntegration
	err = json.Unmarshal(w.Body.Bytes(), &piResponse)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusCreated))
	s.g.Expect(piResponse.ID).Should(Equal(piEntity.ID))
	s.g.Expect(piResponse.Spec).Should(Equal(piEntity.Spec))
}

func (s *PackagingIntegrationRouteSuite) TestCreateDuplicatePackagingIntegration() {
	pi := newPackagingIntegration()
	s.piServiceMock.
		On("CreatePackagingIntegration", pi).
		Return(odahuErrors.AlreadyExistError{})

	piEntityBody, err := json.Marshal(pi)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/packaging/integration", bytes.NewReader(piEntityBody))
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusConflict))
	s.g.Expect(result.Message).Should(ContainSubstring("already exists"))
}

func (s *PackagingIntegrationRouteSuite) TestUpdatePackagingIntegration() {
	piEntity := newPackagingIntegration()

	s.piServiceMock.On("UpdatePackagingIntegration", piEntity).Return(nil)

	piEntityBody, err := json.Marshal(piEntity)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPut, "/packaging/integration", bytes.NewReader(piEntityBody))
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var piResponse packaging.PackagingIntegration
	err = json.Unmarshal(w.Body.Bytes(), &piResponse)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(piResponse.ID).Should(Equal(piEntity.ID))
	s.g.Expect(piResponse.Spec).Should(Equal(piEntity.Spec))
}

func (s *PackagingIntegrationRouteSuite) TestUpdatePackagingIntegrationNotFound() {
	pi := newPackagingIntegration()

	s.piServiceMock.On("UpdatePackagingIntegration", pi).Return(odahuErrors.NotFoundError{})

	piEntityBody, err := json.Marshal(pi)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPut, "/packaging/integration", bytes.NewReader(piEntityBody))
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusNotFound))
	s.g.Expect(result.Message).Should(ContainSubstring("not found"))
}

func (s *PackagingIntegrationRouteSuite) TestDeletePackagingIntegration() {
	pi := newPackagingIntegration()
	s.piServiceMock.On("DeletePackagingIntegration", pi.ID).Return(nil)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/packaging/integration/%s", piID), nil)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(result.Message).Should(ContainSubstring("was deleted"))
}

func (s *PackagingIntegrationRouteSuite) TestDeletePackagingIntegrationNotFound() {
	s.piServiceMock.On("DeletePackagingIntegration", piID).Return(odahuErrors.NotFoundError{})

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/packaging/integration/%s", piID), nil)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusNotFound))
	s.g.Expect(result.Message).Should(ContainSubstring("not found"))
}

// GET methods still work despite disabled API
func (s *PackagingIntegrationRouteSuite) TestDisabledAPIGetPackagingIntegration() {
	packagingConfig := config.NewDefaultModelPackagingConfig()
	packagingConfig.Enabled = false
	s.registerHandlers(packagingConfig)

	pi := newPackagingIntegration()
	s.piServiceMock.On("GetPackagingIntegration", pi.ID).Return(pi, nil)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/packaging/integration/%s", piID), nil)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result packaging.PackagingIntegration
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(result.Spec).Should(Equal(pi.Spec))
}

func (s *PackagingIntegrationRouteSuite) TestDisabledAPIGetAllPackagingIntegration() {
	packagingConfig := config.NewDefaultModelPackagingConfig()
	packagingConfig.Enabled = false
	s.registerHandlers(packagingConfig)

	piList := []packaging.PackagingIntegration{*newPackagingIntegration()}
	s.piServiceMock.On("GetPackagingIntegrationList", mock.Anything, mock.Anything).Return(piList, nil)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/packaging/integration", nil)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result []packaging.PackagingIntegration
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(result).Should(Equal(piList))
}

func (s *PackagingIntegrationRouteSuite) TestDisabledAPICreatePackagingIntegration() {
	packagingConfig := config.NewDefaultModelPackagingConfig()
	packagingConfig.Enabled = false
	s.registerHandlers(packagingConfig)

	pi := newPackagingIntegration()

	piEntityBody, err := json.Marshal(pi)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/packaging/integration", bytes.NewReader(piEntityBody))
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusBadRequest))
	s.g.Expect(result.Message).Should(ContainSubstring(routes.DisabledAPIErrorMessage))
}

func (s *PackagingIntegrationRouteSuite) TestDisabledAPIUpdatePackagingIntegration() {
	packagingConfig := config.NewDefaultModelPackagingConfig()
	packagingConfig.Enabled = false
	s.registerHandlers(packagingConfig)

	pi := newPackagingIntegration()
	piEntityBody, err := json.Marshal(pi)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPut, "/packaging/integration", bytes.NewReader(piEntityBody))
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusBadRequest))
	s.g.Expect(result.Message).Should(ContainSubstring(routes.DisabledAPIErrorMessage))
}

func (s *PackagingIntegrationRouteSuite) TestDisabledAPIDeletePackagingIntegration() {
	packagingConfig := config.NewDefaultModelPackagingConfig()
	packagingConfig.Enabled = false
	s.registerHandlers(packagingConfig)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/packaging/integration/%s", piID), nil)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusBadRequest))
	s.g.Expect(result.Message).Should(ContainSubstring(routes.DisabledAPIErrorMessage))
}
