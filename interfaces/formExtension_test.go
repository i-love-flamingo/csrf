package interfaces_test

import (
	"context"
	"flamingo.me/csrf/interfaces"
	"testing"

	applicationMocks "flamingo.me/csrf/application/mocks"
	"flamingo.me/flamingo/v3/framework/web"
	"flamingo.me/form/domain"
	"github.com/stretchr/testify/suite"
)

type (
	CsrfFormExtensionTestSuite struct {
		suite.Suite

		formExtension *interfaces.CsrfTokenFormExtension
		service       *applicationMocks.Service

		webRequest *web.Request
	}
)

func TestCsrfFormExtensionTestSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, &CsrfFormExtensionTestSuite{})
}

func (t *CsrfFormExtensionTestSuite) SetupSuite() {
	t.webRequest = web.CreateRequest(nil, nil)
}

func (t *CsrfFormExtensionTestSuite) SetupTest() {
	t.service = &applicationMocks.Service{}

	t.formExtension = &interfaces.CsrfTokenFormExtension{}
	t.formExtension.Inject(t.service)
}

func (t *CsrfFormExtensionTestSuite) TearDown() {
	t.service.AssertExpectations(t.T())
	t.service = nil
}

func (t *CsrfFormExtensionTestSuite) TestValidate_WrongToken() {
	t.service.EXPECT().IsValid(t.webRequest).Return(false).Once()

	validationInfo, err := t.formExtension.Validate(context.Background(), t.webRequest, nil, nil)

	t.NoError(err)
	t.True(validationInfo.HasGeneralErrors())
	t.Equal([]domain.Error{
		{
			MessageKey:   "formError.csrfToken.invalid",
			DefaultLabel: "Invalid csrf token.",
		},
	}, validationInfo.GetGeneralErrors())
}

func (t *CsrfFormExtensionTestSuite) TestFilter_Success() {
	t.service.EXPECT().IsValid(t.webRequest).Return(true).Once()

	validationInfo, err := t.formExtension.Validate(context.Background(), t.webRequest, nil, nil)

	t.NoError(err)
	t.False(validationInfo.HasGeneralErrors())
}
