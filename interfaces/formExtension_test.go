package interfaces

import (
	"testing"

	applicationMocks "flamingo.me/csrf/application/mocks"
	"flamingo.me/flamingo/v3/framework/web"
	"flamingo.me/form/domain"
	"github.com/stretchr/testify/suite"
)

type (
	CsrfFormExtensionTestSuite struct {
		suite.Suite

		formExtension *CrsfTokenFormExtension
		service       *applicationMocks.Service

		webRequest *web.Request
	}
)

func TestCsrfFormExtensionTestSuite(t *testing.T) {
	suite.Run(t, &CsrfFormExtensionTestSuite{})
}

func (t *CsrfFormExtensionTestSuite) SetupSuite() {
	t.webRequest = web.CreateRequest(nil, nil)
}

func (t *CsrfFormExtensionTestSuite) SetupTest() {
	t.service = &applicationMocks.Service{}

	t.formExtension = &CrsfTokenFormExtension{}
	t.formExtension.Inject(t.service)
}

func (t *CsrfFormExtensionTestSuite) TearDown() {
	t.service.AssertExpectations(t.T())
	t.service = nil
}

func (t *CsrfFormExtensionTestSuite) TestValidate_WrongToken() {
	t.service.On("IsValid", t.webRequest).Return(false).Once()

	validationInfo, err := t.formExtension.Validate(nil, t.webRequest, nil, nil)

	t.NoError(err)
	t.True(validationInfo.HasGeneralErrors())
	t.Equal([]domain.Error{
		{
			MessageKey:   "formError.crsfToken.invalid",
			DefaultLabel: "Invalid crsf token.",
		},
	}, validationInfo.GetGeneralErrors())
}

func (t *CsrfFormExtensionTestSuite) TestFilter_Success() {
	t.service.On("IsValid", t.webRequest).Return(true).Once()

	validationInfo, err := t.formExtension.Validate(nil, t.webRequest, nil, nil)

	t.NoError(err)
	t.False(validationInfo.HasGeneralErrors())
}
