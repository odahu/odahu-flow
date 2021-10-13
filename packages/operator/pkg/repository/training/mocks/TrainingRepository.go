// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	apistraining "github.com/odahu/odahu-flow/packages/operator/pkg/apis/training"
	filter "github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"

	mock "github.com/stretchr/testify/mock"
)

// TrainingIntegrationRepository is an autogenerated mock type for the TrainingIntegrationRepository type
type TrainingIntegrationRepository struct {
	mock.Mock
}

// DeleteTrainingIntegration provides a mock function with given fields: name
func (_m *TrainingIntegrationRepository) DeleteTrainingIntegration(name string) error {
	ret := _m.Called(name)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetTrainingIntegration provides a mock function with given fields: name
func (_m *TrainingIntegrationRepository) GetTrainingIntegration(name string) (*apistraining.TrainingIntegration, error) {
	ret := _m.Called(name)

	var r0 *apistraining.TrainingIntegration
	if rf, ok := ret.Get(0).(func(string) *apistraining.TrainingIntegration); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*apistraining.TrainingIntegration)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTrainingIntegrationList provides a mock function with given fields: options
func (_m *TrainingIntegrationRepository) GetTrainingIntegrationList(options ...filter.ListOption) ([]apistraining.TrainingIntegration, error) {
	_va := make([]interface{}, len(options))
	for _i := range options {
		_va[_i] = options[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []apistraining.TrainingIntegration
	if rf, ok := ret.Get(0).(func(...filter.ListOption) []apistraining.TrainingIntegration); ok {
		r0 = rf(options...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]apistraining.TrainingIntegration)
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

// SaveTrainingIntegration provides a mock function with given fields: md
func (_m *TrainingIntegrationRepository) SaveTrainingIntegration(md *apistraining.TrainingIntegration) error {
	ret := _m.Called(md)

	var r0 error
	if rf, ok := ret.Get(0).(func(*apistraining.TrainingIntegration) error); ok {
		r0 = rf(md)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateTrainingIntegration provides a mock function with given fields: md
func (_m *TrainingIntegrationRepository) UpdateTrainingIntegration(md *apistraining.TrainingIntegration) error {
	ret := _m.Called(md)

	var r0 error
	if rf, ok := ret.Get(0).(func(*apistraining.TrainingIntegration) error); ok {
		r0 = rf(md)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
