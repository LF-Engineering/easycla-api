package events

import (
	"context"
	"errors"

	"github.com/communitybridge/easycla-api/gen/models"
	"github.com/communitybridge/easycla-api/gen/restapi/operations/events"
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
	SearchEvents(ctx context.Context, params *events.SearchEventsParams) (*models.EventList, error)
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
	return eventType != ""
}

// Create event will create event in database.
// ID passed in input will get ignore.
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

func searchEventsSQLStatement(db *sqlx.DB, params *events.SearchEventsParams) *sqlz.SelectStmt {
	stmt := sqlz.Newx(db).
		Select("*").
		From(EventsTable)

	var conditions []sqlz.WhereCondition
	if params.EventType != nil {
		conditions = append(conditions, sqlz.Eq("event_type", *params.EventType))
	}
	if params.UserID != nil {
		conditions = append(conditions, sqlz.Eq("user_id", *params.UserID))
	}
	if params.ProjectID != nil {
		conditions = append(conditions, sqlz.Eq("project_id", *params.ProjectID))
	}
	if params.CompanyID != nil {
		conditions = append(conditions, sqlz.Eq("company_id", *params.CompanyID))
	}
	if params.Before != nil {
		conditions = append(conditions, sqlz.Lte("event_time", *params.Before))
	}
	if params.After != nil {
		conditions = append(conditions, sqlz.Gte("event_time", *params.After))
	}
	if len(conditions) != 0 {
		stmt = stmt.Where(conditions...)
	}
	if params.Offset != nil {
		stmt = stmt.Offset(*params.Offset)
	}
	if params.PageSize != nil {
		stmt = stmt.Limit(*params.PageSize)
	}
	orderBy := "event_time"
	if params.OrderBy != nil {
		orderBy = *params.OrderBy
	}
	if params.SortOrder != nil && *params.SortOrder == "desc" {
		stmt = stmt.OrderBy(sqlz.Desc(orderBy))
	} else {
		stmt = stmt.OrderBy(sqlz.Asc(orderBy))
	}
	return stmt
}

// SearchEvents returns list of events matching with filter criteria.
func (r *repository) SearchEvents(ctx context.Context, params *events.SearchEventsParams) (*models.EventList, error) {
	var events []SQLEvent
	stmt := searchEventsSQLStatement(r.GetDB(), params)
	err := stmt.GetAll(&events)
	if err != nil {
		return nil, err
	}
	var result models.EventList
	for _, e := range events {
		result.Events = append(result.Events, e.toEvent())
	}
	return &result, nil
}
