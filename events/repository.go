package events

import (
	"errors"

	"github.com/communitybridge/easycla-api/gen/models"
	"github.com/ido50/sqlz"
	"github.com/jmoiron/sqlx"
)

// EventsTable is the name of events table in database
const (
	EventsTable = "cla.events"
)

// ErrInvalidEventType,ErrUserIDRequired are errors
var (
	ErrInvalidEventType = errors.New("invalid event type")
	ErrUserIDRequired   = errors.New("user_id cannot be empty")
)

// Repository interface defines methods of event repository service
type Repository interface {
	CreateEvent(event *models.Event) error
	//	ListEvents(ctx context.Context, params *events.ListEventsParams) (*models.EventList, error)
}

type repository struct {
	db *sqlx.DB
}

// NewRepository creates new instance of audit event repository
func NewRepository(dbConn *sqlx.DB) Repository {
	return &repository{
		db: dbConn,
	}
}

func (r *repository) GetDB() *sqlx.DB {
	return r.db
}

func validEventType(eventType string) bool {
	if eventType == "" {
		return false
	}
	return true
}

func (r *repository) CreateEvent(event *models.Event) error {
	values := make(map[string]interface{})
	if !validEventType(event.EventType) {
		return ErrInvalidEventType
	}
	values["event_type"] = event.EventType
	if event.UserID == "" {
		return ErrUserIDRequired
	}
	values["user_id"] = event.UserID
	if event.ProjectID != "" {
		values["project_id"] = event.ProjectID
	}
	if event.CompanyID != "" {
		values["company_id"] = event.CompanyID
	}
	if event.EventData != "" {
		values["event_data"] = event.EventData
	}
	_, err := sqlz.Newx(r.GetDB()).
		InsertInto(EventsTable).
		ValueMap(values).Exec()
	return err
}
