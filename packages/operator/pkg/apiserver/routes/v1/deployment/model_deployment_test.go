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
	"github.com/gin-gonic/gin"
	odahuflowv1alpha1 "github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/deployment"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/event"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/routes"
	dep_route "github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/routes/v1/deployment"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/routes/v1/deployment/mocks"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	"github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	"github.com/odahu/odahu-flow/packages/operator/pkg/odahuflow"
	dep_post_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/deployment/postgres"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/outbox"
	route_post_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/route/postgres"
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

var (
	mdID                    = "test-model-deployment"
	mdID1                   = "test-model-deployment1"
	mdID2                   = "test-model-deployment2"
	mdRoleName1             = "role-1"
	mdRoleName2             = "role-2"
	mdImage                 = "test/test:123"
	mdMinReplicas           = int32(1)
	mdMaxReplicas           = int32(2)
	mdLivenessInitialDelay  = int32(60)
	mdReadinessInitialDelay = int32(30)
	mdImagePullConnID       = "test-docker-pull-conn-id"
)

var (
	mdAnnotations = map[string]string{"k1": "v1", "k2": "v2"}
	reqMem        = "111Mi"
	reqCPU        = "111m"
	limMem        = "222Mi"
	mdResources   = &odahuflowv1alpha1.ResourceRequirements{
		Limits: &odahuflowv1alpha1.ResourceList{
			CPU:    nil,
			Memory: &limMem,
		},
		Requests: &odahuflowv1alpha1.ResourceList{
			CPU:    &reqCPU,
			Memory: &reqMem,
		},
	}
)

type ModelDeploymentRouteSuite struct {
	suite.Suite
	g              *GomegaWithT
	server         *gin.Engine
	mdService      md_service.Service
	mrService      mr_service.Service
	mdEventsGetter *mocks.ModelDeploymentEventGetter
}

func (s *ModelDeploymentRouteSuite) SetupSuite() {
	s.mdService = md_service.NewService(dep_post_repository.DeploymentRepo{DB: db}, route_post_repository.RouteRepo{
		DB: db,
	}, outbox.EventPublisher{DB: db})
	s.mrService = mr_service.NewService(route_post_repository.RouteRepo{DB: db}, outbox.EventPublisher{DB: db})
	s.mdEventsGetter = &mocks.ModelDeploymentEventGetter{}
}

func (s *ModelDeploymentRouteSuite) registerHTTPHandlers(deploymentConfig config.ModelDeploymentConfig) {
	s.server = gin.Default()
	v1Group := s.server.Group("")
	dep_route.ConfigureRoutes(v1Group, s.mdService, s.mdEventsGetter, s.mrService, nil,
		deploymentConfig, config.NvidiaResourceName)
}

func (s *ModelDeploymentRouteSuite) SetupTest() {
	s.g = NewGomegaWithT(s.T())
	s.registerHTTPHandlers(config.NewDefaultModelDeploymentConfig())
}

func (s *ModelDeploymentRouteSuite) TearDownTest() {
	ctx := context.Background()
	for _, currMdID := range []string{mdID, mdID1, mdID2} {
		if err := s.mdService.DeleteModelDeployment(ctx, currMdID); err != nil && !errors.IsNotFoundError(err) {
			panic(err)
		}
	}
}

func newStubMd() *deployment.ModelDeployment {
	return &deployment.ModelDeployment{
		ID: mdID,
		Spec: odahuflowv1alpha1.ModelDeploymentSpec{
			Image:                      mdImage,
			Predictor:                  odahuflow.OdahuMLServer.ID,
			MinReplicas:                &mdMinReplicas,
			MaxReplicas:                &mdMaxReplicas,
			LivenessProbeInitialDelay:  &mdLivenessInitialDelay,
			ReadinessProbeInitialDelay: &mdReadinessInitialDelay,
			Annotations:                mdAnnotations,
			Resources:                  mdResources,
			RoleName:                   &mdRoleName,
			ImagePullConnectionID:      &mdImagePullConnID,
		},
	}
}

