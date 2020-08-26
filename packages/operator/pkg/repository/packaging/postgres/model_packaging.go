package postgres

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/packaging"
	odahuErrors "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/util/filter"
	utils "github.com/odahu/odahu-flow/packages/operator/pkg/repository/util/postgres"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	ModelPackagingTable = "odahu_operator_packaging"
)

var (
	log = logf.Log.WithName("model-packaging--repository--postgres")
)

type PackagingPostgresRepo struct {
	DB *sql.DB
}

func (repo PackagingPostgresRepo) GetModelPackaging(name string) (*packaging.ModelPackaging, error) {

	mt := new(packaging.ModelPackaging)

	err := repo.DB.QueryRow(
		fmt.Sprintf("SELECT id, spec, status FROM %s WHERE id = $1", ModelPackagingTable),
		name,
	).Scan(&mt.ID, &mt.Spec, &mt.Status)

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

func (repo PackagingPostgresRepo) GetModelPackagingList(options ...filter.ListOption) (
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

	sb := sq.Select("*").From("odahu_operator_packaging").
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
		err := rows.Scan(&mt.ID, &mt.Spec, &mt.Status)
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

func (repo PackagingPostgresRepo) DeleteModelPackaging(name string) error {

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

func (repo PackagingPostgresRepo) UpdateModelPackaging(mt *packaging.ModelPackaging) error {

	// First try to check that row exists otherwise raise exception to fit interface
	_, err := repo.GetModelPackaging(mt.ID)
	if err != nil {
		return err
	}

	sqlStatement := fmt.Sprintf("UPDATE %s SET spec = $1, status = $2 WHERE id = $3", ModelPackagingTable)
	_, err = repo.DB.Exec(sqlStatement, mt.Spec, mt.Status, mt.ID)
	if err != nil {
		return err
	}
	return nil
}

func (repo PackagingPostgresRepo) CreateModelPackaging(mt *packaging.ModelPackaging) error {

	_, err := repo.DB.Exec(
		fmt.Sprintf("INSERT INTO %s (id, spec, status) VALUES($1, $2, $3)", ModelPackagingTable),
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
