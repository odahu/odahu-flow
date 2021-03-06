// Code generated by mockery v2.2.1. DO NOT EDIT.

package mocks

import (
	context "context"

	deployment "github.com/odahu/odahu-flow/packages/operator/pkg/apis/deployment"
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

// CreateModelRoute provides a mock function with given fields: ctx, tx, r
func (_m *Repository) SaveModelRoute(ctx context.Context, tx *sql.Tx, r *deployment.ModelRoute) error {
	ret := _m.Called(ctx, tx, r)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *sql.Tx, *deployment.ModelRoute) error); ok {
		r0 = rf(ctx, tx, r)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DefaultExists provides a mock function with given fields: ctx, mdID, qrr
func (_m *Repository) DefaultExists(ctx context.Context, mdID string, qrr *sql.Tx) (bool, error) {
	ret := _m.Called(ctx, mdID, qrr)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string, *sql.Tx) bool); ok {
		r0 = rf(ctx, mdID, qrr)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, *sql.Tx) error); ok {
		r1 = rf(ctx, mdID, qrr)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteModelRoute provides a mock function with given fields: ctx, tx, name
func (_m *Repository) DeleteModelRoute(ctx context.Context, tx *sql.Tx, name string) error {
	ret := _m.Called(ctx, tx, name)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *sql.Tx, string) error); ok {
		r0 = rf(ctx, tx, name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetModelRoute provides a mock function with given fields: ctx, tx, name
func (_m *Repository) GetModelRoute(ctx context.Context, tx *sql.Tx, name string) (*deployment.ModelRoute, error) {
	ret := _m.Called(ctx, tx, name)

	var r0 *deployment.ModelRoute
	if rf, ok := ret.Get(0).(func(context.Context, *sql.Tx, string) *deployment.ModelRoute); ok {
		r0 = rf(ctx, tx, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*deployment.ModelRoute)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *sql.Tx, string) error); ok {
		r1 = rf(ctx, tx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetModelRouteList provides a mock function with given fields: ctx, tx, options
func (_m *Repository) GetModelRouteList(ctx context.Context, tx *sql.Tx, options ...filter.ListOption) ([]deployment.ModelRoute, error) {
	_va := make([]interface{}, len(options))
	for _i := range options {
		_va[_i] = options[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, tx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []deployment.ModelRoute
	if rf, ok := ret.Get(0).(func(context.Context, *sql.Tx, ...filter.ListOption) []deployment.ModelRoute); ok {
		r0 = rf(ctx, tx, options...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]deployment.ModelRoute)
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

// IsDefault provides a mock function with given fields: ctx, id, tx
func (_m *Repository) IsDefault(ctx context.Context, id string, tx *sql.Tx) (bool, error) {
	ret := _m.Called(ctx, id, tx)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string, *sql.Tx) bool); ok {
		r0 = rf(ctx, id, tx)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, *sql.Tx) error); ok {
		r1 = rf(ctx, id, tx)
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

// UpdateModelRoute provides a mock function with given fields: ctx, tx, md
func (_m *Repository) UpdateModelRoute(ctx context.Context, tx *sql.Tx, md *deployment.ModelRoute) error {
	ret := _m.Called(ctx, tx, md)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *sql.Tx, *deployment.ModelRoute) error); ok {
		r0 = rf(ctx, tx, md)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateModelRouteStatus provides a mock function with given fields: ctx, tx, id, s
func (_m *Repository) UpdateModelRouteStatus(ctx context.Context, tx *sql.Tx, id string, s v1alpha1.ModelRouteStatus) error {
	ret := _m.Called(ctx, tx, id, s)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *sql.Tx, string, v1alpha1.ModelRouteStatus) error); ok {
		r0 = rf(ctx, tx, id, s)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
