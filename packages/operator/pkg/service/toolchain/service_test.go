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

package toolchain_test

import (
	"errors"
	training3 "github.com/odahu/odahu-flow/packages/operator/pkg/apis/training"
	toolchain_mock "github.com/odahu/odahu-flow/packages/operator/pkg/repository/training/mocks"
	"github.com/odahu/odahu-flow/packages/operator/pkg/service/toolchain"
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
	tiService  toolchain.Service
	tiRepoMock *toolchain_mock.ToolchainRepository
}

func (s *TITestSuite) SetupTest() {
	s.tiRepoMock = &toolchain_mock.ToolchainRepository{}
	s.tiService = toolchain.NewService(s.tiRepoMock)
}

// Tests that GetToolchainIntegration method proxies to repo
func (s *TITestSuite) TestGetToolchainIntegration() {
	id := "some-id"

	expectedTI := &training3.ToolchainIntegration{}
	expectedErr := errors.New("some-error")

	s.tiRepoMock.On("GetToolchainIntegration", id).Return(expectedTI, expectedErr)

	pi, err := s.tiService.GetToolchainIntegration(id)

	s.Assert().Equal(pi, expectedTI)
	s.Assert().Equal(err, expectedErr)
}

// Tests that GetToolchainIntegrationList method proxies to repo
func (s *TITestSuite) TestGetToolchainIntegrationList() {
	listOption := filter.Page(1)

	expectedTIList := []training3.ToolchainIntegration{}
	expectedErr := errors.New("some-error")

	s.tiRepoMock.
		On("GetToolchainIntegrationList", mock.AnythingOfType("filter.ListOption")).
		Return(expectedTIList, expectedErr)

	piList, err := s.tiService.GetToolchainIntegrationList(listOption)

	s.Assert().Equal(piList, expectedTIList)
	s.Assert().Equal(err, expectedErr)
}

// Tests that CreateToolchainIntegration method sets CreatedAt/UpdatedAt and proxies to repo
func (s *TITestSuite) TestCreateToolchainIntegration() {
	someTI := &training3.ToolchainIntegration{ID: "some-ti"}

	expectedErr := errors.New("some-error")

	s.tiRepoMock.
		On("SaveToolchainIntegration", someTI).
		Return(expectedErr)

	timeBefore := time.Now()
	err := s.tiService.CreateToolchainIntegration(someTI)

	actualTI := s.tiRepoMock.Calls[0].Arguments.Get(0).(*training3.ToolchainIntegration)

	s.Assert().Equal(someTI, actualTI)
	s.Assert().Equal(err, expectedErr)

	s.Assert().True(someTI.CreatedAt.After(timeBefore))
	s.Assert().Equal(someTI.CreatedAt, someTI.UpdatedAt)
}

// Tests that CreateToolchainIntegration method sets CreatedAt/UpdatedAt and proxies to repo
func (s *TITestSuite) TestUpdateToolchainIntegration() {
	someTI := &training3.ToolchainIntegration{ID: "some-ti", CreatedAt: time.Now(), UpdatedAt: time.Now()}

	expectedErr := errors.New("some-error")

	s.tiRepoMock.
		On("UpdateToolchainIntegration", someTI).
		Return(expectedErr)

	timeBeforeUpdate := time.Now()
	err := s.tiService.UpdateToolchainIntegration(someTI)

	actualPI := s.tiRepoMock.Calls[0].Arguments.Get(0).(*training3.ToolchainIntegration)

	s.Assert().Equal(someTI, actualPI)
	s.Assert().Equal(err, expectedErr)

	s.Assert().True(someTI.UpdatedAt.After(timeBeforeUpdate))
	s.Assert().True(someTI.CreatedAt.Before(timeBeforeUpdate))
}

// Tests that DeleteToolchainIntegration method proxies to repo
func (s *TITestSuite) TestDeleteToolchainIntegration() {
	id := "some-id"
	expectedErr := errors.New("some-error")

	s.tiRepoMock.On("DeleteToolchainIntegration", id).Return(expectedErr)

	err := s.tiService.DeleteToolchainIntegration(id)

	s.Assert().Equal(err, expectedErr)
}
