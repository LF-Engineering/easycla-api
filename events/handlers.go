package events

import (
	"github.com/communitybridge/easycla-api/gen/models"
	"github.com/communitybridge/easycla-api/gen/restapi/operations"
	"github.com/communitybridge/easycla-api/gen/restapi/operations/events"
	"github.com/go-openapi/runtime/middleware"
)

// Configure setups handlers on api with service
func Configure(api *operations.ClaAPI, service Service) {
	api.EventsListEventsHandler = events.ListEventsHandlerFunc(
		func(params events.ListEventsParams) middleware.Responder {
			result, err := service.ListEvents(params.HTTPRequest.Context(), &params)
			if err != nil {
				return events.NewListEventsBadRequest().WithPayload(errorResponse(err))
			}
			return events.NewListEventsOK().WithPayload(result)
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
