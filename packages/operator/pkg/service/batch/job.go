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

package batch

import (
	"context"
	"database/sql"
	"fmt"
	api_types "github.com/odahu/odahu-flow/packages/operator/pkg/apis/batch"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/connection"
	odahuErrors "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"
	"time"
)

type JobRepository interface {
	Create(ctx context.Context, tx *sql.Tx, bij api_types.InferenceJob) (err error)
	UpdateStatus(ctx context.Context, tx *sql.Tx, id string, s api_types.InferenceJobStatus) (err error)
	Delete(ctx context.Context, tx *sql.Tx, id string) (err error)
	List(ctx context.Context, tx *sql.Tx, options ...filter.ListOption) (res []api_types.InferenceJob, err error)
	Get(ctx context.Context, tx *sql.Tx, id string) (res api_types.InferenceJob, err error)
	SetDeletionMark(ctx context.Context, tx *sql.Tx, id string, value bool) error
}

type ServiceRepository interface {
	Get(ctx context.Context, tx *sql.Tx, id string) (res api_types.InferenceService, err error)
}

type ConnectionGetter interface {
	GetConnection(id string, encrypted bool) (*connection.Connection, error)
}

type JobService struct {
	repo JobRepository
	sRepo ServiceRepository
	connGetter ConnectionGetter
}

func NewJobService(repo JobRepository, sRepo ServiceRepository, connGetter ConnectionGetter) *JobService {
	return &JobService{
		repo:       repo,
		sRepo:      sRepo,
		connGetter: connGetter,
	}
}

// Create launches BatchInferenceJob
// Because we ensure immutability of jobs we also generate ID to take this responsibility from client
// Generated ID should be returned to client
func (s *JobService) Create(ctx context.Context, bij *api_types.InferenceJob) (err error) {

	bij.CreatedAt = time.Now().UTC()
	bij.UpdatedAt = time.Now().UTC()

	if errs := ValidateJobInput(*bij); len(errs) > 0 {
		return odahuErrors.InvalidEntityError{
			Entity:           bij.ID,
			ValidationErrors: errs,
		}
	}

	service, err := s.sRepo.Get(ctx, nil, bij.Spec.InferenceServiceID)
	if err != nil {
		if odahuErrors.IsNotFoundError(err) {
			return odahuErrors.InvalidEntityError{
				Entity:           "job",
				ValidationErrors: []error{fmt.Errorf("unable to fetch corresponding service: %s", err)},
			}
		}
		return err
	}

	DefaultJob(bij, service)

	errs, err := ValidateJob(*bij, s.connGetter, service)
	if err != nil {
		return err
	}
	if len(errs) > 0 {
		return odahuErrors.InvalidEntityError{
			Entity:           bij.ID,
			ValidationErrors: errs,
		}
	}

	err = s.repo.Create(ctx, nil, *bij)

	return err
}

func (s *JobService) UpdateStatus(ctx context.Context, id string, status api_types.InferenceJobStatus) error {

	return s.repo.UpdateStatus(ctx, nil, id, status)
}

func (s *JobService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, nil, id)
}
func (s *JobService) List(ctx context.Context, options ...filter.ListOption) ([]api_types.InferenceJob, error) {
	return s.repo.List(ctx, nil, options...)
}

func (s *JobService) Get(ctx context.Context, id string) (api_types.InferenceJob, error) {
	return s.repo.Get(ctx, nil, id)
}

func (s *JobService) SetDeletionMark(ctx context.Context, id string) error {
	return s.repo.SetDeletionMark(ctx, nil, id, true)
}