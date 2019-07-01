package application

import (
	"net/http"
	"net/url"
	"testing"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/stretchr/testify/suite"
)

type (
	ServiceTestSuite struct {
		suite.Suite

		service *ServiceImpl

		webSession *web.Session
		request    *http.Request
	}
)

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, &ServiceTestSuite{})
}

func (t *ServiceTestSuite) SetupTest() {
	t.service = &ServiceImpl{}
	t.service.Inject(flamingo.NullLogger{}, &struct {
		Secret string  `inject:"config:csrf.secret"`
		TTL    float64 `inject:"config:csrf.ttl"`
	}{
		Secret: "6368616e676520746869732070617373776f726420746f206120736563726574",
		TTL:    900,
	})

	t.webSession = web.EmptySession()
	t.request = &http.Request{}
}

func (t *ServiceTestSuite) TearDown() {
	t.webSession = nil
	t.service = nil
	t.request = nil
}

func (t *ServiceTestSuite) TestGenerate_WrongKey() {
	t.service.secret = []byte{}
	t.Empty(t.service.Generate(t.webSession))
}

func (t *ServiceTestSuite) TestGenerate_RightKey() {
	t.NotEmpty(t.service.Generate(t.webSession))

	t.NotEmpty(t.service.Generate(t.webSession))
}

func (t *ServiceTestSuite) TestIsValid_GetRequest() {
	t.request.Method = http.MethodGet
	t.True(t.service.IsValid(web.CreateRequest(t.request, t.webSession)))
}

func (t *ServiceTestSuite) TestIsValid_MalformedToken() {
	t.request.Method = http.MethodPost
	t.False(t.service.IsValid(web.CreateRequest(t.request, t.webSession)))
}

func (t *ServiceTestSuite) TestIsValid_WrongId() {
	t.T().Skip("no session id changes possible right now")

	token := t.service.Generate(t.webSession)

	t.request.Method = http.MethodPost
	t.request.PostForm = url.Values{
		FormTokenName: []string{token},
	}

	t.False(t.service.IsValid(web.CreateRequest(t.request, t.webSession)))
}

func (t *ServiceTestSuite) TestIsValid_WrongTime() {
	t.service.ttl = -100000000

	token := t.service.Generate(t.webSession)
	t.request.Method = http.MethodPost
	t.request.PostForm = url.Values{
		FormTokenName: []string{token},
	}

	t.False(t.service.IsValid(web.CreateRequest(t.request, t.webSession)))
}

func (t *ServiceTestSuite) TestIsValid_FormTokenSuccess() {
	token := t.service.Generate(t.webSession)
	t.request.Method = http.MethodPost
	t.request.PostForm = url.Values{
		FormTokenName: []string{token},
	}

	t.True(t.service.IsValid(web.CreateRequest(t.request, t.webSession)))
}

func (t *ServiceTestSuite) TestIsValid_HeaderTokenSuccess() {
	token := t.service.Generate(t.webSession)
	t.request.Method = http.MethodPost
	t.request.Header = http.Header{
		HeaderTokenName: []string{token},
	}

	t.True(t.service.IsValid(web.CreateRequest(t.request, t.webSession)))
}
