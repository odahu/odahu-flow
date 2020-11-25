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

package route

import (
	"context"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	route "github.com/odahu/odahu-flow/packages/operator/pkg/apis/deployment"
	odahu_errors "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	repo "github.com/odahu/odahu-flow/packages/operator/pkg/repository/route"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"
	hashutil "github.com/odahu/odahu-flow/packages/operator/pkg/utils/hash"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"time"
)

var (
	log = logf.Log.WithName("model-route--service")
)

type Service interface {
	GetModelRoute(ctx context.Context, id string) (*route.ModelRoute, error)
	GetModelRouteList(ctx context.Context, options ...filter.ListOption) ([]route.ModelRoute, error)
	DeleteModelRoute(ctx context.Context, id string) error
	SetDeletionMark(ctx context.Context, id string, value bool) error
	UpdateModelRoute(ctx context.Context, mt *route.ModelRoute) error
	// Try to update status. If spec in storage differs from spec snapshot then update does not happen
	UpdateModelRouteStatus(
		ctx context.Context, id string, status v1alpha1.ModelRouteStatus, spec v1alpha1.ModelRouteSpec) error
	CreateModelRoute(ctx context.Context, mt *route.ModelRoute) error
}

type serviceImpl struct {
	// Repository that has "database/sql" underlying storage
	repo repo.Repository
}

func (s serviceImpl) GetModelRoute(ctx context.Context, id string) (*route.ModelRoute, error) {
	return s.repo.GetModelRoute(ctx, nil, id)
}

func (s serviceImpl) GetModelRouteList(
	ctx context.Context, options ...filter.ListOption,
) ([]route.ModelRoute, error) {
	return s.repo.GetModelRouteList(ctx, nil, options...)
}

func (s serviceImpl) DeleteModelRoute(ctx context.Context, id string) error {
	return s.repo.DeleteModelRoute(ctx, nil, id)
}

func (s serviceImpl) SetDeletionMark(ctx context.Context, id string, value bool) error {
	return s.repo.SetDeletionMark(ctx, nil, id, value)
}

func (s serviceImpl) UpdateModelRoute(ctx context.Context, md *route.ModelRoute) error {
	md.UpdatedAt = time.Now()
	md.DeletionMark = false
	md.Status = v1alpha1.ModelRouteStatus{
	}
	return s.repo.UpdateModelRoute(ctx, nil, md)
}

func (s serviceImpl) UpdateModelRouteStatus(
	ctx context.Context, id string, status v1alpha1.ModelRouteStatus, spec v1alpha1.ModelRouteSpec,
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

	oldMt, err := s.repo.GetModelRoute(ctx, tx, id)
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

	err = s.repo.UpdateModelRouteStatus(ctx, tx, id, status)
	if err != nil {
		return err
	}

	return err
}

func (s serviceImpl) CreateModelRoute(ctx context.Context, md *route.ModelRoute) error {
	md.CreatedAt = time.Now()
	md.UpdatedAt = time.Now()
	md.DeletionMark = false
	md.Status = v1alpha1.ModelRouteStatus{}
	return s.repo.CreateModelRoute(ctx, nil, md)
}

func NewService(repo repo.Repository) Service {
	return &serviceImpl{repo: repo}
}

