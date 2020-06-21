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

package deployment_test

import (
	"bytes"
	"encoding/json"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	odahuflowv1alpha1 "github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/deployment"
	dep_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/deployment"
	dep_k8s_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/deployment/kubernetes"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils"
	"github.com/odahu/odahu-flow/packages/operator/pkg/webserver/routes"
	dep_route "github.com/odahu/odahu-flow/packages/operator/pkg/webserver/routes/v1/deployment"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

const (
	mrID  = "test-mr"
	mrID1 = "test-mr1"
	mrID2 = "test-mr2"
	mrURL = "/test/url"
)

type ModelRouteSuite struct {
	suite.Suite
	g            *GomegaWithT
	server       *gin.Engine
	mdRepository dep_repository.Repository
}

func (s *ModelRouteSuite) SetupSuite() {
	mgr, err := manager.New(cfg, manager.Options{NewClient: utils.NewClient, MetricsBindAddress: "0"})
	if err != nil {
		panic(err)
	}

	s.mdRepository = dep_k8s_repository.NewRepositoryWithOptions(
		testNamespace, mgr.GetClient(), metav1.DeletePropagationBackground,
	)

	err = s.mdRepository.CreateModelDeployment(&deployment.ModelDeployment{
		ID: mdID1,
		Spec: odahuflowv1alpha1.ModelDeploymentSpec{
			Image:                      mdImage,
			MinReplicas:                &mdMinReplicas,
			MaxReplicas:                &mdMaxReplicas,
			LivenessProbeInitialDelay:  &mdLivenessInitialDelay,
			ReadinessProbeInitialDelay: &mdReadinessInitialDelay,
			Annotations:                mdAnnotations,
			Resources:                  mdResources,
			RoleName:                   &mdRoleName,
		},
	})
	if err != nil {
		panic(err)
	}

	err = s.mdRepository.CreateModelDeployment(&deployment.ModelDeployment{
		ID: mdID2,
		Spec: odahuflowv1alpha1.ModelDeploymentSpec{
			Image:                      mdImage,
			MinReplicas:                &mdMinReplicas,
			MaxReplicas:                &mdMaxReplicas,
			LivenessProbeInitialDelay:  &mdLivenessInitialDelay,
			ReadinessProbeInitialDelay: &mdReadinessInitialDelay,
			Annotations:                mdAnnotations,
			Resources:                  mdResources,
			RoleName:                   &mdRoleName,
		},
	})
	if err != nil {
		panic(err)
	}
}

func (s *ModelRouteSuite) TearDownSuite() {
	for _, mdID := range []string{mdID1, mdID2} {
		if err := s.mdRepository.DeleteModelDeployment(mdID); err != nil && !errors.IsNotFound(err) {
			panic(err)
		}
	}
}

func (s *ModelRouteSuite) SetupTest() {
	s.g = NewGomegaWithT(s.T())

	s.registerHTTPHandlers(config.NewDefaultModelDeploymentConfig())
}

func (s *ModelRouteSuite) registerHTTPHandlers(deploymentConfig config.ModelDeploymentConfig) {
	s.server = gin.Default()
	v1Group := s.server.Group("")
	dep_route.ConfigureRoutes(v1Group, s.mdRepository, deploymentConfig, config.NvidiaResourceName)
}

func (s *ModelRouteSuite) TearDownTest() {
	for _, mdID := range []string{mrID, mrID1, mrID2} {
		if err := s.mdRepository.DeleteModelRoute(mdID); err != nil && !errors.IsNotFound(err) {
			panic(err)
		}
	}
}

func newStubMr() *deployment.ModelRoute {
	return &deployment.ModelRoute{
		ID: mrID,
		Spec: odahuflowv1alpha1.ModelRouteSpec{
			URLPrefix: mrURL,
			ModelDeploymentTargets: []odahuflowv1alpha1.ModelDeploymentTarget{
				{
					Name:   mdID1,
					Weight: &dep_route.MaxWeight,
				},
			},
		},
	}
}

func TestModelRouteSuite(t *testing.T) {
	suite.Run(t, new(ModelRouteSuite))
}

func (s *ModelRouteSuite) TestGetMR() {
	mr := newStubMr()
	s.g.Expect(s.mdRepository.CreateModelRoute(mr)).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodGet,
		strings.Replace(dep_route.GetModelRouteURL, ":id", mrID, -1),
		nil,
	)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result deployment.ModelRoute
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(result.Spec).Should(Equal(mr.Spec))
}

func (s *ModelRouteSuite) TestGetMRNotFound() {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodGet,
		strings.Replace(dep_route.GetModelRouteURL, ":id", "not-found", -1),
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

func (s *ModelRouteSuite) TestGetAllModelRoutes() {
	conn := newStubMr()
	s.g.Expect(s.mdRepository.CreateModelRoute(conn)).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodGet,
		dep_route.GetAllModelRouteURL,
		nil,
	)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result []deployment.ModelRoute
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(result).Should(HaveLen(1))
	s.g.Expect(result[0].ID).Should(Equal(conn.ID))
	s.g.Expect(result[0].Spec).Should(Equal(conn.Spec))
}

