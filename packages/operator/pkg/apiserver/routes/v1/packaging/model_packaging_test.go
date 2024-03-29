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
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	odahuflowv1alpha1 "github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/connection"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/packaging"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/training"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/routes"
	pack_route "github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/routes/v1/packaging"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	"github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	kube_client "github.com/odahu/odahu-flow/packages/operator/pkg/kubeclient/packagingclient"
	"github.com/odahu/odahu-flow/packages/operator/pkg/odahuflow"
	conn_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/connection"
	conn_k8s_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/connection/kubernetes"
	mp_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/packaging"
	mp_postgres_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/packaging/postgres"
	mp_service "github.com/odahu/odahu-flow/packages/operator/pkg/service/packaging"
	"github.com/odahu/odahu-flow/packages/operator/pkg/service/packaging_integration"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"
	httputil "github.com/odahu/odahu-flow/packages/operator/pkg/utils/httputil"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"net/http/httptest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
	"testing"
)

var (
	mpIDRoute           = "test-id"
	piIDMpRoute         = "pi-id"
	piEntrypointMpRoute = "/usr/bin/test"
	piImageMpRoute      = "test:image"
)

const (
	testOutConn         = "some-output-connection"
	testOutConnDefault  = "default-output-connection"
	testOutConnNotFound = "out-conn-not-found"
)

type packagingIntegrationService interface {
	GetPackagingIntegration(id string) (*packaging.PackagingIntegration, error)
	GetPackagingIntegrationList(options ...filter.ListOption) ([]packaging.PackagingIntegration, error)
	CreatePackagingIntegration(md *packaging.PackagingIntegration) error
	UpdatePackagingIntegration(md *packaging.PackagingIntegration) error
	DeletePackagingIntegration(id string) error
}

type ModelPackagingRouteSuite struct {
	suite.Suite
	g              *GomegaWithT
	server         *gin.Engine
	packRepo       mp_repository.Repository
	packService    mp_service.Service
	piService      packagingIntegrationService
	kubePackClient kube_client.Client
	connStorage    conn_repository.Repository
	k8sClient      client.Client
}

func (s *ModelPackagingRouteSuite) SetupSuite() {

	s.k8sClient = kubeClient
	s.kubePackClient = kube_client.NewClient(testNamespace, testNamespace, s.k8sClient, cfg)
	piRepo := mp_postgres_repository.PackagingIntegrationRepository{DB: db}
	s.piService = packaging_integration.NewService(&piRepo)
	s.packRepo = mp_postgres_repository.PackagingRepo{DB: db}
	s.packService = mp_service.NewService(s.packRepo)

	err := s.piService.CreatePackagingIntegration(&packaging.PackagingIntegration{
		ID: piIDMpRoute,
		Spec: packaging.PackagingIntegrationSpec{
			Entrypoint:   piEntrypointMpRoute,
			DefaultImage: piImageMpRoute,
			Schema:       packaging.Schema{},
		},
	})
	if err != nil {
		s.T().Fatalf("Cannot create PackagingIntegration: %v", err)
	}

	s.connStorage = conn_k8s_repository.NewRepository(testNamespace, s.k8sClient)
	// Create the connection that will be used as the outputConnection param for a training.
	if err := s.connStorage.SaveConnection(&connection.Connection{
		ID: testOutConn,
		Spec: odahuflowv1alpha1.ConnectionSpec{
			Type: connection.GcsType,
		},
	}); err != nil {
		s.T().Fatalf("Cannot create Connection: %v", err)
	}

	// Create the connection that will be used as the default outputConnection param for a training.
	if err := s.connStorage.SaveConnection(&connection.Connection{
		ID: testOutConnDefault,
		Spec: odahuflowv1alpha1.ConnectionSpec{
			Type: connection.GcsType,
		},
	}); err != nil {
		s.T().Fatalf("Cannot create Connection: %v", err)
	}
}

func (s *ModelPackagingRouteSuite) TearDownSuite() {
	if err := s.piService.DeletePackagingIntegration(piIDMpValid); err != nil {
		s.T().Fatal(err)
	}

}

func (s *ModelPackagingRouteSuite) SetupTest() {
	s.g = NewGomegaWithT(s.T())

	s.registerHandlers(config.NewDefaultModelPackagingConfig())
}

