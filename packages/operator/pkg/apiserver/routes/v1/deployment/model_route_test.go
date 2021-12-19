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
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	odahuflowv1alpha1 "github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/deployment"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/event"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/routes"
	dep_route "github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/routes/v1/deployment"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/routes/v1/deployment/mocks"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	"github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	dep_repository_db "github.com/odahu/odahu-flow/packages/operator/pkg/repository/deployment/postgres"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/outbox"
	route_repository_db "github.com/odahu/odahu-flow/packages/operator/pkg/repository/route/postgres"
	md_service "github.com/odahu/odahu-flow/packages/operator/pkg/service/deployment"
	mr_service "github.com/odahu/odahu-flow/packages/operator/pkg/service/route"
	httputil "github.com/odahu/odahu-flow/packages/operator/pkg/utils/httputil"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const (
	mrID  = "test-mr"
	mrID1 = "test-mr1"
	mrID2 = "test-mr2"
	mrURL = "/custom/test/url"
)

type ModelRouteSuite struct {
	suite.Suite
	g              *GomegaWithT
	server         *gin.Engine
	mdService      md_service.Service
	mrService      mr_service.Service
	mrRepo         route_repository_db.RouteRepo
	mrEventsGetter *mocks.RoutesEventGetter
}

func (s *ModelRouteSuite) SetupSuite() {

	s.mrRepo = route_repository_db.RouteRepo{DB: db}
	s.mdService = md_service.NewService(
		dep_repository_db.DeploymentRepo{DB: db},
		route_repository_db.RouteRepo{DB: db},
		outbox.EventPublisher{DB: db})
	s.mrService = mr_service.NewService(s.mrRepo, outbox.EventPublisher{DB: db})
	s.mrEventsGetter = &mocks.RoutesEventGetter{}

	err := s.mdService.CreateModelDeployment(context.Background(), &deployment.ModelDeployment{
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

	err = s.mdService.CreateModelDeployment(context.Background(), &deployment.ModelDeployment{
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
		if err := s.mdService.DeleteModelDeployment(context.Background(), mdID); err != nil && !errors.IsNotFoundError(err) {
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
	dep_route.ConfigureRoutes(v1Group, s.mdService, nil, s.mrService, s.mrEventsGetter,
		deploymentConfig, config.NvidiaResourceName)
}

func (s *ModelRouteSuite) TearDownTest() {
	for _, mdID := range []string{mrID, mrID1, mrID2} {
		if err := s.mrService.DeleteModelRoute(context.Background(), mdID); err != nil && !errors.IsNotFoundError(err) {
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
	s.g.Expect(s.mrService.CreateModelRoute(context.Background(), mr)).NotTo(HaveOccurred())

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

	var result httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusNotFound))
	s.g.Expect(result.Message).Should(ContainSubstring("not found"))
}

func (s *ModelRouteSuite) TestGetAllModelRoutes() {
	conn := newStubMr()
	s.g.Expect(s.mrService.CreateModelRoute(context.Background(), conn)).NotTo(HaveOccurred())

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
	s.g.Expect(result).Should(HaveLen(3)) // two defaults and one that we created

	ids := make([]string, len(result))
	specs := make([]odahuflowv1alpha1.ModelRouteSpec, len(result))
	for i, v := range result {
		ids[i] = v.ID
		specs[i] = v.Spec
	}

	s.g.Expect(ids).Should(ContainElement(conn.ID))
	s.g.Expect(specs).Should(ContainElement(conn.Spec))
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
	s.g.Expect(result).Should(HaveLen(2)) // only suite deployments default routes
}

func (s *ModelRouteSuite) TestGetAllModelRoutesPaging() {
	mr1 := newStubMr()
	mr1.ID = mrID1
	s.g.Expect(s.mrService.CreateModelRoute(context.Background(), mr1)).NotTo(HaveOccurred())

	mr2 := newStubMr()
	mr2.ID = mrID2
	s.g.Expect(s.mrService.CreateModelRoute(context.Background(), mr2)).NotTo(HaveOccurred())

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
	query.Set("page", "4")
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

	mr, err := s.mrService.GetModelRoute(context.Background(), mrID)
	s.g.Expect(err).ShouldNot(HaveOccurred())
	s.g.Expect(mr.ID).To(Equal(mrEntity.ID))
	s.g.Expect(mr.Spec).To(Equal(mrEntity.Spec))
}

func (s *ModelRouteSuite) TestCreateDuplicateMR() {
	mr := newStubMr()
	s.g.Expect(s.mrService.CreateModelRoute(context.Background(), mr)).NotTo(HaveOccurred())

	mrEntityBody, err := json.Marshal(mr)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, dep_route.CreateModelRouteURL, bytes.NewReader(mrEntityBody))
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result httputil.HTTPResult
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

	var result httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusBadRequest))
	s.g.Expect(result.Message).Should(ContainSubstring(dep_route.URLPrefixEmptyErrorMessage))
}

func (s *ModelRouteSuite) TestUpdateMR() {
	mr := newStubMr()
	s.g.Expect(s.mrService.CreateModelRoute(context.Background(), mr)).NotTo(HaveOccurred())

	newURL := "/custom/new/url"
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

	mr, err = s.mrService.GetModelRoute(context.Background(), mrID)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.g.Expect(mr.ID).To(Equal(mrEntity.ID))
	s.g.Expect(mr.Spec).To(Equal(mrEntity.Spec))
}

func (s *ModelRouteSuite) TestUpdateDefaultRoute() {

	ctx := context.Background()

	r, err := md_service.GetDefaultModelRoute(ctx, nil, mdID1, s.mrRepo)
	s.g.Expect(err).NotTo(HaveOccurred())

	// API request

	newURL := "/custom/new/url"
	mrEntity := newStubMr()
	mrEntity.Spec.URLPrefix = newURL
	mrEntity.ID = r
	mrEntityBody, err := json.Marshal(mrEntity)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodPut,
		dep_route.UpdateModelRouteURL,
		bytes.NewReader(mrEntityBody),
	)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	s.g.Expect(w.Code).Should(Equal(403))

	var result httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())
	errMsg := fmt.Sprintf("access forbidden: unable to update default route with ID \"%v\"", r)
	s.g.Expect(result.Message).Should(Equal(errMsg))

}

