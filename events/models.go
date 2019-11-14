package events

import (
	"database/sql"

	"github.com/communitybridge/easycla-api/gen/models"
)

// SQLEvent struct represent row of sql.events table
type SQLEvent struct {
	ID        sql.NullString `db:"id"`
	EventType sql.NullString `db:"event_type"`
	UserID    sql.NullString `db:"user_id"`
	ProjectID sql.NullString `db:"project_id"`
	CompanyID sql.NullString `db:"company_id"`
	EventTime sql.NullInt64  `db:"event_time"`
	EventData sql.NullString `db:"event_data"`
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
