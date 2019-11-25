package cla_groups_test

import (
	"context"
	"github.com/communitybridge/easycla-api/cla_groups"
	"github.com/communitybridge/easycla-api/gen/models"
	params "github.com/communitybridge/easycla-api/gen/restapi/operations/cla_groups"
	"github.com/communitybridge/easycla-api/gen/restapi/operations/events"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_UpdateCLAGroup(t *testing.T) {
	prepareTestDatabase()
	totalClaGroups := numberOfCLAGroups()
	claGroupID := "d9dc5834-3d9a-4d04-abb6-4a36ed378304" // cla_group id (of Dataloader)in fixtures
	assert.Equal(t, 3, int(totalClaGroups))
	claGroupName := "cncf cla"
	tests := []struct {
		name    string
		args    *params.UpdateCLAGroupParams
		want    error
		wantErr bool
	}{
		{
			name: "success",
			args: &params.UpdateCLAGroupParams{
				ClaGroup : &models.UpdateClaGroup{
					CclaEnabled:  false,
					IclaEnabled:  true,
					ClaGroupName: claGroupName,
				},
				ClaGroupID:claGroupID,

			},
			want: nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := claGroupsService.UpdateCLAGroup(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateCLAGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, 3, int(totalClaGroups),"number of cla_groups should not change") {
				t.Fail()
			}
			res,err := getCLAGroup(tt.args.ClaGroupID)
			if !assert.Nil(t,err) {
				t.Fail()
			}
			if !assert.Equal(t, res.CCLAEnabled.Bool, tt.args.ClaGroup.CclaEnabled) {
				t.Fail()
			}
			if !assert.Equal(t, res.ICLAEnabled.Bool, tt.args.ClaGroup.IclaEnabled) {
				t.Fail()
			}
			if !assert.Equal(t, res.CLAGroupName.String, tt.args.ClaGroup.ClaGroupName) {
				t.Fail()
			}
			list, err := eventsService.SearchEvents(context.TODO(), &events.SearchEventsParams{})
			if !assert.Nil(t, err) {
				t.Fail()
			}
			if !assert.Equal(t, 1, len(list.Events)) {
				t.Fail()
			}
			if !assert.Equal(t, cla_groups.CLAGroupDeleted, list.Events[0].EventType) {
				t.Fail()
			}
		})
	}
}
