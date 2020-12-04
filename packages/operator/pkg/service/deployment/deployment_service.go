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
	"fmt"
	"github.com/google/uuid"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/deployment"
	odahu_errors "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	repo "github.com/odahu/odahu-flow/packages/operator/pkg/repository/deployment"
	mrRepo "github.com/odahu/odahu-flow/packages/operator/pkg/repository/route"
	db_utils "github.com/odahu/odahu-flow/packages/operator/pkg/utils/db"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"
	hashutil "github.com/odahu/odahu-flow/packages/operator/pkg/utils/hash"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"time"
)

var (
	log = logf.Log.WithName("model-deployment--service")
	defaultWeight            = int32(100)
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
	GetDefaultModelRoute(ctx context.Context, mdID string) (*deployment.ModelRoute, error)
}

type serviceImpl struct {
	// Repository that has "database/sql" underlying storage
	repo repo.Repository
	mrRepo mrRepo.Repository
}

func (s serviceImpl) GetModelDeployment(ctx context.Context, id string) (*deployment.ModelDeployment, error) {
	return s.repo.GetModelDeployment(ctx, nil, id)
}

func (s serviceImpl) GetModelDeploymentList(
	ctx context.Context, options ...filter.ListOption,
) ([]deployment.ModelDeployment, error) {
	return s.repo.GetModelDeploymentList(ctx, nil, options...)
}


func GetDefaultModelRoute(ctx context.Context, tx *sql.Tx, mdID string, repository mrRepo.Repository) (string, error) {

	mrs, err := repository.GetModelRouteList(ctx, tx, filter.ListFilter(&mrRepo.Filter{
		Default: []bool{true},
		MdID: []string{mdID},
	}))
	if err != nil {
		return "", err
	}
	if len(mrs) > 1 {
		return "", fmt.Errorf("model deployment must have one default route, but have %v", len(mrs))
	}
	if len(mrs) == 0 {
		return "", nil
	}
	return mrs[0].ID, nil
}

func (s serviceImpl) GetDefaultModelRoute(ctx context.Context, mdID string) (*deployment.ModelRoute, error) {
	tx, err := s.mrRepo.BeginTransaction(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		db_utils.FinishTx(tx, err, log)
	}()
	id, err := GetDefaultModelRoute(ctx, tx, mdID, s.mrRepo)
	if err != nil {
		return nil, err
	}
	return s.mrRepo.GetModelRoute(ctx, tx, id)
}

func (s serviceImpl) DeleteModelDeployment(ctx context.Context, id string) (err error) {
	tx, err := s.repo.BeginTransaction(ctx)
	if err != nil {
		return err
	}
	defer func() {
		db_utils.FinishTx(tx, err, log)
	}()

	mrDefaultID, err := GetDefaultModelRoute(ctx, tx, id, s.mrRepo)
	if err != nil {
		return err
	}

	if mrDefaultID != "" {
		err = s.mrRepo.DeleteModelRoute(ctx, tx, mrDefaultID)
		if err != nil {
			return
		}
	}

	err = s.repo.DeleteModelDeployment(ctx, tx, id)
	return
}

func (s serviceImpl) SetDeletionMark(ctx context.Context, id string, value bool) error {
	return s.repo.SetDeletionMark(ctx, nil, id, value)
}

func (s serviceImpl) UpdateModelDeployment(ctx context.Context, md *deployment.ModelDeployment) error {
	md.UpdatedAt = time.Now()
	md.DeletionMark = false
	md.Status = v1alpha1.ModelDeploymentStatus{}
	return s.repo.UpdateModelDeployment(ctx, nil, md)
}

func (s serviceImpl) UpdateModelDeploymentStatus(
	ctx context.Context, id string, status v1alpha1.ModelDeploymentStatus, spec v1alpha1.ModelDeploymentSpec,
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


func constructDefaultRoute(mdID string) deployment.ModelRoute {
	return deployment.ModelRoute{
		ID:           mdID + "-" + uuid.New().String()[:5],
		Default:      true,
		DeletionMark: false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Spec:         v1alpha1.ModelRouteSpec{
			URLPrefix:              fmt.Sprintf("/model/%s", mdID),
			ModelDeploymentTargets: []v1alpha1.ModelDeploymentTarget{
				{
					Name:   mdID,
					Weight: &defaultWeight,
				},
			},
		},
	}
}

func (s serviceImpl) CreateModelDeployment(ctx context.Context, md *deployment.ModelDeployment) (err error) {

	var tx *sql.Tx

	tx, err = s.repo.BeginTransaction(ctx)
	if err != nil {
		return err
	}
	defer func() {
		db_utils.FinishTx(tx, err, log)
	}()

	md.CreatedAt = time.Now()
	md.UpdatedAt = time.Now()
	md.DeletionMark = false
	md.Status = v1alpha1.ModelDeploymentStatus{}
	err = s.repo.CreateModelDeployment(ctx, tx, md)
	if err != nil {
		return
	}


	exists, err := s.mrRepo.DefaultExists(ctx, md.ID, tx)
	if err != nil || exists {
		return err
	}
	// Every Model deployment must have a default HTTP route that sends 100% of traffic to the model
	defRoute := constructDefaultRoute(md.ID)
	err = s.mrRepo.CreateModelRoute(ctx, tx, &defRoute)
	if err != nil {
		return fmt.Errorf("unable to create default ModelRoute: %v", err)
	}


	return err
}

func NewService(repo repo.Repository, mrRepo mrRepo.Repository) Service {
	return &serviceImpl{repo: repo, mrRepo: mrRepo}
}

