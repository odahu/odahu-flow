// Code generated by mockery v2.9.4. DO NOT EDIT.

package training_integration

import (
	filter "github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"
	mock "github.com/stretchr/testify/mock"

	training "github.com/odahu/odahu-flow/packages/operator/pkg/apis/training"
)

// MockTrainingIntegrationService is an autogenerated mock type for the trainingIntegrationService type
type MockTrainingIntegrationService struct {
	mock.Mock
}

// CreateTrainingIntegration provides a mock function with given fields: md
func (_m *MockTrainingIntegrationService) CreateTrainingIntegration(md *training.TrainingIntegration) error {
	ret := _m.Called(md)

	var r0 error
	if rf, ok := ret.Get(0).(func(*training.TrainingIntegration) error); ok {
		r0 = rf(md)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteTrainingIntegration provides a mock function with given fields: name
func (_m *MockTrainingIntegrationService) DeleteTrainingIntegration(name string) error {
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
func (_m *MockTrainingIntegrationService) GetTrainingIntegration(name string) (*training.TrainingIntegration, error) {
	ret := _m.Called(name)

	var r0 *training.TrainingIntegration
	if rf, ok := ret.Get(0).(func(string) *training.TrainingIntegration); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*training.TrainingIntegration)
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
func (_m *MockTrainingIntegrationService) GetTrainingIntegrationList(options ...filter.ListOption) ([]training.TrainingIntegration, error) {
	_va := make([]interface{}, len(options))
	for _i := range options {
		_va[_i] = options[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []training.TrainingIntegration
	if rf, ok := ret.Get(0).(func(...filter.ListOption) []training.TrainingIntegration); ok {
		r0 = rf(options...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]training.TrainingIntegration)
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

// UpdateTrainingIntegration provides a mock function with given fields: md
func (_m *MockTrainingIntegrationService) UpdateTrainingIntegration(md *training.TrainingIntegration) error {
	ret := _m.Called(md)

	var r0 error
	if rf, ok := ret.Get(0).(func(*training.TrainingIntegration) error); ok {
		r0 = rf(md)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
