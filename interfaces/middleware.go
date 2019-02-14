package interfaces

import (
	"context"

	"flamingo.me/csrf/application"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/pkg/errors"
)

type (
	CsrfMiddleware struct {
		responder *web.Responder
		service   application.Service
	}
)

func (m *CsrfMiddleware) Inject(r *web.Responder, s application.Service) {
	m.responder = r
	m.service = s
}

func (m *CsrfMiddleware) Secured(action web.Action) web.Action {
	return func(ctx context.Context, r *web.Request) web.Result {
		if !m.service.IsValid(r) {
			return m.responder.Forbidden(errors.New("csrf_token is not valid"))
		}

		return action(ctx, r)
	}
}
