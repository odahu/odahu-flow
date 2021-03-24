/*
 *
 *     Copyright 2021 EPAM Systems
 *
 *     Licensed under the Apache License, Version 2.0 (the "License");
 *     you may not use this file except in compliance with the License.
 *     You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 *     Unless required by applicable law or agreed to in writing, software
 *     distributed under the License is distributed on an "AS IS" BASIS,
 *     WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *     See the License for the specific language governing permissions and
 *     limitations under the License.
 */

package postgres

import (
	"context"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	api_types "github.com/odahu/odahu-flow/packages/operator/pkg/apis/batch"
	odahuErrors "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	utils "github.com/odahu/odahu-flow/packages/operator/pkg/repository/util/postgres"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	BatchInferenceJobTable = "odahu_batch_inference_job"
	uniqueViolationPostgresCode = pq.ErrorCode("23505") // unique_violation
	foreignKeyViolationCode = pq.ErrorCode("23503")

	odahuJobServiceFKConstraint = "odahu_bij_bis_fk"
)

var (
	log = logf.Log.WithName("batch-inference--repository--postgres")
	MaxSize   = 500
	FirstPage = 0
)

// InferenceJob persistence repository
type BIJRepo struct {
	DB *sql.DB
}

func (r BIJRepo) Create(ctx context.Context, tx *sql.Tx, bij api_types.InferenceJob) (err error) {
	var qrr utils.Querier
	qrr = r.DB
	if tx != nil {
		qrr = tx
	}

	stmt, args, err := sq.
		Insert(BatchInferenceJobTable).
		Columns("id", "spec", "status", "created", "updated", "service").
		Values(bij.ID, bij.Spec, bij.Status, bij.CreatedAt, bij.UpdatedAt, bij.Spec.InferenceServiceID).
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
		if ok {
			switch {
			case pqError.Code == uniqueViolationPostgresCode:
				return odahuErrors.AlreadyExistError{Entity: bij.ID}
			case pqError.Code == foreignKeyViolationCode && pqError.Constraint == odahuJobServiceFKConstraint:
				return odahuErrors.CreatingJobServiceNotFound{
					Entity:  bij.ID,
					Service: bij.Spec.InferenceServiceID,
				}
			default:
				return err
			}
		}
		return err
	}
	return nil
}

func (r BIJRepo) UpdateStatus(
	ctx context.Context, tx *sql.Tx, id string, s api_types.InferenceJobStatus) (err error) {

	var qrr utils.Querier
	qrr = r.DB
	if tx != nil {
		qrr = tx
	}

	stmt, args, err := sq.Update(BatchInferenceJobTable).
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

func (r BIJRepo) Delete(ctx context.Context, tx *sql.Tx, id string) (err error) {
	var qrr utils.Querier
	qrr = r.DB
	if tx != nil {
		qrr = tx
	}

	stmt, args, err := sq.Delete(BatchInferenceJobTable).Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Dollar).ToSql()
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

func (r BIJRepo) List(
	ctx context.Context, tx *sql.Tx, options ...filter.ListOption) (res []api_types.InferenceJob, err error) {
	var qrr utils.Querier
	qrr = r.DB
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

	sb := sq.Select("id, spec, status, deletionmark, created, updated").From(BatchInferenceJobTable).
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

	// To avoid nil
	res = make([]api_types.InferenceJob, 0)
	for rows.Next() {
		j := api_types.InferenceJob{}
		err := rows.Scan(&j.ID, &j.Spec, &j.Status, &j.DeletionMark, &j.CreatedAt, &j.UpdatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, j)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func (r BIJRepo) Get(ctx context.Context, tx *sql.Tx, id string) (res api_types.InferenceJob, err error) {
	var qrr utils.Querier
	qrr = r.DB
	if tx != nil {
		qrr = tx
	}

	query, args, err := sq.
		Select("id", "spec", "status", "deletionmark", "created", "updated").
		From(BatchInferenceJobTable).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return res, err
	}

	err = qrr.QueryRowContext(
		ctx,
		query,
		args...,
	).Scan(&res.ID, &res.Spec, &res.Status, &res.DeletionMark, &res.CreatedAt, &res.UpdatedAt)

	switch {
	case err == sql.ErrNoRows:
		return res, odahuErrors.NotFoundError{Entity: id}
	case err != nil:
		log.Error(err, "error during sql query")
		return res, err
	default:
		return res, nil
	}
}

func (r BIJRepo) SetDeletionMark(ctx context.Context, tx *sql.Tx, id string, value bool) error {

	var qrr utils.Querier
	qrr = r.DB
	if tx != nil {
		qrr = tx
	}

	return utils.SetDeletionMark(ctx, qrr, BatchInferenceJobTable, id, value)
}