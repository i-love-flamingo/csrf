package csrf

import (
	"flamingo.me/csrf/application"
	"flamingo.me/csrf/interfaces"
	"flamingo.me/csrf/interfaces/templatefunctions"

	"fmt"

	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
	"flamingo.me/form/domain"
)

const csrfTTL = 900.0

// Module for core/csrfPreventionFilter
type Module struct {
	All bool `inject:"config:csrf.all"`
}

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	injector.Bind((*application.Service)(nil)).To(application.ServiceImpl{})
	flamingo.BindTemplateFunc(injector, "csrfToken", new(templatefunctions.CsrfTokenFunc))
	flamingo.BindTemplateFunc(injector, "csrfInput", new(templatefunctions.CsrfInputFunc))

	injector.BindMap((*domain.FormExtension)(nil), "formExtension.csrfToken").To(interfaces.CsrfTokenFormExtension{})

	if m.All {
		injector.BindMulti((*web.Filter)(nil)).To(interfaces.CsrfFilter{})
	}
}

// CueConfig for the module
func (m *Module) CueConfig() string {
	// language=cue
	return fmt.Sprintf(`
csrf: {
	all: bool | *false
	secret: string | *"somethingSuperSecret"
	ttl: number | *%f
}
`, csrfTTL)
}