func (s *ModelPackagingRouteSuite) registerHandlers(packagingConfig config.ModelPackagingConfig) {
	s.server = gin.Default()
	v1Group := s.server.Group("")
	packGroup := v1Group.Group("", routes.DisableAPIMiddleware(packagingConfig.Enabled))

	pack_route.ConfigureRoutes(
		packGroup, s.kubePackClient, s.packService,
		s.piService, s.connStorage, packagingConfig,
		config.NvidiaResourceName,
	)
}

func (s *ModelPackagingRouteSuite) TearDownTest() {
	for _, mpID := range []string{mpIDRoute, testMpID1, testMpID2} {
		if err := s.packService.DeleteModelPackaging(
			context.Background(), mpID); err != nil && !errors.IsNotFoundError(err) {
			// If a model packaging is not found then it was not created during a test case
			s.T().Fatalf("Cannot delete ModelPackaging: %v", err)
		}
	}
}

func newModelPackaging() *packaging.ModelPackaging {
	resources := config.NewDefaultModelTrainingConfig().DefaultResources
	return &packaging.ModelPackaging{
		ID: mpIDRoute,
		Spec: packaging.ModelPackagingSpec{
			ArtifactName:     mpArtifactName,
			IntegrationName:  piIDMpRoute,
			Image:            mpImage,
			Resources:        &resources,
			OutputConnection: testOutConn,
		},
	}
}

func (s *ModelPackagingRouteSuite) createModelPackagings() {
	mp1 := newModelPackaging()
	mp1.ID = testMpID1
	s.g.Expect(s.packService.CreateModelPackaging(context.Background(), mp1)).NotTo(HaveOccurred())

	mp2 := newModelPackaging()
	mp2.ID = testMpID2
	s.g.Expect(s.packService.CreateModelPackaging(context.Background(), mp2)).NotTo(HaveOccurred())
}

func TestModelPackagingRouteSuite(t *testing.T) {
	suite.Run(t, new(ModelPackagingRouteSuite))
}

func (s *ModelPackagingRouteSuite) TestGetMP() {
	mp := newModelPackaging()
	s.g.Expect(s.packService.CreateModelPackaging(context.Background(), mp)).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodGet,
		strings.Replace(pack_route.GetModelPackagingURL, ":id", mp.ID, -1),
		nil,
	)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result packaging.ModelPackaging
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(result.Spec).Should(Equal(mp.Spec))
}

func (s *ModelPackagingRouteSuite) TestGetMPNotFound() {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodGet,
		strings.Replace(pack_route.GetModelPackagingURL, ":id", "not-found", -1),
		nil,
	)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusNotFound))
	s.g.Expect(result.Message).Should(ContainSubstring("not found"))
}

func (s *ModelPackagingRouteSuite) TestGetAllMPEmptyResult() {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, pack_route.GetAllModelPackagingURL, nil)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var mpResponse []packaging.ModelPackaging
	err = json.Unmarshal(w.Body.Bytes(), &mpResponse)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(mpResponse).Should(HaveLen(0))
}

func (s *ModelPackagingRouteSuite) TestGetAllMP() {
	s.createModelPackagings()

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, pack_route.GetAllModelPackagingURL, nil)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result []packaging.ModelPackaging
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(result).Should(HaveLen(2))

	for _, mp := range result {
		s.g.Expect(mp.ID).To(Or(Equal(testMpID1), Equal(testMpID2)))
	}
}

