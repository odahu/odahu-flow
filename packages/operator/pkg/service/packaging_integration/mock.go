/*
 * Copyright 2021 EPAM Systems
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package packaging_integration

import (
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/packaging"
	filter "github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"
	mock "github.com/stretchr/testify/mock"
)

// MockService is an autogenerated mock type for the Service type
type MockService struct {
	mock.Mock
}

// CreatePackagingIntegration provides a mock function with given fields: md
func (_m *MockService) CreatePackagingIntegration(md *packaging.PackagingIntegration) error {
	ret := _m.Called(md)

	var r0 error
	if rf, ok := ret.Get(0).(func(*packaging.PackagingIntegration) error); ok {
		r0 = rf(md)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeletePackagingIntegration provides a mock function with given fields: id
func (_m *MockService) DeletePackagingIntegration(id string) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetPackagingIntegration provides a mock function with given fields: id
func (_m *MockService) GetPackagingIntegration(id string) (*packaging.PackagingIntegration, error) {
	ret := _m.Called(id)

	var r0 *packaging.PackagingIntegration
	if rf, ok := ret.Get(0).(func(string) *packaging.PackagingIntegration); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*packaging.PackagingIntegration)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPackagingIntegrationList provides a mock function with given fields: options
func (_m *MockService) GetPackagingIntegrationList(options ...filter.ListOption) ([]packaging.PackagingIntegration, error) {
	_va := make([]interface{}, len(options))
	for _i := range options {
		_va[_i] = options[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []packaging.PackagingIntegration
	if rf, ok := ret.Get(0).(func(...filter.ListOption) []packaging.PackagingIntegration); ok {
		r0 = rf(options...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]packaging.PackagingIntegration)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(...filter.ListOption) error); ok {
		r1 = rf(options...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdatePackagingIntegration provides a mock function with given fields: md
func (_m *MockService) UpdatePackagingIntegration(md *packaging.PackagingIntegration) error {
	ret := _m.Called(md)

	var r0 error
	if rf, ok := ret.Get(0).(func(*packaging.PackagingIntegration) error); ok {
		r0 = rf(md)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
