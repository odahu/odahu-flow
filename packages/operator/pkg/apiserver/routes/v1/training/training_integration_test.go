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
	"github.com/gin-gonic/gin"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/training"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/routes"
	train_route "github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/routes/v1/training"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	odahuErrors "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	"github.com/odahu/odahu-flow/packages/operator/pkg/service/training_integration"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"
	httputil "github.com/odahu/odahu-flow/packages/operator/pkg/utils/httputil"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
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
	g                              *GomegaWithT
	server                         *gin.Engine
	trainingIntegrationServiceMock training_integration.MockTrainingIntegrationService
}

func TestTIGenericRouteSuite(t *testing.T) {
	suite.Run(t, new(TIGenericRouteSuite))
}

func (s *TIGenericRouteSuite) SetupTest() {
	s.g = NewGomegaWithT(s.T())
	s.trainingIntegrationServiceMock = training_integration.MockTrainingIntegrationService{}
	s.registerHandlers(config.NewDefaultModelTrainingConfig())
}

func (s *TIGenericRouteSuite) registerHandlers(trainingConfig config.ModelTrainingConfig) {
	s.server = gin.Default()
	v1Group := s.server.Group("")
	trainingGroup := v1Group.Group("", routes.DisableAPIMiddleware(trainingConfig.Enabled))
	train_route.ConfigureTrainingIntegrationRoutes(trainingGroup, &s.trainingIntegrationServiceMock)
}

func newTiStub() *training.TrainingIntegration {
	return &training.TrainingIntegration{
		ID: testTrainingIntegrationID,
		Spec: v1alpha1.TrainingIntegrationSpec{
			Entrypoint:             tiEntrypoint,
			DefaultImage:           tiDefaultImage,
			AdditionalEnvironments: tiAdditionalEnvironments,
		},
	}
}

func (s *TIGenericRouteSuite) TestGetTrainingIntegration() {
	ti := newTiStub()

	s.trainingIntegrationServiceMock.On("GetTrainingIntegration", ti.ID).Return(ti, nil)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, strings.Replace(
		train_route.GetTrainingIntegrationURL, ":id", ti.ID, -1,
	), nil)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var tiResponse training.TrainingIntegration
	err = json.Unmarshal(w.Body.Bytes(), &tiResponse)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(tiResponse.ID).Should(Equal(ti.ID))
	s.g.Expect(tiResponse.Spec).Should(Equal(ti.Spec))
}

func (s *TIGenericRouteSuite) TestGetTrainingIntegrationNotFound() {
	trainingIntegrationID := "not-found"

	s.trainingIntegrationServiceMock.
		On("GetTrainingIntegration", trainingIntegrationID).
		Return(nil, odahuErrors.NotFoundError{})

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, strings.Replace(
		train_route.GetTrainingIntegrationURL, ":id", trainingIntegrationID, -1,
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
	s.trainingIntegrationServiceMock.
		On("GetTrainingIntegrationList", mock.Anything, mock.Anything).
		Return([]training.TrainingIntegration{}, nil)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, train_route.GetAllTrainingIntegrationURL, nil)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var tiResponse []training.TrainingIntegration
	err = json.Unmarshal(w.Body.Bytes(), &tiResponse)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(tiResponse).Should(HaveLen(0))
}

func (s *TIGenericRouteSuite) TestGetAllTi() {
	s.trainingIntegrationServiceMock.
		On("GetTrainingIntegrationList", mock.Anything, mock.Anything).
		Return([]training.TrainingIntegration{
			{ID: testTrainingIntegrationID1},
			{ID: testTrainingIntegrationID2},
		}, nil)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, train_route.GetAllTrainingIntegrationURL, nil)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result []training.TrainingIntegration
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(result).Should(HaveLen(2))

	for _, ti := range result {
		s.g.Expect(ti.ID).To(Or(Equal(testTrainingIntegrationID1), Equal(testTrainingIntegrationID2)))
	}
}

