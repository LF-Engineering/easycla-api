package cla_groups

import (
	"encoding/json"

	"github.com/communitybridge/easycla-api/gen/models"
	log "github.com/communitybridge/easycla-api/logging"
)

// Events
const (
	CLAGroupCreated = "cla_group_created"
	CLAGroupUpdated = "cla_group_updated"
	CLAGroupDeleted = "cla_group_deleted"
)

// CLAGroupCreatedEvent is data for cla_group_created event
type CLAGroupCreatedEvent struct {
	ClaGroupID   string `json:"cla_group_id,omitempty"`
	ClaGroupName string `json:"cla_group_name,omitempty"`
	CclaEnabled  bool   `json:"ccla_enabled,omitempty"`
	IclaEnabled  bool   `json:"icla_enabled,omitempty"`
	FoundationID string `json:"foundation_id,omitempty"`
}

// CLAGroupUpdatedEvent is data for cla_group_updated event
type CLAGroupUpdatedEvent struct {
	ClaGroupID   string `json:"cla_group_id,omitempty"`
	ClaGroupName string `json:"cla_group_name,omitempty"`
	CclaEnabled  bool   `json:"ccla_enabled,omitempty"`
	IclaEnabled  bool   `json:"icla_enabled,omitempty"`
}

// CLAGroupDeletedEvent is data for cla_group_deleted event
type CLAGroupDeletedEvent struct {
	ClaGroupID string `json:"cla_group_id,omitempty"`
}

func (s *service) createEvent(userID, companyID, projectID, eventType string, data interface{}) {
	event := models.Event{
		UserID:    userID,
		CompanyID: companyID,
		ProjectID: projectID,
		EventType: eventType,
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
