package training

import (
	"context"
	"database/sql"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/training"
	odahu_errors "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"
	repo"github.com/odahu/odahu-flow/packages/operator/pkg/repository/training"
	hashutil "github.com/odahu/odahu-flow/packages/operator/pkg/utils/hash"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	TagKey = "name"
)

var (
	txOptions = &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  false,
	}
	log = logf.Log.WithName("model-training--service")
)

type Service interface {
	GetModelTraining(id string) (*training.ModelTraining, error)
	GetModelTrainingList(options ...filter.ListOption) ([]training.ModelTraining, error)
	DeleteModelTraining(id string) error
	SetDeletionMark(id string, value bool) error
	UpdateModelTraining(mt *training.ModelTraining) error
	// Try to update status. If spec in storage differs from spec snapshot then update does not happen
	UpdateModelTrainingStatus(id string, status v1alpha1.ModelTrainingStatus, spec v1alpha1.ModelTrainingSpec) error
	CreateModelTraining(mt *training.ModelTraining) error
}

type MTFilter struct {
	Toolchain    []string `name:"toolchain" postgres:"spec->>'toolchain'"`
	ModelName    []string `name:"model_name" postgres:"spec->'model'->>'name'"`
	ModelVersion []string `name:"model_version" postgres:"spec->'model'->>'version'"`
}

type serviceImpl struct {
	db   *sql.DB
	// Repository that has "database/sql" underlying storage
	repo repo.Repository
}

func (s serviceImpl) GetModelTraining(id string) (*training.ModelTraining, error) {
	ctx := context.TODO()
	return s.repo.GetModelTraining(ctx, s.db, id)
}

func (s serviceImpl) GetModelTrainingList(options ...filter.ListOption) ([]training.ModelTraining, error) {
	return s.repo.GetModelTrainingList(nil, nil)
}

func (s serviceImpl) DeleteModelTraining(id string) error {
	return s.repo.DeleteModelTraining(nil, nil, id)
}

func (s serviceImpl) SetDeletionMark(id string, value bool) error {
	return s.repo.SetDeletionMark(nil, nil, id, value)
}

func (s serviceImpl) UpdateModelTraining(mt *training.ModelTraining) error {
	return s.repo.UpdateModelTraining(nil, nil, mt)
}

func (s serviceImpl) UpdateModelTrainingStatus(id string, status v1alpha1.ModelTrainingStatus, spec v1alpha1.ModelTrainingSpec) error {

	ctx := context.TODO()

	tx, err := s.db.BeginTx(ctx, txOptions)
	defer func() {
		err := tx.Rollback()
		if err != nil {
			log.Error(err, "Error while rollback transaction")
		}
	}()

	if err != nil {
		return err
	}

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
		return odahu_errors.SpecWasChangedError{Entity: id}
	}

	err = s.repo.UpdateModelTrainingStatus(ctx, tx, id, status)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s serviceImpl) CreateModelTraining(mt *training.ModelTraining) error {
	return s.repo.CreateModelTraining(nil, nil, mt)
}

func NewService(repo repo.Repository, db *sql.DB) Service {
	return &serviceImpl{repo: repo, db: db}
}

