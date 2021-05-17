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

package packaging_test

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	apis "github.com/odahu/odahu-flow/packages/operator/pkg/apis/packaging"
	odahu_errs "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/packaging/mocks"
	service "github.com/odahu/odahu-flow/packages/operator/pkg/service/packaging"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
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
	service  service.Service
	db       *sql.DB
	dbMock   sqlmock.Sqlmock
	as       *assert.Assertions
	nilTx    *sql.Tx
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
	s.service = service.NewService(mockRepo)
}

func (s *TestSuite) TestGetModelPackaging() {
	as := assert.New(s.T())

	en := newStubMT()
	ctx := context.Background()
	s.mockRepo.On("GetModelPackaging", ctx, s.nilTx, enID).Return(en, nil)

	actEn, err := s.service.GetModelPackaging(ctx, enID)
	as.NoError(err)
	as.Equal(en, actEn)
	s.mockRepo.AssertExpectations(s.T())
}

func (s *TestSuite) TestGetModelPackagingList() {
	as := assert.New(s.T())

	ens := []apis.ModelPackaging{*newStubMT()}
	ctx := context.Background()
	stubFilter := newStubFilter()
	s.mockRepo.
		On("GetModelPackagingList", ctx, s.nilTx, mock.AnythingOfType("filter.ListOption")).
		Return(ens, nil)

	actualEns, err := s.service.GetModelPackagingList(ctx, stubFilter)
	as.NoError(err)
	as.Equal(ens, actualEns)
	s.mockRepo.AssertExpectations(s.T())
}

func (s *TestSuite) TestGetModelPackagingList_Error() {
	as := assert.New(s.T())

	ctx := context.Background()
	stubFilter := newStubFilter()
	anyError := errors.New("any error")
	s.mockRepo.
		On("GetModelPackagingList", ctx, s.nilTx, mock.AnythingOfType("filter.ListOption")).
		Return(nil, anyError)

	actualEns, err := s.service.GetModelPackagingList(ctx, stubFilter)
	as.Error(err)
	as.Equal(anyError, err)
	as.Nil(actualEns)
	s.mockRepo.AssertExpectations(s.T())
}

func (s *TestSuite) TestDeleteModelPackaging() {
	as := assert.New(s.T())

	ctx := context.Background()
	s.mockRepo.On("DeleteModelPackaging", ctx, s.nilTx, enID).Return(nil)

	as.NoError(s.service.DeleteModelPackaging(ctx, enID))
	s.mockRepo.AssertExpectations(s.T())
}

func (s *TestSuite) TestSetDeletionMark() {
	as := assert.New(s.T())

	ctx := context.Background()
	s.mockRepo.On("SetDeletionMark", ctx, s.nilTx, enID, true).Return(nil)

	as.NoError(s.service.SetDeletionMark(ctx, enID, true))
	s.mockRepo.AssertExpectations(s.T())
}

func (s *TestSuite) TestUpdateModelPackaging() {
	as := assert.New(s.T())

	ctx := context.Background()
	en := newStubMT()
	s.mockRepo.On("UpdateModelPackaging", ctx, s.nilTx, en).Return(nil)

	as.NoError(s.service.UpdateModelPackaging(ctx, en))
	s.mockRepo.AssertExpectations(s.T())
}

func (s *TestSuite) TestUpdateModelPackaging_UpdatedAt() {
	as := assert.New(s.T())

	ctx := context.Background()
	en := newStubMT()
	s.mockRepo.On("UpdateModelPackaging", ctx, s.nilTx, en).Return(nil)

	timeBeforeCall := time.Now()
	as.NoError(s.service.UpdateModelPackaging(ctx, en))
	// UpdatedAt field must be updated on now during the invocation
	as.True(timeBeforeCall.Before(en.UpdatedAt))
	// UpdatedAt field must be not updated on now during the invocation
	as.True(timeBeforeCall.After(en.CreatedAt))
	s.mockRepo.AssertExpectations(s.T())
}

