package events_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"

	"github.com/stretchr/testify/assert"

	"github.com/communitybridge/easycla-api/events"

	"github.com/communitybridge/easycla-api/gen/models"
	params "github.com/communitybridge/easycla-api/gen/restapi/operations/events"
)

func Test_CreateEvent(t *testing.T) {
	testEvent := models.Event{
		CompanyID: "company",
		EventData: "{\"users\": 1}",
		EventType: "event",
		ProjectID: "project",
		UserID:    "user",
	}
	minimalEvent := models.Event{
		EventType: "event",
		UserID:    "user",
	}
	tests := []struct {
		name     string
		args     models.Event
		want     error
		expected models.Event
	}{
		{
			name: "empty event",
			args: models.Event{
				EventType: "",
			},
			want: events.ErrInvalidEventType,
		},
		{
			name: "without user_id",
			args: models.Event{
				EventType: "UserAdded",
				UserID:    "",
			},
			want: events.ErrUserIDRequired,
		},
		{
			name:     "valid event",
			args:     testEvent,
			want:     nil,
			expected: testEvent,
		},
		{
			name:     "minimal event",
			args:     minimalEvent,
			want:     nil,
			expected: minimalEvent,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prepareTestDatabase()
			currentTime := time.Now().Unix()
			err := eventsService.CreateEvent(tt.args)
			if err != tt.want {
				t.Errorf("CreateEvent() error = %v, want %v", err, tt.want)
				return
			}
			if err != nil {
				return
			}
			list, err := eventsService.SearchEvents(context.TODO(), &params.SearchEventsParams{After: newInt64(currentTime)})
			if err != nil {
				t.Fatal(err)
			}
			if !assert.Equal(t, 1, len(list.Events)) {
				t.Fail()
			}
			got := *list.Events[0]
			if !assert.Equal(t, true, strfmt.IsUUID4(got.ID)) {
				t.Fail()
			}
			if !assert.Equal(t, true, got.EventTime >= currentTime) {
				t.Fail()
			}
			expected := tt.expected
			expected.ID = got.ID
			expected.EventTime = got.EventTime
			if !assert.Equal(t, expected, got) {
				t.Fail()
			}
		})
	}
}
