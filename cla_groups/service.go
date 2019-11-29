package cla_groups

import (
	"github.com/communitybridge/easycla-api/events"
	"github.com/communitybridge/easycla-api/gen/models"
	"github.com/communitybridge/easycla-api/gen/restapi/operations/cla_groups"
)

// Service interface defines methods of cla_groups service
type Service interface {
	CreateCLAGroup(in *cla_groups.CreateCLAGroupParams) (*models.ClaGroup, error)
	DeleteCLAGroup(in *cla_groups.DeleteCLAGroupParams) error
	UpdateCLAGroup(in *cla_groups.UpdateCLAGroupParams) error
	ListCLAGroups(in *cla_groups.ListCLAGroupsParams) (*models.ClaGroupList, error)
}

type service struct {
	repo   Repository
	events events.Service
}

// NewService creates new instance of event service
func NewService(repo Repository, eventService events.Service) Service {
	return &service{
		repo:   repo,
		events: eventService,
	}
}

func (s *service) CreateCLAGroup(in *cla_groups.CreateCLAGroupParams) (*models.ClaGroup, error) {
	result, err := s.repo.CreateCLAGroup(in.ClaGroup)
	if err != nil {
		return nil, err
	}
	s.createEvent("user", "", "", CLAGroupCreated, CLAGroupCreatedEvent{
		ClaGroupID:   result.ID,
		ClaGroupName: result.ClaGroupName,
		ProjectID:    result.ProjectID,
		CclaEnabled:  result.CclaEnabled,
		IclaEnabled:  result.IclaEnabled,
	})
	return result, err
}

func (s *service) DeleteCLAGroup(in *cla_groups.DeleteCLAGroupParams) error {
	err := s.repo.DeleteCLAGroup(in.ClaGroupID)
	if err != nil {
		return err
	}
	s.createEvent("user", "", "", CLAGroupDeleted, CLAGroupDeletedEvent{
		ClaGroupID: in.ClaGroupID,
	})
	return nil
}

func (s *service) UpdateCLAGroup(in *cla_groups.UpdateCLAGroupParams) error {
	err := s.repo.UpdateCLAGroup(in.ClaGroupID, in.ClaGroup)
	if err != nil {
		return err
	}
	s.createEvent("user", "", "", CLAGroupUpdated, CLAGroupUpdatedEvent{
		ClaGroupID:   in.ClaGroupID,
		ClaGroupName: in.ClaGroup.ClaGroupName,
		CclaEnabled:  in.ClaGroup.CclaEnabled,
		IclaEnabled:  in.ClaGroup.IclaEnabled,
	})
	return nil
}

func (s *service) ListCLAGroups(in *cla_groups.ListCLAGroupsParams) (*models.ClaGroupList, error) {
	return s.repo.ListCLAGroups(in)
}
