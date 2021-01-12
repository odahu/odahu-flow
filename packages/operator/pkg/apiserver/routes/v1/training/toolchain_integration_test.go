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

package training_test

import (
	"bytes"
	"encoding/json"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	httputil "github.com/odahu/odahu-flow/packages/operator/pkg/utils/httputil"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/training"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/routes"
	train_route "github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/routes/v1/training"
	odahuErrors "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	mt_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/training"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"
)

const (
	tiEntrypoint   = "test-entrypoint"
	tiDefaultImage = "test:image"
)

var (
	tiAdditionalEnvironments = map[string]string{
		"name-123": "value-456",
	}
)

type TIGenericRouteSuite struct {
	suite.Suite
	g            *GomegaWithT
	server       *gin.Engine
	mtRepository mt_repository.ToolchainRepository
}

func (s *TIGenericRouteSuite) TearDownTest() {
	for _, mpID := range []string{
		testToolchainIntegrationID,
		testToolchainIntegrationID1,
		testToolchainIntegrationID2,
	} {
		if err := s.mtRepository.DeleteToolchainIntegration(mpID); err != nil && !odahuErrors.IsNotFoundError(err) {
			// If a model training is not found then it was not created during a test case
			// All other errors propagate as a panic
			panic(err)
		}
	}
}

func (s *TIGenericRouteSuite) SetupTest() {
	s.g = NewGomegaWithT(s.T())

	s.registerHandlers(config.NewDefaultModelTrainingConfig())
}

func (s *TIGenericRouteSuite) registerHandlers(trainingConfig config.ModelTrainingConfig) {
	s.server = gin.Default()
	v1Group := s.server.Group("")
	trainingGroup := v1Group.Group("", routes.DisableAPIMiddleware(trainingConfig.Enabled))
	train_route.ConfigureToolchainRoutes(trainingGroup, s.mtRepository)
}

func newTiStub() *training.ToolchainIntegration {
	return &training.ToolchainIntegration{
		ID: testToolchainIntegrationID,
		Spec: v1alpha1.ToolchainIntegrationSpec{
			Entrypoint:             tiEntrypoint,
			DefaultImage:           tiDefaultImage,
			AdditionalEnvironments: tiAdditionalEnvironments,
		},
	}
}

func (s *TIGenericRouteSuite) newMultipleTiStubs() {
	ti1 := &training.ToolchainIntegration{
		ID: testToolchainIntegrationID1,
		Spec: v1alpha1.ToolchainIntegrationSpec{
			DefaultImage: testToolchainMtImage,
		},
	}
	s.g.Expect(s.mtRepository.CreateToolchainIntegration(ti1)).NotTo(HaveOccurred())

	ti2 := &training.ToolchainIntegration{
		ID: testToolchainIntegrationID2,
		Spec: v1alpha1.ToolchainIntegrationSpec{
			DefaultImage: testToolchainMtImage,
		},
	}
	s.g.Expect(s.mtRepository.CreateToolchainIntegration(ti2)).NotTo(HaveOccurred())
}

func (s *TIGenericRouteSuite) TestGetToolchainIntegration() {
	ti := newTiStub()
	s.g.Expect(s.mtRepository.CreateToolchainIntegration(ti)).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, strings.Replace(
		train_route.GetToolchainIntegrationURL, ":id", ti.ID, -1,
	), nil)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var tiResponse training.ToolchainIntegration
	err = json.Unmarshal(w.Body.Bytes(), &tiResponse)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(tiResponse.ID).Should(Equal(ti.ID))
	s.g.Expect(tiResponse.Spec).Should(Equal(ti.Spec))
}

func (s *TIGenericRouteSuite) TestGetToolchainIntegrationNotFound() {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, strings.Replace(
		train_route.GetToolchainIntegrationURL, ":id", "not-found", -1,
	), nil)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusNotFound))
	s.g.Expect(result.Message).Should(ContainSubstring("not found"))
}

