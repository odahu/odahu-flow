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
)

const (
	BatchInferenceServiceTable = "odahu_batch_inference_service"
)

// BatchInferenceService persistence repository
type BISRepo struct {
	DB *sql.DB
}

func (r BISRepo) Create(ctx context.Context, tx *sql.Tx, bis api_types.InferenceService) (err error) {
	var qrr utils.Querier
	qrr = r.DB
	if tx != nil {
		qrr = tx
	}

	stmt, args, err := sq.
		Insert(BatchInferenceServiceTable).
		Columns("id", "spec", "status", "created", "updated").
		Values(bis.ID, bis.Spec, bis.Status, bis.CreatedAt, bis.UpdatedAt).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	_, err = qrr.ExecContext(
		ctx,
		stmt,
		args...,
	)
	if err != nil {
		pqError, ok := err.(*pq.Error)
		if ok && pqError.Code == uniqueViolationPostgresCode {
			return odahuErrors.AlreadyExistError{Entity: bis.ID}
		}
		return err
	}
	return nil
}

func (r BISRepo) Update(
	ctx context.Context, tx *sql.Tx, id string, bis api_types.InferenceService) (err error) {

	var qrr utils.Querier
	qrr = r.DB
	if tx != nil {
		qrr = tx
	}

	stmt, args, err := sq.Update(BatchInferenceServiceTable).
		Set("deletionmark", bis.DeletionMark).
		Set("spec", bis.Spec).
		Set("created", bis.CreatedAt).
		Set("updated", bis.UpdatedAt).
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

func (r BISRepo) Delete(ctx context.Context, tx *sql.Tx, id string) (err error) {
	var qrr utils.Querier
	qrr = r.DB
	if tx != nil {
		qrr = tx
	}

	stmt, args, err := sq.Delete(BatchInferenceServiceTable).Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return err
	}

	result, err := qrr.ExecContext(ctx, stmt, args...)

	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch {
			case pqErr.Code == foreignKeyViolationCode && pqErr.Constraint == odahuJobServiceFKConstraint:
				return odahuErrors.DeletingServiceHasJobs{
					Entity: id,
				}
			default:
				return err
			}
		}
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

func (r BISRepo) List(
	ctx context.Context, tx *sql.Tx, options ...filter.ListOption) (res []api_types.InferenceService, err error) {

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

	sb := sq.Select("id, spec, deletionmark, created, updated").From(BatchInferenceServiceTable).
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
	res = make([]api_types.InferenceService, 0)

	for rows.Next() {
		s := api_types.InferenceService{}
		err := rows.Scan(&s.ID, &s.Spec, &s.DeletionMark, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func (r BISRepo) Get(ctx context.Context, tx *sql.Tx, id string) (res api_types.InferenceService, err error) {
	var qrr utils.Querier
	qrr = r.DB
	if tx != nil {
		qrr = tx
	}

	query, args, err := sq.
		Select("id", "spec", "deletionmark", "created", "updated").
		From(BatchInferenceServiceTable).
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
	).Scan(&res.ID, &res.Spec, &res.DeletionMark, &res.CreatedAt, &res.UpdatedAt)

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
