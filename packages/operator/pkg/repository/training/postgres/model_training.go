package postgres

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/training"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	odahuErrors "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	utils "github.com/odahu/odahu-flow/packages/operator/pkg/repository/util/postgres"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	ModelTrainingTable = "odahu_operator_training"
)

var (
	log = logf.Log.WithName("model-training--repository--postgres")
)

type TrainingRepo struct {
	DB *sql.DB
}

func (repo TrainingRepo) GetModelTraining(name string) (*training.ModelTraining, error) {

	mt := new(training.ModelTraining)

	err := repo.DB.QueryRow(
		fmt.Sprintf("SELECT id, spec, status, deletionmark FROM %s WHERE id = $1", ModelTrainingTable),
		name,
	).Scan(&mt.ID, &mt.Spec, &mt.Status, &mt.DeletionMark)

	switch {
	case err == sql.ErrNoRows:

		return nil, odahuErrors.NotFoundError{Entity: name}
	case err != nil:
		log.Error(err, "error during sql query")

		return nil, err
	default:

		return mt, nil
	}

}

func (repo TrainingRepo) GetModelTrainingList(options ...filter.ListOption) (
	[]training.ModelTraining, error,
) {

	listOptions := &filter.ListOptions{
		Filter: nil,
		Page:   &FirstPage,
		Size:   &MaxSize,
	}
	for _, option := range options {
		option(listOptions)
	}

	offset := *listOptions.Size * (*listOptions.Page)

	sb := sq.Select("id, spec, status, deletionmark").From("odahu_operator_training").
		OrderBy("id").
		Offset(uint64(offset)).
		Limit(uint64(*listOptions.Size)).PlaceholderFormat(sq.Dollar)

	sb = utils.TransformFilter(sb, listOptions.Filter)
	stmt, args, err := sb.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := repo.DB.Query(stmt, args...)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Error(err, "error during rows.Close()")
		}
		if err := rows.Err(); err != nil {
			log.Error(err, "error during rows iterating")
		}
	}()

	mts := make([]training.ModelTraining, 0)

	for rows.Next() {
		mt := new(training.ModelTraining)
		err := rows.Scan(&mt.ID, &mt.Spec, &mt.Status, &mt.DeletionMark)
		if err != nil {
			return nil, err
		}
		mts = append(mts, *mt)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return mts, nil

}

func (repo TrainingRepo) DeleteModelTraining(name string) error {

	// First try to check that row exists otherwise raise exception to fit interface
	_, err := repo.GetModelTraining(name)
	if err != nil {
		return err
	}

	// If exists, delete it

	sqlStatement := fmt.Sprintf("DELETE FROM %s WHERE id = $1", ModelTrainingTable)
	_, err = repo.DB.Exec(sqlStatement, name)
	if err != nil {
		return err
	}
	return nil
}

func (repo TrainingRepo) SetDeletionMark(id string, value bool) error {
	return utils.SetDeletionMark(repo.DB, ModelTrainingTable, id, value)
}

func (repo TrainingRepo) UpdateModelTraining(mt *training.ModelTraining) error {

	mt.Status.State = ""

	stmt, args, err := sq.Update(ModelTrainingTable).
		Set("spec", mt.Spec).
		Set("status", mt.Status).
		Where(sq.Eq{"id": mt.ID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	result, err := repo.DB.Exec(stmt, args...)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return odahuErrors.NotFoundError{Entity: mt.ID}
	}

	return nil
}

func (repo TrainingRepo) UpdateModelTrainingStatus(id string, s v1alpha1.ModelTrainingStatus) error {

	stmt, args, err := sq.Update(ModelTrainingTable).
		Set("status", s).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	result, err := repo.DB.Exec(stmt, args...)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return odahuErrors.NotFoundError{Entity: id}
	}

	return nil
}

func (repo TrainingRepo) CreateModelTraining(mt *training.ModelTraining) error {

	_, err := repo.DB.Exec(
		fmt.Sprintf("INSERT INTO %s (id, spec, status) VALUES($1, $2, $3)", ModelTrainingTable),
		mt.ID, mt.Spec, mt.Status,
	)
	if err != nil {
		pqError, ok := err.(*pq.Error)
		if ok && pqError.Code == uniqueViolationPostgresCode {
			return odahuErrors.AlreadyExistError{Entity: mt.ID}
		}
		return err
	}
	return nil

}
