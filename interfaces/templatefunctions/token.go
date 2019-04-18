package templatefunctions

import (
	"context"

	"flamingo.me/csrf/application"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// CsrfTokenFunc delivers template function to render only csrf token value.
	CsrfTokenFunc struct {
		service application.Service
		logger  flamingo.Logger
	}
)

// Inject dependencies
func (f *CsrfTokenFunc) Inject(s application.Service, l flamingo.Logger) {
	f.service = s
	f.logger = l
}

// Func delivers template function.
func (f *CsrfTokenFunc) Func(ctx context.Context) interface{} {
	return func() interface{} {
		s := web.SessionFromContext(ctx)
		if s == nil {
			f.logger.WithField("csrf", "templateFunc").Error("can't find session")
			return ""
		}

		return f.service.Generate(s)
	}
}
