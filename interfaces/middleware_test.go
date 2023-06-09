package interfaces_test

import (
	"context"
	"flamingo.me/csrf/interfaces"
	"net/http"
	"testing"

	applicationMocks "flamingo.me/csrf/application/mocks"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/stretchr/testify/suite"
)

type (
	CsrfMiddlewareTestSuite struct {
		suite.Suite

		middleware *interfaces.CsrfMiddleware
		service    *applicationMocks.Service

		action     web.Action
		context    context.Context
		webRequest *web.Request
	}
)

func TestCsrfMiddlewareTestSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, &CsrfMiddlewareTestSuite{})
}

func (t *CsrfMiddlewareTestSuite) SetupSuite() {
	t.context = context.Background()
	t.action = func(ctx context.Context, req *web.Request) web.Result {
		return &web.Response{}
	}
	t.webRequest = web.CreateRequest(nil, nil)
}

func (t *CsrfMiddlewareTestSuite) SetupTest() {
	t.service = &applicationMocks.Service{}

	t.middleware = &interfaces.CsrfMiddleware{}
	t.middleware.Inject(&web.Responder{}, t.service)
}

func (t *CsrfMiddlewareTestSuite) TearDown() {
	t.service.AssertExpectations(t.T())
	t.service = nil
}

func (t *CsrfMiddlewareTestSuite) TestMiddleware_WrongToken() {
	t.service.EXPECT().IsValidPost(t.webRequest).Return(false).Once()

	handler := t.middleware.Secured(t.action)
	response := handler(t.context, t.webRequest)
	forbidden, ok := response.(*web.ServerErrorResponse)
	t.True(ok)
	t.Equal(uint(http.StatusForbidden), forbidden.Response.Status)
}

func (t *CsrfMiddlewareTestSuite) TestMiddleware_Success() {
	t.service.EXPECT().IsValidPost(t.webRequest).Return(true).Once()

	handler := t.middleware.Secured(t.action)
	response := handler(t.context, t.webRequest)
	t.IsType(&web.Response{}, response)
}

func (t *CsrfMiddlewareTestSuite) TestMiddleware_HeaderWrongToken() {
	t.service.EXPECT().IsValidHeader(t.webRequest).Return(false).Once()

	handler := t.middleware.SecuredHeader(t.action)
	response := handler(t.context, t.webRequest)
	forbidden, ok := response.(*web.ServerErrorResponse)
	t.True(ok)
	t.Equal(uint(http.StatusForbidden), forbidden.Response.Status)
}

func (t *CsrfMiddlewareTestSuite) TestMiddleware_HeaderSuccess() {
	t.service.EXPECT().IsValidHeader(t.webRequest).Return(true).Once()

	handler := t.middleware.SecuredHeader(t.action)
	response := handler(t.context, t.webRequest)
	t.IsType(&web.Response{}, response)
}