func (s *ModelRouteSuite) TestUpdateMRNotFound() {
	mrEntity := newStubMr()

	mrEntityBody, err := json.Marshal(mrEntity)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPut, dep_route.UpdateModelRouteURL, bytes.NewReader(mrEntityBody))
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusNotFound))
	s.g.Expect(result.Message).Should(ContainSubstring("not found"))
}

func (s *ModelRouteSuite) TestValidateUpdateMR() {
	mr := newStubMr()
	s.g.Expect(s.mrService.CreateModelRoute(context.Background(), mr)).NotTo(HaveOccurred())

	mr.Spec.URLPrefix = ""
	connEntityBody, err := json.Marshal(mr)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPut, dep_route.UpdateModelRouteURL, bytes.NewReader(connEntityBody))
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusBadRequest))
	s.g.Expect(result.Message).Should(ContainSubstring(dep_route.URLPrefixEmptyErrorMessage))
}

func (s *ModelRouteSuite) TestDeleteMR() {
	mr := newStubMr()
	s.g.Expect(s.mrService.CreateModelRoute(context.Background(), mr)).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodDelete,
		strings.Replace(dep_route.DeleteModelRouteURL, ":id", mrID, -1),
		nil,
	)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(result.Message).Should(ContainSubstring("was deleted"))

	mrList, err := s.mrService.GetModelRouteList(context.Background())
	s.g.Expect(err).NotTo(HaveOccurred())
	s.g.Expect(mrList).To(HaveLen(2)) // only suite default routes
}

func (s *ModelRouteSuite) TestDeleteDefaultRoute() {

	ctx := context.Background()

	r, err := md_service.GetDefaultModelRoute(ctx, nil, mdID1, s.mrRepo)
	s.g.Expect(err).NotTo(HaveOccurred())

	// API request
	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodDelete,
		strings.Replace(dep_route.DeleteModelRouteURL, ":id", r, -1),
		nil,
	)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	s.g.Expect(w.Code).Should(Equal(403))

	var result httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())
	errMsg := fmt.Sprintf("access forbidden: unable to delete default route with ID \"%v\"", r)
	s.g.Expect(result.Message).Should(Equal(errMsg))

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

	var result httputil.HTTPResult
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
	s.g.Expect(s.mrService.CreateModelRoute(context.Background(), mr)).NotTo(HaveOccurred())

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
	s.g.Expect(result).Should(HaveLen(2)) // only suite deployments default routes
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

	var result httputil.HTTPResult
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

	var result httputil.HTTPResult
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

	var result httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusBadRequest))
	s.g.Expect(result.Message).Should(ContainSubstring(routes.DisabledAPIErrorMessage))
}

func (s *ModelRouteSuite) TestGetRouteEventsIncorrectCursor() {

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodGet,
		dep_route.EventsModelRouteURL+"?cursor=not-number",
		nil,
	)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)
	s.g.Expect(w.Code).Should(Equal(400))

	var result httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.g.Expect(result.Message).Should(ContainSubstring("Incorrect \"cursor\" query parameter"))
}

func (s *ModelRouteSuite) TestGetRouteEvents() {

	events := []event.RouteEvent{
		{
			Payload:   deployment.ModelRoute{ID: "route-1"},
			EventType: event.ModelRouteCreatedEventType,
			Datetime:  time.Time{},
		},
		{
			EntityID:  "route-2",
			EventType: event.ModelRouteDeletedEventType,
			Datetime:  time.Time{},
		},
	}

	s.mrEventsGetter.On("Get", mock.Anything, 0).Return(events, 5, nil)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodGet,
		dep_route.EventsModelRouteURL,
		nil,
	)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)
	s.g.Expect(w.Code).Should(Equal(200))

	var result event.LatestRouteEvents
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(result.Cursor).Should(Equal(5))
	s.g.Expect(result.Events).Should(Equal(events))
	s.g.Expect(err).NotTo(HaveOccurred())
}

func (s *ModelRouteSuite) TestGetRouteEventsWithCursor() {

	events := []event.RouteEvent{
		{
			Payload:   deployment.ModelRoute{ID: "route-1"},
			EventType: event.ModelRouteCreatedEventType,
			Datetime:  time.Time{},
		},
		{
			EntityID:  "route-2",
			EventType: event.ModelRouteDeletedEventType,
			Datetime:  time.Time{},
		},
	}

	s.mrEventsGetter.On("Get", mock.Anything, 6).Return(events, 9, nil)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodGet,
		dep_route.EventsModelRouteURL+"?cursor=6",
		nil,
	)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)
	s.g.Expect(w.Code).Should(Equal(200))

	var result event.LatestRouteEvents
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(result.Cursor).Should(Equal(9))
	s.g.Expect(result.Events).Should(Equal(events))
	s.g.Expect(err).NotTo(HaveOccurred())
}
