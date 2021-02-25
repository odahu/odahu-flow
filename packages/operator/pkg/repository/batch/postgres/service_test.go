/*
 *
 *     Copyright 2021 EPAM Systems
 *
 *     Licensed under the Apache License, Version 2.0 (the "License");
 *     you may not use this file except in compliance with the License.
 *     You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 *     Unless required by applicable law or agreed to in writing, software
 *     distributed under the License is distributed on an "AS IS" BASIS,
 *     WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *     See the License for the specific language governing permissions and
 *     limitations under the License.
 */

package postgres_test

import (
	"context"
	"fmt"
	api_types "github.com/odahu/odahu-flow/packages/operator/pkg/apis/batch"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/batch/postgres"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"
	hashutil "github.com/odahu/odahu-flow/packages/operator/pkg/utils/hash"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
	"time"
)

var testCreateServiceCases = []struct {
	testName string
	bis    api_types.InferenceService // ID will be assigned automatically as "bij{testNumber}"
	expectedErrString string
}{
	{
		testName: "ok",
		bis: api_types.InferenceService{
			DeletionMark: true,
			CreatedAt:    time.Now().UTC(),
			UpdatedAt:    time.Now().UTC(),
			Spec:         api_types.InferenceServiceSpec{
				Image: "some:image",
			},
			Status:       api_types.InferenceServiceStatus{
			},
		},
	},
}

func TestCreateService(t *testing.T) {

	r := postgres.BISRepo{DB: db}
	for i, test := range testCreateServiceCases {
		i := i
		test := test
		t.Run(fmt.Sprintf("TestCreateService#%s", test.testName), func(t *testing.T) {

			req := require.New(t)
			defer func() {
				req.NoError(cleanupJobsServices())
			}()

			bis := test.bis

			// generate ID
			bisID := "bis-" + strconv.Itoa(i)
			bis.ID =bisID
			err := r.Create(context.TODO(), nil, bis)

			// Check expectation about returned error
			if len(test.expectedErrString) > 0 {
				req.Error(err)
				req.Equal(test.expectedErrString, err.Error())
			} else {
				req.Nil(err)
			}

			// Check that retrieved object equals to created one
			res, err := r.Get(context.TODO(), nil, bisID)
			req.Nil(err)

			hashutil.Equal(bis, res)
		})
	}
}


var testUpdateServiceCases = []struct {
	testName          string
	bisID             string
	updated            api_types.InferenceService
	fixtures          []api_types.InferenceService // Entities that should be created in DB before test logic
	expectedErrString string
}{
	{
		testName:          "not found",
		bisID:             "entity",
		updated:            api_types.InferenceService{ID: "entity", Spec: api_types.InferenceServiceSpec{
			Image: "image",
		}},
		fixtures:          nil,
		expectedErrString: "entity \"entity\" is not found",
	},
	{
		testName: "ok",
		bisID:    "entity",
		updated:            api_types.InferenceService{ID: "entity", Spec: api_types.InferenceServiceSpec{
			Image: "image",
		}},
		fixtures:         []api_types.InferenceService{{ID: "entity"}},
	},
}

func TestUpdateService(t *testing.T) {

	r := postgres.BISRepo{DB: db}
	for _, test := range testUpdateServiceCases {
		test := test
		t.Run(fmt.Sprintf("TestUpdateService#%s", test.testName), func(t *testing.T) {
			req := require.New(t)
			defer func() {
				req.NoError(cleanupJobsServices())
			}()

			// Create fixtures if required
			for _, bis := range test.fixtures {
				err := r.Create(context.TODO(), nil, bis)
				req.NoError(err)
			}

			// Try to update status
			err := r.Update(context.TODO(), nil, test.bisID, test.updated)

			// Check expectation about returned error
			if len(test.expectedErrString) > 0 {
				req.EqualError(err, test.expectedErrString)
			} else {
				req.NoError(err)
				// Check that status is updated
				res, err := r.Get(context.TODO(), nil, test.bisID)
				req.NoError(err)

				hashutil.Equal(res, test.updated)
			}
		})
	}
}


var testDeleteServiceCases = []struct {
	testName          string
	id                string
	fixtures          []api_types.InferenceService // Entities that should be created in DB before test logic
	jobFixtures          []api_types.InferenceJob // Entities that should be created in DB before test logic
	expectedErrString string
}{
	{
		testName:          "not found",
		id:                "entity",
		fixtures:          nil,
		expectedErrString: "entity \"entity\" is not found",
	},
	{
		testName: "ok",
		id:       "entity",
		fixtures: []api_types.InferenceService{{ID: "entity"}},
	},
	{
		testName: "error because some job exist",
		id:       "entity",
		fixtures: []api_types.InferenceService{{ID: "entity"}},
		jobFixtures: []api_types.InferenceJob{{ID: "job", Spec: api_types.InferenceJobSpec{
			Service: "entity",
		}}},
		expectedErrString: `Unable to delete service: "entity". Cause: there are child jobs`,
	},
}

