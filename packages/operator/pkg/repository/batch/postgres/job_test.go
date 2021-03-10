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
	sq "github.com/Masterminds/squirrel"
	api_types "github.com/odahu/odahu-flow/packages/operator/pkg/apis/batch"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/batch/postgres"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"
	hashutil "github.com/odahu/odahu-flow/packages/operator/pkg/utils/hash"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var testCreateCases = []struct {
	testName string
	bij    api_types.InferenceJob // ID will be assigned automatically as "bij{testNumber}"
	bisFixtures []api_types.InferenceService
	expectedErrString string
}{
	{
		testName: "ok",
		bij: api_types.InferenceJob{
			ID: "bij",
			DeletionMark: true,
			CreatedAt:    time.Now().UTC(),
			UpdatedAt:    time.Now().UTC(),
			Spec:         api_types.InferenceJobSpec{
				InferenceServiceID: "service",
				BatchRequestID:     "some-id",
			},
			Status:       api_types.InferenceJobStatus{
				State: "Creating with prefilled reason",
			},
		},
		bisFixtures: []api_types.InferenceService{{
			ID: "service",
		}},
	},
	{
		testName: "service not found",
		bij: api_types.InferenceJob{
			ID: "bij",
			DeletionMark: true,
			CreatedAt:    time.Now().UTC(),
			UpdatedAt:    time.Now().UTC(),
			Spec:         api_types.InferenceJobSpec{
				InferenceServiceID: "service",
				BatchRequestID:     "some-id",
			},
			Status:       api_types.InferenceJobStatus{
				State: "Creating with prefilled reason",
			},
		},
		expectedErrString: `Unable to create job: "bij". There is no service with ID: service`,
	},
}

func TestCreate(t *testing.T) {

	r := postgres.BIJRepo{DB: db}
	rs := postgres.BISRepo{DB: db}
	for _, test := range testCreateCases {
		t.Run(fmt.Sprintf("TestCreate#%s", test.testName), func(t *testing.T) {

			req := require.New(t)
			defer func() {
				req.NoError(cleanupJobsServices())
			}()

			bij := test.bij

			// Create fixtures if required
			for _, bis := range test.bisFixtures {
				err := rs.Create(context.TODO(), nil, bis)
				req.NoError(err)
			}

			// generate ID
			err := r.Create(context.TODO(), nil, bij)

			// Check expectation about returned error
			if len(test.expectedErrString) > 0 {
				req.Error(err)
				req.Equal(test.expectedErrString, err.Error())
			} else {
				req.Nil(err)
				// Check that retrieved object equals to created one
				res, err := r.Get(context.TODO(), nil, bij.ID)
				req.Nil(err)
				hashutil.Equal(bij, res)
			}

		})
	}
}


var testUpdateStatusCases = []struct {
	testName string
	bijID string
	bijs    api_types.InferenceJobStatus
	bisFixtures []api_types.InferenceService
	fixtures []api_types.InferenceJob // Jobs that should be created in DB before test logic
	expectedErrString string
}{
	{
		testName: "not found",
		bijID: "job",
		bijs:              api_types.InferenceJobStatus{},
		fixtures:          nil,
		expectedErrString: "entity \"job\" is not found",
	},
	{
		testName: "ok",
		bijID: "job",
		bijs:              api_types.InferenceJobStatus{
			State: "New state",
		},
		fixtures:         []api_types.InferenceJob{{ID: "job", Spec: api_types.InferenceJobSpec{
			InferenceServiceID: "service",
		}}},
		bisFixtures: []api_types.InferenceService{{
			ID: "service",
		}},
	},
}

func TestUpdateStatus(t *testing.T) {

	r := postgres.BIJRepo{DB: db}
	rs := postgres.BISRepo{DB: db}
	for _, test := range testUpdateStatusCases {
		t.Run(fmt.Sprintf("TestUpdateStatus#%s", test.testName), func(t *testing.T) {
			req := require.New(t)
			defer func() {
				req.NoError(cleanupJobsServices())
			}()

			// Create fixtures if required
			for _, bis := range test.bisFixtures {
				err := rs.Create(context.TODO(), nil, bis)
				req.NoError(err)
			}

			for _, bij := range test.fixtures {
				err := r.Create(context.TODO(), nil, bij)
				req.NoError(err)
			}

			// Try to update status
			err := r.UpdateStatus(context.TODO(), nil, test.bijID, test.bijs)

			// Check expectation about returned error
			if len(test.expectedErrString) > 0 {
				req.EqualError(err, test.expectedErrString)
			} else {
				req.NoError(err)
				// Check that status is updated
				res, err := r.Get(context.TODO(), nil, test.bijID)
				req.NoError(err)

				hashutil.Equal(res.Status, test.bijs)
			}
		})
	}
}