func (s *TestSuite) TestUpdateModelPackagingStatus() {
	as := assert.New(s.T())

	// Assume transaction commit
	s.dbMock.ExpectBegin()
	s.dbMock.ExpectCommit()

	// Assume entity exists in repository
	ctx := context.Background()
	mockTx, err := s.db.Begin()
	if err != nil {
		s.T().Fatal(err)
	}
	repoEn := newStubMT()
	s.mockRepo.On("GetModelPackaging", ctx, mockTx, enID).Return(repoEn, nil)
	s.mockRepo.On("BeginTransaction", ctx).Return(mockTx, nil)

	// Assume that repository return no error while set new status with not touched spec snapshot
	newStatus := repoEn.Status
	newStatus.PodName = "new Pod name"
	s.mockRepo.
		On("UpdateModelPackagingStatus", ctx, mockTx, enID, newStatus).
		Return(nil)

	// Call service with the same spec snapshot as in repository and new status
	specSnapshot := repoEn.Spec
	as.NoError(s.service.UpdateModelPackagingStatus(ctx, enID, newStatus, specSnapshot))
	as.NoError(s.dbMock.ExpectationsWereMet())
	s.mockRepo.AssertExpectations(s.T())
}

func (s *TestSuite) TestUpdateModelPackagingStatusSpecTouched() {
	as := assert.New(s.T())

	// Assume transaction rollback
	s.dbMock.ExpectBegin()
	s.dbMock.ExpectRollback()

	// Assume entity exists in repository
	ctx := context.Background()
	mockTx, err := s.db.Begin()
	if err != nil {
		s.T().Fatal(err)
	}
	repoEn := newStubMT()
	s.mockRepo.On("GetModelPackaging", ctx, mockTx, enID).Return(repoEn, nil)
	s.mockRepo.On("BeginTransaction", ctx).Return(mockTx, nil)

	// Assume that repository return no error while set new status with not touched spec snapshot
	newStatus := repoEn.Status
	newStatus.PodName = "new Pod name"

	// Call service with the same spec snapshot as in repository and new status
	specSnapshot := repoEn.Spec
	specSnapshot.Image = "image in spec was changed"
	err = s.service.UpdateModelPackagingStatus(ctx, enID, newStatus, specSnapshot)
	as.Error(err)

	// Error about spec was touched must be raised
	as.True(odahu_errs.IsSpecWasTouchedError(err))
	as.Equal(odahu_errs.SpecWasTouched{Entity: enID}, err)

	as.NoError(s.dbMock.ExpectationsWereMet())
	s.mockRepo.AssertExpectations(s.T())
	// Update in repo should not be called
	s.mockRepo.AssertNotCalled(s.T(), "UpdateModelPackagingStatus")
}

func (s *TestSuite) TestCreateModelPackaging() {
	as := assert.New(s.T())

	en := newStubMT()
	ctx := context.Background()
	s.mockRepo.On("SaveModelPackaging", ctx, s.nilTx, en).Return(nil)

	as.NoError(s.service.CreateModelPackaging(ctx, en))
	s.mockRepo.AssertExpectations(s.T())
}
func (s *TestSuite) TestCreateModelPackaging_CreatedAt() {
	as := assert.New(s.T())

	en := newStubMT()
	ctx := context.Background()
	s.mockRepo.On("SaveModelPackaging", ctx, s.nilTx, en).Return(nil)

	timeBeforeCall := time.Now()
	as.NoError(s.service.CreateModelPackaging(ctx, en))
	// CreatedAt, UpdatedAt fields must be updated on now during the invocation
	as.True(timeBeforeCall.Before(en.CreatedAt))
	as.True(timeBeforeCall.Before(en.UpdatedAt))
	s.mockRepo.AssertExpectations(s.T())
}

func (s *TestSuite) TestCreateModelPackaging_Error() {
	as := assert.New(s.T())

	en := newStubMT()
	ctx := context.Background()
	anyError := errors.New("any error")
	s.mockRepo.On("SaveModelPackaging", ctx, s.nilTx, en).Return(anyError)

	as.Error(s.service.CreateModelPackaging(ctx, en))
	s.mockRepo.AssertExpectations(s.T())
}

// Helpers

func newStubFilter() filter.ListOption {
	return func(options *filter.ListOptions) {
	}
}

func newStubMT() *apis.ModelPackaging {
	return &apis.ModelPackaging{
		ID:           enID,
		DeletionMark: false,
		Spec:         apis.ModelPackagingSpec{},
		Status:       v1alpha1.ModelPackagingStatus{},
	}
}
