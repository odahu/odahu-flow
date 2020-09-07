package postgres

import (
	"context"
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	odahuErrors "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	"reflect"
)

const (
	tagKey = "postgres"
	deletionMarkColumn = "deletionmark"
)

func TransformFilter(sqlBuilder sq.SelectBuilder, filter interface{}) sq.SelectBuilder {
	if filter == nil {
		return sqlBuilder
	}

	var conditions sq.And

	elem := reflect.ValueOf(filter).Elem()
	for i := 0; i < elem.NumField(); i++ {
		value, ok := elem.Field(i).Interface().([]string)

		if !ok {
			continue
		}
		if len(value) == 0 {
			continue
		}
		if len(value) == 1 && (value[0] == "*") {
			continue
		}

		field := elem.Type().Field(i).Tag.Get(tagKey)

		conditions = append(conditions, sq.Eq{field: value})
	}

	if len(conditions) > 0 {
		newSQLBuilder := sqlBuilder.Where(conditions)
		return newSQLBuilder
	}
	return sqlBuilder
}

func SetDeletionMark(ctx context.Context, qrr Querier, tableName string, id string, value bool) error {
	stmt, args, err := sq.
		Update(tableName).
		Set(deletionMarkColumn, value).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}
	res, err := qrr.ExecContext(ctx, stmt, args...)
	if err != nil {
		return err
	}
	rowsAffected, err:= res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return odahuErrors.NotFoundError{Entity: id}
	}
	if rowsAffected > 1 {
		return fmt.Errorf("more that one rows found for ID %s", id)
	}
	return nil
}


type Querier interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}