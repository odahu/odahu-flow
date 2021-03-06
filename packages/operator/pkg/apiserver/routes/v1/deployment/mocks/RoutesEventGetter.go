// Code generated by mockery v2.2.1. DO NOT EDIT.

package mocks

import (
	context "context"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/event"

	mock "github.com/stretchr/testify/mock"
)

// RoutesEventGetter is an autogenerated mock type for the RoutesEventGetter type
type RoutesEventGetter struct {
	mock.Mock
}

// Get provides a mock function with given fields: ctx, cursor
func (_m *RoutesEventGetter) Get(ctx context.Context, cursor int) ([]event.RouteEvent, int, error) {
	ret := _m.Called(ctx, cursor)

	var r0 []event.RouteEvent
	if rf, ok := ret.Get(0).(func(context.Context, int) []event.RouteEvent); ok {
		r0 = rf(ctx, cursor)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]event.RouteEvent)
		}
	}

	var r1 int
	if rf, ok := ret.Get(1).(func(context.Context, int) int); ok {
		r1 = rf(ctx, cursor)
	} else {
		r1 = ret.Get(1).(int)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, int) error); ok {
		r2 = rf(ctx, cursor)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}
