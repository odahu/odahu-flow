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
	api_types "github.com/odahu/odahu-flow/packages/operator/pkg/apis/batch"
	odahuErrs "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"
	"time"
)

type InferenceServiceRepo interface {
	Create(ctx context.Context, tx *sql.Tx, bis api_types.InferenceService) (err error)
	Get(ctx context.Context, tx *sql.Tx, id string) (res api_types.InferenceService, err error)
	Update(ctx context.Context, tx *sql.Tx, id string, bis api_types.InferenceService) (err error)
	List(ctx context.Context, tx *sql.Tx, options ...filter.ListOption) (res []api_types.InferenceService, err error)
	Delete(ctx context.Context, tx *sql.Tx, id string) (err error)
}

type InferenceServiceService struct {
	repo InferenceServiceRepo
}

func NewInferenceServiceService(repo InferenceServiceRepo) *InferenceServiceService {
	return &InferenceServiceService{repo: repo}
}

// Create creates api_types.InferenceService
func (s *InferenceServiceService) Create(
	ctx context.Context, bis *api_types.InferenceService) (err error) {

	// Set fields that managed by platform. Cannot be overridden by user
	bis.DeletionMark = false
	bis.CreatedAt = time.Now().UTC()
	bis.UpdatedAt = time.Now().UTC()
	bis.Status = api_types.InferenceServiceStatus{}

	// Defaulting
	DefaultCreate(bis)

	// Validation
	if errs := ValidateCreateUpdate(*bis); len(errs) > 0 {
		return odahuErrs.InvalidEntityError{
			Entity:           bis.ID,
			ValidationErrors: errs,
		}
	}

	err = s.repo.Create(ctx, nil, *bis)
	return err
}

// Update updates api_types.InferenceService
func (s *InferenceServiceService) Update(
	ctx context.Context, id string, bis *api_types.InferenceService) (err error) {

	bis.DeletionMark = false
	bis.UpdatedAt = time.Now().UTC()

	// Validation
	if errs := ValidateCreateUpdate(*bis); len(errs) > 0 {
		return odahuErrs.InvalidEntityError{
			Entity:           bis.ID,
			ValidationErrors: errs,
		}
	}

	old, err := s.repo.Get(ctx, nil, id)
	if err != nil {
		return err
	}
	bis.CreatedAt = old.CreatedAt

	return s.repo.Update(ctx, nil, id, *bis)
}

func (s *InferenceServiceService) Delete(ctx context.Context, id string) (err error) {
	return s.repo.Delete(ctx, nil, id)
}

func (s *InferenceServiceService) Get(ctx context.Context, id string) (res api_types.InferenceService, err error) {
	return s.repo.Get(ctx, nil, id)
}

func (s *InferenceServiceService) List(
	ctx context.Context, options ...filter.ListOption) (res []api_types.InferenceService, err error) {
	return s.repo.List(ctx, nil, options...)
}