var testDeleteCases = []struct {
	testName string
	bijID string
	fixtures []api_types.InferenceJob // Jobs that should be created in DB before test logic
	bisFixtures []api_types.InferenceService
	expectedErrString string
}{
	{
		testName: "not found",
		bijID: "job",
		fixtures:          nil,
		expectedErrString: "entity \"job\" is not found",
	},
	{
		testName: "ok",
		bijID: "job",
		fixtures:         []api_types.InferenceJob{{ID: "job", Spec: api_types.InferenceJobSpec{
			InferenceServiceID: "service",
		}}},
		bisFixtures: []api_types.InferenceService{{
			ID: "service",
		}},
	},
}

func TestDelete(t *testing.T) {

	r := postgres.BIJRepo{DB: db}
	rs := postgres.BISRepo{DB: db}
	for _, test := range testDeleteCases {
		t.Run(fmt.Sprintf("TestDelete#%s", test.testName), func(t *testing.T) {
			req := require.New(t)
			defer func() {
				req.NoError(cleanupJobsServices())
			}()

			// Create fixtures if required
			for _, bis := range test.bisFixtures {
				err := rs.Create(context.TODO(), nil, bis)
				req.NoError(err)
			}
			for _, bij := range test.fixtures {
				err := r.Create(context.TODO(), nil, bij)
				req.NoError(err)
			}

			// Try to delete
			err := r.Delete(context.TODO(), nil, test.bijID)

			// Check expectation about returned error
			if len(test.expectedErrString) > 0 {
				req.EqualError(err, test.expectedErrString)
			} else {
				// If deletion was successful then we expect not found while retrieving
				req.NoError(err)
				_, err := r.Get(context.TODO(), nil, test.bijID)
				req.EqualError(err, "entity \"job\" is not found")
			}
		})
	}
}

func job(jobID string) api_types.InferenceJob {
	return api_types.InferenceJob{
		ID:   jobID,
		Spec: api_types.InferenceJobSpec{
			InferenceServiceID: "service",
		},
	}
}

var testListCases = []struct {
	testName string
	bisFixtures []api_types.InferenceService
	fixtures []api_types.InferenceJob // Jobs that should be created in DB before test logic
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
		fixtures:        []api_types.InferenceJob{job("job"), job("job2"), job("job3")},
		bisFixtures: []api_types.InferenceService{{ID: "service"}},
		expectedLen: 3,
	},
	{
		testName: "first page with size=1",
		fixtures:        []api_types.InferenceJob{job("job"), job("job2"), job("job3")},
		bisFixtures: []api_types.InferenceService{{ID: "service"}},
		expectedLen: 1,
		filter: []filter.ListOption{filter.Page(0), filter.Size(1)},
	},
	{
		testName: "second page with size=1",
		fixtures:        []api_types.InferenceJob{job("job"), job("job2"), job("job3")},
		bisFixtures: []api_types.InferenceService{{ID: "service"}},
		expectedLen: 1,
		filter: []filter.ListOption{filter.Page(1), filter.Size(1)},
	},
	{
		testName: "forth page with size=1",
		fixtures:        []api_types.InferenceJob{job("job"), job("job2"), job("job3")},
		bisFixtures: []api_types.InferenceService{{ID: "service"}},
		expectedLen: 0,
		filter: []filter.ListOption{filter.Page(3), filter.Size(1)},
	},
	{
		testName: "first page with size=2",
		fixtures:        []api_types.InferenceJob{job("job"), job("job2"), job("job3")},
		bisFixtures: []api_types.InferenceService{{ID: "service"}},
		expectedLen: 2,  // all except last one
		filter: []filter.ListOption{filter.Page(0), filter.Size(2)},
	},
}

