package templatefunctions_test

import (
	"context"
	"flamingo.me/csrf/interfaces/templatefunctions"
	"testing"

	applicationMocks "flamingo.me/csrf/application/mocks"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/stretchr/testify/suite"
)

type (
	CsrfTokenFuncTestSuite struct {
		suite.Suite

		csrfFunc *templatefunctions.CsrfTokenFunc

		service *applicationMocks.Service

		context context.Context
		session *web.Session
	}
)

func TestCsrfTokenFuncTestSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, &CsrfTokenFuncTestSuite{})
}

func (t *CsrfTokenFuncTestSuite) SetupSuite() {
	t.session = web.EmptySession()
	t.context = web.ContextWithSession(context.Background(), t.session)
}

func (t *CsrfTokenFuncTestSuite) SetupTest() {
	t.service = &applicationMocks.Service{}

	t.csrfFunc = &templatefunctions.CsrfTokenFunc{}
	t.csrfFunc.Inject(t.service, flamingo.NullLogger{})
}

func (t *CsrfTokenFuncTestSuite) TearDown() {
	t.service.AssertExpectations(t.T())
	t.service = nil
	t.csrfFunc = nil
}

func (t *CsrfTokenFuncTestSuite) TestFunc() {
	t.service.EXPECT().Generate(t.session).Return("token").Once()

	function := t.csrfFunc.Func(t.context)
	csrfFunc, ok := function.(func() interface{})
	t.True(ok)

	content := csrfFunc()
	t.Equal("token", content)
}
