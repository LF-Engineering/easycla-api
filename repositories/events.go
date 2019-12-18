package repositories

import (
	"encoding/json"

	"github.com/communitybridge/easycla-api/gen/models"
	log "github.com/communitybridge/easycla-api/logging"
)

// EventType is the type for holding type of the event
type EventType string

// Various Events
const (
	EventTypeRepositoriesCreated EventType = "repositories_created"
	EventTypeRepositoriesDeleted EventType = "repositories_deleted"
)

// CreatedEventData is data for repositories_created event
type CreatedEventData struct {
	Input               *models.CreateRepositoriesInput `json:"input"`
	CreatedRepositories *models.RepositoryList          `json:"created_repositories"`
}

// DeletedEventData is data for repositories_deleted event
type DeletedEventData struct {
	Input *models.DeleteRepositoriesInput `json:"input"`
}

func (s *service) createEvent(userID, projectID string, eventType EventType, data interface{}) {
	event := models.Event{
		UserID:    userID,
		ProjectID: projectID,
		EventType: string(eventType),
	}
	if data != nil {
		eventData, err := json.Marshal(data)
		if err != nil {
			log.Error(log.Trace(), err)
		}
		event.EventData = string(eventData)
	}
	err := s.events.CreateEvent(event)
	if err != nil {
		log.Error(log.Trace(), err)
	}
}
