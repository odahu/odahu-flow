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

package deployment_test

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	apis "github.com/odahu/odahu-flow/packages/operator/pkg/apis/deployment"
	odahu_errs "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/deployment/mocks"
	route_mocks "github.com/odahu/odahu-flow/packages/operator/pkg/repository/route/mocks"
	service "github.com/odahu/odahu-flow/packages/operator/pkg/service/deployment"
	event_pub_mocks "github.com/odahu/odahu-flow/packages/operator/pkg/service/deployment/mocks"
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
	mockRepo  *mocks.Repository
	rMockRepo *route_mocks.Repository
	eMockPub  *event_pub_mocks.EventPublisher
	service   service.Service
	db        *sql.DB
	dbMock    sqlmock.Sqlmock
	as        *assert.Assertions
	nilTx     *sql.Tx
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
	rMockRepo := &route_mocks.Repository{}
	eMockPub := &event_pub_mocks.EventPublisher{}

	s.mockRepo = mockRepo
	s.rMockRepo = rMockRepo
	s.eMockPub = eMockPub
	s.db = db
	s.dbMock = dbMock
	s.service = service.NewService(mockRepo, rMockRepo, eMockPub)
}

func (s *TestSuite) TestGetModelDeployment() {
	as := assert.New(s.T())

	en := newStubMT()
	ctx := context.Background()
	s.mockRepo.On("GetModelDeployment", ctx, s.nilTx, enID).Return(en, nil)

	actEn, err := s.service.GetModelDeployment(ctx, enID)
	as.NoError(err)
	as.Equal(en, actEn)
	s.mockRepo.AssertExpectations(s.T())
}

func (s *TestSuite) TestGetModelDeploymentList() {
	as := assert.New(s.T())

	ens := []apis.ModelDeployment{*newStubMT()}
	ctx := context.Background()
	stubFilter := newStubFilter()
	s.mockRepo.
		On("GetModelDeploymentList", ctx, s.nilTx, mock.AnythingOfType("filter.ListOption")).
		Return(ens, nil)

	actualEns, err := s.service.GetModelDeploymentList(ctx, stubFilter)
	as.NoError(err)
	as.Equal(ens, actualEns)
	s.mockRepo.AssertExpectations(s.T())
}

func (s *TestSuite) TestGetModelDeploymentList_Error() {
	as := assert.New(s.T())

	ctx := context.Background()
	stubFilter := newStubFilter()
	anyError := errors.New("any error")
	s.mockRepo.
		On("GetModelDeploymentList", ctx, s.nilTx, mock.AnythingOfType("filter.ListOption")).
		Return(nil, anyError)

	actualEns, err := s.service.GetModelDeploymentList(ctx, stubFilter)
	as.Error(err)
	as.Equal(anyError, err)
	as.Nil(actualEns)
	s.mockRepo.AssertExpectations(s.T())
}

func (s *TestSuite) TestDeleteModelDeployment() {
	as := assert.New(s.T())

	ctx := context.Background()

	// Assume transaction commit
	s.dbMock.ExpectBegin()
	s.dbMock.ExpectCommit()
	mockTx, err := s.db.Begin()
	as.NoError(err)

	defaultRouteID := "defaultRoute"
	s.mockRepo.On("BeginTransaction", ctx).Return(mockTx, nil)
	s.rMockRepo.On("GetModelRouteList", ctx, mockTx, mock.Anything).Return([]apis.ModelRoute{
		{
			ID:      defaultRouteID,
			Default: true,
		},
	}, nil)
	s.rMockRepo.On("DeleteModelRoute", ctx, mockTx, defaultRouteID).Return(nil)
	s.mockRepo.On("DeleteModelDeployment", ctx, mockTx, enID).Return(nil)

	s.eMockPub.On("PublishEvent", ctx, mockTx, mock.Anything).Return(nil)

	as.NoError(s.service.DeleteModelDeployment(ctx, enID))
	s.mockRepo.AssertExpectations(s.T())
	as.NoError(s.dbMock.ExpectationsWereMet())
}

func (s *TestSuite) TestSetDeletionMark() {
	as := assert.New(s.T())
	ctx := context.Background()

	// Assume transaction commit
	s.dbMock.ExpectBegin()
	s.dbMock.ExpectCommit()
	mockTx, err := s.db.Begin()
	as.NoError(err)

	s.mockRepo.On("BeginTransaction", ctx).Return(mockTx, nil)
	s.mockRepo.On("SetDeletionMark", ctx, mockTx, enID, true).Return(nil)
	s.eMockPub.On("PublishEvent", ctx, mockTx, mock.Anything).Return(nil)

	as.NoError(s.service.SetDeletionMark(ctx, enID, true))
	s.mockRepo.AssertExpectations(s.T())
	as.NoError(s.dbMock.ExpectationsWereMet())
}

func (s *TestSuite) TestUpdateModelDeployment() {
	as := assert.New(s.T())
	ctx := context.Background()

	// Assume transaction commit
	s.dbMock.ExpectBegin()
	s.dbMock.ExpectCommit()
	mockTx, err := s.db.Begin()
	as.NoError(err)

	en := newStubMT()
	s.mockRepo.On("BeginTransaction", ctx).Return(mockTx, nil)
	s.mockRepo.On("UpdateModelDeployment", ctx, mockTx, en).Return(nil)
	s.eMockPub.On("PublishEvent", ctx, mockTx, mock.Anything).Return(nil)

	timeBeforeCall := time.Now()
	as.NoError(s.service.UpdateModelDeployment(ctx, en))
	s.mockRepo.AssertExpectations(s.T())
	as.NoError(s.dbMock.ExpectationsWereMet())

	// UpdatedAt field must be updated on now during the invocation
	as.True(timeBeforeCall.Before(en.UpdatedAt))
	// UpdatedAt field must be not updated on now during the invocation
	as.True(timeBeforeCall.After(en.CreatedAt))
	s.mockRepo.AssertExpectations(s.T())
	as.NoError(s.dbMock.ExpectationsWereMet())
}