func TestDeleteService(t *testing.T) {

	r := postgres.BISRepo{DB: db}
	jR := postgres.BIJRepo{DB: db}
	for _, test := range testDeleteServiceCases {
		test := test
		t.Run(fmt.Sprintf("TestDelete#%s", test.testName), func(t *testing.T) {
			req := require.New(t)
			defer func() {
				req.NoError(cleanupJobsServices())
			}()

			// Create fixtures if required
			for _, bis := range test.fixtures {
				err := r.Create(context.TODO(), nil, bis)
				req.NoError(err)
			}
			for _, bij := range test.jobFixtures {
				err := jR.Create(context.TODO(), nil, bij)
				req.NoError(err)
			}

			// Try to delete
			err := r.Delete(context.TODO(), nil, test.id)

			// Check expectation about returned error
			if len(test.expectedErrString) > 0 {
				req.EqualError(err, test.expectedErrString)
			} else {
				// If deletion was successful then we expect not found while retrieving
				req.NoError(err)
				_, err := r.Get(context.TODO(), nil, test.id)
				req.EqualError(err, "entity \"entity\" is not found")
			}
		})
	}
}


var testListServiceCases = []struct {
	testName string
	fixtures []api_types.InferenceService // Entities that should be created in DB before test logic
	expectedLen int
	filter []filter.ListOption
}{
	{
		testName: "empty",
		fixtures:          nil,
		expectedLen: 0,
	},
	{
		testName: "full",
		fixtures:         []api_types.InferenceService{{ID: "entity"},{ID: "entity2"},{ID: "entity3"}},
		expectedLen: 3,
	},
	{
		testName: "first page with size=1",
		fixtures:         []api_types.InferenceService{{ID: "entity"},{ID: "entity2"},{ID: "entity3"}},
		expectedLen: 1,
		filter: []filter.ListOption{filter.Page(0), filter.Size(1)},
	},
	{
		testName: "second page with size=1",
		fixtures:         []api_types.InferenceService{{ID: "entity"},{ID: "entity2"},{ID: "entity3"}},
		expectedLen: 1,
		filter: []filter.ListOption{filter.Page(1), filter.Size(1)},
	},
	{
		testName: "forth page with size=1",
		fixtures:         []api_types.InferenceService{{ID: "entity"},{ID: "entity2"},{ID: "entity3"}},
		expectedLen: 0,
		filter: []filter.ListOption{filter.Page(3), filter.Size(1)},
	},
	{
		testName: "first page with size=2",
		fixtures:         []api_types.InferenceService{{ID: "entity"},{ID: "entity2"},{ID: "entity3"}},
		expectedLen: 2,  // all except last one
		filter: []filter.ListOption{filter.Page(0), filter.Size(2)},
	},
}

func TestListService(t *testing.T) {

	r := postgres.BISRepo{DB: db}
	for _, test := range testListServiceCases {
		test := test
		t.Run(fmt.Sprintf("TestListService#%s", test.testName), func(t *testing.T) {
			req := require.New(t)
			defer func() {
				req.NoError(cleanupJobsServices())
			}()

			// Create fixtures if required
			for _, bij := range test.fixtures {
				err := r.Create(context.TODO(), nil, bij)
				req.NoError(err)
			}

			// Try to List
			res, err := r.List(context.TODO(), nil, test.filter...)

			req.NoError(err)
			req.Len(res, test.expectedLen)
		})
	}
}


var testGetServiceCases = []struct {
	testName          string
	id                string
	fixtures          []api_types.InferenceService // Services that should be created in DB before test logic
	expectedErrString string
}{
	{
		testName:          "not found",
		id:                "entity",
		fixtures:          nil,
		expectedErrString: "entity \"entity\" is not found",
	},
	{
		testName: "ok",
		id:       "entity",
		fixtures: []api_types.InferenceService{{ID: "entity"}},
	},
}

func TestGetService(t *testing.T) {

	r := postgres.BISRepo{DB: db}
	for _, test := range testGetServiceCases {
		test := test
		t.Run(fmt.Sprintf("TestGetService#%s", test.testName), func(t *testing.T) {
			req := require.New(t)
			defer func() {
				req.NoError(cleanupJobsServices())
			}()

			// Create fixtures if required
			for _, bij := range test.fixtures {
				err := r.Create(context.TODO(), nil, bij)
				req.NoError(err)
			}

			// Try to Get
			_, err := r.Get(context.TODO(), nil, test.id)

			// Check expectation about returned error
			if len(test.expectedErrString) > 0 {
				req.EqualError(err, test.expectedErrString)
			} else {
				req.NoError(err)
			}
		})
	}
}


