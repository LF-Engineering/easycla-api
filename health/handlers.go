// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package health

import (
	"github.com/communitybridge/easycla-api/gen/models"
	"github.com/communitybridge/easycla-api/gen/restapi/operations"
	"github.com/communitybridge/easycla-api/gen/restapi/operations/health"

	"github.com/go-openapi/runtime/middleware"
)

// Configure setups handlers on api with service
func Configure(api *operations.ClaAPI, service Service) {

	api.HealthGetHealthHandler = health.GetHealthHandlerFunc(func(params health.GetHealthParams) middleware.Responder {
		result, err := service.HealthCheck(params.HTTPRequest.Context(), params)
		if err != nil {
			return health.NewGetHealthBadRequest().WithPayload(errorResponse(err))
		}

		return health.NewGetHealthOK().WithPayload(result)
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
