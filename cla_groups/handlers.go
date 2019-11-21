package cla_groups

import (
	"strconv"

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

	api.ClaGroupsDeleteCLAGroupHandler = cla_groups.DeleteCLAGroupHandlerFunc(
		func(params cla_groups.DeleteCLAGroupParams) middleware.Responder {
			err := service.DeleteCLAGroup(&params)
			if err != nil {
				return cla_groups.NewDeleteCLAGroupBadRequest().WithPayload(errorResponse(err))
			}
			return cla_groups.NewDeleteCLAGroupOK()
		})

	api.ClaGroupsUpdateCLAGroupHandler = cla_groups.UpdateCLAGroupHandlerFunc(
		func(params cla_groups.UpdateCLAGroupParams) middleware.Responder {
			err := service.UpdateCLAGroup(&params)
			if err != nil {
				if err == ErrClaGroupNotFound {
					return cla_groups.NewDeleteCLAGroupNotFound().WithPayload(&models.ErrorResponse{
						Code:    strconv.Itoa(cla_groups.UpdateCLAGroupNotFoundCode),
						Message: err.Error(),
					})
				}
				return cla_groups.NewUpdateCLAGroupBadRequest().WithPayload(errorResponse(err))
			}
			return cla_groups.NewUpdateCLAGroupOK()
		})

	api.ClaGroupsListCLAGroupsHandler = cla_groups.ListCLAGroupsHandlerFunc(
		func(params cla_groups.ListCLAGroupsParams) middleware.Responder {
			response, err := service.ListCLAGroups(&params)
			if err != nil {
				return cla_groups.NewListCLAGroupsBadRequest().WithPayload(errorResponse(err))
			}
			return cla_groups.NewListCLAGroupsOK().WithPayload(response)
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
