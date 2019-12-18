package repositories

import (
	"strconv"

	"github.com/LF-Engineering/lfx-kit/auth"
	"github.com/communitybridge/easycla-api/gen/models"
	"github.com/communitybridge/easycla-api/gen/restapi/operations"
	"github.com/communitybridge/easycla-api/gen/restapi/operations/repositories"
	"github.com/go-openapi/runtime/middleware"
)

// Configure setups handlers on api with service
func Configure(api *operations.ClaAPI, service Service) {
	api.RepositoriesCreateRepositoriesHandler = repositories.CreateRepositoriesHandlerFunc(
		func(params repositories.CreateRepositoriesParams, user *auth.User) middleware.Responder {
			if !(user.Admin || user.IsUserAuthorizedForAll(auth.Project) || user.IsUserAuthorized(auth.Project, *params.RepositoriesInput.ProjectID)) {
				return repositories.NewCreateRepositoriesUnauthorized()
			}
			result, err := service.CreateRepositories(user, &params)
			if err != nil {
				return repositories.NewCreateRepositoriesBadRequest().WithPayload(errorResponse(err))
			}
			return repositories.NewCreateRepositoriesOK().WithPayload(result)
		})
	api.RepositoriesDeleteRepositoriesHandler = repositories.DeleteRepositoriesHandlerFunc(
		func(params repositories.DeleteRepositoriesParams, user *auth.User) middleware.Responder {
			if !(user.Admin || user.IsUserAuthorizedForAll(auth.Project) || user.IsUserAuthorized(auth.Project, *params.RepositoriesInput.ProjectID)) {
				return repositories.NewDeleteRepositoriesUnauthorized()
			}
			err := service.DeleteRepositories(user, &params)
			if err != nil {
				return repositories.NewDeleteRepositoriesBadRequest().WithPayload(&models.ErrorResponse{
					Code:    strconv.Itoa(repositories.DeleteRepositoriesBadRequestCode),
					Message: err.Error(),
				})
			}
			return repositories.NewDeleteRepositoriesOK()
		})
	api.RepositoriesListRepositoriesHandler = repositories.ListRepositoriesHandlerFunc(
		func(params repositories.ListRepositoriesParams, user *auth.User) middleware.Responder {
			if !(user.Admin || user.IsUserAuthorizedForAll(auth.Project) || user.IsUserAuthorized(auth.Project, params.ProjectID)) {
				return repositories.NewCreateRepositoriesUnauthorized()
			}
			result, err := service.ListRepositories(&params)
			if err != nil {
				return repositories.NewListRepositoriesBadRequest().WithPayload(errorResponse(err))
			}
			return repositories.NewListRepositoriesOK().WithPayload(result)
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
