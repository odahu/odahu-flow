package postgres

import (
	"context"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/packaging"
	odahuErrors "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	utils "github.com/odahu/odahu-flow/packages/operator/pkg/repository/util/postgres"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	ModelPackagingTable = "odahu_operator_packaging"
)

var (
	log = logf.Log.WithName("model-packaging--repository--postgres")
	txOptions = &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  false,
	}
)

type PackagingRepo struct {
	DB *sql.DB
}

func (repo PackagingRepo) GetModelPackaging(
	ctx context.Context, tx *sql.Tx, id string) (*packaging.ModelPackaging, error) {

	var qrr utils.Querier
	qrr = repo.DB
	if tx != nil {
		qrr = tx
	}

	mt := new(packaging.ModelPackaging)

	q, args, err := sq.
		Select("id", "spec", "status", "deletionmark").
		From(ModelPackagingTable).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	err = qrr.QueryRowContext(ctx, q, args...).Scan(&mt.ID, &mt.Spec, &mt.Status, &mt.DeletionMark)

	switch {
	case err == sql.ErrNoRows:
		return nil, odahuErrors.NotFoundError{Entity: id}
	case err != nil:
		log.Error(err, "error during sql query")
		return nil, err
	default:
		return mt, nil
	}
}

func (repo PackagingRepo) GetModelPackagingList(
	ctx context.Context, tx *sql.Tx, options ...filter.ListOption) ([]packaging.ModelPackaging, error) {

	var qrr utils.Querier
	qrr = repo.DB
	if tx != nil {
		qrr = tx
	}

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

	rows, err := qrr.QueryContext(ctx, stmt, args...)
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

func (repo PackagingRepo) DeleteModelPackaging(ctx context.Context, tx *sql.Tx, id string) error {

	var qrr utils.Querier
	qrr = repo.DB
	if tx != nil {
		qrr = tx
	}

	stmt, args, err := sq.Delete(ModelPackagingTable).Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return err
	}

	result, err := qrr.ExecContext(ctx, stmt, args...)

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

func (repo PackagingRepo) SetDeletionMark(ctx context.Context, tx *sql.Tx, id string, value bool) error {

	var qrr utils.Querier
	qrr = repo.DB
	if tx != nil {
		qrr = tx
	}

	return utils.SetDeletionMark(ctx, qrr, ModelPackagingTable, id, value)
}

func (repo PackagingRepo) UpdateModelPackaging(ctx context.Context, tx *sql.Tx, mp *packaging.ModelPackaging) error {

	var qrr utils.Querier
	qrr = repo.DB
	if tx != nil {
		qrr = tx
	}

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

	result, err := qrr.ExecContext(ctx, stmt, args...)
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

func (repo PackagingRepo) UpdateModelPackagingStatus(
	ctx context.Context, tx *sql.Tx, id string, s v1alpha1.ModelPackagingStatus) error {

	var qrr utils.Querier
	qrr = repo.DB
	if tx != nil {
		qrr = tx
	}

	stmt, args, err := sq.Update(ModelPackagingTable).
		Set("status", s).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	result, err := qrr.ExecContext(ctx, stmt, args...)
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

func (repo PackagingRepo) CreateModelPackaging(ctx context.Context, tx *sql.Tx, mp *packaging.ModelPackaging) error {

	var qrr utils.Querier
	qrr = repo.DB
	if tx != nil {
		qrr = tx
	}

	stmt, args, err := sq.
		Insert(ModelPackagingTable).
		Columns("id", "spec", "status").
		Values(mp.ID, mp.Spec, mp.Status).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	_, err = qrr.ExecContext(
		ctx,
		stmt,
		args...
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

func (repo PackagingRepo) BeginTransaction(ctx context.Context) (*sql.Tx, error) {
	return repo.DB.BeginTx(ctx, txOptions)
}