func (s *ModelRouteSuite) TestGetAllEmptyModelRoutes() {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodGet,
		dep_route.GetAllModelRouteURL,
		nil,
	)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result []deployment.ModelRoute
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(result).Should(HaveLen(0))
}

func (s *ModelRouteSuite) TestGetAllModelRoutesPaging() {
	mr1 := newStubMr()
	mr1.ID = mrID1
	s.g.Expect(s.mdRepository.CreateModelRoute(mr1)).NotTo(HaveOccurred())

	mr2 := newStubMr()
	mr2.ID = mrID2
	s.g.Expect(s.mdRepository.CreateModelRoute(mr2)).NotTo(HaveOccurred())

	connNames := map[string]interface{}{mrID1: nil, mrID2: nil}

	// Return first page
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, dep_route.GetAllModelRouteURL, nil)
	s.g.Expect(err).NotTo(HaveOccurred())

	query := req.URL.Query()
	query.Set("size", "1")
	query.Set("page", "0")
	req.URL.RawQuery = query.Encode()

	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result []deployment.ModelRoute
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(result).Should(HaveLen(1))
	delete(connNames, result[0].ID)

	// Return second page
	w = httptest.NewRecorder()
	req, err = http.NewRequest(http.MethodGet, dep_route.GetAllModelRouteURL, nil)
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
	req, err = http.NewRequest(http.MethodGet, dep_route.GetAllModelRouteURL, nil)
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

func (s *ModelRouteSuite) TestCreateMR() {
	mrEntity := newStubMr()

	mrEntityBody, err := json.Marshal(mrEntity)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, dep_route.CreateModelRouteURL, bytes.NewReader(mrEntityBody))
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var mrResponse deployment.ModelRoute
	err = json.Unmarshal(w.Body.Bytes(), &mrResponse)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusCreated))
	s.g.Expect(mrResponse.ID).To(Equal(mrEntity.ID))
	s.g.Expect(mrResponse.Spec).To(Equal(mrEntity.Spec))

	mr, err := s.mdRepository.GetModelRoute(mrID)
	s.g.Expect(err).ShouldNot(HaveOccurred())
	s.g.Expect(mr.ID).To(Equal(mrEntity.ID))
	s.g.Expect(mr.Spec).To(Equal(mrEntity.Spec))
}

// CreatedAt and UpdatedAt field should automatically be updated after create request
func (s *ModelRouteSuite) TestCreateMRModifiable() {
	newResource := newStubMr()

	newResourceBody, err := json.Marshal(newResource)
	s.g.Expect(err).NotTo(HaveOccurred())

	reqTime := routes.GetTimeNowTruncatedToSeconds()
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, dep_route.CreateModelRouteURL, bytes.NewReader(newResourceBody))
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var resp deployment.ModelRoute
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

func (s *ModelRouteSuite) TestCreateDuplicateMR() {
	mr := newStubMr()
	s.g.Expect(s.mdRepository.CreateModelRoute(mr)).NotTo(HaveOccurred())

	mrEntityBody, err := json.Marshal(mr)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, dep_route.CreateModelRouteURL, bytes.NewReader(mrEntityBody))
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result routes.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusConflict))
	s.g.Expect(result.Message).Should(ContainSubstring("already exists"))
}

func (s *ModelRouteSuite) TestValidateCreateMR() {
	mr := newStubMr()
	mr.Spec.URLPrefix = ""

	mrEntity, err := json.Marshal(mr)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, dep_route.CreateModelRouteURL, bytes.NewReader(mrEntity))
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result routes.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusBadRequest))
	s.g.Expect(result.Message).Should(ContainSubstring(dep_route.URLPrefixEmptyErrorMessage))
}

func (s *ModelRouteSuite) TestUpdateMR() {
	mr := newStubMr()
	s.g.Expect(s.mdRepository.CreateModelRoute(mr)).NotTo(HaveOccurred())

	newURL := "/new/url"
	mrEntity := newStubMr()
	mrEntity.Spec.URLPrefix = newURL

	mrEntityBody, err := json.Marshal(mrEntity)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPut, dep_route.UpdateModelRouteURL, bytes.NewReader(mrEntityBody))
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var mrResponse deployment.ModelRoute
	err = json.Unmarshal(w.Body.Bytes(), &mrResponse)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(mrResponse.ID).To(Equal(mrEntity.ID))
	s.g.Expect(mrResponse.Spec).To(Equal(mrEntity.Spec))

	mr, err = s.mdRepository.GetModelRoute(mrID)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.g.Expect(mr.ID).To(Equal(mrEntity.ID))
	s.g.Expect(mr.Spec).To(Equal(mrEntity.Spec))
}

