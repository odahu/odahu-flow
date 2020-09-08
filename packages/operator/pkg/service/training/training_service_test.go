/*
 * Copyright 2020 EPAM Systems
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package training_test

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	apis "github.com/odahu/odahu-flow/packages/operator/pkg/apis/training"
	odahu_errs "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/training/mocks"
	service "github.com/odahu/odahu-flow/packages/operator/pkg/service/training"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

const (
	enID = "entity-id"
)

func TestSuiteRun(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

type TestSuite struct {
	suite.Suite
	mockRepo *mocks.Repository
	service service.Service
	db *sql.DB
	dbMock sqlmock.Sqlmock
	as *assert.Assertions
}


func (s *TestSuite) SetupSuite() {
	s.as = assert.New(s.T())
}

func (s *TestSuite) SetupTest() {
	db, dbMock, err := sqlmock.New()
	if err != nil {
		s.T().Fatal("Unable initialize sql mock")
	}
	mockRepo := &mocks.Repository{}

	s.mockRepo = mockRepo
	s.db = db
	s.dbMock = dbMock
	s.service = service.NewService(mockRepo, db)
}

func (s *TestSuite) TestGetModelTraining() {
	as := assert.New(s.T())

	en := newStubMT()
	ctx := context.Background()
	s.mockRepo.On("GetModelTraining", ctx, s.db, enID).Return(en, nil)

	actEn, err := s.service.GetModelTraining(ctx, enID)
	as.NoError(err)
	as.Equal(en, actEn)
	s.mockRepo.AssertExpectations(s.T())
}

func (s *TestSuite) TestGetModelTrainingList() {
	as := assert.New(s.T())

	ens := []apis.ModelTraining{*newStubMT()}
	ctx := context.Background()
	stubFilter := newStubFilter()
	s.mockRepo.
		On("GetModelTrainingList", ctx, s.db, mock.AnythingOfType("filter.ListOption")).
		Return(ens, nil)

	actualEns, err := s.service.GetModelTrainingList(ctx, stubFilter)
	as.NoError(err)
	as.Equal(ens, actualEns)
	s.mockRepo.AssertExpectations(s.T())
}

func (s *TestSuite) TestGetModelTrainingList_Error() {
	as := assert.New(s.T())

	ctx := context.Background()
	stubFilter := newStubFilter()
	anyError := errors.New("any error")
	s.mockRepo.
		On("GetModelTrainingList", ctx, s.db, mock.AnythingOfType("filter.ListOption")).
		Return(nil, anyError)

	actualEns, err := s.service.GetModelTrainingList(ctx, stubFilter)
	as.Error(err)
	as.Equal(anyError, err)
	as.Nil(actualEns)
	s.mockRepo.AssertExpectations(s.T())
}

func (s *TestSuite) TestDeleteModelTraining() {
	as := assert.New(s.T())

	ctx := context.Background()
	s.mockRepo.On("DeleteModelTraining", ctx, s.db, enID).Return(nil)

	as.NoError(s.service.DeleteModelTraining(ctx, enID))
	s.mockRepo.AssertExpectations(s.T())
}

func (s *TestSuite) TestSetDeletionMark() {
	as := assert.New(s.T())

	ctx := context.Background()
	s.mockRepo.On("SetDeletionMark", ctx, s.db, enID, true).Return(nil)

	as.NoError(s.service.SetDeletionMark(ctx, enID, true))
	s.mockRepo.AssertExpectations(s.T())
}

func (s *TestSuite) TestUpdateModelTraining() {
	as := assert.New(s.T())

	ctx := context.Background()
	en := newStubMT()
	s.mockRepo.On("UpdateModelTraining", ctx, s.db, en).Return(nil)

	as.NoError(s.service.UpdateModelTraining(ctx, en))
	s.mockRepo.AssertExpectations(s.T())
}

func (s *TestSuite) TestUpdateModelTrainingStatus() {
	as := assert.New(s.T())

	// Assume entity exists in repository
	ctx := context.Background()
	mockTx := mock.AnythingOfType("*sql.Tx")
	repoEn := newStubMT()
	s.mockRepo.On("GetModelTraining", ctx, mockTx, enID).Return(repoEn, nil)

	// Assume transaction commit
	s.dbMock.ExpectBegin()
	s.dbMock.ExpectCommit()

	// Assume that repository return no error while set new status with not touched spec snapshot
	newStatus := repoEn.Status
	newStatus.PodName = "new Pod name"
	s.mockRepo.
		On("UpdateModelTrainingStatus", ctx, mockTx, enID, newStatus).
		Return(nil)

	// Call service with the same spec snapshot as in repository and new status
	specSnapshot := repoEn.Spec
	as.NoError(s.service.UpdateModelTrainingStatus(ctx, enID, newStatus, specSnapshot))
	as.NoError(s.dbMock.ExpectationsWereMet())
	s.mockRepo.AssertExpectations(s.T())
}

func (s *TestSuite) TestUpdateModelTrainingStatusSpecTouched() {
	as := assert.New(s.T())

	// Assume entity exists in repository
	ctx := context.Background()
	mockTx := mock.AnythingOfType("*sql.Tx")
	repoEn := newStubMT()
	s.mockRepo.On("GetModelTraining", ctx, mockTx, enID).Return(repoEn, nil)

	// Assume transaction rollback
	s.dbMock.ExpectBegin()
	s.dbMock.ExpectRollback()

	// Assume that repository return no error while set new status with not touched spec snapshot
	newStatus := repoEn.Status
	newStatus.PodName = "new Pod name"

	// Call service with the same spec snapshot as in repository and new status
	specSnapshot := repoEn.Spec
	specSnapshot.Image = "image in spec was changed"
	err := s.service.UpdateModelTrainingStatus(ctx, enID, newStatus, specSnapshot)
	as.Error(err)

	// Error about spec was touched must be raised
	as.True(odahu_errs.IsSpecWasTouchedError(err))
	as.Equal(odahu_errs.SpecWasTouched{Entity: enID}, err)

	as.NoError(s.dbMock.ExpectationsWereMet())
	s.mockRepo.AssertExpectations(s.T())
	// Update in repo should not be called
	s.mockRepo.AssertNotCalled(s.T(), "UpdateModelTrainingStatus")
}

func (s *TestSuite) TestCreateModelTraining() {
	as := assert.New(s.T())

	en := newStubMT()
	ctx := context.Background()
	s.mockRepo.On("CreateModelTraining", ctx, s.db, en).Return(nil)

	as.NoError(s.service.CreateModelTraining(ctx, en))
	s.mockRepo.AssertExpectations(s.T())
}

func (s *TestSuite) TestCreateModelTraining_Error() {
	as := assert.New(s.T())

	en := newStubMT()
	ctx := context.Background()
	anyError := errors.New("any error")
	s.mockRepo.On("CreateModelTraining", ctx, s.db, en).Return(anyError)

	as.Error(s.service.CreateModelTraining(ctx, en))
	s.mockRepo.AssertExpectations(s.T())
}

// Helpers

func newStubFilter() filter.ListOption {
	return func(options *filter.ListOptions) {
	}
}

func newStubMT() *apis.ModelTraining {
	return &apis.ModelTraining{
		ID:           enID,
		DeletionMark: false,
		Spec:         v1alpha1.ModelTrainingSpec{},
		Status:       v1alpha1.ModelTrainingStatus{},
	}
}