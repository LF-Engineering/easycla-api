package cla_templates

import (
	"github.com/LF-Engineering/lfx-kit/auth"
	"github.com/communitybridge/easycla-api/gen/models"
	"github.com/communitybridge/easycla-api/gen/restapi/operations"
	"github.com/communitybridge/easycla-api/gen/restapi/operations/cla_templates"
	"github.com/go-openapi/runtime/middleware"
)

// Configure setups handlers on api with service
func Configure(api *operations.ClaAPI, service Service) {
	api.ClaTemplatesCreateCLATemplateHandler = cla_templates.CreateCLATemplateHandlerFunc(
		func(params cla_templates.CreateCLATemplateParams, user *auth.User) middleware.Responder {
			result, err := service.CreateCLATemplate(&params)
			if err != nil {
				return cla_templates.NewCreateCLATemplateBadRequest().WithPayload(errorResponse(err))
			}
			return cla_templates.NewCreateCLATemplateCreated().WithPayload(result)
		})
}

type codedResponse interface {
	Code() string
}

func errorResponse(err error) *models.ErrorResponse {
	code := ""
	if e, ok := err.(codedResponse); ok {
		code = e.Code()
	}

	e := models.ErrorResponse{
		Code:    code,
		Message: err.Error(),
	}

	return &e
}