func (s *ModelDeploymentRouteSuite) newMultipleMds() []*deployment.ModelDeployment {
	ctx := context.Background()
	md1 := newStubMd()
	md1.ID = mdID1
	md1.Spec.RoleName = &mdRoleName1
	s.g.Expect(s.mdService.CreateModelDeployment(ctx, md1)).NotTo(HaveOccurred())

	md2 := newStubMd()
	md2.ID = mdID2
	md2.Spec.RoleName = &mdRoleName2
	s.g.Expect(s.mdService.CreateModelDeployment(ctx, md2)).NotTo(HaveOccurred())

	return []*deployment.ModelDeployment{md1, md2}
}

func TestModelDeploymentRouteSuite(t *testing.T) {
	suite.Run(t, new(ModelDeploymentRouteSuite))
}

func (s *ModelDeploymentRouteSuite) TestGetMD() {
	ctx := context.Background()
	md := newStubMd()
	s.g.Expect(s.mdService.CreateModelDeployment(ctx, md)).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodGet,
		strings.Replace(dep_route.GetModelDeploymentURL, ":id", mdID, -1),
		nil,
	)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result deployment.ModelDeployment
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(result.ID).Should(Equal(md.ID))
	s.g.Expect(result.Spec).Should(Equal(md.Spec))

	s.g.Expect(result.Status.AvailableReplicas).Should(Equal(md.Status.AvailableReplicas))
	s.g.Expect(result.Status.Deployment).Should(Equal(md.Status.Deployment))
	s.g.Expect(result.Status.LastCredsUpdatedTime).Should(Equal(md.Status.LastCredsUpdatedTime))
	s.g.Expect(result.Status.LastRevisionName).Should(Equal(md.Status.LastRevisionName))
	s.g.Expect(result.Status.Replicas).Should(Equal(md.Status.Replicas))
	s.g.Expect(result.Status.Service).Should(Equal(md.Status.Service))
	s.g.Expect(result.Status.ServiceURL).Should(Equal(md.Status.ServiceURL))
	s.g.Expect(result.Status.State).Should(Equal(md.Status.State))
}

func (s *ModelDeploymentRouteSuite) TestGetMDNotFound() {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodGet,
		strings.Replace(dep_route.GetModelDeploymentURL, ":id", "not-found", -1),
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

func (s *ModelDeploymentRouteSuite) TestGetAllMDEmptyResult() {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, dep_route.GetAllModelDeploymentURL, nil)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var mdResponse []deployment.ModelDeployment
	err = json.Unmarshal(w.Body.Bytes(), &mdResponse)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(mdResponse).Should(HaveLen(0))
}

func (s *ModelDeploymentRouteSuite) TestGetAllMD() {
	s.newMultipleMds()

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, dep_route.GetAllModelDeploymentURL, nil)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result []deployment.ModelDeployment
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(result).Should(HaveLen(2))

	for _, md := range result {
		s.g.Expect(md.ID).To(Or(Equal(mdID1), Equal(mdID2)))
	}
}

func (s *ModelDeploymentRouteSuite) TestGetAllMdByRole() {
	mds := s.newMultipleMds()

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, dep_route.GetAllModelDeploymentURL, nil)
	s.g.Expect(err).NotTo(HaveOccurred())

	query := req.URL.Query()
	query.Set("roleName", mdRoleName2)
	req.URL.RawQuery = query.Encode()

	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result []deployment.ModelDeployment
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(result).Should(HaveLen(1))
	s.g.Expect(result[0].ID).Should(Equal(mdID2))
	s.g.Expect(result[0].Spec).Should(Equal(mds[1].Spec))
}

func (s *ModelDeploymentRouteSuite) TestGetAllMdMultipleFiltersByRole() {
	s.newMultipleMds()

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, dep_route.GetAllModelDeploymentURL, nil)
	s.g.Expect(err).NotTo(HaveOccurred())

	query := req.URL.Query()
	query.Set("roleName", mdRoleName1)
	query.Add("roleName", mdRoleName2)
	req.URL.RawQuery = query.Encode()

	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result []deployment.ModelDeployment
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(result).Should(HaveLen(2))
}

