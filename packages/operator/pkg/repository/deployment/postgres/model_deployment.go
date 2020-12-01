package postgres

import (
	"context"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/deployment"
	odahuErrors "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	utils "github.com/odahu/odahu-flow/packages/operator/pkg/repository/util/postgres"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	ModelDeploymentTable = "odahu_operator_deployment"
	ModelRouteTable = "odahu_operator_route"
	uniqueViolationPostgresCode = pq.ErrorCode("23505") // unique_violation
)

var (
	log = logf.Log.WithName("model-deployment--repository--postgres")
	MaxSize   = 500
	FirstPage = 0
	txOptions = &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  false,
	}
)

type DeploymentRepo struct {
	DB *sql.DB
}

func (repo DeploymentRepo) GetModelDeployment(
	ctx context.Context, tx *sql.Tx, id string) (*deployment.ModelDeployment, error) {

	var qrr utils.Querier
	qrr = repo.DB
	if tx != nil {
		qrr = tx
	}

	mt := new(deployment.ModelDeployment)

	q, args, err := sq.
		Select("id", "spec", "status", "deletionmark", "created", "updated").
		From(ModelDeploymentTable).
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

func (repo DeploymentRepo) GetModelDeploymentList(
	ctx context.Context, tx *sql.Tx, options ...filter.ListOption) ([]deployment.ModelDeployment, error) {

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
		From("odahu_operator_deployment").
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

	mts := make([]deployment.ModelDeployment, 0)

	for rows.Next() {
		mt := new(deployment.ModelDeployment)
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

func (repo DeploymentRepo) DeleteModelDeployment(ctx context.Context, tx *sql.Tx, id string) error {

	var qrr utils.Querier
	qrr = repo.DB
	if tx != nil {
		qrr = tx
	}

	stmt, args, err := sq.Delete(ModelDeploymentTable).Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Dollar).ToSql()
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

func (repo DeploymentRepo) SetDeletionMark(ctx context.Context, tx *sql.Tx, id string, value bool) error {

	var qrr utils.Querier
	qrr = repo.DB
	if tx != nil {
		qrr = tx
	}

	return utils.SetDeletionMark(ctx, qrr, ModelDeploymentTable, id, value)
}

func (repo DeploymentRepo) UpdateModelDeployment(
	ctx context.Context, tx *sql.Tx, md *deployment.ModelDeployment) error {

	var qrr utils.Querier
	qrr = repo.DB
	if tx != nil {
		qrr = tx
	}

	md.Status.State = ""

	stmt, args, err := sq.Update(ModelDeploymentTable).
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

func (repo DeploymentRepo) UpdateModelDeploymentStatus(
	ctx context.Context, tx *sql.Tx, id string, s v1alpha1.ModelDeploymentStatus) error {

	var qrr utils.Querier
	qrr = repo.DB
	if tx != nil {
		qrr = tx
	}

	stmt, args, err := sq.Update(ModelDeploymentTable).
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

func (repo DeploymentRepo) CreateModelDeployment(
	ctx context.Context, tx *sql.Tx, md *deployment.ModelDeployment) error {

	var qrr utils.Querier
	qrr = repo.DB
	if tx != nil {
		qrr = tx
	}

	stmt, args, err := sq.
		Insert(ModelDeploymentTable).
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

func (repo DeploymentRepo) BeginTransaction(ctx context.Context) (*sql.Tx, error) {
	return repo.DB.BeginTx(ctx, txOptions)
}
