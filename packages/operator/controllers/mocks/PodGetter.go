// Code generated by mockery v2.2.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	v1 "k8s.io/api/core/v1"
)

// PodGetter is an autogenerated mock type for the PodGetter type
type PodGetter struct {
	mock.Mock
}

// GetPod provides a mock function with given fields: ctx, name, namespace
func (_m *PodGetter) GetPod(ctx context.Context, name string, namespace string) (v1.Pod, error) {
	ret := _m.Called(ctx, name, namespace)

	var r0 v1.Pod
	if rf, ok := ret.Get(0).(func(context.Context, string, string) v1.Pod); ok {
		r0 = rf(ctx, name, namespace)
	} else {
		r0 = ret.Get(0).(v1.Pod)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, name, namespace)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
