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

package deployment

import (
	"context"
	"database/sql"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/deployment"
	odahu_errors "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"
	repo "github.com/odahu/odahu-flow/packages/operator/pkg/repository/deployment"
	hashutil "github.com/odahu/odahu-flow/packages/operator/pkg/utils/hash"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var (
	txOptions = &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  false,
	}
	log = logf.Log.WithName("model-deployment--service")
)

type Service interface {
	GetModelDeployment(ctx context.Context, id string) (*deployment.ModelDeployment, error)
	GetModelDeploymentList(ctx context.Context, options ...filter.ListOption) ([]deployment.ModelDeployment, error)
	DeleteModelDeployment(ctx context.Context, id string) error
	SetDeletionMark(ctx context.Context, id string, value bool) error
	UpdateModelDeployment(ctx context.Context, mt *deployment.ModelDeployment) error
	// Try to update status. If spec in storage differs from spec snapshot then update does not happen
	UpdateModelDeploymentStatus(
		ctx context.Context, id string, status v1alpha1.ModelDeploymentStatus, spec v1alpha1.ModelDeploymentSpec) error
	CreateModelDeployment(ctx context.Context, mt *deployment.ModelDeployment) error
}

type serviceImpl struct {
	db   *sql.DB
	// Repository that has "database/sql" underlying storage
	repo repo.Repository
}

func (s serviceImpl) GetModelDeployment(ctx context.Context, id string) (*deployment.ModelDeployment, error) {
	return s.repo.GetModelDeployment(ctx, s.db, id)
}

func (s serviceImpl) GetModelDeploymentList(
	ctx context.Context, options ...filter.ListOption,
) ([]deployment.ModelDeployment, error) {
	return s.repo.GetModelDeploymentList(ctx, s.db, options...)
}

func (s serviceImpl) DeleteModelDeployment(ctx context.Context, id string) error {
	return s.repo.DeleteModelDeployment(ctx, s.db, id)
}

func (s serviceImpl) SetDeletionMark(ctx context.Context, id string, value bool) error {
	return s.repo.SetDeletionMark(ctx, s.db, id, value)
}

func (s serviceImpl) UpdateModelDeployment(ctx context.Context, mt *deployment.ModelDeployment) error {
	return s.repo.UpdateModelDeployment(ctx, s.db, mt)
}

func (s serviceImpl) UpdateModelDeploymentStatus(
	ctx context.Context, id string, status v1alpha1.ModelDeploymentStatus, spec v1alpha1.ModelDeploymentSpec,
) (err error) {

	tx, err := s.db.BeginTx(ctx, txOptions)
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

	oldMt, err := s.repo.GetModelDeployment(ctx, tx, id)
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

	err = s.repo.UpdateModelDeploymentStatus(ctx, tx, id, status)
	if err != nil {
		return err
	}

	return err
}

func (s serviceImpl) CreateModelDeployment(ctx context.Context, mt *deployment.ModelDeployment) error {
	return s.repo.CreateModelDeployment(ctx, s.db, mt)
}

func NewService(repo repo.Repository, db *sql.DB) Service {
	return &serviceImpl{repo: repo, db: db}
}

