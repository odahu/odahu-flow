// Code generated by mockery v2.2.1. DO NOT EDIT.

package mocks

import (
	context "context"

	apispackaging "github.com/odahu/odahu-flow/packages/operator/pkg/apis/packaging"

	filter "github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"

	mock "github.com/stretchr/testify/mock"

	sql "database/sql"

	v1alpha1 "github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// BeginTransaction provides a mock function with given fields: ctx
func (_m *Repository) BeginTransaction(ctx context.Context) (*sql.Tx, error) {
	ret := _m.Called(ctx)

	var r0 *sql.Tx
	if rf, ok := ret.Get(0).(func(context.Context) *sql.Tx); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sql.Tx)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateModelPackaging provides a mock function with given fields: ctx, tx, mp
func (_m *Repository) SaveModelPackaging(ctx context.Context, tx *sql.Tx, mp *apispackaging.ModelPackaging) error {
	ret := _m.Called(ctx, tx, mp)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *sql.Tx, *apispackaging.ModelPackaging) error); ok {
		r0 = rf(ctx, tx, mp)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteModelPackaging provides a mock function with given fields: ctx, tx, id
func (_m *Repository) DeleteModelPackaging(ctx context.Context, tx *sql.Tx, id string) error {
	ret := _m.Called(ctx, tx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *sql.Tx, string) error); ok {
		r0 = rf(ctx, tx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetModelPackaging provides a mock function with given fields: ctx, tx, id
func (_m *Repository) GetModelPackaging(ctx context.Context, tx *sql.Tx, id string) (*apispackaging.ModelPackaging, error) {
	ret := _m.Called(ctx, tx, id)

	var r0 *apispackaging.ModelPackaging
	if rf, ok := ret.Get(0).(func(context.Context, *sql.Tx, string) *apispackaging.ModelPackaging); ok {
		r0 = rf(ctx, tx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*apispackaging.ModelPackaging)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *sql.Tx, string) error); ok {
		r1 = rf(ctx, tx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetModelPackagingList provides a mock function with given fields: ctx, tx, options
func (_m *Repository) GetModelPackagingList(ctx context.Context, tx *sql.Tx, options ...filter.ListOption) ([]apispackaging.ModelPackaging, error) {
	_va := make([]interface{}, len(options))
	for _i := range options {
		_va[_i] = options[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, tx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []apispackaging.ModelPackaging
	if rf, ok := ret.Get(0).(func(context.Context, *sql.Tx, ...filter.ListOption) []apispackaging.ModelPackaging); ok {
		r0 = rf(ctx, tx, options...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]apispackaging.ModelPackaging)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *sql.Tx, ...filter.ListOption) error); ok {
		r1 = rf(ctx, tx, options...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetDeletionMark provides a mock function with given fields: ctx, tx, id, value
func (_m *Repository) SetDeletionMark(ctx context.Context, tx *sql.Tx, id string, value bool) error {
	ret := _m.Called(ctx, tx, id, value)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *sql.Tx, string, bool) error); ok {
		r0 = rf(ctx, tx, id, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateModelPackaging provides a mock function with given fields: ctx, tx, mp
func (_m *Repository) UpdateModelPackaging(ctx context.Context, tx *sql.Tx, mp *apispackaging.ModelPackaging) error {
	ret := _m.Called(ctx, tx, mp)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *sql.Tx, *apispackaging.ModelPackaging) error); ok {
		r0 = rf(ctx, tx, mp)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateModelPackagingStatus provides a mock function with given fields: ctx, tx, id, s
func (_m *Repository) UpdateModelPackagingStatus(ctx context.Context, tx *sql.Tx, id string, s v1alpha1.ModelPackagingStatus) error {
	ret := _m.Called(ctx, tx, id, s)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *sql.Tx, string, v1alpha1.ModelPackagingStatus) error); ok {
		r0 = rf(ctx, tx, id, s)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
