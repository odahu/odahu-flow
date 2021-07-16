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

package packaging

import (
	"context"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/packaging"
	odahu_errors "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	repo "github.com/odahu/odahu-flow/packages/operator/pkg/repository/packaging"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"
	hashutil "github.com/odahu/odahu-flow/packages/operator/pkg/utils/hash"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"time"
)

var (
	log = logf.Log.WithName("model-packaging--service")
)

type Service interface {
	GetModelPackaging(ctx context.Context, id string) (*packaging.ModelPackaging, error)
	GetModelPackagingList(ctx context.Context, options ...filter.ListOption) ([]packaging.ModelPackaging, error)
	DeleteModelPackaging(ctx context.Context, id string) error
	SetDeletionMark(ctx context.Context, id string, value bool) error
	UpdateModelPackaging(ctx context.Context, mt *packaging.ModelPackaging) error
	// Try to update status. If spec in storage differs from spec snapshot then update does not happen
	UpdateModelPackagingStatus(
		ctx context.Context, id string, status v1alpha1.ModelPackagingStatus, spec packaging.ModelPackagingSpec) error
	CreateModelPackaging(ctx context.Context, mt *packaging.ModelPackaging) error
}

type serviceImpl struct {
	// Repository that has "database/sql" underlying storage
	repo repo.Repository
}

func (s serviceImpl) GetModelPackaging(ctx context.Context, id string) (*packaging.ModelPackaging, error) {
	return s.repo.GetModelPackaging(ctx, nil, id)
}

func (s serviceImpl) GetModelPackagingList(
	ctx context.Context, options ...filter.ListOption,
) ([]packaging.ModelPackaging, error) {
	return s.repo.GetModelPackagingList(ctx, nil, options...)
}

func (s serviceImpl) DeleteModelPackaging(ctx context.Context, id string) error {
	return s.repo.DeleteModelPackaging(ctx, nil, id)
}

func (s serviceImpl) SetDeletionMark(ctx context.Context, id string, value bool) error {
	return s.repo.SetDeletionMark(ctx, nil, id, value)
}

func (s serviceImpl) UpdateModelPackaging(ctx context.Context, mp *packaging.ModelPackaging) error {
	mp.UpdatedAt = time.Now()
	oldMp, err := s.GetModelPackaging(ctx, mp.ID)
	if err != nil {
		return err
	}
	mp.CreatedAt = oldMp.CreatedAt
	mp.DeletionMark = false
	mp.Status = v1alpha1.ModelPackagingStatus{
		State: v1alpha1.ModelPackagingUnknown,
	}
	return s.repo.UpdateModelPackaging(ctx, nil, mp)
}

func (s serviceImpl) UpdateModelPackagingStatus(
	ctx context.Context, id string, status v1alpha1.ModelPackagingStatus, spec packaging.ModelPackagingSpec,
) (err error) {

	tx, err := s.repo.BeginTransaction(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err == nil {
			if err := tx.Commit(); err != nil {
				log.Error(err, "Error while commit transaction")
			}
		} else {
			if err := tx.Rollback(); err != nil {
				log.Error(err, "Error while rollback transaction")
			}
		}
	}()

	oldMt, err := s.repo.GetModelPackaging(ctx, tx, id)
	if err != nil {
		return err
	}

	oldHash, err := hashutil.Hash(oldMt.Spec)
	if err != nil {
		return err
	}

	specHash, err := hashutil.Hash(spec)
	if err != nil {
		return err
	}

	if oldHash != specHash {
		return odahu_errors.SpecWasTouched{Entity: id}
	}

	err = s.repo.UpdateModelPackagingStatus(ctx, tx, id, status)
	if err != nil {
		return err
	}

	return err
}

func (s serviceImpl) CreateModelPackaging(ctx context.Context, mp *packaging.ModelPackaging) error {
	mp.CreatedAt = time.Now()
	mp.UpdatedAt = time.Now()
	mp.DeletionMark = false
	mp.Status = v1alpha1.ModelPackagingStatus{
		State: v1alpha1.ModelPackagingUnknown,
	}
	return s.repo.SaveModelPackaging(ctx, nil, mp)
}

func NewService(repo repo.Repository) Service {
	return &serviceImpl{repo: repo}
}
