package orgs

import (
	"github.com/LF-Engineering/lfx-kit/auth"
	"github.com/communitybridge/easycla-api/gen/models"
	"github.com/communitybridge/easycla-api/gen/restapi/operations"
	"github.com/communitybridge/easycla-api/gen/restapi/operations/organization"
	"github.com/go-openapi/runtime/middleware"
)

// Configure setups handlers on api with service
func Configure(api *operations.ClaAPI, service Service) {
	api.OrganizationGetOrgFoundationsHandler = organization.GetOrgFoundationsHandlerFunc(
		func(params organization.GetOrgFoundationsParams, user *auth.User) middleware.Responder {
			result, err := service.GetOrgFoundations(params.HTTPRequest.Context(), &params)
			if err != nil {
				return organization.NewGetOrgFoundationsBadRequest().WithPayload(errorResponse(err))
			}
			return organization.NewGetOrgFoundationsOK().WithPayload(result)
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
