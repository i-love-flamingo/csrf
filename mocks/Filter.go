package mocks

import http "net/http"
import mock "github.com/stretchr/testify/mock"
import router "flamingo.me/flamingo/framework/router"
import web "flamingo.me/flamingo/framework/web"

// Filter is an autogenerated mock type for the Filter type
type Filter struct {
	mock.Mock
}

// Filter provides a mock function with given fields: ctx, w, fc
func (_m *Filter) Filter(ctx web.Context, w http.ResponseWriter, fc *router.FilterChain) web.Response {
	ret := _m.Called(ctx, w, fc)

	var r0 web.Response
	if rf, ok := ret.Get(0).(func(web.Context, http.ResponseWriter, *router.FilterChain) web.Response); ok {
		r0 = rf(ctx, w, fc)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(web.Response)
		}
	}

	return r0
}