func (s *ModelDeploymentRouteSuite) TestGetAllMdPaging() {
	s.newMultipleMds()

	mdNames := map[string]interface{}{mdID1: nil, mdID2: nil}

	// Return first page
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, dep_route.GetAllModelDeploymentURL, nil)
	s.g.Expect(err).NotTo(HaveOccurred())

	query := req.URL.Query()
	query.Set("size", "1")
	query.Set("page", "0")
	req.URL.RawQuery = query.Encode()

	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result []deployment.ModelDeployment
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(result).Should(HaveLen(1))
	delete(mdNames, result[0].ID)

	// Return second page
	w = httptest.NewRecorder()
	req, err = http.NewRequest(http.MethodGet, dep_route.GetAllModelDeploymentURL, nil)
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
	delete(mdNames, result[0].ID)

	// Return third empty page
	w = httptest.NewRecorder()
	req, err = http.NewRequest(http.MethodGet, dep_route.GetAllModelDeploymentURL, nil)
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

func (s *ModelDeploymentRouteSuite) TestCreateMD() {
	ctx := context.Background()
	mdEntity := newStubMd()
	mdEntityBody, err := json.Marshal(mdEntity)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, dep_route.CreateModelDeploymentURL, bytes.NewReader(mdEntityBody))
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var mdResponse deployment.ModelDeployment
	err = json.Unmarshal(w.Body.Bytes(), &mdResponse)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusCreated))
	s.g.Expect(mdEntity.ID).To(Equal(mdResponse.ID))
	s.g.Expect(mdEntity.Spec).To(Equal(mdResponse.Spec))

	md, err := s.mdService.GetModelDeployment(ctx, mdID)
	s.g.Expect(err).ShouldNot(HaveOccurred())
	s.g.Expect(md.Spec).To(Equal(mdEntity.Spec))
}

func (s *ModelDeploymentRouteSuite) TestCreateMDValidation() {
	mdEntity := newStubMd()
	mdEntity.Spec.Image = ""
	mdEntityBody, err := json.Marshal(mdEntity)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, dep_route.CreateModelDeploymentURL, bytes.NewReader(mdEntityBody))
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusBadRequest))
	s.g.Expect(result.Message).To(ContainSubstring(dep_route.EmptyImageErrorMessage))
}

func (s *ModelDeploymentRouteSuite) TestCreateDuplicateMD() {
	ctx := context.Background()
	md := newStubMd()

	s.g.Expect(s.mdService.CreateModelDeployment(ctx, md)).NotTo(HaveOccurred())

	mdEntityBody, err := json.Marshal(md)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, dep_route.CreateModelDeploymentURL, bytes.NewReader(mdEntityBody))
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusConflict))
	s.g.Expect(result.Message).Should(ContainSubstring("already exists"))
}

func (s *ModelDeploymentRouteSuite) TestUpdateMD() {
	ctx := context.Background()
	md := newStubMd()
	s.g.Expect(s.mdService.CreateModelDeployment(ctx, md)).NotTo(HaveOccurred())

	newMaxReplicas := mdMaxReplicas + 1
	newMDLivenessIninitialDelay := mdLivenessInitialDelay + 1
	newMDReadinessInitialDelay := mdReadinessInitialDelay + 1

	mdEntity := newStubMd()
	mdEntity.Spec.MaxReplicas = &newMaxReplicas
	mdEntity.Spec.LivenessProbeInitialDelay = &newMDLivenessIninitialDelay
	mdEntity.Spec.ReadinessProbeInitialDelay = &newMDReadinessInitialDelay

	mdEntityBody, err := json.Marshal(mdEntity)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPut, dep_route.UpdateModelDeploymentURL, bytes.NewReader(mdEntityBody))
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var mdResponse deployment.ModelDeployment
	err = json.Unmarshal(w.Body.Bytes(), &mdResponse)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(mdEntity.ID).To(Equal(mdResponse.ID))
	s.g.Expect(mdEntity.Spec).To(Equal(mdResponse.Spec))

	md, err = s.mdService.GetModelDeployment(ctx, mdID)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.g.Expect(mdEntity.ID).To(Equal(md.ID))
	s.g.Expect(mdEntity.Spec).To(Equal(md.Spec))
}

