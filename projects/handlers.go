package projects

import (
	"github.com/LF-Engineering/lfx-kit/auth"
	"github.com/communitybridge/easycla-api/gen/models"
	"github.com/communitybridge/easycla-api/gen/restapi/operations"
	"github.com/communitybridge/easycla-api/gen/restapi/operations/project"
	"github.com/go-openapi/runtime/middleware"
)

// Configure setups handlers on api with service
func Configure(api *operations.ClaAPI, service Service) {
	api.ProjectGetProjectHandler = project.GetProjectHandlerFunc(
		func(params project.GetProjectParams, user *auth.User) middleware.Responder {
			result, err := service.GetProject(params.HTTPRequest.Context(), &params)
			if err != nil {
				return project.NewGetProjectBadRequest().WithPayload(errorResponse(err))
			}
			return project.NewGetProjectOK().WithPayload(result)
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
