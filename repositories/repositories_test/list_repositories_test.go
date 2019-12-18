package repositories_test

import (
	"fmt"
	"testing"

	"github.com/communitybridge/easycla-api/repositories"

	"github.com/communitybridge/easycla-api/gen/models"
	params "github.com/communitybridge/easycla-api/gen/restapi/operations/repositories"
	"github.com/stretchr/testify/assert"
)

func Test_ListRepositories(t *testing.T) {
	prepareTestDatabase()
	tests := []struct {
		name    string
		args    *params.ListRepositoriesParams
		want    *models.RepositoryList
		wantErr bool
	}{
		{
			name: "filter project",
			args: &params.ListRepositoriesParams{
				ProjectID: PrometheusProjectID,
			},
			want:    &models.RepositoryList{Repositories: []*models.Repository{ClientJava, HAProxyExportter, Prometheus, NodeExporter}},
			wantErr: false,
		},
		{
			name: "filter project and cla_group",
			args: &params.ListRepositoriesParams{
				ProjectID:  PrometheusProjectID,
				ClaGroupID: newString(PrometheusClaGroupID),
			},
			want:    &models.RepositoryList{Repositories: []*models.Repository{ClientJava, HAProxyExportter, Prometheus}},
			wantErr: false,
		},
		{
			name: "filter type of repo",
			args: &params.ListRepositoriesParams{
				ProjectID:      PrometheusProjectID,
				RepositoryType: newString(string(repositories.RepositoryTypeGerrit)),
			},
			want:    &models.RepositoryList{Repositories: []*models.Repository{NodeExporter}},
			wantErr: false,
		},
		{
			name: "filter by organization",
			args: &params.ListRepositoriesParams{
				ProjectID:              PrometheusProjectID,
				RepositoryOrganization: newString("prometheut"),
			},
			want:    &models.RepositoryList{Repositories: []*models.Repository{NodeExporter}},
			wantErr: false,
		},
		{
			name: "pagination",
			args: &params.ListRepositoriesParams{
				ProjectID: PrometheusProjectID,
				Offset:    newInt64(1),
				PageSize:  newInt64(2),
			},
			want:    &models.RepositoryList{Repositories: []*models.Repository{HAProxyExportter, Prometheus}},
			wantErr: false,
		},
		{
			name: "sort_order",
			args: &params.ListRepositoriesParams{
				ProjectID: PrometheusProjectID,
				OrderBy:   newString("organization_name"),
				SortOrder: newString("desc"),
			},
			want:    &models.RepositoryList{Repositories: []*models.Repository{NodeExporter, Prometheus, HAProxyExportter, ClientJava}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repositoriesService.ListRepositories(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListCLAGroups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, len(tt.want.Repositories), len(got.Repositories)) {
				return
			}
			for i := range got.Repositories {
				if !assert.Equal(t,
					tt.want.Repositories[i], got.Repositories[i],
					expectedRepositoriesOrderMsg(got, tt.want)) {
					return
				}
			}
		})
	}
}

func expectedRepositoriesOrderMsg(got *models.RepositoryList, want *models.RepositoryList) string {
	var gotOrder, wantOrder []string
	for _, v := range got.Repositories {
		gotOrder = append(gotOrder, v.Name)
	}
	for _, v := range want.Repositories {
		wantOrder = append(wantOrder, v.Name)
	}
	return fmt.Sprintf("got [%#v], expected [%#v]", gotOrder, wantOrder)
}