func (s *TIGenericRouteSuite) TestGetAllTiEmptyResult() {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, train_route.GetAllToolchainIntegrationURL, nil)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var tiResponse []training.ToolchainIntegration
	err = json.Unmarshal(w.Body.Bytes(), &tiResponse)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(tiResponse).Should(HaveLen(0))
}

func (s *TIGenericRouteSuite) TestGetAllTi() {
	s.newMultipleTiStubs()

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, train_route.GetAllToolchainIntegrationURL, nil)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result []training.ToolchainIntegration
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(result).Should(HaveLen(2))

	for _, ti := range result {
		s.g.Expect(ti.ID).To(Or(Equal(testToolchainIntegrationID1), Equal(testToolchainIntegrationID2)))
	}
}

func (s *TIGenericRouteSuite) TestGetAllTiPaging() {
	s.newMultipleTiStubs()

	toolchainsNames := map[string]interface{}{testToolchainIntegrationID1: nil, testToolchainIntegrationID2: nil}

	// Return first page
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, train_route.GetAllToolchainIntegrationURL, nil)
	s.g.Expect(err).NotTo(HaveOccurred())

	query := req.URL.Query()
	query.Set("size", "1")
	query.Set("page", "0")
	req.URL.RawQuery = query.Encode()

	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var toolchains []training.ToolchainIntegration
	err = json.Unmarshal(w.Body.Bytes(), &toolchains)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(toolchains).Should(HaveLen(1))
	delete(toolchainsNames, toolchains[0].ID)

	// Return second page
	w = httptest.NewRecorder()
	req, err = http.NewRequest(http.MethodGet, train_route.GetAllToolchainIntegrationURL, nil)
	s.g.Expect(err).NotTo(HaveOccurred())

	query = req.URL.Query()
	query.Set("size", "1")
	query.Set("page", "1")
	req.URL.RawQuery = query.Encode()

	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	err = json.Unmarshal(w.Body.Bytes(), &toolchains)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(toolchains).Should(HaveLen(1))
	delete(toolchainsNames, toolchains[0].ID)

	// Return third empty page
	w = httptest.NewRecorder()
	req, err = http.NewRequest(http.MethodGet, train_route.GetAllToolchainIntegrationURL, nil)
	s.g.Expect(err).NotTo(HaveOccurred())

	query = req.URL.Query()
	query.Set("size", "1")
	query.Set("page", "2")
	req.URL.RawQuery = query.Encode()

	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	err = json.Unmarshal(w.Body.Bytes(), &toolchains)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(toolchains).Should(HaveLen(0))
	s.g.Expect(toolchains).Should(BeEmpty())
}

func (s *TIGenericRouteSuite) TestCreateToolchainIntegration() {
	tiEntity := newTiStub()

	tiEntityBody, err := json.Marshal(tiEntity)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodPost, train_route.CreateToolchainIntegrationURL, bytes.NewReader(tiEntityBody),
	)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var tiResponse training.ToolchainIntegration
	err = json.Unmarshal(w.Body.Bytes(), &tiResponse)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusCreated))
	s.g.Expect(tiResponse.ID).Should(Equal(tiEntity.ID))
	s.g.Expect(tiResponse.Spec).Should(Equal(tiEntity.Spec))

	ti, err := s.mtRepository.GetToolchainIntegration(testToolchainIntegrationID)
	s.g.Expect(err).ShouldNot(HaveOccurred())
	s.g.Expect(ti.Spec).To(Equal(tiEntity.Spec))
}

// CreatedAt and UpdatedAt field should automatically be updated after create request
func (s *TIGenericRouteSuite) TestCreateToolchainIntegrationModifiable() {
	newResource := newTiStub()

	newResourceBody, err := json.Marshal(newResource)
	s.g.Expect(err).NotTo(HaveOccurred())

	reqTime := routes.GetTimeNowTruncatedToSeconds()
	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodPost, train_route.CreateToolchainIntegrationURL, bytes.NewReader(newResourceBody),
	)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var resp training.ToolchainIntegration
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

