package cla_groups

import (
	"github.com/communitybridge/easycla-api/gen/models"
	"github.com/communitybridge/easycla-api/gen/restapi/operations"
	"github.com/communitybridge/easycla-api/gen/restapi/operations/cla_groups"
	"github.com/go-openapi/runtime/middleware"
)

// Configure setups handlers on api with service
func Configure(api *operations.ClaAPI, service Service) {
	api.ClaGroupsCreateCLAGroupHandler = cla_groups.CreateCLAGroupHandlerFunc(
		func(params cla_groups.CreateCLAGroupParams) middleware.Responder {
			result, err := service.CreateCLAGroup(&params)
			if err != nil {
				return cla_groups.NewCreateCLAGroupBadRequest().WithPayload(errorResponse(err))
			}
			return cla_groups.NewCreateCLAGroupCreated().WithPayload(result)
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
