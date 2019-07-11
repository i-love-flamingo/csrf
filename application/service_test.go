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

func (t *ServiceTestSuite) TestIsValidPost_GetRequest() {
	t.request.Method = http.MethodGet
	t.True(t.service.IsValidPost(web.CreateRequest(t.request, t.webSession)))
}

func (t *ServiceTestSuite) TestIsValidPost_MalformedToken() {
	t.request.Method = http.MethodPost
	t.False(t.service.IsValidPost(web.CreateRequest(t.request, t.webSession)))
}

func (t *ServiceTestSuite) TestIsValidPost_WrongId() {
	t.T().Skip("no session id changes possible right now")

	token := t.service.Generate(t.webSession)

	t.request.Method = http.MethodPost
	t.request.PostForm = url.Values{
		FormTokenName: []string{token},
	}

	t.False(t.service.IsValidPost(web.CreateRequest(t.request, t.webSession)))
}

func (t *ServiceTestSuite) TestIsValidPost_WrongTime() {
	t.service.ttl = -100000000

	token := t.service.Generate(t.webSession)
	t.request.Method = http.MethodPost
	t.request.PostForm = url.Values{
		FormTokenName: []string{token},
	}

	t.False(t.service.IsValidPost(web.CreateRequest(t.request, t.webSession)))
}

func (t *ServiceTestSuite) TestIsValidPost_Success() {
	token := t.service.Generate(t.webSession)
	t.request.Method = http.MethodPost
	t.request.PostForm = url.Values{
		FormTokenName: []string{token},
	}

	t.True(t.service.IsValidPost(web.CreateRequest(t.request, t.webSession)))
}

func (t *ServiceTestSuite) TestIsValidHeader_NoToken() {
	t.service.Generate(t.webSession)
	t.request.Method = http.MethodGet
	t.False(t.service.IsValidHeader(web.CreateRequest(t.request, t.webSession)))
}

func (t *ServiceTestSuite) TestIsValidHeader_EmptyToken() {
	t.service.Generate(t.webSession)
	t.request.Header = map[string][]string{}
	t.request.Header.Set(HeaderTokenName, "")

	t.False(t.service.IsValidHeader(web.CreateRequest(t.request, t.webSession)))
}

func (t *ServiceTestSuite) TestIsValidHeader_BadToken() {
	t.service.Generate(t.webSession)
	t.request.Header = map[string][]string{}
	t.request.Header.Set(HeaderTokenName, "a88f88e883")

	t.False(t.service.IsValidHeader(web.CreateRequest(t.request, t.webSession)))
}

func (t *ServiceTestSuite) TestIsValidHeader_Success() {
	token := t.service.Generate(t.webSession)
	t.request.Header = map[string][]string{}
	t.request.Header.Set(HeaderTokenName, token)

	t.True(t.service.IsValidHeader(web.CreateRequest(t.request, t.webSession)))
}