func (s *TIGenericRouteSuite) TestGetAllTiPaging() {
	expectedTrainingIntegrations := []training.TrainingIntegration{
		{ID: testTrainingIntegrationID1},
		{ID: testTrainingIntegrationID2},
	}
	s.trainingIntegrationServiceMock.
		On("GetTrainingIntegrationList", mock.Anything, mock.Anything).
		Return(expectedTrainingIntegrations, nil)

	// Return first page
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, train_route.GetAllTrainingIntegrationURL, nil)
	s.g.Expect(err).NotTo(HaveOccurred())

	query := req.URL.Query()

	expectedSize := 2
	expectedPage := 0
	query.Set("size", strconv.Itoa(expectedSize))
	query.Set("page", strconv.Itoa(expectedPage))
	req.URL.RawQuery = query.Encode()

	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	// Check that query params properly converted to list options
	sizeFunc := s.trainingIntegrationServiceMock.Calls[0].Arguments.Get(0).(filter.ListOption)
	pageFunc := s.trainingIntegrationServiceMock.Calls[0].Arguments.Get(1).(filter.ListOption)

	listOptions := &filter.ListOptions{}
	sizeFunc(listOptions)
	actualSize := *listOptions.Size
	pageFunc(listOptions)
	actualPage := *listOptions.Page

	s.g.Expect(actualPage).To(Equal(expectedPage))
	s.g.Expect(actualSize).To(Equal(expectedSize))

	var actualTrainingIntegrations []training.TrainingIntegration
	err = json.Unmarshal(w.Body.Bytes(), &actualTrainingIntegrations)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(actualTrainingIntegrations).Should(HaveLen(2))
	s.g.Expect(actualTrainingIntegrations).To(Equal(expectedTrainingIntegrations))
}

func (s *TIGenericRouteSuite) TestCreateTrainingIntegration() {
	expectedTI := newTiStub()

	s.trainingIntegrationServiceMock.
		On("CreateTrainingIntegration", expectedTI).
		Return(nil)

	tiEntityBody, err := json.Marshal(expectedTI)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodPost, train_route.CreateTrainingIntegrationURL, bytes.NewReader(tiEntityBody),
	)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var tiResponse training.TrainingIntegration
	err = json.Unmarshal(w.Body.Bytes(), &tiResponse)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusCreated))
	s.g.Expect(tiResponse.ID).Should(Equal(expectedTI.ID))
	s.g.Expect(tiResponse.Spec).Should(Equal(expectedTI.Spec))
}

func (s *TIGenericRouteSuite) TestCreateTrainingIntegrationValidation() {
	tiEntity := newTiStub()
	tiEntity.Spec.Entrypoint = ""

	tiEntityBody, err := json.Marshal(tiEntity)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodPost, train_route.CreateTrainingIntegrationURL, bytes.NewReader(tiEntityBody),
	)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var resultResponse httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &resultResponse)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusBadRequest))
	s.g.Expect(resultResponse.Message).Should(ContainSubstring(train_route.ValidationTiErrorMessage))
}

func (s *TIGenericRouteSuite) TestCreateDuplicateTrainingIntegration() {
	ti := newTiStub()

	s.trainingIntegrationServiceMock.
		On("CreateTrainingIntegration", ti).
		Return(odahuErrors.AlreadyExistError{})

	tiEntityBody, err := json.Marshal(ti)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodPost, train_route.CreateTrainingIntegrationURL, bytes.NewReader(tiEntityBody),
	)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusConflict))
	s.g.Expect(result.Message).Should(ContainSubstring("already exists"))
}

func (s *TIGenericRouteSuite) TestUpdateTrainingIntegration() {
	updatedTi := newTiStub()
	s.trainingIntegrationServiceMock.
		On("UpdateTrainingIntegration", updatedTi).
		Return(nil)

	tiEntityBody, err := json.Marshal(updatedTi)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodPut, train_route.UpdateTrainingIntegrationURL, bytes.NewReader(tiEntityBody),
	)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var tiResponse training.TrainingIntegration
	err = json.Unmarshal(w.Body.Bytes(), &tiResponse)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(tiResponse.ID).Should(Equal(updatedTi.ID))
	s.g.Expect(tiResponse.Spec).Should(Equal(updatedTi.Spec))
}

func (s *TIGenericRouteSuite) TestUpdateTrainingIntegrationValidation() {
	updatedTi := newTiStub()
	updatedTi.Spec.Entrypoint = ""

	tiEntityBody, err := json.Marshal(updatedTi)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodPut, train_route.UpdateTrainingIntegrationURL, bytes.NewReader(tiEntityBody),
	)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var resultResponse httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &resultResponse)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusBadRequest))
	s.g.Expect(resultResponse.Message).Should(ContainSubstring(train_route.ValidationTiErrorMessage))
}

func (s *TIGenericRouteSuite) TestUpdateTrainingIntegrationNotFound() {
	ti := newTiStub()

	s.trainingIntegrationServiceMock.
		On("UpdateTrainingIntegration", ti).
		Return(odahuErrors.NotFoundError{})

	tiBody, err := json.Marshal(ti)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPut, train_route.UpdateTrainingIntegrationURL, bytes.NewReader(tiBody))
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var response httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &response)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusNotFound))
	s.g.Expect(response.Message).Should(ContainSubstring("not found"))
}

