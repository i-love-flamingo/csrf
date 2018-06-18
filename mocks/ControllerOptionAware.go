package mocks

import mock "github.com/stretchr/testify/mock"
import router "flamingo.me/flamingo/framework/router"

// ControllerOptionAware is an autogenerated mock type for the ControllerOptionAware type
type ControllerOptionAware struct {
	mock.Mock
}

// CheckOption provides a mock function with given fields: option
func (_m *ControllerOptionAware) CheckOption(option router.ControllerOption) bool {
	ret := _m.Called(option)

	var r0 bool
	if rf, ok := ret.Get(0).(func(router.ControllerOption) bool); ok {
		r0 = rf(option)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}
