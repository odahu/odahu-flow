/*
 * Copyright 2021 EPAM Systems
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package training_integration_test

import (
	"errors"
	training3 "github.com/odahu/odahu-flow/packages/operator/pkg/apis/training"
	training_integration_mock "github.com/odahu/odahu-flow/packages/operator/pkg/repository/training/mocks"
	"github.com/odahu/odahu-flow/packages/operator/pkg/service/training_integration"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

func TestSuiteRun(t *testing.T) {
	suite.Run(t, new(TITestSuite))
}

type TITestSuite struct {
	suite.Suite
	tiService  *training_integration.Service
	tiRepoMock *training_integration_mock.TrainingIntegrationRepository
}

func (s *TITestSuite) SetupTest() {
	s.tiRepoMock = &training_integration_mock.TrainingIntegrationRepository{}
	s.tiService = training_integration.NewService(s.tiRepoMock)
}

// Tests that GetTrainingIntegration method proxies to repo
func (s *TITestSuite) TestGetTrainingIntegration() {
	id := "some-id"

	expectedTI := &training3.TrainingIntegration{}
	expectedErr := errors.New("some-error")

	s.tiRepoMock.On("GetTrainingIntegration", id).Return(expectedTI, expectedErr)

	pi, err := s.tiService.GetTrainingIntegration(id)

	s.Assert().Equal(pi, expectedTI)
	s.Assert().Equal(err, expectedErr)
}

// Tests that GetTrainingIntegrationList method proxies to repo
func (s *TITestSuite) TestGetTrainingIntegrationList() {
	listOption := filter.Page(1)

	expectedTIList := []training3.TrainingIntegration{}
	expectedErr := errors.New("some-error")

	s.tiRepoMock.
		On("GetTrainingIntegrationList", mock.AnythingOfType("filter.ListOption")).
		Return(expectedTIList, expectedErr)

	piList, err := s.tiService.GetTrainingIntegrationList(listOption)

	s.Assert().Equal(piList, expectedTIList)
	s.Assert().Equal(err, expectedErr)
}

// Tests that CreateTrainingIntegration method sets CreatedAt/UpdatedAt and proxies to repo
func (s *TITestSuite) TestCreateTrainingIntegration() {
	someTI := &training3.TrainingIntegration{ID: "some-ti"}

	expectedErr := errors.New("some-error")

	s.tiRepoMock.
		On("SaveTrainingIntegration", someTI).
		Return(expectedErr)

	timeBefore := time.Now()
	err := s.tiService.CreateTrainingIntegration(someTI)

	actualTI := s.tiRepoMock.Calls[0].Arguments.Get(0).(*training3.TrainingIntegration)

	s.Assert().Equal(someTI, actualTI)
	s.Assert().Equal(err, expectedErr)

	s.Assert().True(someTI.CreatedAt.After(timeBefore))
	s.Assert().Equal(someTI.CreatedAt, someTI.UpdatedAt)
}

// Tests that CreateTrainingIntegration method sets CreatedAt/UpdatedAt and proxies to repo
func (s *TITestSuite) TestUpdateTrainingIntegration() {
	someTI := &training3.TrainingIntegration{ID: "some-ti", CreatedAt: time.Now(), UpdatedAt: time.Now()}

	expectedErr := errors.New("some-error")

	s.tiRepoMock.
		On("UpdateTrainingIntegration", someTI).
		Return(expectedErr)

	s.tiRepoMock.
		On("GetTrainingIntegration", someTI.ID).
		Return(someTI, nil)

	timeBeforeUpdate := time.Now()
	err := s.tiService.UpdateTrainingIntegration(someTI)

	actualPI := s.tiRepoMock.Calls[1].Arguments.Get(0).(*training3.TrainingIntegration)

	s.Assert().Equal(someTI, actualPI)
	s.Assert().Equal(err, expectedErr)

	s.Assert().True(someTI.UpdatedAt.After(timeBeforeUpdate))
	s.Assert().True(someTI.CreatedAt.Before(timeBeforeUpdate))
}

// Tests that DeleteTrainingIntegration method proxies to repo
func (s *TITestSuite) TestDeleteTrainingIntegration() {
	id := "some-id"
	expectedErr := errors.New("some-error")

	s.tiRepoMock.On("DeleteTrainingIntegration", id).Return(expectedErr)

	err := s.tiService.DeleteTrainingIntegration(id)

	s.Assert().Equal(err, expectedErr)
}
