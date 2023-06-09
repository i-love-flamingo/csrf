package interfaces_test

import (
	"context"
	"flamingo.me/csrf/interfaces"
	"net/http"
	"net/http/httptest"
	"testing"

	applicationMocks "flamingo.me/csrf/application/mocks"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/stretchr/testify/suite"
)

type (
	CsrfFilterTestSuite struct {
		suite.Suite

		filter      *interfaces.CsrfFilter
		service     *applicationMocks.Service
		nextFilter  *MockFilter
		filterChain *web.FilterChain

		context        context.Context
		webRequest     *web.Request
		responseWriter http.ResponseWriter
	}

	MockFilter struct{}
)

func (fnc MockFilter) Filter(context.Context, *web.Request, http.ResponseWriter, *web.FilterChain) web.Result {
	return &web.Response{}
}

func TestCsrfFilterTestSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, &CsrfFilterTestSuite{})
}

func (t *CsrfFilterTestSuite) SetupSuite() {
	t.context = context.Background()
	t.responseWriter = httptest.NewRecorder()
	t.webRequest = web.CreateRequest(nil, nil)
}

func (t *CsrfFilterTestSuite) SetupTest() {
	t.service = &applicationMocks.Service{}

	t.filter = &interfaces.CsrfFilter{}
	t.filter.Inject(&web.Responder{}, t.service)

	t.nextFilter = &MockFilter{}
	t.filterChain = web.NewFilterChain(nil, t.nextFilter)
}

func (t *CsrfFilterTestSuite) TearDown() {
	t.service.AssertExpectations(t.T())
	t.service = nil
}

func (t *CsrfFilterTestSuite) TestFilter_WrongToken() {
	t.service.EXPECT().IsValid(t.webRequest).Return(false).Once()

	response := t.filter.Filter(t.context, t.webRequest, t.responseWriter, t.filterChain)
	forbidden, ok := response.(*web.ServerErrorResponse)
	t.True(ok)
	t.Equal(uint(http.StatusForbidden), forbidden.Response.Status)
}

func (t *CsrfFilterTestSuite) TestFilter_Success() {
	t.service.EXPECT().IsValid(t.webRequest).Return(true).Once()

	response := t.filter.Filter(t.context, t.webRequest, t.responseWriter, t.filterChain)
	t.IsType(&web.Response{}, response)
}
