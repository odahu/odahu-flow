package postgres

import (
	"context"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	route "github.com/odahu/odahu-flow/packages/operator/pkg/apis/deployment"
	odahuErrors "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	utils "github.com/odahu/odahu-flow/packages/operator/pkg/repository/util/postgres"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	ModelRouteTable = "odahu_operator_route"
	uniqueViolationPostgresCode = pq.ErrorCode("23505") // unique_violation
)

var (
	log = logf.Log.WithName("model-route--repository--postgres")
	MaxSize   = 500
	FirstPage = 0
	txOptions = &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  false,
	}
)

type RouteRepo struct {
	DB *sql.DB
}

func (repo RouteRepo) GetModelRoute(
	ctx context.Context, tx *sql.Tx, id string) (*route.ModelRoute, error) {

	var qrr utils.Querier
	qrr = repo.DB
	if tx != nil {
		qrr = tx
	}

	mt := new(route.ModelRoute)

	q, args, err := sq.
		Select("id", "spec", "status", "deletionmark", "created", "updated").
		From(ModelRouteTable).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	err = qrr.QueryRowContext(ctx, q, args...).
		Scan(&mt.ID, &mt.Spec, &mt.Status, &mt.DeletionMark, &mt.CreatedAt, &mt.UpdatedAt)

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

func (repo RouteRepo) GetModelRouteList(
	ctx context.Context, tx *sql.Tx, options ...filter.ListOption) ([]route.ModelRoute, error) {

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

	sb := sq.
		Select("id, spec, status, deletionmark, created, updated").
		From("odahu_operator_route").
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

	mts := make([]route.ModelRoute, 0)

	for rows.Next() {
		mt := new(route.ModelRoute)
		err := rows.Scan(&mt.ID, &mt.Spec, &mt.Status, &mt.DeletionMark, &mt.CreatedAt, &mt.UpdatedAt)
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

func (repo RouteRepo) DeleteModelRoute(ctx context.Context, tx *sql.Tx, id string) error {

	var qrr utils.Querier
	qrr = repo.DB
	if tx != nil {
		qrr = tx
	}

	stmt, args, err := sq.Delete(ModelRouteTable).Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Dollar).ToSql()
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

func (repo RouteRepo) SetDeletionMark(ctx context.Context, tx *sql.Tx, id string, value bool) error {

	var qrr utils.Querier
	qrr = repo.DB
	if tx != nil {
		qrr = tx
	}

	return utils.SetDeletionMark(ctx, qrr, ModelRouteTable, id, value)
}

func (repo RouteRepo) UpdateModelRoute(
	ctx context.Context, tx *sql.Tx, md *route.ModelRoute) error {

	var qrr utils.Querier
	qrr = repo.DB
	if tx != nil {
		qrr = tx
	}

	md.Status.State = ""

	stmt, args, err := sq.Update(ModelRouteTable).
		Set("spec", md.Spec).
		Set("status", md.Status).
		Set("updated", md.UpdatedAt).
		Where(sq.Eq{"id": md.ID}).
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
		return odahuErrors.NotFoundError{Entity: md.ID}
	}

	return nil
}

func (repo RouteRepo) UpdateModelRouteStatus(
	ctx context.Context, tx *sql.Tx, id string, s v1alpha1.ModelRouteStatus) error {

	var qrr utils.Querier
	qrr = repo.DB
	if tx != nil {
		qrr = tx
	}

	stmt, args, err := sq.Update(ModelRouteTable).
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

func (repo RouteRepo) CreateModelRoute(
	ctx context.Context, tx *sql.Tx, md *route.ModelRoute) error {

	var qrr utils.Querier
	qrr = repo.DB
	if tx != nil {
		qrr = tx
	}

	stmt, args, err := sq.
		Insert(ModelRouteTable).
		Columns("id", "spec", "status", "created", "updated").
		Values(md.ID, md.Spec, md.Status, md.CreatedAt, md.UpdatedAt).
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
			return odahuErrors.AlreadyExistError{Entity: md.ID}
		}
		return err
	}
	return nil

}


func (repo RouteRepo) BeginTransaction(ctx context.Context) (*sql.Tx, error) {
	return repo.DB.BeginTx(ctx, txOptions)
}