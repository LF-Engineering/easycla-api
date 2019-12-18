package repositories_test

import (
	"context"
	"testing"
	"time"

	"github.com/LF-Engineering/lfx-kit/auth"
	"github.com/communitybridge/easycla-api/gen/models"
	"github.com/communitybridge/easycla-api/gen/restapi/operations/events"
	params "github.com/communitybridge/easycla-api/gen/restapi/operations/repositories"
	"github.com/communitybridge/easycla-api/repositories"
	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
)

func Test_CreateRepositories(t *testing.T) {
	t.Run("valid input", func(t *testing.T) {
		prepareTestDatabase()
		deleteAllRepositories()
		timeBeforeTest := time.Now().Unix()
		args := []*models.CreateRepositoriesInputRepositoriesItems0{
			{
				Enabled:          newBool(NodeExporter.Enabled),
				ExternalID:       newString(NodeExporter.ExternalID),
				Name:             newString(NodeExporter.Name),
				OrganizationName: newString(NodeExporter.OrganizationName),
				RepositoryType:   NodeExporter.RepositoryType,
				URL:              newString(NodeExporter.URL),
			},
			{
				Enabled:          newBool(Prometheus.Enabled),
				ExternalID:       newString(Prometheus.ExternalID),
				Name:             newString(Prometheus.Name),
				OrganizationName: newString(Prometheus.OrganizationName),
				RepositoryType:   Prometheus.RepositoryType,
				URL:              newString(Prometheus.URL),
			},
		}
		result, err := repositoriesService.CreateRepositories(&auth.User{}, &params.CreateRepositoriesParams{
			RepositoriesInput: &models.CreateRepositoriesInput{
				ClaGroupID:   newString(PrometheusClaGroupID),
				ProjectID:    newString(PrometheusProjectID),
				Repositories: args,
			},
		})
		if !assert.Nil(t, err, "error should be nil") {
			return
		}
		if !assert.NotNil(t, result, "result should not be nil") {
			return
		}
		if !assert.Equal(t, 2, len(result.Repositories), "result should return newly created repositories") {
			return
		}
		for i, got := range result.Repositories {
			if !assert.Equal(t, true, got.CreatedAt >= timeBeforeTest, "created_at should be valid") {
				return
			}
			if !assert.Equal(t, true, strfmt.IsUUID4(got.ID), "id should be valid") {
				return
			}
			if !assert.Equal(t, got.Enabled, *args[i].Enabled) {
				return
			}
			if !assert.Equal(t, got.ExternalID, *args[i].ExternalID) {
				return
			}
			if !assert.Equal(t, got.Name, *args[i].Name) {
				return
			}
			if !assert.Equal(t, got.OrganizationName, *args[i].OrganizationName) {
				return
			}
			if !assert.Equal(t, got.RepositoryType, args[i].RepositoryType) {
				return
			}
			if !assert.Equal(t, got.URL, *args[i].URL) {
				return
			}
		}
		if !assert.Equal(t, int64(2), numberOfRepositories(), "repository not created") {
			return
		}
		elist, err := eventsService.SearchEvents(context.TODO(), &events.SearchEventsParams{})
		if !assert.Nil(t, err) {
			return
		}
		if !assert.Equal(t, 1, len(elist.Events), "event not created") {
			return
		}
		if !assert.Equal(t, string(repositories.EventTypeRepositoriesCreated), elist.Events[0].EventType, "invalid event type") {
			return
		}
		if !assert.Equal(t, PrometheusProjectID, elist.Events[0].ProjectID, "invalid project_id") {
			return
		}

	})
	t.Run("invalid input", func(t *testing.T) {
		prepareTestDatabase()
		deleteAllRepositories()
		args := []*models.CreateRepositoriesInputRepositoriesItems0{
			{
				Enabled:          newBool(Prometheus.Enabled),
				ExternalID:       newString(Prometheus.ExternalID),
				Name:             newString(Prometheus.Name),
				OrganizationName: newString(Prometheus.OrganizationName),
				RepositoryType:   Prometheus.RepositoryType,
				URL:              newString(Prometheus.URL),
			},
		}
		result, err := repositoriesService.CreateRepositories(&auth.User{}, &params.CreateRepositoriesParams{
			RepositoriesInput: &models.CreateRepositoriesInput{
				ClaGroupID:   newString("00000000-0000-0000-0000-000000000000"),
				ProjectID:    newString(PrometheusProjectID),
				Repositories: args,
			},
		})
		if !assert.NotNil(t, err, "error should not be nil") {
			return
		}
		if !assert.Equal(t, err, repositories.ErrInvalidClgGroupAndProjectID, "got wrong error") {
			return
		}
		if !assert.Nil(t, result, "result should  be nil") {
			return
		}
		if !assert.Equal(t, int64(0), numberOfRepositories(), "repository should not be created") {
			return
		}
	})
}
