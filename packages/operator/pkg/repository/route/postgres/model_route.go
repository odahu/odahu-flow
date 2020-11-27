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

	ClID = "id"
	ClSpec = "spec"
	ClStatus = "status"
	ClDelMark = "deletionmark"
	ClCreated = "created"
	ClUpdated = "updated"
	ClIsDefault = "is_default"
	ClFirstMDName = "spec->'modelDeployments'->0->>'mdName'"
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
		Select(ClID, ClSpec, ClStatus, ClDelMark, ClCreated, ClUpdated, ClIsDefault).
		From(ModelRouteTable).
		Where(sq.Eq{ClID: id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	err = qrr.QueryRowContext(ctx, q, args...).
		Scan(&mt.ID, &mt.Spec, &mt.Status, &mt.DeletionMark, &mt.CreatedAt, &mt.UpdatedAt, &mt.Default)

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
		Select(ClID, ClSpec, ClStatus, ClDelMark, ClCreated, ClUpdated, ClIsDefault).
		From(ModelRouteTable).
		OrderBy(ClID).
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
		err := rows.Scan(&mt.ID, &mt.Spec, &mt.Status, &mt.DeletionMark, &mt.CreatedAt, &mt.UpdatedAt, &mt.Default)
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

	stmt, args, err := sq.Delete(ModelRouteTable).Where(sq.Eq{ClID: id}).PlaceholderFormat(sq.Dollar).ToSql()
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
		Set(ClSpec, md.Spec).
		Set(ClStatus, md.Status).
		Set(ClUpdated, md.UpdatedAt).
		Set(ClIsDefault, md.Default).
		Where(sq.Eq{ClID: md.ID}).
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
		Set(ClStatus, s).
		Where(sq.Eq{ClID: id}).
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
		Columns(ClID, ClSpec, ClStatus, ClCreated, ClUpdated, ClIsDefault).
		Values(md.ID, md.Spec, md.Status, md.CreatedAt, md.UpdatedAt, md.Default).
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

// Check whether the default ModelRoute for ModelDeployment with id=mdID is existed
func (repo RouteRepo) DefaultExists(ctx context.Context, mdID string, tx *sql.Tx) (bool, error) {

	var qrr utils.Querier
	qrr = repo.DB
	if tx != nil {
		qrr = tx
	}

	stmt, args, err := sq.
		Select("count(id)").
		From(ModelRouteTable).
		Where(sq.Eq{ClIsDefault: true, ClFirstMDName: mdID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return false, err
	}

	row := qrr.QueryRowContext(ctx, stmt, args...)
	var count int
	if err := row.Scan(&count); err != nil {
		return false, err
	}
	return count >= 1, nil
}

// Check whether the the route is default route
func (repo RouteRepo) IsDefault(ctx context.Context, id string, tx *sql.Tx) (bool, error) {

	var qrr utils.Querier
	qrr = repo.DB
	if tx != nil {
		qrr = tx
	}

	stmt, args, err := sq.
		Select(ClIsDefault).
		From(ModelRouteTable).
		Where(sq.Eq{ClID: id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return false, err
	}

	rows, err := qrr.QueryContext(ctx, stmt, args...)
	if err != nil {
		return false, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Error(err, "error during rows.Close()")
		}
		if err := rows.Err(); err != nil {
			log.Error(err, "error during rows iterating")
		}
	}()
	var isDefault bool
	for rows.Next() {
		if err := rows.Scan(&isDefault); err != nil {
			return false, err
		}
		return isDefault, nil  //nolint
	}

	return false, odahuErrors.NotFoundError{Entity: id}
}