func TestList(t *testing.T) {

	r := postgres.BIJRepo{DB: db}
	rs := postgres.BISRepo{DB: db}
	for _, test := range testListCases {
		t.Run(fmt.Sprintf("TestList#%s", test.testName), func(t *testing.T) {
			req := require.New(t)
			defer func() {
				req.NoError(cleanupJobsServices())
			}()

			// Create fixtures if required
			for _, bis := range test.bisFixtures {
				err := rs.Create(context.TODO(), nil, bis)
				req.NoError(err)
			}

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


var testGetCases = []struct {
	testName string
	bijID string
	fixtures []api_types.InferenceJob // Jobs that should be created in DB before test logic
	bisFixtures []api_types.InferenceService
	expectedErrString string
}{
	{
		testName: "not found",
		bijID: "job",
		fixtures:          nil,
		expectedErrString: "entity \"job\" is not found",
	},
	{
		testName: "ok",
		bijID: "job",
		fixtures:         []api_types.InferenceJob{{ID: "job", Spec: api_types.InferenceJobSpec{
			InferenceServiceID: "service",
		}}},
		bisFixtures: []api_types.InferenceService{{ID: "service"}},
	},
}

func TestGet(t *testing.T) {

	r := postgres.BIJRepo{DB: db}
	rs := postgres.BISRepo{DB: db}
	for _, test := range testGetCases {
		t.Run(fmt.Sprintf("TestGet#%s", test.testName), func(t *testing.T) {
			req := require.New(t)
			defer func() {
				req.NoError(cleanupJobsServices())
			}()

			// Create fixtures if required
			for _, bis := range test.bisFixtures {
				err := rs.Create(context.TODO(), nil, bis)
				req.NoError(err)
			}

			for _, bij := range test.fixtures {
				err := r.Create(context.TODO(), nil, bij)
				req.NoError(err)
			}

			// Try to Get
			_, err := r.Get(context.TODO(), nil, test.bijID)

			// Check expectation about returned error
			if len(test.expectedErrString) > 0 {
				req.EqualError(err, test.expectedErrString)
			} else {
				req.NoError(err)
			}
		})
	}
}


var testSetDeletionMarkCases = []struct {
	testName string
	bijID string
	deletionMark bool
	fixtures []api_types.InferenceJob // Jobs that should be created in DB before test logic
	bisFixtures []api_types.InferenceService
	expectedErrString string
}{
	{
		testName: "not found",
		bijID: "job",
		fixtures:          nil,
		expectedErrString: "entity \"job\" is not found",
	},
	{
		testName: "set true",
		bijID: "job",
		deletionMark: true,
		fixtures:         []api_types.InferenceJob{{ID: "job", DeletionMark: false, Spec: api_types.InferenceJobSpec{
			InferenceServiceID: "service",
		}}},
		bisFixtures: []api_types.InferenceService{{ID: "service"}},
	},
	{
		testName: "set false",
		bijID: "job",
		deletionMark: false,
		fixtures:         []api_types.InferenceJob{{ID: "job", DeletionMark: true, Spec: api_types.InferenceJobSpec{
			InferenceServiceID: "service",
		}}},
		bisFixtures: []api_types.InferenceService{{ID: "service"}},
	},
}

func TestSetDeletionMark(t *testing.T) {

	r := postgres.BIJRepo{DB: db}
	rs := postgres.BISRepo{DB: db}
	for _, test := range testSetDeletionMarkCases {
		t.Run(fmt.Sprintf("TestSetDeletionMark#%s", test.testName), func(t *testing.T) {
			req := require.New(t)
			defer func() {
				req.NoError(cleanupJobsServices())
			}()

			// Create fixtures if required
			for _, bis := range test.bisFixtures {
				err := rs.Create(context.TODO(), nil, bis)
				req.NoError(err)
			}

			for _, bij := range test.fixtures {
				err := r.Create(context.TODO(), nil, bij)
				req.NoError(err)
			}

			// Try to set deletion mark
			err := r.SetDeletionMark(context.TODO(), nil, test.bijID, test.deletionMark)

			// Check expectation about returned error
			if len(test.expectedErrString) > 0 {
				req.EqualError(err, test.expectedErrString)
			} else {
				// If deletion mark was successful then we expect to see new flag
				req.NoError(err)
				res, err := r.Get(context.TODO(), nil, test.bijID)
				req.NoError(err)
				req.Equal(res.DeletionMark, test.deletionMark)
			}
		})
	}
}


// Utils

func cleanupJobsServices() error {
	query, _, err := sq.Delete(postgres.BatchInferenceJobTable).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return err
	}
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	query, _, err = sq.Delete(postgres.BatchInferenceServiceTable).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return err
	}
	_, err = db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