func (s *TIGenericRouteSuite) TestCreateToolchainIntegrationValidation() {
	tiEntity := newTiStub()
	tiEntity.Spec.Entrypoint = ""

	tiEntityBody, err := json.Marshal(tiEntity)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodPost, train_route.CreateToolchainIntegrationURL, bytes.NewReader(tiEntityBody),
	)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var resultResponse httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &resultResponse)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusBadRequest))
	s.g.Expect(resultResponse.Message).Should(ContainSubstring(train_route.ValidationTiErrorMessage))
}

func (s *TIGenericRouteSuite) TestCreateDuplicateToolchainIntegration() {
	ti := newTiStub()

	s.g.Expect(s.mtRepository.CreateToolchainIntegration(ti)).NotTo(HaveOccurred())

	tiEntityBody, err := json.Marshal(ti)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodPost, train_route.CreateToolchainIntegrationURL, bytes.NewReader(tiEntityBody),
	)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusConflict))
	s.g.Expect(result.Message).Should(ContainSubstring("already exists"))
}

func (s *TIGenericRouteSuite) TestUpdateToolchainIntegration() {
	ti := newTiStub()
	s.g.Expect(s.mtRepository.CreateToolchainIntegration(ti)).NotTo(HaveOccurred())

	updatedTi := newTiStub()
	updatedTi.Spec.Entrypoint = "new-entrypoint"

	tiEntityBody, err := json.Marshal(updatedTi)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodPut, train_route.UpdateToolchainIntegrationURL, bytes.NewReader(tiEntityBody),
	)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var tiResponse training.ToolchainIntegration
	err = json.Unmarshal(w.Body.Bytes(), &tiResponse)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(tiResponse.ID).Should(Equal(updatedTi.ID))
	s.g.Expect(tiResponse.Spec).Should(Equal(updatedTi.Spec))

	ti, err = s.mtRepository.GetToolchainIntegration(testToolchainIntegrationID)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.g.Expect(ti.Spec).To(Equal(updatedTi.Spec))
}

// UpdatedAt field should automatically be updated after update request
func (s *TIGenericRouteSuite) TestUpdateToolchainIntegrationModifiable() {
	resource := newTiStub()
	s.g.Expect(s.mtRepository.CreateToolchainIntegration(resource)).NotTo(HaveOccurred())

	time.Sleep(1 * time.Second)

	newResource := newTiStub()

	newResourceBody, err := json.Marshal(newResource)
	s.g.Expect(err).NotTo(HaveOccurred())

	reqTime := routes.GetTimeNowTruncatedToSeconds()
	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodPut, train_route.UpdateToolchainIntegrationURL, bytes.NewReader(newResourceBody),
	)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var respResource training.ToolchainIntegration
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

func (s *TIGenericRouteSuite) TestUpdateToolchainIntegrationValidation() {
	ti := newTiStub()
	s.g.Expect(s.mtRepository.CreateToolchainIntegration(ti)).NotTo(HaveOccurred())

	updatedTi := newTiStub()
	updatedTi.Spec.Entrypoint = ""

	tiEntityBody, err := json.Marshal(updatedTi)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodPut, train_route.UpdateToolchainIntegrationURL, bytes.NewReader(tiEntityBody),
	)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var resultResponse httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &resultResponse)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusBadRequest))
	s.g.Expect(resultResponse.Message).Should(ContainSubstring(train_route.ValidationTiErrorMessage))
}

func (s *TIGenericRouteSuite) TestUpdateToolchainIntegrationNotFound() {
	ti := newTiStub()

	tiBody, err := json.Marshal(ti)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPut, train_route.UpdateToolchainIntegrationURL, bytes.NewReader(tiBody))
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var response httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &response)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusNotFound))
	s.g.Expect(response.Message).Should(ContainSubstring("not found"))
}

