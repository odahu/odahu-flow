package postgres

import (
	sq "github.com/Masterminds/squirrel"
	"reflect"
)

const tagKey = "postgres"

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
