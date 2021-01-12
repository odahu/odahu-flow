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
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	httputil "github.com/odahu/odahu-flow/packages/operator/pkg/utils/httputil"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/connection"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/packaging"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/routes"
	pack_route "github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/routes/v1/packaging"
	odahuErrors "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	mp_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/packaging"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"
	"k8s.io/apimachinery/pkg/api/errors"
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
	g            *GomegaWithT
	server       *gin.Engine
	mpRepository mp_repository.PackagingIntegrationRepository
}

func (s *PackagingIntegrationRouteSuite) SetupTest() {
	s.g = NewGomegaWithT(s.T())
	s.registerHandlers(config.NewDefaultModelPackagingConfig())
}

func (s *PackagingIntegrationRouteSuite) registerHandlers(packagingConfig config.ModelPackagingConfig) {
	s.server = gin.Default()
	v1Group := s.server.Group("")
	packGroup := v1Group.Group("", routes.DisableAPIMiddleware(packagingConfig.Enabled))
	pack_route.ConfigurePiRoutes(packGroup, s.mpRepository)
}

func (s *PackagingIntegrationRouteSuite) TearDownTest() {
	for _, piID := range []string{piID} {
		err := s.mpRepository.DeletePackagingIntegration(piID)
		if err != nil && !errors.IsNotFound(err) && !odahuErrors.IsNotFoundError(err) {
			s.T().Fail()
		}
	}
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
	s.g.Expect(s.mpRepository.CreatePackagingIntegration(pi)).NotTo(HaveOccurred())

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
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/packaging/integration/not-found", nil)
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

	pi, err := s.mpRepository.GetPackagingIntegration(piID)
	s.g.Expect(err).ShouldNot(HaveOccurred())
	s.g.Expect(pi.Spec).To(Equal(piEntity.Spec))
}

// CreatedAt and UpdatedAt field should automatically be updated after create request
func (s *PackagingIntegrationRouteSuite) TestCreatePackagingIntegrationModifiable() {
	newResource := newPackagingIntegration()

	newResourceBody, err := json.Marshal(newResource)
	s.g.Expect(err).NotTo(HaveOccurred())

	reqTime := routes.GetTimeNowTruncatedToSeconds()
	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodPost, pack_route.CreatePackagingIntegrationURL, bytes.NewReader(newResourceBody),
	)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var resp packaging.ModelPackaging
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

func (s *PackagingIntegrationRouteSuite) TestCreateDuplicatePackagingIntegration() {
	pi := newPackagingIntegration()
	s.g.Expect(s.mpRepository.CreatePackagingIntegration(pi)).NotTo(HaveOccurred())

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
	pi := newPackagingIntegration()
	s.g.Expect(s.mpRepository.CreatePackagingIntegration(pi)).NotTo(HaveOccurred())

	piEntity := newPackagingIntegration()
	piEntity.Spec.Entrypoint = "new-entrypoint"

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

	pi, err = s.mpRepository.GetPackagingIntegration(piID)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.g.Expect(pi.Spec).To(Equal(piEntity.Spec))
}

// UpdatedAt field should automatically be updated after update request
func (s *PackagingIntegrationRouteSuite) TestUpdatePackagingIntegrationModifiable() {
	resource := newPackagingIntegration()
	s.g.Expect(s.mpRepository.CreatePackagingIntegration(resource)).NotTo(HaveOccurred())

	time.Sleep(1 * time.Second)

	newResource := newPackagingIntegration()

	newResourceBody, err := json.Marshal(newResource)
	s.g.Expect(err).NotTo(HaveOccurred())

	reqTime := routes.GetTimeNowTruncatedToSeconds()
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPut, pack_route.UpdatePackagingIntegrationURL, bytes.NewReader(newResourceBody))
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var respResource packaging.ModelPackaging
	err = json.Unmarshal(w.Body.Bytes(), &respResource)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(respResource.Status.CreatedAt).NotTo(BeNil())
	createdAtWasNotUpdated := reqTime.After(respResource.Status.CreatedAt.Time)
	s.g.Expect(createdAtWasNotUpdated).Should(Equal(true))
	s.g.Expect(respResource.Status.UpdatedAt).NotTo(BeNil())
	updatedAtWasUpdated := reqTime.Before(respResource.Status.UpdatedAt) || reqTime.Equal(respResource.Status.UpdatedAt)
	s.g.Expect(updatedAtWasUpdated).Should(Equal(true))
}

func (s *PackagingIntegrationRouteSuite) TestUpdatePackagingIntegrationNotFound() {
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

	s.g.Expect(w.Code).Should(Equal(http.StatusNotFound))
	s.g.Expect(result.Message).Should(ContainSubstring("not found"))
}

func (s *PackagingIntegrationRouteSuite) TestDeletePackagingIntegration() {
	pi := newPackagingIntegration()
	s.g.Expect(s.mpRepository.CreatePackagingIntegration(pi)).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/packaging/integration/%s", piID), nil)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(result.Message).Should(ContainSubstring("was deleted"))

	piList, err := s.mpRepository.GetPackagingIntegrationList()
	s.g.Expect(err).NotTo(HaveOccurred())
	s.g.Expect(piList).To(HaveLen(0))
}

func (s *PackagingIntegrationRouteSuite) TestDeletePackagingIntegrationNotFound() {
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

func (s *PackagingIntegrationRouteSuite) TestDisabledAPIGetPackagingIntegration() {
	packagingConfig := config.NewDefaultModelPackagingConfig()
	packagingConfig.Enabled = false
	s.registerHandlers(packagingConfig)

	pi := newPackagingIntegration()
	s.g.Expect(s.mpRepository.CreatePackagingIntegration(pi)).NotTo(HaveOccurred())

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

	pi := newPackagingIntegration()
	s.g.Expect(s.mpRepository.CreatePackagingIntegration(pi)).NotTo(HaveOccurred())

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

func (s *PackagingIntegrationRouteSuite) TestDisabledAPICreatePackagingIntegration() {
	packagingConfig := config.NewDefaultModelPackagingConfig()
	packagingConfig.Enabled = false
	s.registerHandlers(packagingConfig)

	pi := newPackagingIntegration()
	s.g.Expect(s.mpRepository.CreatePackagingIntegration(pi)).NotTo(HaveOccurred())

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
