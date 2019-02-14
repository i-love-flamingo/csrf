package templatefunctions

import (
	"context"
	"testing"

	applicationMocks "flamingo.me/csrf/application/mocks"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/stretchr/testify/suite"
)

type (
	CsrfTokenFuncTestSuite struct {
		suite.Suite

		csrfFunc *CsrfTokenFunc

		service *applicationMocks.Service

		context context.Context
		session *web.Session
	}
)

func TestCsrfTokenFuncTestSuite(t *testing.T) {
	suite.Run(t, &CsrfTokenFuncTestSuite{})
}

func (t *CsrfTokenFuncTestSuite) SetupSuite() {
	t.session = web.EmptySession()
	t.context = web.ContextWithSession(context.Background(), t.session)
}

func (t *CsrfTokenFuncTestSuite) SetupTest() {
	t.service = &applicationMocks.Service{}

	t.csrfFunc = &CsrfTokenFunc{}
	t.csrfFunc.Inject(t.service, flamingo.NullLogger{})
}

func (t *CsrfTokenFuncTestSuite) TearDown() {
	t.service.AssertExpectations(t.T())
	t.service = nil
	t.csrfFunc = nil
}

func (t *CsrfTokenFuncTestSuite) TestFunc() {
	t.service.On("Generate", t.session).Return("token").Once()

	function := t.csrfFunc.Func(t.context)
	csrfFunc, ok := function.(func() interface{})
	t.True(ok)
	content := csrfFunc()
	t.Equal("token", content)
}
