package repositories

import (
	"github.com/LF-Engineering/lfx-kit/auth"
	"github.com/communitybridge/easycla-api/events"
	"github.com/communitybridge/easycla-api/gen/models"
	params "github.com/communitybridge/easycla-api/gen/restapi/operations/repositories"
)

// Service interface defines methods of repository service
type Service interface {
	CreateRepositories(user *auth.User, in *params.CreateRepositoriesParams) (*models.RepositoryList, error)
	DeleteRepositories(user *auth.User, in *params.DeleteRepositoriesParams) error
	ListRepositories(in *params.ListRepositoriesParams) (*models.RepositoryList, error)
}

type service struct {
	repo   Repository
	events events.Service
}

// NewService creates new instance of repository service
func NewService(repo Repository, eventService events.Service) Service {
	return &service{
		repo:   repo,
		events: eventService,
	}
}
func (s *service) CreateRepositories(user *auth.User, in *params.CreateRepositoriesParams) (*models.RepositoryList, error) {
	result, err := s.repo.CreateRepositories(in.RepositoriesInput)
	if err != nil {
		return nil, err
	}
	s.createEvent("user", *in.RepositoriesInput.ProjectID, EventTypeRepositoriesCreated, &CreatedEventData{
		Input:               in.RepositoriesInput,
		CreatedRepositories: result,
	})
	return result, nil
}

func (s *service) DeleteRepositories(user *auth.User, in *params.DeleteRepositoriesParams) error {
	err := s.repo.DeleteRepositories(in.RepositoriesInput)
	if err != nil {
		return err
	}
	s.createEvent("user", *in.RepositoriesInput.ProjectID, EventTypeRepositoriesDeleted, &DeletedEventData{
		Input: in.RepositoriesInput,
	})
	return nil
}

func (s *service) ListRepositories(in *params.ListRepositoriesParams) (*models.RepositoryList, error) {
	return s.repo.ListRepositories(in)
}
