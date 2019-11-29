package events_test

import (
	"context"
	"testing"

	"github.com/communitybridge/easycla-api/gen/models"

	params "github.com/communitybridge/easycla-api/gen/restapi/operations/events"
	"github.com/stretchr/testify/assert"
)

func Test_SearchEvents(t *testing.T) {
	prepareTestDatabase()
	tests := []struct {
		name    string
		args    *params.SearchEventsParams
		want    *models.EventList
		wantErr bool
	}{
		{
			name:    "default", // page_size 100, sortby event_time asc
			args:    &params.SearchEventsParams{PageSize: newInt64(100)},
			want:    &models.EventList{Events: []*models.Event{Event1, Event2, Event3, Event4, Event5, Event6}},
			wantErr: false,
		},
		{
			name:    "filter by event_type",
			args:    &params.SearchEventsParams{EventType: newString(EventUserUpdated)},
			want:    &models.EventList{Events: []*models.Event{Event2, Event6}},
			wantErr: false,
		},
		{
			name:    "filter by user_id",
			args:    &params.SearchEventsParams{UserID: newString(UserMani)},
			want:    &models.EventList{Events: []*models.Event{Event1, Event3, Event5}},
			wantErr: false,
		},
		{
			name:    "filter by project_id",
			args:    &params.SearchEventsParams{ProjectID: newString(ProjectPrometheus)},
			want:    &models.EventList{Events: []*models.Event{Event2, Event5}},
			wantErr: false,
		},
		{
			name:    "filter by before",
			args:    &params.SearchEventsParams{Before: newInt64(3)},
			want:    &models.EventList{Events: []*models.Event{Event1, Event2, Event3}},
			wantErr: false,
		},
		{
			name:    "filter by after",
			args:    &params.SearchEventsParams{After: newInt64(4)},
			want:    &models.EventList{Events: []*models.Event{Event4, Event5, Event6}},
			wantErr: false,
		},
		{
			name: "page_size and offset",
			args: &params.SearchEventsParams{
				PageSize: newInt64(2),
				Offset:   newInt64(3),
			},
			want:    &models.EventList{Events: []*models.Event{Event4, Event5}},
			wantErr: false,
		},
		{
			name: "order by event_type asc", // default sort order is asc
			args: &params.SearchEventsParams{
				OrderBy: newString("event_type"),
			},
			want:    &models.EventList{Events: []*models.Event{Event3, Event4, Event1, Event5, Event2, Event6}},
			wantErr: false,
		},
		{
			name: "order by event_type desc",
			args: &params.SearchEventsParams{
				OrderBy:   newString("event_type"),
				SortOrder: newString("desc"),
			},
			want:    &models.EventList{Events: []*models.Event{Event2, Event6, Event1, Event5, Event4, Event3}},
			wantErr: false,
		},
		{
			name: "filter by user_id and project_id",
			args: &params.SearchEventsParams{
				UserID:    newString(UserMani),
				ProjectID: newString(ProjectKubernetes),
			},
			want:    &models.EventList{Events: []*models.Event{Event3}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := eventsService.SearchEvents(context.TODO(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchEvents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, result) {
				t.Fail()
			}
		})
	}
}