func (s *TestSuite) TestUpdateModelDeploymentStatus() {
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
	s.mockRepo.On("GetModelDeployment", ctx, mockTx, enID).Return(repoEn, nil)
	s.mockRepo.On("BeginTransaction", ctx).Return(mockTx, nil)
	s.eMockPub.On("PublishEvent", ctx, mockTx, mock.Anything).Return(nil)

	// Assume that repository return no error while set new status with not touched spec snapshot
	newStatus := repoEn.Status
	newStatus.Replicas = 3
	s.mockRepo.
		On("UpdateModelDeploymentStatus", ctx, mockTx, enID, newStatus).
		Return(nil)

	// Call service with the same spec snapshot as in repository and new status
	specSnapshot := repoEn.Spec
	as.NoError(s.service.UpdateModelDeploymentStatus(ctx, enID, newStatus, specSnapshot))
	as.NoError(s.dbMock.ExpectationsWereMet())
	s.mockRepo.AssertExpectations(s.T())
	as.NoError(s.dbMock.ExpectationsWereMet())
}

func (s *TestSuite) TestUpdateModelDeploymentStatusSpecTouched() {
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
	s.mockRepo.On("GetModelDeployment", ctx, mockTx, enID).Return(repoEn, nil)
	s.mockRepo.On("BeginTransaction", ctx).Return(mockTx, nil)
	s.eMockPub.On("PublishEvent", ctx, mockTx, mock.Anything).Return(nil)

	// Assume that repository return no error while set new status with not touched spec snapshot
	newStatus := repoEn.Status
	newStatus.Replicas = 3

	// Call service with the same spec snapshot as in repository and new status
	specSnapshot := repoEn.Spec
	specSnapshot.Image = "image in spec was changed"
	err = s.service.UpdateModelDeploymentStatus(ctx, enID, newStatus, specSnapshot)
	as.Error(err)

	// Error about spec was touched must be raised
	as.True(odahu_errs.IsSpecWasTouchedError(err))
	as.Equal(odahu_errs.SpecWasTouched{Entity: enID}, err)

	as.NoError(s.dbMock.ExpectationsWereMet())
	s.mockRepo.AssertExpectations(s.T())
	// Update in repo should not be called
	s.mockRepo.AssertNotCalled(s.T(), "UpdateModelDeploymentStatus")
}

func (s *TestSuite) TestCreateModelDeployment() {
	as := assert.New(s.T())

	en := newStubMT()
	ctx := context.Background()

	// Assume transaction commit
	s.dbMock.ExpectBegin()
	s.dbMock.ExpectCommit()
	mockTx, err := s.db.Begin()
	if err != nil {
		s.T().Fatal(err)
	}
	s.mockRepo.On("BeginTransaction", ctx).Return(mockTx, nil)
	s.mockRepo.On("SaveModelDeployment", ctx, mockTx, en).Return(nil)
	s.rMockRepo.On("DefaultExists", ctx, enID, mockTx).Return(false, nil)
	s.rMockRepo.On("SaveModelRoute", ctx, mockTx, mock.Anything).
		Return(nil)
	s.eMockPub.On("PublishEvent", ctx, mockTx, mock.Anything).Return(nil)

	timeBeforeCall := time.Now()

	as.NoError(s.service.CreateModelDeployment(ctx, en))
	s.mockRepo.AssertExpectations(s.T())
	as.NoError(s.dbMock.ExpectationsWereMet())

	// CreatedAt, UpdatedAt fields must be updated on now during the invocation
	as.True(timeBeforeCall.Before(en.CreatedAt))
	as.True(timeBeforeCall.Before(en.UpdatedAt))
}

func (s *TestSuite) TestCreateModelDeployment_Error() {
	as := assert.New(s.T())

	en := newStubMT()

	ctx := context.Background()
	anyError := errors.New("any error")

	// Assume transaction commit
	s.dbMock.ExpectBegin()
	s.dbMock.ExpectRollback()
	mockTx, err := s.db.Begin()
	if err != nil {
		s.T().Fatal(err)
	}
	s.mockRepo.On("BeginTransaction", ctx).Return(mockTx, nil)
	s.mockRepo.On("SaveModelDeployment", ctx, mockTx, en).Return(anyError)

	as.Error(s.service.CreateModelDeployment(ctx, en))
	s.mockRepo.AssertExpectations(s.T())
	s.rMockRepo.AssertNotCalled(s.T(), "DefaultExists")
	s.rMockRepo.AssertNotCalled(s.T(), "SaveModelRoute")
	as.NoError(s.dbMock.ExpectationsWereMet())
}

// Helpers

func newStubFilter() filter.ListOption {
	return func(options *filter.ListOptions) {
	}
}

func newStubMT() *apis.ModelDeployment {
	return &apis.ModelDeployment{
		ID:           enID,
		DeletionMark: false,
		Spec:         v1alpha1.ModelDeploymentSpec{},
		Status:       v1alpha1.ModelDeploymentStatus{},
	}
}
