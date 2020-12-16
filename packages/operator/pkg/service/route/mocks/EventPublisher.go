// Code generated by mockery v2.2.1. DO NOT EDIT.

package mocks

import (
	context "context"

	event "github.com/odahu/odahu-flow/packages/operator/pkg/apis/event"
	mock "github.com/stretchr/testify/mock"

	sql "database/sql"
)

// EventPublisher is an autogenerated mock type for the EventPublisher type
type EventPublisher struct {
	mock.Mock
}

// PublishEvent provides a mock function with given fields: ctx, tx, _a2
func (_m *EventPublisher) PublishEvent(ctx context.Context, tx *sql.Tx, _a2 event.Event) error {
	ret := _m.Called(ctx, tx, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *sql.Tx, event.Event) error); ok {
		r0 = rf(ctx, tx, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
