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
	repo_dep "github.com/odahu/odahu-flow/packages/operator/pkg/repository/deployment/postgres"
	repo_route "github.com/odahu/odahu-flow/packages/operator/pkg/repository/route/postgres"
	route_interface "github.com/odahu/odahu-flow/packages/operator/pkg/repository/route"
	service "github.com/odahu/odahu-flow/packages/operator/pkg/service/deployment"
	route_service "github.com/odahu/odahu-flow/packages/operator/pkg/service/route"
	"github.com/odahu/odahu-flow/packages/operator/pkg/testhelpers/testenvs"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)


func TestIntegrationSuiteRun(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

type IntegrationTestSuite struct {
	suite.Suite
	as *assert.Assertions

	DB      *sql.DB
	closeDB func() error
	service service.Service
	repo repo_dep.DeploymentRepo
	routeRepo repo_route.RouteRepo
	routeService route_service.Service
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.as = assert.New(s.T())
	var err error
	s.DB, _, s.closeDB, err = testenvs.SetupTestDB()
	if err != nil {
		s.FailNow("Unable to init real database")
	}
	s.repo = repo_dep.DeploymentRepo{DB: s.DB}
	s.routeRepo = repo_route.RouteRepo{DB: s.DB}
	s.service = service.NewService(s.repo, s.routeRepo)
	s.routeService = route_service.NewService(s.routeRepo)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	if err := s.closeDB(); err != nil {
		s.T().Fatal("Error during release test DB resources")
	}
}


func (s *IntegrationTestSuite) TestFullCase() {
	as := assert.New(s.T())

	en := newStubMT()
	en2 := newStubMT()
	en2.ID = "newID"
	ctx := context.Background()

	as.NoError(s.service.CreateModelDeployment(ctx, en))
	deps, err := s.repo.GetModelDeploymentList(ctx, nil)
	as.NoError(err)
	as.Len(deps, 1)
	routes, err := s.routeRepo.GetModelRouteList(ctx, nil)
	as.NoError(err)
	as.Len(routes, 1)
	r := routes[0]

	// Check default route
	route, err := s.routeRepo.GetModelRoute(ctx, nil, r.ID)
	as.NoError(err)
	as.Equal(route.Spec.ModelDeploymentTargets[0].Name, en.ID)
	as.Equal(route.Default, true)

	// Check default route via API of route repository
	isDef, err := s.routeRepo.IsDefault(ctx, route.ID, nil)
	as.NoError(err)
	as.True(isDef)

	defExists, err := s.routeRepo.DefaultExists(ctx, en.ID, nil)
	as.NoError(err)
	as.True(defExists)


	// It is impossible to delete default route via route Service (only via deployment service by deletion
	// corresponding ModelDeployment)
	err = s.routeService.DeleteModelRoute(ctx, route.ID)
	as.Errorf(err, "unable to delete default route with ID \"%v\"", route.ID)

	// It is impossible to update default route via route Service
	err = s.routeService.UpdateModelRoute(ctx, route)
	as.Errorf(err, "unable to update default route with ID \"%v\"", route.ID)

	// Add yet another deployment
	as.NoError(s.service.CreateModelDeployment(ctx, en2))

	// Check count of deployments
	deps, err = s.repo.GetModelDeploymentList(ctx, nil)
	as.NoError(err)
	as.Len(deps, 2)

	// Check count of routes without filters
	routes, err = s.routeRepo.GetModelRouteList(ctx, nil)
	as.NoError(err)
	as.Len(routes, 2)

	// Check that there are no NOT default routes
	routes, err = s.routeRepo.GetModelRouteList(ctx, nil, filter.ListFilter(&route_interface.Filter{Default: []bool{
		false,
	}}))
	as.NoError(err)
	as.Len(routes, 0)

	// Try to select route of first deployment
	routes, err = s.routeRepo.GetModelRouteList(ctx, nil, filter.ListFilter(&route_interface.Filter{
		MdID: []string{en.ID},
		Default: []bool{true, false},
	}))
	as.NoError(err)
	as.Len(routes, 1)

	// Delete first deployment
	as.NoError(s.service.DeleteModelDeployment(ctx, en.ID))
	// There are not default routes of first deployment
	routes, err = s.routeRepo.GetModelRouteList(ctx, nil, filter.ListFilter(&route_interface.Filter{
		MdID: []string{en.ID},
		Default: []bool{true},
	}))
	as.NoError(err)
	as.Len(routes, 0)

	// Check the same via api of route repo
	defExists, err = s.routeRepo.DefaultExists(ctx, en.ID, nil)
	as.NoError(err)
	as.False(defExists)

}

