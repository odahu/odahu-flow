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

package route_test

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	apis "github.com/odahu/odahu-flow/packages/operator/pkg/apis/deployment"
	odahu_errs "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/route/mocks"
	event_pub_mocks "github.com/odahu/odahu-flow/packages/operator/pkg/service/deployment/mocks"
	service "github.com/odahu/odahu-flow/packages/operator/pkg/service/route"
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
	eMockPub *event_pub_mocks.EventPublisher
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
	eMockPub := &event_pub_mocks.EventPublisher{}

	s.mockRepo = mockRepo
	s.db = db
	s.dbMock = dbMock
	s.eMockPub = eMockPub
	s.service = service.NewService(mockRepo, eMockPub)
}

func (s *TestSuite) TestGetModelRoute() {
	as := assert.New(s.T())

	en := newStubMT()
	ctx := context.Background()
	s.mockRepo.On("GetModelRoute", ctx, s.nilTx, enID).Return(en, nil)

	actEn, err := s.service.GetModelRoute(ctx, enID)
	as.NoError(err)
	as.Equal(en, actEn)
	s.mockRepo.AssertExpectations(s.T())
}

func (s *TestSuite) TestGetModelRouteList() {
	as := assert.New(s.T())

	ens := []apis.ModelRoute{*newStubMT()}
	ctx := context.Background()
	stubFilter := newStubFilter()
	s.mockRepo.
		On("GetModelRouteList", ctx, s.nilTx, mock.AnythingOfType("filter.ListOption")).
		Return(ens, nil)

	actualEns, err := s.service.GetModelRouteList(ctx, stubFilter)
	as.NoError(err)
	as.Equal(ens, actualEns)
	s.mockRepo.AssertExpectations(s.T())
}

func (s *TestSuite) TestGetModelRouteList_Error() {
	as := assert.New(s.T())

	ctx := context.Background()
	stubFilter := newStubFilter()
	anyError := errors.New("any error")
	s.mockRepo.
		On("GetModelRouteList", ctx, s.nilTx, mock.AnythingOfType("filter.ListOption")).
		Return(nil, anyError)

	actualEns, err := s.service.GetModelRouteList(ctx, stubFilter)
	as.Error(err)
	as.Equal(anyError, err)
	as.Nil(actualEns)
	s.mockRepo.AssertExpectations(s.T())
}

func (s *TestSuite) TestDeleteModelRoute() {
	as := assert.New(s.T())

	ctx := context.Background()

	// Assume transaction commit
	s.dbMock.ExpectBegin()
	s.dbMock.ExpectCommit()
	mockTx, err := s.db.Begin()
	if err != nil {
		s.T().Fatal(err)
	}
	s.mockRepo.On("BeginTransaction", ctx).Return(mockTx, nil)
	s.mockRepo.On("IsDefault", ctx, enID, mockTx).Return(false, nil)
	s.mockRepo.On("DeleteModelRoute", ctx, mockTx, enID).Return(nil)
	s.eMockPub.On("PublishEvent", ctx, mockTx, mock.Anything).Return(nil)

	as.NoError(s.service.DeleteModelRoute(ctx, enID))
	s.mockRepo.AssertExpectations(s.T())
	as.NoError(s.dbMock.ExpectationsWereMet())
}

func (s *TestSuite) TestSetDeletionMark() {
	as := assert.New(s.T())

	// Assume transaction commit
	s.dbMock.ExpectBegin()
	s.dbMock.ExpectCommit()
	mockTx, err := s.db.Begin()
	as.NoError(err)

	ctx := context.Background()
	s.mockRepo.On("BeginTransaction", ctx).Return(mockTx, nil)
	s.mockRepo.On("SetDeletionMark", ctx, mockTx, enID, true).Return(nil)
	s.eMockPub.On("PublishEvent", ctx, mockTx, mock.Anything).Return(nil)

	as.NoError(s.service.SetDeletionMark(ctx, enID, true))
	s.mockRepo.AssertExpectations(s.T())
	as.NoError(s.dbMock.ExpectationsWereMet())
}

func (s *TestSuite) TestUpdateModelRoute() {
	as := assert.New(s.T())

	ctx := context.Background()
	en := newStubMT()

	// Assume transaction commit
	s.dbMock.ExpectBegin()
	s.dbMock.ExpectCommit()
	mockTx, err := s.db.Begin()
	if err != nil {
		s.T().Fatal(err)
	}
	s.mockRepo.On("BeginTransaction", ctx).Return(mockTx, nil)
	s.mockRepo.On("IsDefault", ctx, enID, mockTx).Return(false, nil)
	s.mockRepo.On("UpdateModelRoute", ctx, mockTx, en).Return(nil)
	s.eMockPub.On("PublishEvent", ctx, mockTx, mock.Anything).Return(nil)

	as.NoError(s.service.UpdateModelRoute(ctx, en))
	s.mockRepo.AssertExpectations(s.T())
	as.NoError(s.dbMock.ExpectationsWereMet())
}

func (s *TestSuite) TestUpdateModelRoute_UpdatedAt() {
	as := assert.New(s.T())

	ctx := context.Background()
	en := newStubMT()

	// Assume transaction commit
	s.dbMock.ExpectBegin()
	s.dbMock.ExpectCommit()
	mockTx, err := s.db.Begin()
	if err != nil {
		s.T().Fatal(err)
	}
	s.mockRepo.On("BeginTransaction", ctx).Return(mockTx, nil)
	s.mockRepo.On("IsDefault", ctx, enID, mockTx).Return(false, nil)
	s.mockRepo.On("UpdateModelRoute", ctx, mockTx, en).Return(nil)
	s.eMockPub.On("PublishEvent", ctx, mockTx, mock.Anything).Return(nil)

	timeBeforeCall := time.Now()
	as.NoError(s.service.UpdateModelRoute(ctx, en))
	// UpdatedAt field must be updated on now during the invocation
	as.True(timeBeforeCall.Before(en.UpdatedAt))
	// UpdatedAt field must be not updated on now during the invocation
	as.True(timeBeforeCall.After(en.CreatedAt))
	s.mockRepo.AssertExpectations(s.T())
	as.NoError(s.dbMock.ExpectationsWereMet())
}

