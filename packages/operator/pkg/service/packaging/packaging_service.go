package packaging

import (
	"context"
	"database/sql"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/packaging"
	odahu_errors "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"
	repo "github.com/odahu/odahu-flow/packages/operator/pkg/repository/packaging"
	hashutil "github.com/odahu/odahu-flow/packages/operator/pkg/utils/hash"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var (
	txOptions = &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  false,
	}
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
	db   *sql.DB
	// Repository that has "database/sql" underlying storage
	repo repo.Repository
}

func (s serviceImpl) GetModelPackaging(ctx context.Context, id string) (*packaging.ModelPackaging, error) {
	return s.repo.GetModelPackaging(ctx, s.db, id)
}

func (s serviceImpl) GetModelPackagingList(
	ctx context.Context, options ...filter.ListOption,
) ([]packaging.ModelPackaging, error) {
	return s.repo.GetModelPackagingList(ctx, s.db, options...)
}

func (s serviceImpl) DeleteModelPackaging(ctx context.Context, id string) error {
	return s.repo.DeleteModelPackaging(ctx, s.db, id)
}

func (s serviceImpl) SetDeletionMark(ctx context.Context, id string, value bool) error {
	return s.repo.SetDeletionMark(ctx, s.db, id, value)
}

func (s serviceImpl) UpdateModelPackaging(ctx context.Context, mt *packaging.ModelPackaging) error {
	return s.repo.UpdateModelPackaging(ctx, s.db, mt)
}

func (s serviceImpl) UpdateModelPackagingStatus(
	ctx context.Context, id string, status v1alpha1.ModelPackagingStatus, spec packaging.ModelPackagingSpec,
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

func (s serviceImpl) CreateModelPackaging(ctx context.Context, mt *packaging.ModelPackaging) error {
	return s.repo.CreateModelPackaging(ctx, s.db, mt)
}

func NewService(repo repo.Repository, db *sql.DB) Service {
	return &serviceImpl{repo: repo, db: db}
}

