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

package training

import (
	"context"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/training"
	odahu_errors "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	repo "github.com/odahu/odahu-flow/packages/operator/pkg/repository/training"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"
	hashutil "github.com/odahu/odahu-flow/packages/operator/pkg/utils/hash"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"time"
)

var (
	log = logf.Log.WithName("model-training--service")
)

type Service interface {
	GetModelTraining(ctx context.Context, id string) (*training.ModelTraining, error)
	GetModelTrainingList(ctx context.Context, options ...filter.ListOption) ([]training.ModelTraining, error)
	DeleteModelTraining(ctx context.Context, id string) error
	SetDeletionMark(ctx context.Context, id string, value bool) error
	UpdateModelTraining(ctx context.Context, mt *training.ModelTraining) error
	// Try to update status. If spec in storage differs from spec snapshot then update does not happen
	UpdateModelTrainingStatus(
		ctx context.Context, id string, status v1alpha1.ModelTrainingStatus, spec v1alpha1.ModelTrainingSpec) error
	CreateModelTraining(ctx context.Context, mt *training.ModelTraining) error
}

type serviceImpl struct {
	// Repository that has "database/sql" underlying storage
	repo repo.Repository
}

func (s serviceImpl) GetModelTraining(ctx context.Context, id string) (*training.ModelTraining, error) {
	return s.repo.GetModelTraining(ctx, nil, id)
}

func (s serviceImpl) GetModelTrainingList(
	ctx context.Context, options ...filter.ListOption,
) ([]training.ModelTraining, error) {
	return s.repo.GetModelTrainingList(ctx, nil, options...)
}

func (s serviceImpl) DeleteModelTraining(ctx context.Context, id string) error {
	return s.repo.DeleteModelTraining(ctx, nil, id)
}

func (s serviceImpl) SetDeletionMark(ctx context.Context, id string, value bool) error {
	return s.repo.SetDeletionMark(ctx, nil, id, value)
}

func (s serviceImpl) UpdateModelTraining(ctx context.Context, mt *training.ModelTraining) error {
	mt.UpdatedAt = time.Now()
	oldMt, err := s.GetModelTraining(ctx, mt.ID)
	if err != nil {
		return err
	}
	mt.CreatedAt = oldMt.CreatedAt
	mt.DeletionMark = false
	mt.Status = v1alpha1.ModelTrainingStatus{
		State: v1alpha1.ModelTrainingUnknown,
	}
	return s.repo.UpdateModelTraining(ctx, nil, mt)
}

func (s serviceImpl) UpdateModelTrainingStatus(
	ctx context.Context, id string, status v1alpha1.ModelTrainingStatus, spec v1alpha1.ModelTrainingSpec,
) (err error) {

	log.Info("!!!DEBUG2!!!")

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

	oldMt, err := s.repo.GetModelTraining(ctx, tx, id)
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

	err = s.repo.UpdateModelTrainingStatus(ctx, tx, id, status)
	if err != nil {
		return err
	}

	log.Info("!!!DEBUG3!!!")

	return err
}

func (s serviceImpl) CreateModelTraining(ctx context.Context, mt *training.ModelTraining) error {
	mt.CreatedAt = time.Now()
	mt.UpdatedAt = time.Now()
	mt.DeletionMark = false
	mt.Status = v1alpha1.ModelTrainingStatus{
		State: v1alpha1.ModelTrainingUnknown,
	}
	return s.repo.SaveModelTraining(ctx, nil, mt)
}

func NewService(repo repo.Repository) Service {
	return &serviceImpl{repo: repo}
}
