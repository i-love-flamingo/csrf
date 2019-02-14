package interfaces

import (
	"context"

	"flamingo.me/csrf/application"
	"flamingo.me/flamingo/v3/framework/web"
	"flamingo.me/form/domain"
)

type (
	CrsfTokenFormExtension struct {
		service application.Service
	}
)

func (f *CrsfTokenFormExtension) Inject(service application.Service) {
	f.service = service
}

func (f *CrsfTokenFormExtension) Validate(_ context.Context, req *web.Request, _ domain.ValidatorProvider, _ interface{}) (*domain.ValidationInfo, error) {
	validationInfo := domain.ValidationInfo{}

	if !f.service.IsValid(req) {
		validationInfo.AddGeneralError("formError.crsfToken.invalid", "Invalid crsf token.")
	}

	return &validationInfo, nil
}
