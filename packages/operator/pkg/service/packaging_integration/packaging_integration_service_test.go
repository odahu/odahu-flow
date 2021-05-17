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

package packaging_integration

import (
	"errors"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/packaging"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/packaging/mocks"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

func TestSuiteRun(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

type TestSuite struct {
	suite.Suite
	piService  *service
	piRepoMock *mocks.PackagingIntegrationRepository
}

func (s *TestSuite) SetupTest() {
	s.piRepoMock = &mocks.PackagingIntegrationRepository{}
	s.piService = NewService(s.piRepoMock)
}

// Tests that GetPackagingIntegration method proxies to repo
func (s *TestSuite) TestGetPackagingIntegration() {
	id := "some-id"

	expectedPI := &packaging.PackagingIntegration{}
	expectedErr := errors.New("some-error")

	s.piRepoMock.On("GetPackagingIntegration", id).Return(expectedPI, expectedErr)

	pi, err := s.piService.GetPackagingIntegration(id)

	s.Assert().Equal(pi, expectedPI)
	s.Assert().Equal(err, expectedErr)
}

// Tests that GetPackagingIntegrationList method proxies to repo
func (s *TestSuite) TestGetPackagingIntegrationList() {
	listOption := filter.Page(1)

	expectedPIList := []packaging.PackagingIntegration{}
	expectedErr := errors.New("some-error")

	s.piRepoMock.
		On("GetPackagingIntegrationList", mock.AnythingOfType("filter.ListOption")).
		Return(expectedPIList, expectedErr)

	piList, err := s.piService.GetPackagingIntegrationList(listOption)

	s.Assert().Equal(piList, expectedPIList)
	s.Assert().Equal(err, expectedErr)
}

// Tests that CreatePackagingIntegration method sets CreatedAt/UpdatedAt and proxies to repo
func (s *TestSuite) TestCreatePackagingIntegration() {
	somePI := &packaging.PackagingIntegration{ID: "some-pi"}

	expectedErr := errors.New("some-error")

	s.piRepoMock.
		On("SavePackagingIntegration", somePI).
		Return(expectedErr)

	timeBefore := time.Now()
	err := s.piService.CreatePackagingIntegration(somePI)

	actualPI := s.piRepoMock.Calls[0].Arguments.Get(0).(*packaging.PackagingIntegration)

	s.Assert().Equal(somePI, actualPI)
	s.Assert().Equal(err, expectedErr)

	s.Assert().True(somePI.CreatedAt.After(timeBefore))
	s.Assert().Equal(somePI.CreatedAt, somePI.UpdatedAt)
}

// Tests that CreatePackagingIntegration method sets CreatedAt/UpdatedAt and proxies to repo
func (s *TestSuite) TestUpdatePackagingIntegration() {
	somePI := &packaging.PackagingIntegration{ID: "some-pi", CreatedAt: time.Now(), UpdatedAt: time.Now()}

	expectedErr := errors.New("some-error")

	s.piRepoMock.
		On("UpdatePackagingIntegration", somePI).
		Return(expectedErr)

	timeBeforeUpdate := time.Now()
	err := s.piService.UpdatePackagingIntegration(somePI)

	actualPI := s.piRepoMock.Calls[0].Arguments.Get(0).(*packaging.PackagingIntegration)

	s.Assert().Equal(somePI, actualPI)
	s.Assert().Equal(err, expectedErr)

	s.Assert().True(somePI.UpdatedAt.After(timeBeforeUpdate))
	s.Assert().True(somePI.CreatedAt.Before(timeBeforeUpdate))
}

// Tests that DeletePackagingIntegration method proxies to repo
func (s *TestSuite) TestDeletePackagingIntegration() {
	id := "some-id"
	expectedErr := errors.New("some-error")

	s.piRepoMock.On("DeletePackagingIntegration", id).Return(expectedErr)

	err := s.piService.DeletePackagingIntegration(id)

	s.Assert().Equal(err, expectedErr)
}
