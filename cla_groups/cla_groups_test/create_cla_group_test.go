package cla_groups_test

import (
	"context"
	"github.com/communitybridge/easycla-api/cla_groups"
	"github.com/communitybridge/easycla-api/gen/restapi/operations/events"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/communitybridge/easycla-api/gen/models"
	params "github.com/communitybridge/easycla-api/gen/restapi/operations/cla_groups"
	"github.com/go-openapi/strfmt"
)

func Test_CreateCLAGroup(t *testing.T) {
	prepareTestDatabase()
	totalClaGroups := numberOfCLAGroups()
	assert.Equal(t, 3, int(totalClaGroups))
	foundationId := "CNCFCreateTest"
	claGroupName := "cncf cla"
	projectManager := "413f4711-a3c3-4635-9dad-a0ba58694205"
	tests := []struct {
		name    string
		args    *models.CreateClaGroup
		want    *models.ClaGroup
		wantErr bool
	}{
		{
			name: "success",
			args: &models.CreateClaGroup{
				CclaEnabled:     newBool(true),
				IclaEnabled:     newBool(false),
				FoundationID:    newString(foundationId),
				ClaGroupName:    newString(claGroupName),
				ProjectManagers: []strfmt.UUID{strfmt.UUID(projectManager)},
			},
			want: &models.ClaGroup{
				CclaEnabled:     true,
				IclaEnabled:     false,
				FoundationID:    foundationId,
				ClaGroupName:    claGroupName,
				ProjectManagers: []strfmt.UUID{strfmt.UUID(projectManager)},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			currentTime := time.Now().Unix()
			res, err := claGroupsService.CreateCLAGroup(&params.CreateCLAGroupParams{
				ClaGroup: tt.args,
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateCLAGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, totalClaGroups+1, numberOfCLAGroups(), "total number of cla_groups do not match") {
				t.Fail()
			}
			if !assert.Equal(t, tt.want.CclaEnabled, res.CclaEnabled, "ccla_enabled do not match") {
				t.Fail()
			}
			if !assert.Equal(t, tt.want.IclaEnabled, res.IclaEnabled, "icla_enabled do not match") {
				t.Fail()
			}
			if !assert.Equal(t, tt.want.ClaGroupName, res.ClaGroupName, "cla_group_name do not match") {
				t.Fail()
			}
			if !assert.Equal(t, tt.want.FoundationID, res.FoundationID, "foundation_id do not match") {
				t.Fail()
			}
			if !assert.Equal(t, tt.want.ProjectManagers, res.ProjectManagers, "project managers do not match") {
				t.Fail()
			}
			if !assert.Equal(t, len(tt.want.ProjectManagers), int(numberOfProjectManagers(res.ID)), "project managers do not match") {
				t.Fail()
			}
			if !assert.Equal(t, true, strfmt.IsUUID4(res.ID), "id is not valid") {
				t.Fail()
			}
			if !assert.Equal(t, true, res.CreatedAt >= currentTime, "created_at is not valid") {
				t.Fail()
			}
			if !assert.Equal(t, true, res.UpdatedAt >= currentTime, "updated_at is not valid") {
				t.Fail()
			}

			list, err := claGroupsService.ListCLAGroups(&params.ListCLAGroupsParams{FoundationID:&foundationId})
			if !assert.Nil(t,err,"get cla group list failed") {
				t.Fail()
			}
			if !assert.Equal(t, len(list.ClaGroups),int(1)) {
				t.Fail()
			}
			if !assert.Equal(t, list.ClaGroups[0],res){
				t.Fail()
			}
			elist, err := eventsService.SearchEvents(context.TODO(), &events.SearchEventsParams{})
			if !assert.Nil(t, err) {
				t.Fail()
			}
			if !assert.Equal(t, 1, len(elist.Events)) {
				t.Fail()
			}
			if !assert.Equal(t, cla_groups.CLAGroupUpdated, elist.Events[0].EventType) {
				t.Fail()
			}
		})
	}
}