func (s *ModelPackagingRouteSuite) TestGetAllMTPaging() {
	s.createModelPackagings()

	mpNames := map[string]interface{}{testMpID1: nil, testMpID2: nil}

	// Return first page
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, pack_route.GetAllModelPackagingURL, nil)
	s.g.Expect(err).NotTo(HaveOccurred())

	query := req.URL.Query()
	query.Set("size", "1")
	query.Set("page", "0")
	req.URL.RawQuery = query.Encode()

	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var trainings []training.ModelTraining
	err = json.Unmarshal(w.Body.Bytes(), &trainings)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(trainings).Should(HaveLen(1))
	delete(mpNames, trainings[0].ID)

	// Return second page
	w = httptest.NewRecorder()
	req, err = http.NewRequest(http.MethodGet, pack_route.GetAllModelPackagingURL, nil)
	s.g.Expect(err).NotTo(HaveOccurred())

	query = req.URL.Query()
	query.Set("size", "1")
	query.Set("page", "1")
	req.URL.RawQuery = query.Encode()

	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	err = json.Unmarshal(w.Body.Bytes(), &trainings)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(trainings).Should(HaveLen(1))
	delete(mpNames, trainings[0].ID)

	// Return third empty page
	w = httptest.NewRecorder()
	req, err = http.NewRequest(http.MethodGet, pack_route.GetAllModelPackagingURL, nil)
	s.g.Expect(err).NotTo(HaveOccurred())

	query = req.URL.Query()
	query.Set("size", "1")
	query.Set("page", "2")
	req.URL.RawQuery = query.Encode()

	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	err = json.Unmarshal(w.Body.Bytes(), &trainings)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(trainings).Should(HaveLen(0))
	s.g.Expect(trainings).Should(BeEmpty())
}

func (s *ModelPackagingRouteSuite) TestCreateMP() {
	mpEntity := newModelPackaging()

	mpEntityBody, err := json.Marshal(mpEntity)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, pack_route.CreateModelPackagingURL, bytes.NewReader(mpEntityBody))
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var mpResponse packaging.ModelPackaging
	err = json.Unmarshal(w.Body.Bytes(), &mpResponse)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusCreated))
	s.g.Expect(mpResponse.ID).Should(Equal(mpEntity.ID))
	s.g.Expect(mpResponse.Spec).Should(Equal(mpEntity.Spec))

	mp, err := s.packService.GetModelPackaging(context.Background(), mpIDRoute)
	s.g.Expect(err).ShouldNot(HaveOccurred())
	s.g.Expect(mp.Spec).To(Equal(mpEntity.Spec))
}

func (s *ModelPackagingRouteSuite) TestCreateDuplicateMP() {
	mp := newModelPackaging()
	s.g.Expect(s.packService.CreateModelPackaging(context.Background(), mp)).NotTo(HaveOccurred())

	mpEntityBody, err := json.Marshal(mp)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, pack_route.CreateModelPackagingURL, bytes.NewReader(mpEntityBody))
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusConflict))
	s.g.Expect(result.Message).Should(ContainSubstring("already exists"))
}

func (s *ModelPackagingRouteSuite) TestUpdateMP() {
	mp := newModelPackaging()
	s.g.Expect(s.packService.CreateModelPackaging(context.Background(), mp)).NotTo(HaveOccurred())

	updatedMp := newModelPackaging()
	updatedMp.Spec.Image += "123"

	mpEntityBody, err := json.Marshal(updatedMp)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPut, pack_route.UpdateModelPackagingURL, bytes.NewReader(mpEntityBody))
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var mpResponse packaging.ModelPackaging
	err = json.Unmarshal(w.Body.Bytes(), &mpResponse)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(mpResponse.ID).Should(Equal(updatedMp.ID))
	s.g.Expect(mpResponse.Spec).Should(Equal(updatedMp.Spec))

	mp, err = s.packService.GetModelPackaging(context.Background(), mpIDRoute)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.g.Expect(mp.Spec).To(Equal(updatedMp.Spec))
}

func (s *ModelPackagingRouteSuite) TestUpdateMPNotFound() {
	mpEntity := newModelPackaging()

	mpEntityBody, err := json.Marshal(mpEntity)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPut, pack_route.UpdateModelPackagingURL, bytes.NewReader(mpEntityBody))
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusNotFound))
	s.g.Expect(result.Message).Should(ContainSubstring("not found"))
}

func (s *ModelPackagingRouteSuite) TestDeleteMP() {
	mp := newModelPackaging()
	s.g.Expect(s.packService.CreateModelPackaging(context.Background(), mp)).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodDelete,
		strings.Replace(pack_route.DeleteModelPackagingURL, ":id", mp.ID, -1),
		nil,
	)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(result.Message).Should(ContainSubstring("was deleted"))

	mpList, err := s.packService.GetModelPackagingList(context.Background())
	s.g.Expect(err).NotTo(HaveOccurred())
	s.g.Expect(mpList).To(HaveLen(0))
}