func (s *ModelDeploymentRouteSuite) TestUpdateMDValidation() {
	ctx := context.Background()
	md := newStubMd()
	s.g.Expect(s.mdService.CreateModelDeployment(ctx, md)).NotTo(HaveOccurred())

	mdEntity := newStubMd()
	mdEntity.Spec.Image = ""

	mdEntityBody, err := json.Marshal(mdEntity)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPut, dep_route.UpdateModelDeploymentURL, bytes.NewReader(mdEntityBody))
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusBadRequest))
	s.g.Expect(result.Message).To(ContainSubstring(dep_route.EmptyImageErrorMessage))
}

func (s *ModelDeploymentRouteSuite) TestUpdateMDNotFound() {
	mdEntity := newStubMd()

	mdEntityBody, err := json.Marshal(mdEntity)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPut, dep_route.UpdateModelDeploymentURL, bytes.NewReader(mdEntityBody))
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusNotFound))
	s.g.Expect(result.Message).To(ContainSubstring("not found"))
}

func (s *ModelDeploymentRouteSuite) TestDeleteMD() {
	ctx := context.Background()
	md := newStubMd()
	s.g.Expect(s.mdService.CreateModelDeployment(ctx, md)).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodDelete,
		strings.Replace(dep_route.DeleteModelDeploymentURL, ":id", md.ID, -1),
		nil,
	)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(result.Message).Should(ContainSubstring("was deleted"))

	mdList, err := s.mdService.GetModelDeploymentList(ctx)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.g.Expect(mdList).To(HaveLen(1))
	fetchedMd, err := s.mdService.GetModelDeployment(ctx, mdID)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.g.Expect(fetchedMd.DeletionMark).Should(BeTrue())
}

func (s *ModelDeploymentRouteSuite) TestDeleteMDNotFound() {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodDelete,
		strings.Replace(dep_route.DeleteModelDeploymentURL, ":id", "not-found", -1),
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

func (s *ModelDeploymentRouteSuite) TestDisabledAPIGetMD() {
	ctx := context.Background()
	deploymentConfig := config.NewDefaultModelDeploymentConfig()
	deploymentConfig.Enabled = false
	s.registerHTTPHandlers(deploymentConfig)

	md := newStubMd()
	s.g.Expect(s.mdService.CreateModelDeployment(ctx, md)).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodGet,
		strings.Replace(dep_route.GetModelDeploymentURL, ":id", mdID, -1),
		nil,
	)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result deployment.ModelDeployment
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(result.ID).Should(Equal(md.ID))
	s.g.Expect(result.Spec).Should(Equal(md.Spec))

	s.g.Expect(result.Status.AvailableReplicas).Should(Equal(md.Status.AvailableReplicas))
	s.g.Expect(result.Status.Deployment).Should(Equal(md.Status.Deployment))
	s.g.Expect(result.Status.LastCredsUpdatedTime).Should(Equal(md.Status.LastCredsUpdatedTime))
	s.g.Expect(result.Status.LastRevisionName).Should(Equal(md.Status.LastRevisionName))
	s.g.Expect(result.Status.Replicas).Should(Equal(md.Status.Replicas))
	s.g.Expect(result.Status.Service).Should(Equal(md.Status.Service))
	s.g.Expect(result.Status.ServiceURL).Should(Equal(md.Status.ServiceURL))
	s.g.Expect(result.Status.State).Should(Equal(md.Status.State))
}

func (s *ModelDeploymentRouteSuite) TestDisabledAPIGetAllMD() {
	deploymentConfig := config.NewDefaultModelDeploymentConfig()
	deploymentConfig.Enabled = false
	s.registerHTTPHandlers(deploymentConfig)

	s.newMultipleMds()

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, dep_route.GetAllModelDeploymentURL, nil)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result []deployment.ModelDeployment
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(result).Should(HaveLen(2))

	for _, md := range result {
		s.g.Expect(md.ID).To(Or(Equal(mdID1), Equal(mdID2)))
	}
}

