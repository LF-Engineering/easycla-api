package events

import (
	"github.com/LF-Engineering/lfx-kit/auth"
	"github.com/communitybridge/easycla-api/gen/models"
	"github.com/communitybridge/easycla-api/gen/restapi/operations"
	"github.com/communitybridge/easycla-api/gen/restapi/operations/events"
	"github.com/go-openapi/runtime/middleware"
)

// Configure setups handlers on api with service
func Configure(api *operations.ClaAPI, service Service) {
	api.EventsSearchEventsHandler = events.SearchEventsHandlerFunc(
		func(params events.SearchEventsParams, user *auth.User) middleware.Responder {
			result, err := service.SearchEvents(params.HTTPRequest.Context(), &params)
			if err != nil {
				return events.NewSearchEventsBadRequest().WithPayload(errorResponse(err))
			}
			return events.NewSearchEventsOK().WithPayload(result)
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