func (s *ModelPackagingRouteSuite) TestDeleteMPNotFound() {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodDelete,
		strings.Replace(pack_route.DeleteModelPackagingURL, ":id", "some-mp-id", -1),
		nil,
	)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusNotFound))
	s.g.Expect(result.Message).Should(ContainSubstring("not found"))
}

func (s *ModelPackagingRouteSuite) TestSavingMPResult() {
	resultCM := &corev1.ConfigMap{
		ObjectMeta: v1.ObjectMeta{
			Name:      odahuflow.GeneratePackageResultCMName(mpIDRoute),
			Namespace: testNamespace,
		},
	}
	s.g.Expect(s.k8sClient.Create(context.TODO(), resultCM)).NotTo(HaveOccurred())
	defer s.k8sClient.Delete(context.TODO(), resultCM)

	expectedMPResult := []odahuflowv1alpha1.ModelPackagingResult{
		{
			Name:  "test-name-1",
			Value: "test-value-1",
		},
		{
			Name:  "test-name-2",
			Value: "test-value-2",
		},
	}
	expectedMPResultBody, err := json.Marshal(expectedMPResult)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodPut,
		strings.Replace(pack_route.SaveModelPackagingResultURL, ":id", mpIDRoute, -1),
		bytes.NewReader(expectedMPResultBody),
	)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	result := []odahuflowv1alpha1.ModelPackagingResult{}
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.g.Expect(expectedMPResult).Should(Equal(result))

	result, err = s.kubePackClient.GetModelPackagingResult(mpIDRoute)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.g.Expect(expectedMPResult).To(Equal(result))
}

func (s *ModelPackagingRouteSuite) TestDisabledAPIDeleteMP() {
	packagingConfig := config.NewDefaultModelPackagingConfig()
	packagingConfig.Enabled = false
	s.registerHandlers(packagingConfig)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodDelete,
		strings.Replace(pack_route.DeleteModelPackagingURL, ":id", "some-mp-id", -1),
		nil,
	)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusBadRequest))
	s.g.Expect(result.Message).Should(ContainSubstring(routes.DisabledAPIErrorMessage))
}

func (s *ModelPackagingRouteSuite) TestDisabledAPIUpdateMP() {
	packagingConfig := config.NewDefaultModelPackagingConfig()
	packagingConfig.Enabled = false
	s.registerHandlers(packagingConfig)

	mpEntity := newModelPackaging()

	mpEntityBody, err := json.Marshal(mpEntity)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPut, pack_route.UpdateModelPackagingURL, bytes.NewReader(mpEntityBody))
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusBadRequest))
	s.g.Expect(result.Message).Should(ContainSubstring(routes.DisabledAPIErrorMessage))
}

func (s *ModelPackagingRouteSuite) TestDisabledAPICreateMP() {
	packagingConfig := config.NewDefaultModelPackagingConfig()
	packagingConfig.Enabled = false
	s.registerHandlers(packagingConfig)

	mpEntity := newModelPackaging()

	mpEntityBody, err := json.Marshal(mpEntity)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, pack_route.CreateModelPackagingURL, bytes.NewReader(mpEntityBody))
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusBadRequest))
	s.g.Expect(result.Message).Should(ContainSubstring(routes.DisabledAPIErrorMessage))
}

func (s *ModelPackagingRouteSuite) TestDisabledAPIGetMP() {
	packagingConfig := config.NewDefaultModelPackagingConfig()
	packagingConfig.Enabled = false
	s.registerHandlers(packagingConfig)

	mp := newModelPackaging()
	s.g.Expect(s.packService.CreateModelPackaging(context.Background(), mp)).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodGet,
		strings.Replace(pack_route.GetModelPackagingURL, ":id", mp.ID, -1),
		nil,
	)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result packaging.ModelPackaging
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(result.Spec).Should(Equal(mp.Spec))
}

func (s *ModelPackagingRouteSuite) TestDisabledAPIGetAllMP() {
	packagingConfig := config.NewDefaultModelPackagingConfig()
	packagingConfig.Enabled = false
	s.registerHandlers(packagingConfig)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, pack_route.GetAllModelPackagingURL, nil)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var mpResponse []packaging.ModelPackaging
	err = json.Unmarshal(w.Body.Bytes(), &mpResponse)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(mpResponse).Should(HaveLen(0))
}