func (s *TIGenericRouteSuite) TestDeleteToolchainIntegration() {
	ti := newTiStub()
	s.g.Expect(s.mtRepository.CreateToolchainIntegration(ti)).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodDelete, strings.Replace(
		train_route.DeleteToolchainIntegrationURL, ":id", ti.ID, -1,
	), nil)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(result.Message).Should(ContainSubstring("was deleted"))

	tiList, err := s.mtRepository.GetToolchainIntegrationList()
	s.g.Expect(err).NotTo(HaveOccurred())
	s.g.Expect(tiList).To(HaveLen(0))
}

func (s *TIGenericRouteSuite) TestDeleteToolchainIntegrationNotFound() {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodDelete, strings.Replace(
		train_route.DeleteToolchainIntegrationURL, ":id", "not-found", -1,
	), nil)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusNotFound))
	s.g.Expect(result.Message).Should(ContainSubstring("not found"))
}

func (s *TIGenericRouteSuite) TestDisabledAPIDeleteToolchainIntegration() {
	trainingConfig := config.NewDefaultModelTrainingConfig()
	trainingConfig.Enabled = false
	s.registerHandlers(trainingConfig)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodDelete, strings.Replace(
		train_route.DeleteToolchainIntegrationURL, ":id", "not-found", -1,
	), nil)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusBadRequest))
	s.g.Expect(result.Message).Should(ContainSubstring(routes.DisabledAPIErrorMessage))
}

func (s *TIGenericRouteSuite) TestDisabledAPIUpdateToolchainIntegration() {
	trainingConfig := config.NewDefaultModelTrainingConfig()
	trainingConfig.Enabled = false
	s.registerHandlers(trainingConfig)

	ti := newTiStub()

	tiBody, err := json.Marshal(ti)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPut, train_route.UpdateToolchainIntegrationURL, bytes.NewReader(tiBody))
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusBadRequest))
	s.g.Expect(result.Message).Should(ContainSubstring(routes.DisabledAPIErrorMessage))
}

func (s *TIGenericRouteSuite) TestDisabledAPICreateToolchainIntegrationValidation() {
	trainingConfig := config.NewDefaultModelTrainingConfig()
	trainingConfig.Enabled = false
	s.registerHandlers(trainingConfig)

	tiEntity := newTiStub()
	tiEntity.Spec.Entrypoint = ""

	tiEntityBody, err := json.Marshal(tiEntity)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodPost, train_route.CreateToolchainIntegrationURL, bytes.NewReader(tiEntityBody),
	)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)
	var result httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusBadRequest))
	s.g.Expect(result.Message).Should(ContainSubstring(routes.DisabledAPIErrorMessage))
}

func (s *TIGenericRouteSuite) TestDisabledAPIGetAllTi() {
	trainingConfig := config.NewDefaultModelTrainingConfig()
	trainingConfig.Enabled = false
	s.registerHandlers(trainingConfig)

	s.newMultipleTiStubs()

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, train_route.GetAllToolchainIntegrationURL, nil)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result []training.ToolchainIntegration
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(result).Should(HaveLen(2))

	for _, ti := range result {
		s.g.Expect(ti.ID).To(Or(Equal(testToolchainIntegrationID1), Equal(testToolchainIntegrationID2)))
	}
}

func (s *TIGenericRouteSuite) TestDisabledAPIGetToolchainIntegration() {
	trainingConfig := config.NewDefaultModelTrainingConfig()
	trainingConfig.Enabled = false
	s.registerHandlers(trainingConfig)

	ti := newTiStub()
	s.g.Expect(s.mtRepository.CreateToolchainIntegration(ti)).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, strings.Replace(
		train_route.GetToolchainIntegrationURL, ":id", ti.ID, -1,
	), nil)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var tiResponse training.ToolchainIntegration
	err = json.Unmarshal(w.Body.Bytes(), &tiResponse)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(tiResponse.ID).Should(Equal(ti.ID))
	s.g.Expect(tiResponse.Spec).Should(Equal(ti.Spec))
}
