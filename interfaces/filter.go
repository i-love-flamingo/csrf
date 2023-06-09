package interfaces

import (
	"context"
	"net/http"

	"flamingo.me/flamingo/v3/framework/web"

	"flamingo.me/csrf/application"
)

type (
	// CsrfFilter is used for all request if it's defined by configuration.
	CsrfFilter struct {
		responder *web.Responder
		service   application.Service
	}
)

// Inject dependencies
func (f *CsrfFilter) Inject(r *web.Responder, s application.Service) {
	f.responder = r
	f.service = s
}

// Filter is used on each request, and it calls csrf Service to validate token from request.
func (f *CsrfFilter) Filter(ctx context.Context, r *web.Request, w http.ResponseWriter, chain *web.FilterChain) web.Result {
	if !f.service.IsValid(r) {
		return f.responder.Forbidden(ErrInvalidToken)
	}

	return chain.Next(ctx, r, w)
}
