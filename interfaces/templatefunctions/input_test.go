package templatefunctions_test

import (
	"context"
	"testing"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/stretchr/testify/suite"

	applicationMocks "flamingo.me/csrf/application/mocks"
	"flamingo.me/csrf/interfaces/templatefunctions"
)

type (
	CsrfInputFuncTestSuite struct {
		suite.Suite
		csrfFunc *templatefunctions.CsrfInputFunc
		service  *applicationMocks.Service
		context  context.Context
		session  *web.Session
	}
)

func TestCsrfInputFuncTestSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, &CsrfInputFuncTestSuite{})
}

func (t *CsrfInputFuncTestSuite) SetupSuite() {
	t.session = web.EmptySession()
	t.context = web.ContextWithSession(context.Background(), t.session)
}

func (t *CsrfInputFuncTestSuite) SetupTest() {
	t.service = &applicationMocks.Service{}

	t.csrfFunc = &templatefunctions.CsrfInputFunc{}
	t.csrfFunc.Inject(t.service, flamingo.NullLogger{})
}

func (t *CsrfInputFuncTestSuite) TearDown() {
	t.service.AssertExpectations(t.T())
	t.service = nil
	t.csrfFunc = nil
}

func (t *CsrfInputFuncTestSuite) TestFunc() {
	t.service.EXPECT().Generate(t.session).Return("token").Once()

	function := t.csrfFunc.Func(t.context)
	csrfFunc, ok := function.(func() interface{})
	t.True(ok)

	content := csrfFunc()
	t.Equal(`<input type="hidden" name="csrftoken" value="token" data-qa="csrfInput" />`, content)
}