func (s *ModelDeploymentRouteSuite) TestDisabledAPICreateMD() {
	ctx := context.Background()
	deploymentConfig := config.NewDefaultModelDeploymentConfig()
	deploymentConfig.Enabled = false
	s.registerHTTPHandlers(deploymentConfig)
	md := newStubMd()

	s.g.Expect(s.mdService.CreateModelDeployment(ctx, md)).NotTo(HaveOccurred())

	mdEntityBody, err := json.Marshal(md)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, dep_route.CreateModelDeploymentURL, bytes.NewReader(mdEntityBody))
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusBadRequest))
	s.g.Expect(result.Message).Should(ContainSubstring(routes.DisabledAPIErrorMessage))
}

func (s *ModelDeploymentRouteSuite) TestDisabledAPIUpdateMD() {
	deploymentConfig := config.NewDefaultModelDeploymentConfig()
	deploymentConfig.Enabled = false
	s.registerHTTPHandlers(deploymentConfig)
	mdEntity := newStubMd()

	mdEntityBody, err := json.Marshal(mdEntity)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPut, dep_route.UpdateModelDeploymentURL, bytes.NewReader(mdEntityBody))
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusBadRequest))
	s.g.Expect(result.Message).Should(ContainSubstring(routes.DisabledAPIErrorMessage))
}

func (s *ModelDeploymentRouteSuite) TestDisabledAPIDeleteMD() {
	deploymentConfig := config.NewDefaultModelDeploymentConfig()
	deploymentConfig.Enabled = false
	s.registerHTTPHandlers(deploymentConfig)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodDelete,
		strings.Replace(dep_route.DeleteModelDeploymentURL, ":id", "12345", -1),
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

func (s *ModelDeploymentRouteSuite) TestGetDefaultRoute() {
	s.newMultipleMds()

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodGet,
		strings.Replace(dep_route.GetModelDeploymentDefaultRouteURL, ":id", mdID1, -1),
		nil,
	)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result deployment.ModelRoute
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(result.Default).Should(BeTrue())
}

func (s *ModelDeploymentRouteSuite) TestGetDeploymentEventsIncorrectCursor() {

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodGet,
		dep_route.EventsModelDeploymentURL+"?cursor=not-number",
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

func (s *ModelDeploymentRouteSuite) TestGetDeploymentEvents() {

	events := []event.DeploymentEvent{
		{
			Payload:   deployment.ModelDeployment{ID: "deployment-1"},
			EventType: event.ModelDeploymentCreatedEventType,
			Datetime:  time.Time{},
		},
		{
			EntityID:  "deployment-2",
			EventType: event.ModelDeploymentDeletedEventType,
			Datetime:  time.Time{},
		},
	}

	s.mdEventsGetter.On("Get", mock.Anything, 0).Return(events, 5, nil)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodGet,
		dep_route.EventsModelDeploymentURL,
		nil,
	)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)
	s.g.Expect(w.Code).Should(Equal(200))

	var result event.LatestDeploymentEvents
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(result.Cursor).Should(Equal(5))
	s.g.Expect(result.Events).Should(Equal(events))
	s.g.Expect(err).NotTo(HaveOccurred())
}

func (s *ModelDeploymentRouteSuite) TestGetDeploymentEventsWithCursor() {

	events := []event.DeploymentEvent{
		{
			Payload:   deployment.ModelDeployment{ID: "deployment-1"},
			EventType: event.ModelRouteCreatedEventType,
			Datetime:  time.Time{},
		},
		{
			EntityID:  "deployment-2",
			EventType: event.ModelRouteDeletedEventType,
			Datetime:  time.Time{},
		},
	}

	s.mdEventsGetter.On("Get", mock.Anything, 6).Return(events, 9, nil)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodGet,
		dep_route.EventsModelDeploymentURL+"?cursor=6",
		nil,
	)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)
	s.g.Expect(w.Code).Should(Equal(200))

	var result event.LatestDeploymentEvents
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(result.Cursor).Should(Equal(9))
	s.g.Expect(result.Events).Should(Equal(events))
	s.g.Expect(err).NotTo(HaveOccurred())
}
