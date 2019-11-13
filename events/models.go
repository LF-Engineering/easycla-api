package events

import (
	"database/sql"

	"github.com/communitybridge/easycla-api/gen/models"
)

// SQLEvent struct represent row of sql.events table
type SQLEvent struct {
	ID        sql.NullString `json:"id,omitempty"`
	EventType sql.NullString `json:"event_type,omitempty"`
	UserID    sql.NullString `json:"user_id,omitempty"`
	ProjectID sql.NullString `json:"project_id,omitempty"`
	CompanyID sql.NullString `json:"company_id,omitempty"`
	EventTime sql.NullInt64  `json:"event_time,omitempty"`
	EventData sql.NullString `json:"event_data,omitempty"`
}

func (e *SQLEvent) toEvent() *models.Event {
	return &models.Event{
		ID:        e.ID.String,
		EventType: e.EventType.String,
		UserID:    e.UserID.String,
		ProjectID: e.ProjectID.String,
		CompanyID: e.CompanyID.String,
		EventTime: e.EventTime.Int64,
		EventData: e.EventData.String,
	}
}
