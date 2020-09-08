// Code generated by mockery v2.2.1. DO NOT EDIT.

package mocks

import (
	context "context"

	filter "github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"
	mock "github.com/stretchr/testify/mock"

	packaging "github.com/odahu/odahu-flow/packages/operator/pkg/apis/packaging"

	postgres "github.com/odahu/odahu-flow/packages/operator/pkg/repository/util/postgres"

	v1alpha1 "github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// CreateModelPackaging provides a mock function with given fields: ctx, qrr, mp
func (_m *Repository) CreateModelPackaging(ctx context.Context, qrr postgres.Querier, mp *packaging.ModelPackaging) error {
	ret := _m.Called(ctx, qrr, mp)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, postgres.Querier, *packaging.ModelPackaging) error); ok {
		r0 = rf(ctx, qrr, mp)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteModelPackaging provides a mock function with given fields: ctx, qrr, id
func (_m *Repository) DeleteModelPackaging(ctx context.Context, qrr postgres.Querier, id string) error {
	ret := _m.Called(ctx, qrr, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, postgres.Querier, string) error); ok {
		r0 = rf(ctx, qrr, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetModelPackaging provides a mock function with given fields: ctx, qrr, id
func (_m *Repository) GetModelPackaging(ctx context.Context, qrr postgres.Querier, id string) (*packaging.ModelPackaging, error) {
	ret := _m.Called(ctx, qrr, id)

	var r0 *packaging.ModelPackaging
	if rf, ok := ret.Get(0).(func(context.Context, postgres.Querier, string) *packaging.ModelPackaging); ok {
		r0 = rf(ctx, qrr, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*packaging.ModelPackaging)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, postgres.Querier, string) error); ok {
		r1 = rf(ctx, qrr, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetModelPackagingList provides a mock function with given fields: ctx, qrr, options
func (_m *Repository) GetModelPackagingList(ctx context.Context, qrr postgres.Querier, options ...filter.ListOption) ([]packaging.ModelPackaging, error) {
	_va := make([]interface{}, len(options))
	for _i := range options {
		_va[_i] = options[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, qrr)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []packaging.ModelPackaging
	if rf, ok := ret.Get(0).(func(context.Context, postgres.Querier, ...filter.ListOption) []packaging.ModelPackaging); ok {
		r0 = rf(ctx, qrr, options...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]packaging.ModelPackaging)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, postgres.Querier, ...filter.ListOption) error); ok {
		r1 = rf(ctx, qrr, options...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetDeletionMark provides a mock function with given fields: ctx, qrr, id, value
func (_m *Repository) SetDeletionMark(ctx context.Context, qrr postgres.Querier, id string, value bool) error {
	ret := _m.Called(ctx, qrr, id, value)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, postgres.Querier, string, bool) error); ok {
		r0 = rf(ctx, qrr, id, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateModelPackaging provides a mock function with given fields: ctx, qrr, mp
func (_m *Repository) UpdateModelPackaging(ctx context.Context, qrr postgres.Querier, mp *packaging.ModelPackaging) error {
	ret := _m.Called(ctx, qrr, mp)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, postgres.Querier, *packaging.ModelPackaging) error); ok {
		r0 = rf(ctx, qrr, mp)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateModelPackagingStatus provides a mock function with given fields: ctx, qrr, id, s
func (_m *Repository) UpdateModelPackagingStatus(ctx context.Context, qrr postgres.Querier, id string, s v1alpha1.ModelPackagingStatus) error {
	ret := _m.Called(ctx, qrr, id, s)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, postgres.Querier, string, v1alpha1.ModelPackagingStatus) error); ok {
		r0 = rf(ctx, qrr, id, s)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
