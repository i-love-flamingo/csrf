package csrf

import (
	"flamingo.me/csrf/application"
	"flamingo.me/csrf/interfaces"
	"flamingo.me/csrf/interfaces/templatefunctions"
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
	"flamingo.me/form/domain"
)

// Module for core/csrfPreventionFilter
type Module struct {
	All bool `inject:"config:csrf.all"`
}

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	injector.Bind((*application.Service)(nil)).To(application.ServiceImpl{})
	flamingo.BindTemplateFunc(injector, "csrfToken", new(templatefunctions.CsrfTokenFunc))
	flamingo.BindTemplateFunc(injector, "csrfInput", new(templatefunctions.CsrfInputFunc))

	injector.BindMap((*domain.FormExtension)(nil), "formExtension.csrfToken").To(interfaces.CrsfTokenFormExtension{})

	if m.All {
		injector.BindMulti((*web.Filter)(nil)).To(interfaces.CsrfFilter{})
	}
}

// DefaultConfig for this module
func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"csrf.all":    false,
		"csrf.secret": "somethingSuperSecret",
		"csrf.ttl":    900.0,
	}
}
