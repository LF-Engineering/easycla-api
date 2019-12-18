package repositories_test

import (
	"context"
	"testing"

	"github.com/communitybridge/easycla-api/repositories"

	"github.com/stretchr/testify/assert"

	"github.com/communitybridge/easycla-api/gen/models"

	"github.com/LF-Engineering/lfx-kit/auth"
	"github.com/communitybridge/easycla-api/gen/restapi/operations/events"
	params "github.com/communitybridge/easycla-api/gen/restapi/operations/repositories"
)

func Test_DeleteRepositories(t *testing.T) {
	prepareTestDatabase()
	count := numberOfRepositories()
	isExist := isRepositoryPresent(Minikube.ID)
	if !assert.Equal(t, true, isExist) {
		return
	}
	err := repositoriesService.DeleteRepositories(&auth.User{}, &params.DeleteRepositoriesParams{
		RepositoriesInput: &models.DeleteRepositoriesInput{
			ClaGroupID:    newString(KubernetesClaGroupID),
			ProjectID:     newString(KubernetesProjectID),
			RepositoryIds: []string{Minikube.ID},
		},
	})
	if !assert.Nil(t, err, "error should be nil") {
		return
	}
	if !assert.Equal(t, count-1, numberOfRepositories(), "repository not deleted") {
		return
	}
	if !assert.Equal(t, false, isRepositoryPresent(Minikube.ID), "correct repository not deleted") {
		return
	}
	elist, err := eventsService.SearchEvents(context.TODO(), &events.SearchEventsParams{})
	if !assert.Nil(t, err) {
		return
	}
	if !assert.Equal(t, 1, len(elist.Events), "event not created") {
		return
	}
	if !assert.Equal(t, string(repositories.EventTypeRepositoriesDeleted), elist.Events[0].EventType, "invalid event type") {
		return
	}
	if !assert.Equal(t, Minikube.ProjectID, elist.Events[0].ProjectID, "invalid project_id") {
		return
	}
}