func (s *TIGenericRouteSuite) TestDeleteTrainingIntegration() {
	id := "id"
	s.trainingIntegrationServiceMock.
		On("DeleteTrainingIntegration", id).
		Return(nil)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodDelete, strings.Replace(
		train_route.DeleteTrainingIntegrationURL, ":id", id, -1,
	), nil)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(result.Message).Should(ContainSubstring("was deleted"))
}

func (s *TIGenericRouteSuite) TestDeleteTrainingIntegrationNotFound() {
	notFoundID := "not-found"
	s.trainingIntegrationServiceMock.
		On("DeleteTrainingIntegration", notFoundID).
		Return(odahuErrors.NotFoundError{})

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodDelete, strings.Replace(
		train_route.DeleteTrainingIntegrationURL, ":id", notFoundID, -1,
	), nil)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusNotFound))
	s.g.Expect(result.Message).Should(ContainSubstring("not found"))
}

func (s *TIGenericRouteSuite) TestDisabledAPIDeleteTrainingIntegration() {
	trainingConfig := config.NewDefaultModelTrainingConfig()
	trainingConfig.Enabled = false
	s.registerHandlers(trainingConfig)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodDelete, strings.Replace(
		train_route.DeleteTrainingIntegrationURL, ":id", "not-found", -1,
	), nil)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusBadRequest))
	s.g.Expect(result.Message).Should(ContainSubstring(routes.DisabledAPIErrorMessage))
}

func (s *TIGenericRouteSuite) TestDisabledAPIUpdateTrainingIntegration() {
	trainingConfig := config.NewDefaultModelTrainingConfig()
	trainingConfig.Enabled = false
	s.registerHandlers(trainingConfig)

	ti := newTiStub()

	tiBody, err := json.Marshal(ti)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPut, train_route.UpdateTrainingIntegrationURL, bytes.NewReader(tiBody))
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusBadRequest))
	s.g.Expect(result.Message).Should(ContainSubstring(routes.DisabledAPIErrorMessage))
}

func (s *TIGenericRouteSuite) TestDisabledAPICreateTrainingIntegrationValidation() {
	trainingConfig := config.NewDefaultModelTrainingConfig()
	trainingConfig.Enabled = false
	s.registerHandlers(trainingConfig)

	tiEntity := newTiStub()
	tiEntity.Spec.Entrypoint = ""

	tiEntityBody, err := json.Marshal(tiEntity)
	s.g.Expect(err).NotTo(HaveOccurred())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodPost, train_route.CreateTrainingIntegrationURL, bytes.NewReader(tiEntityBody),
	)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)
	var result httputil.HTTPResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusBadRequest))
	s.g.Expect(result.Message).Should(ContainSubstring(routes.DisabledAPIErrorMessage))
}

// GET requests are served despite API is disabled
func (s *TIGenericRouteSuite) TestDisabledAPIGetAllTi() {
	trainingConfig := config.NewDefaultModelTrainingConfig()
	trainingConfig.Enabled = false
	s.registerHandlers(trainingConfig)

	expectedTrainingIntegrations := []training.TrainingIntegration{
		{ID: testTrainingIntegrationID1},
		{ID: testTrainingIntegrationID2},
	}
	s.trainingIntegrationServiceMock.
		On("GetTrainingIntegrationList", mock.Anything, mock.Anything).
		Return(expectedTrainingIntegrations, nil)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, train_route.GetAllTrainingIntegrationURL, nil)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result []training.TrainingIntegration
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(result).Should(HaveLen(2))

	for _, ti := range result {
		s.g.Expect(ti.ID).To(Or(Equal(testTrainingIntegrationID1), Equal(testTrainingIntegrationID2)))
	}
}

func (s *TIGenericRouteSuite) TestDisabledAPIGetTrainingIntegration() {
	trainingConfig := config.NewDefaultModelTrainingConfig()
	trainingConfig.Enabled = false
	s.registerHandlers(trainingConfig)

	ti := newTiStub()
	s.trainingIntegrationServiceMock.
		On("GetTrainingIntegration", ti.ID).
		Return(ti, nil)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, strings.Replace(
		train_route.GetTrainingIntegrationURL, ":id", ti.ID, -1,
	), nil)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var tiResponse training.TrainingIntegration
	err = json.Unmarshal(w.Body.Bytes(), &tiResponse)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(tiResponse.ID).Should(Equal(ti.ID))
	s.g.Expect(tiResponse.Spec).Should(Equal(ti.Spec))
}