func (s *TestSuite) TestUpdateModelRouteStatus() {
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
	s.mockRepo.On("GetModelRoute", ctx, mockTx, enID).Return(repoEn, nil)
	s.mockRepo.On("BeginTransaction", ctx).Return(mockTx, nil)
	s.eMockPub.On("PublishEvent", ctx, mockTx, mock.Anything).Return(nil)

	// Assume that repository return no error while set new status with not touched spec snapshot
	newStatus := repoEn.Status
	newStatus.EdgeURL = "old"
	s.mockRepo.
		On("UpdateModelRouteStatus", ctx, mockTx, enID, newStatus).
		Return(nil)

	// Call service with the same spec snapshot as in repository and new status
	specSnapshot := repoEn.Spec
	as.NoError(s.service.UpdateModelRouteStatus(ctx, enID, newStatus, specSnapshot))
	as.NoError(s.dbMock.ExpectationsWereMet())
	s.mockRepo.AssertExpectations(s.T())
}

func (s *TestSuite) TestUpdateModelRouteStatusSpecTouched() {
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
	s.mockRepo.On("GetModelRoute", ctx, mockTx, enID).Return(repoEn, nil)
	s.mockRepo.On("BeginTransaction", ctx).Return(mockTx, nil)
	s.eMockPub.On("PublishEvent", ctx, mockTx, mock.Anything).Return(nil)

	// Assume that repository return no error while set new status with not touched spec snapshot
	newStatus := repoEn.Status
	newStatus.EdgeURL = "old"

	// Call service with the same spec snapshot as in repository and new status
	specSnapshot := repoEn.Spec
	specSnapshot.URLPrefix = "prefix was changed"
	err = s.service.UpdateModelRouteStatus(ctx, enID, newStatus, specSnapshot)
	as.Error(err)

	// Error about spec was touched must be raised
	as.True(odahu_errs.IsSpecWasTouchedError(err))
	as.Equal(odahu_errs.SpecWasTouched{Entity: enID}, err)

	as.NoError(s.dbMock.ExpectationsWereMet())
	s.mockRepo.AssertExpectations(s.T())
	// Update in repo should not be called
	s.mockRepo.AssertNotCalled(s.T(), "UpdateModelRouteStatus")
}

func (s *TestSuite) TestCreateModelRoute() {
	as := assert.New(s.T())

	s.dbMock.ExpectBegin()
	s.dbMock.ExpectCommit()
	mockTx, err := s.db.Begin()
	as.NoError(err)

	en := newStubMT()
	ctx := context.Background()
	s.mockRepo.On("BeginTransaction", ctx).Return(mockTx, nil)
	s.mockRepo.On("SaveModelRoute", ctx, mockTx, en).Return(nil)
	s.eMockPub.On("PublishEvent", ctx, mockTx, mock.Anything).Return(nil)

	as.NoError(s.service.CreateModelRoute(ctx, en))
	s.mockRepo.AssertExpectations(s.T())
	as.NoError(s.dbMock.ExpectationsWereMet())
}

func (s *TestSuite) TestCreateModelRoute_CreatedAt() {
	as := assert.New(s.T())

	en := newStubMT()

	s.dbMock.ExpectBegin()
	s.dbMock.ExpectCommit()
	mockTx, err := s.db.Begin()
	as.NoError(err)

	ctx := context.Background()
	s.mockRepo.On("BeginTransaction", ctx).Return(mockTx, nil)
	s.mockRepo.On("SaveModelRoute", ctx, mockTx, en).Return(nil)
	s.eMockPub.On("PublishEvent", ctx, mockTx, mock.Anything).Return(nil)

	timeBeforeCall := time.Now()
	as.NoError(s.service.CreateModelRoute(ctx, en))
	// CreatedAt, UpdatedAt fields must be updated on now during the invocation
	as.True(timeBeforeCall.Before(en.CreatedAt))
	as.True(timeBeforeCall.Before(en.UpdatedAt))
	s.mockRepo.AssertExpectations(s.T())
	as.NoError(s.dbMock.ExpectationsWereMet())
}

func (s *TestSuite) TestCreateModelRoute_Error() {
	as := assert.New(s.T())

	s.dbMock.ExpectBegin()
	s.dbMock.ExpectRollback()
	mockTx, err := s.db.Begin()
	as.NoError(err)

	en := newStubMT()
	ctx := context.Background()
	anyError := errors.New("any error")
	s.mockRepo.On("BeginTransaction", ctx).Return(mockTx, nil)
	s.mockRepo.On("SaveModelRoute", ctx, mockTx, en).Return(anyError)
	s.eMockPub.On("PublishEvent", ctx, mockTx, mock.Anything).Return(nil)

	as.Error(s.service.CreateModelRoute(ctx, en))
	s.mockRepo.AssertExpectations(s.T())
	as.NoError(s.dbMock.ExpectationsWereMet())
}

// Helpers

func newStubFilter() filter.ListOption {
	return func(options *filter.ListOptions) {
	}
}

func newStubMT() *apis.ModelRoute {
	return &apis.ModelRoute{
		ID:           enID,
		DeletionMark: false,
		Spec:         v1alpha1.ModelRouteSpec{},
		Status:       v1alpha1.ModelRouteStatus{},
	}
}
