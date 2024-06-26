package interfaces

import (
	"context"

	"flamingo.me/flamingo/v3/framework/web"

	"flamingo.me/csrf/application"
)

type (
	// CsrfMiddleware is middleware which can be attached to particular action from controller to validate csrf token.
	CsrfMiddleware struct {
		responder *web.Responder
		service   application.Service
	}
)

// Inject dependencies
func (m *CsrfMiddleware) Inject(r *web.Responder, s application.Service) {
	m.responder = r
	m.service = s
}

// Secured validates csrf token by using csrf Service if controller action is wrapped with this middleware.
func (m *CsrfMiddleware) Secured(action web.Action) web.Action {
	return func(ctx context.Context, r *web.Request) web.Result {
		if !m.service.IsValidPost(r) {
			return m.responder.Forbidden(ErrInvalidToken)
		}

		return action(ctx, r)
	}
}

// SecuredHeader validates csrf token from header field by using csrf Service if controller action is wrapped with this middleware.
func (m *CsrfMiddleware) SecuredHeader(action web.Action) web.Action {
	return func(ctx context.Context, r *web.Request) web.Result {
		if !m.service.IsValidHeader(r) {
			return m.responder.Forbidden(ErrInvalidToken)
		}

		return action(ctx, r)
	}
}
