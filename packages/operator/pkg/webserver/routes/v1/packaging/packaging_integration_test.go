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
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	odahuflow_apis "github.com/odahu/odahu-flow/packages/operator/pkg/apis"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/connection"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/odahuflow/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/packaging"
	mp_config "github.com/odahu/odahu-flow/packages/operator/pkg/config/packaging"
	conn_k8s_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/connection/kubernetes"
	mp_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/packaging"
	mp_k8s_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/packaging/kubernetes"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils"
	"github.com/odahu/odahu-flow/packages/operator/pkg/webserver/routes"
	pack_route "github.com/odahu/odahu-flow/packages/operator/pkg/webserver/routes/v1/packaging"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/manager"
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

type mpiValidationSuite struct {
	suite.Suite
	g              *GomegaWithT
	server         *gin.Engine
	k8sEnvironment *envtest.Environment
	mpRepository   mp_repository.Repository
}

func (s *mpiValidationSuite) SetupSuite() {
	var cfg *rest.Config

	s.k8sEnvironment = &envtest.Environment{
		CRDDirectoryPaths: []string{filepath.Join("..", "..", "..", "..", "..", "config", "crds")},
	}

	err := odahuflow_apis.AddToScheme(scheme.Scheme)
	if err != nil {
		s.T().Fatalf("Cannot setup the odahuflow schema: %v", err)
	}

	cfg, err = s.k8sEnvironment.Start()
	if err != nil {
		s.T().Fatalf("Cannot setup the test k8s api: %v", err)
	}

	mgr, err := manager.New(cfg, manager.Options{NewClient: utils.NewClient})
	if err != nil {
		s.T().Fatalf("Cannot setup the test k8s manager: %v", err)
	}

	s.server = gin.Default()
	s.mpRepository = mp_k8s_repository.NewRepository(testNamespace, testNamespace, mgr.GetClient(), nil)
	pack_route.ConfigureRoutes(s.server.Group(""), s.mpRepository, conn_k8s_repository.NewRepository(
		testNamespace, mgr.GetClient(),
	))
}

func (s *mpiValidationSuite) TearDownSuite() {
	if err := s.k8sEnvironment.Stop(); err != nil {
		s.T().Fatal("Cannot stop the test k8s api")
	}
}

func (s *mpiValidationSuite) SetupTest() {
	s.g = NewGomegaWithT(s.T())
}

func (s *mpiValidationSuite) TearDownTest() {
	viper.Set(mp_config.Enabled, true)

	for _, piID := range []string{piID} {
		if err := s.mpRepository.DeletePackagingIntegration(piID); err != nil && !errors.IsNotFound(err) {
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

func TestModelPackagingIntegrationSuite(t *testing.T) {
	suite.Run(t, new(mpiValidationSuite))
}

func (s *mpiValidationSuite) TestGetPackagingIntegration() {
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

func (s *mpiValidationSuite) TestGetPackagingIntegrationNotFound() {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/packaging/integration/not-found", nil)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result routes.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusNotFound))
	s.g.Expect(result.Message).Should(ContainSubstring("not found"))
}

func (s *mpiValidationSuite) TestCreatePackagingIntegration() {
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
func (s *mpiValidationSuite) TestCreatePackagingIntegrationModifiable(){
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

func (s *mpiValidationSuite) TestCreateDuplicatePackagingIntegration() {
	pi := newPackagingIntegration()
	s.g.Expect(s.mpRepository.CreatePackagingIntegration(pi)).NotTo(HaveOccurred())

	piEntityBody, err := json.Marshal(pi)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/packaging/integration", bytes.NewReader(piEntityBody))
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result routes.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusConflict))
	s.g.Expect(result.Message).Should(ContainSubstring("already exists"))
}

func (s *mpiValidationSuite) TestUpdatePackagingIntegration() {
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
func (s *mpiValidationSuite) TestUpdatePackagingIntegrationModifiable(){
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

func (s *mpiValidationSuite) TestUpdatePackagingIntegrationNotFound() {
	pi := newPackagingIntegration()
	piEntityBody, err := json.Marshal(pi)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPut, "/packaging/integration", bytes.NewReader(piEntityBody))
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result routes.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusNotFound))
	s.g.Expect(result.Message).Should(ContainSubstring("not found"))
}

func (s *mpiValidationSuite) TestDeletePackagingIntegration() {
	pi := newPackagingIntegration()
	s.g.Expect(s.mpRepository.CreatePackagingIntegration(pi)).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/packaging/integration/%s", piID), nil)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result routes.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(result.Message).Should(ContainSubstring("was deleted"))

	piList, err := s.mpRepository.GetPackagingIntegrationList()
	s.g.Expect(err).NotTo(HaveOccurred())
	s.g.Expect(piList).To(HaveLen(0))
}

func (s *mpiValidationSuite) TestDeletePackagingIntegrationNotFound() {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/packaging/integration/%s", piID), nil)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result routes.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusNotFound))
	s.g.Expect(result.Message).Should(ContainSubstring("not found"))
}

func (s *mpiValidationSuite) TestDisabledAPIGetPackagingIntegration() {
	viper.Set(mp_config.Enabled, false)

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

func (s *mpiValidationSuite) TestDisabledAPIGetAllPackagingIntegration() {
	viper.Set(mp_config.Enabled, false)

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

func (s *mpiValidationSuite) TestDisabledAPICreatePackagingIntegration() {
	viper.Set(mp_config.Enabled, false)

	pi := newPackagingIntegration()
	s.g.Expect(s.mpRepository.CreatePackagingIntegration(pi)).NotTo(HaveOccurred())

	piEntityBody, err := json.Marshal(pi)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/packaging/integration", bytes.NewReader(piEntityBody))
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result routes.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusBadRequest))
	s.g.Expect(result.Message).Should(ContainSubstring(routes.DisabledAPIErrorMessage))
}

func (s *mpiValidationSuite) TestDisabledAPIUpdatePackagingIntegration() {
	viper.Set(mp_config.Enabled, false)

	pi := newPackagingIntegration()
	piEntityBody, err := json.Marshal(pi)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPut, "/packaging/integration", bytes.NewReader(piEntityBody))
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result routes.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusBadRequest))
	s.g.Expect(result.Message).Should(ContainSubstring(routes.DisabledAPIErrorMessage))
}

func (s *mpiValidationSuite) TestDisabledAPIDeletePackagingIntegration() {
	viper.Set(mp_config.Enabled, false)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/packaging/integration/%s", piID), nil)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result routes.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusBadRequest))
	s.g.Expect(result.Message).Should(ContainSubstring(routes.DisabledAPIErrorMessage))
}
