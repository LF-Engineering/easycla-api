package cla_groups_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/communitybridge/easycla-api/cla_groups"
	params "github.com/communitybridge/easycla-api/gen/restapi/operations/cla_groups"
	"github.com/communitybridge/easycla-api/gen/restapi/operations/events"
	"github.com/stretchr/testify/assert"
)

const (
	DataloaderClaGroupID   = "d9dc5834-3d9a-4d04-abb6-4a36ed378304"
	DataloaderRepositoryID = "ae85c1d2-1e93-4e4d-b6d5-216b1a41ff17"
)

func Test_DeleteCLAGroup(t *testing.T) {
	prepareTestDatabase()
	claGroupID := DataloaderClaGroupID
	nonExistentClaGroupID := "e9dc5834-3d9a-4d04-abb6-4a36ed378300"
	/* Check if db is in required state */
	if !assert.Equal(t, true, isCLAGroupPresent(claGroupID),
		fmt.Sprintf("cla_group with id : %s must be present", claGroupID)) {
		t.Fail()
	}
	if !assert.Equal(t, 2, int(numberOfProjectManagers(claGroupID)),
		fmt.Sprintf("cla_group %s must have 2 project managers", claGroupID)) {
		t.Fail()
	}
	if !assert.Equal(t, int64(2), numberOfRepositories(),
		"total number of repositories must be 2") {
		t.Fail()
	}
	if !assert.Equal(t, true, isRepositoryPresent(DataloaderRepositoryID),
		fmt.Sprintf("repository with ID %s must be present", DataloaderRepositoryID)) {
		t.Fail()
	}
	list, err := eventsService.SearchEvents(context.TODO(), &events.SearchEventsParams{})
	if !assert.Nil(t, err) {
		t.Fail()
	}
	if !assert.Equal(t, 0, len(list.Events)) {
		t.Fail()
	}

	tests := []struct {
		name    string
		args    *params.DeleteCLAGroupParams
		want    error
		wantErr bool
	}{
		{
			name: "cla_group exist",
			args: &params.DeleteCLAGroupParams{
				ClaGroupID: claGroupID,
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "cla_group not exist",
			args: &params.DeleteCLAGroupParams{
				ClaGroupID: nonExistentClaGroupID,
			},
			want:    cla_groups.ErrClaGroupNotFound,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := claGroupsService.DeleteCLAGroup(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteCLAGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.want.Error() != err.Error() {
				t.Errorf("DeleteCLAGroup() error = %v, expected %v", err, tt.want)
				return
			}
			if !tt.wantErr {
				if !assert.Equal(t, false, isCLAGroupPresent(tt.args.ClaGroupID)) {
					t.Fail()
				}
				if !assert.Equal(t, 0, int(numberOfProjectManagers(tt.args.ClaGroupID))) {
					t.Fail()
				}
				if !assert.Equal(t, false, isRepositoryPresent(DataloaderRepositoryID),
					"repository associated with cla_group must be deleted") {
					t.Fail()
				}
				if !assert.Equal(t, int64(1), numberOfRepositories(),
					"other repositories must remain unaffected") {
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
			}
		})
	}
}
