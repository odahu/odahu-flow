package postgres

import (
	"context"
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/packaging"
	odahuErrors "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"
	utils "github.com/odahu/odahu-flow/packages/operator/pkg/repository/util/postgres"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	ModelPackagingTable = "odahu_operator_packaging"
)

var (
	log = logf.Log.WithName("model-packaging--repository--postgres")
)

type PackagingRepo struct {
	DB *sql.DB
}

func (repo PackagingRepo) GetModelPackaging(name string) (*packaging.ModelPackaging, error) {

	mt := new(packaging.ModelPackaging)

	err := repo.DB.QueryRow(
		fmt.Sprintf("SELECT id, spec, status, deletionmark FROM %s WHERE id = $1", ModelPackagingTable),
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

func (repo PackagingRepo) GetModelPackagingList(options ...filter.ListOption) (
	[]packaging.ModelPackaging, error,
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

	sb := sq.Select("id, spec, status, deletionmark").From("odahu_operator_packaging").
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

	mts := make([]packaging.ModelPackaging, 0)

	for rows.Next() {
		mt := new(packaging.ModelPackaging)
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

func (repo PackagingRepo) DeleteModelPackaging(name string) error {

	// First try to check that row exists otherwise raise exception to fit interface
	_, err := repo.GetModelPackaging(name)
	if err != nil {
		return err
	}

	// If exists, delete it

	sqlStatement := fmt.Sprintf("DELETE FROM %s WHERE id = $1", ModelPackagingTable)
	_, err = repo.DB.Exec(sqlStatement, name)
	if err != nil {
		return err
	}
	return nil
}

func (repo PackagingRepo) SetDeletionMark(id string, value bool) error {
	return utils.SetDeletionMark(context.TODO(), repo.DB, ModelPackagingTable, id, value)
}

func (repo PackagingRepo) UpdateModelPackaging(mp *packaging.ModelPackaging) error {

	mp.Status.State = ""

	stmt, args, err := sq.Update(ModelPackagingTable).
		Set("spec", mp.Spec).
		Set("status", mp.Status).
		Where(sq.Eq{"id": mp.ID}).
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
		return odahuErrors.NotFoundError{Entity: mp.ID}
	}

	return nil
}

func (repo PackagingRepo) UpdateModelPackagingStatus(id string, s v1alpha1.ModelPackagingStatus) error {

	stmt, args, err := sq.Update(ModelPackagingTable).
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

func (repo PackagingRepo) CreateModelPackaging(mp *packaging.ModelPackaging) error {

	_, err := repo.DB.Exec(
		fmt.Sprintf("INSERT INTO %s (id, spec, status) VALUES($1, $2, $3)", ModelPackagingTable),
		mp.ID, mp.Spec, mp.Status,
	)
	if err != nil {
		pqError, ok := err.(*pq.Error)
		if ok && pqError.Code == uniqueViolationPostgresCode {
			return odahuErrors.AlreadyExistError{Entity: mp.ID}
		}
		return err
	}
	return nil

}
