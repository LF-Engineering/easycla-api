package cla_templates

import (
	"strconv"

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
			if !(user.Admin || user.IsUserAuthorizedForAll(auth.Project)) {
				return cla_templates.NewCreateCLATemplateUnauthorized()
			}
			result, err := service.CreateCLATemplate(&params)
			if err != nil {
				return cla_templates.NewCreateCLATemplateBadRequest().WithPayload(errorResponse(err))
			}
			return cla_templates.NewCreateCLATemplateCreated().WithPayload(result)
		})

	api.ClaTemplatesGetCLATemplateHandler = cla_templates.GetCLATemplateHandlerFunc(
		func(params cla_templates.GetCLATemplateParams, user *auth.User) middleware.Responder {
			if len(user.ResourceIDsByType(auth.Project)) == 0 {
				return cla_templates.NewGetCLATemplateUnauthorized()
			}
			result, err := service.GetCLATemplate(&params)
			if err != nil {
				if err == ErrClaTemplateNotFound {
					return cla_templates.NewGetCLATemplateNotFound().WithPayload(&models.ErrorResponse{
						Code:    strconv.Itoa(cla_templates.DeleteCLATemplateNotFoundCode),
						Message: err.Error(),
					})
				}
				return cla_templates.NewGetCLATemplateBadRequest().WithPayload(errorResponse(err))
			}
			return cla_templates.NewGetCLATemplateOK().WithPayload(result)
		})

	api.ClaTemplatesUpdateCLATemplateHandler = cla_templates.UpdateCLATemplateHandlerFunc(
		func(params cla_templates.UpdateCLATemplateParams, user *auth.User) middleware.Responder {
			if !(user.Admin || user.IsUserAuthorizedForAll(auth.Project)) {
				return cla_templates.NewUpdateCLATemplateUnauthorized()
			}
			result, err := service.UpdateCLATemplate(&params)
			if err != nil {
				if err == ErrClaTemplateNotFound {
					return cla_templates.NewUpdateCLATemplateNotFound().WithPayload(&models.ErrorResponse{
						Code:    strconv.Itoa(cla_templates.DeleteCLATemplateNotFoundCode),
						Message: err.Error(),
					})
				}
				return cla_templates.NewUpdateCLATemplateBadRequest().WithPayload(errorResponse(err))
			}
			return cla_templates.NewUpdateCLATemplateOK().WithPayload(result)
		})

	api.ClaTemplatesDeleteCLATemplateHandler = cla_templates.DeleteCLATemplateHandlerFunc(
		func(params cla_templates.DeleteCLATemplateParams, user *auth.User) middleware.Responder {
			if !(user.Admin || user.IsUserAuthorizedForAll(auth.Project)) {
				return cla_templates.NewDeleteCLATemplateUnauthorized()
			}
			err := service.DeleteCLATemplate(&params)
			if err != nil {
				if err == ErrClaTemplateNotFound {
					return cla_templates.NewDeleteCLATemplateNotFound().WithPayload(&models.ErrorResponse{
						Code:    strconv.Itoa(cla_templates.DeleteCLATemplateNotFoundCode),
						Message: err.Error(),
					})
				}
				return cla_templates.NewDeleteCLATemplateBadRequest().WithPayload(errorResponse(err))
			}
			return cla_templates.NewDeleteCLATemplateOK()
		})

	api.ClaTemplatesListCLATemplatesHandler = cla_templates.ListCLATemplatesHandlerFunc(
		func(params cla_templates.ListCLATemplatesParams, user *auth.User) middleware.Responder {
			if len(user.ResourceIDsByType(auth.Project)) == 0 {
				return cla_templates.NewListCLATemplatesUnauthorized()
			}
			result, err := service.ListCLATemplate(&params)
			if err != nil {
				return cla_templates.NewListCLATemplatesBadRequest().WithPayload(errorResponse(err))
			}
			return cla_templates.NewListCLATemplatesOK().WithPayload(result)
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
