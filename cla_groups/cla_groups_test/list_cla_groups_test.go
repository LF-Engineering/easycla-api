package cla_groups_test

import (
	"fmt"
	"testing"

	"github.com/communitybridge/easycla-api/gen/models"
	"github.com/communitybridge/easycla-api/gen/restapi/operations/cla_groups"
	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
)

var Dataloader = &models.ClaGroup{
	CclaEnabled:     true,
	ClaGroupName:    "Dataloader",
	CreatedAt:       2,
	ProjectID:       "GraphQL",
	IclaEnabled:     false,
	ID:              "d9dc5834-3d9a-4d04-abb6-4a36ed378304",
	ProjectManagers: []strfmt.UUID{"413f4711-a3c3-4635-9dad-a0ba58694203", "413f4711-a3c3-4635-9dad-a0ba58694204"},
	UpdatedAt:       10,
}

var Kubernetes = &models.ClaGroup{
	CclaEnabled:     true,
	ClaGroupName:    "Kubernetes",
	CreatedAt:       1,
	ProjectID:       "CNCF",
	IclaEnabled:     true,
	ID:              "73007448-6192-403b-86f2-9ee00ea07060",
	ProjectManagers: []strfmt.UUID{"413f4711-a3c3-4635-9dad-a0ba58694201"},
	UpdatedAt:       30,
}

var Prometheus = &models.ClaGroup{
	CclaEnabled:     true,
	ClaGroupName:    "Prometheus",
	CreatedAt:       3,
	ProjectID:       "CNCF",
	IclaEnabled:     false,
	ID:              "086bb357-349a-492b-b64d-230a357f3712",
	ProjectManagers: []strfmt.UUID{"413f4711-a3c3-4635-9dad-a0ba58694201", "413f4711-a3c3-4635-9dad-a0ba58694202"},
	UpdatedAt:       20,
}

func Test_ListCLAGroups(t *testing.T) {
	prepareTestDatabase()
	var cncf = "CNCF"
	var projectManagerID = "413f4711-a3c3-4635-9dad-a0ba58694202"
	var pageSize int64 = 1
	var offset int64 = 1
	var orderByUpdatedAt = "updated_at"
	var sortOrderDesc = "desc"

	tests := []struct {
		name    string
		args    *cla_groups.ListCLAGroupsParams
		want    *models.ClaGroupList
		wantErr bool
	}{
		{
			name: "default",
			args: &cla_groups.ListCLAGroupsParams{},
			want: &models.ClaGroupList{
				ClaGroups: []*models.ClaGroup{Dataloader, Kubernetes, Prometheus},
			},
			wantErr: false,
		},
		{
			name: "filter by project_id",
			args: &cla_groups.ListCLAGroupsParams{
				ProjectID: &cncf,
			},
			want: &models.ClaGroupList{
				ClaGroups: []*models.ClaGroup{Kubernetes, Prometheus},
			},
			wantErr: false,
		},
		{
			name: "filter by project_manager_id",
			args: &cla_groups.ListCLAGroupsParams{
				ProjectManagerID: &projectManagerID,
			},
			want: &models.ClaGroupList{
				ClaGroups: []*models.ClaGroup{Prometheus},
			},
			wantErr: false,
		},
		{
			name: "filter by project_id, project_manager_id",
			args: &cla_groups.ListCLAGroupsParams{
				ProjectID:        &cncf,
				ProjectManagerID: &projectManagerID,
			},
			want: &models.ClaGroupList{
				ClaGroups: []*models.ClaGroup{Prometheus},
			},
			wantErr: false,
		},
		{
			name: "pagination",
			args: &cla_groups.ListCLAGroupsParams{
				PageSize: &pageSize,
				Offset:   &offset,
			},
			want: &models.ClaGroupList{
				ClaGroups: []*models.ClaGroup{Kubernetes},
			},
			wantErr: false,
		},
		{
			name: "order_by",
			args: &cla_groups.ListCLAGroupsParams{
				OrderBy: &orderByUpdatedAt,
			},
			want: &models.ClaGroupList{
				ClaGroups: []*models.ClaGroup{Dataloader, Prometheus, Kubernetes},
			},
			wantErr: false,
		},
		{
			name: "sort_order",
			args: &cla_groups.ListCLAGroupsParams{
				SortOrder: &sortOrderDesc,
			},
			want: &models.ClaGroupList{
				ClaGroups: []*models.ClaGroup{Prometheus, Kubernetes, Dataloader},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := claGroupsService.ListCLAGroups(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListCLAGroups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, len(tt.want.ClaGroups), len(got.ClaGroups)) {
				return
			}
			for i := range got.ClaGroups {
				if !assert.Equal(t,
					tt.want.ClaGroups[i], got.ClaGroups[i],
					expectedClaGroupOrderMsg(got.ClaGroups, tt.want.ClaGroups)) {
					return
				}
			}
		})
	}
}

func expectedClaGroupOrderMsg(got []*models.ClaGroup, want []*models.ClaGroup) string {
	var gotOrder, wantOrder []string
	for _, v := range got {
		gotOrder = append(gotOrder, v.ClaGroupName)
	}
	for _, v := range want {
		wantOrder = append(wantOrder, v.ClaGroupName)
	}
	return fmt.Sprintf("got [%#v], expected [%#v]", gotOrder, wantOrder)
}