// UpdatedAt field should automatically be updated after update request
func (s *ModelRouteSuite) TestUpdateMRModifiable() {
	resource := newStubMr()
	s.g.Expect(s.mdRepository.CreateModelRoute(resource)).NotTo(HaveOccurred())

	time.Sleep(1 * time.Second)

	newResource := newStubMr()

	newResourceBody, err := json.Marshal(newResource)
	s.g.Expect(err).NotTo(HaveOccurred())

	reqTime := routes.GetTimeNowTruncatedToSeconds()
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPut, dep_route.UpdateModelRouteURL, bytes.NewReader(newResourceBody))
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var respResource deployment.ModelRoute
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

func (s *ModelRouteSuite) TestUpdateMRNotFound() {
	mrEntity := newStubMr()

	mrEntityBody, err := json.Marshal(mrEntity)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPut, dep_route.UpdateModelRouteURL, bytes.NewReader(mrEntityBody))
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result routes.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusNotFound))
	s.g.Expect(result.Message).Should(ContainSubstring("not found"))
}

func (s *ModelRouteSuite) TestValidateUpdateMR() {
	mr := newStubMr()
	s.g.Expect(s.mdRepository.CreateModelRoute(mr)).NotTo(HaveOccurred())

	mr.Spec.URLPrefix = ""
	connEntityBody, err := json.Marshal(mr)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPut, dep_route.UpdateModelRouteURL, bytes.NewReader(connEntityBody))
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result routes.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusBadRequest))
	s.g.Expect(result.Message).Should(ContainSubstring(dep_route.URLPrefixEmptyErrorMessage))
}

func (s *ModelRouteSuite) TestDeleteMR() {
	mr := newStubMr()
	s.g.Expect(s.mdRepository.CreateModelRoute(mr)).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodDelete,
		strings.Replace(dep_route.DeleteModelRouteURL, ":id", mrID, -1),
		nil,
	)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result routes.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(result.Message).Should(ContainSubstring("was deleted"))

	mrList, err := s.mdRepository.GetModelRouteList()
	s.g.Expect(err).NotTo(HaveOccurred())
	s.g.Expect(mrList).To(HaveLen(0))
}

func (s *ModelRouteSuite) TestDeleteMRNotFound() {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodDelete,
		strings.Replace(dep_route.DeleteModelRouteURL, ":id", mrID, -1),
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

func (s *ModelRouteSuite) TestDisabledAPIGetMR() {
	deploymentConfig := config.NewDefaultModelDeploymentConfig()
	deploymentConfig.Enabled = false
	s.registerHTTPHandlers(deploymentConfig)

	mr := newStubMr()
	s.g.Expect(s.mdRepository.CreateModelRoute(mr)).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodGet,
		strings.Replace(dep_route.GetModelRouteURL, ":id", mrID, -1),
		nil,
	)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result deployment.ModelRoute
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(result.Spec).Should(Equal(mr.Spec))
}

func (s *ModelRouteSuite) TestDisabledAPIGetAllModelRoutes() {
	deploymentConfig := config.NewDefaultModelDeploymentConfig()
	deploymentConfig.Enabled = false
	s.registerHTTPHandlers(deploymentConfig)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodGet,
		dep_route.GetAllModelRouteURL,
		nil,
	)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result []deployment.ModelRoute
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(result).Should(HaveLen(0))
}

func (s *ModelRouteSuite) TestDisabledAPICreateMR() {
	deploymentConfig := config.NewDefaultModelDeploymentConfig()
	deploymentConfig.Enabled = false
	s.registerHTTPHandlers(deploymentConfig)

	mr := newStubMr()
	mrEntityBody, err := json.Marshal(mr)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, dep_route.CreateModelRouteURL, bytes.NewReader(mrEntityBody))
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result routes.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusBadRequest))
	s.g.Expect(result.Message).Should(ContainSubstring(routes.DisabledAPIErrorMessage))
}

func (s *ModelRouteSuite) TestDisabledAPIUpdateMR() {
	deploymentConfig := config.NewDefaultModelDeploymentConfig()
	deploymentConfig.Enabled = false
	s.registerHTTPHandlers(deploymentConfig)

	mrEntity := newStubMr()

	mrEntityBody, err := json.Marshal(mrEntity)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPut, dep_route.UpdateModelRouteURL, bytes.NewReader(mrEntityBody))
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result routes.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusBadRequest))
	s.g.Expect(result.Message).Should(ContainSubstring(routes.DisabledAPIErrorMessage))
}

func (s *ModelRouteSuite) TestDisabledAPIDeleteMR() {
	deploymentConfig := config.NewDefaultModelDeploymentConfig()
	deploymentConfig.Enabled = false
	s.registerHTTPHandlers(deploymentConfig)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodDelete,
		strings.Replace(dep_route.DeleteModelRouteURL, ":id", mrID, -1),
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
