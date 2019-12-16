package repositories

import (
	"github.com/communitybridge/easycla-api/gen/models"
	params "github.com/communitybridge/easycla-api/gen/restapi/operations/repositories"
)

// Service interface defines methods of repository service
type Service interface {
	CreateRepositories(in *params.CreateRepositoriesParams) (*models.RepositoryList, error)
	DeleteRepositories(in *params.DeleteRepositoriesParams) error
	ListRepositories(in *params.ListRepositoriesParams) (*models.RepositoryList, error)
}

type service struct {
	repo Repository
}

// NewService creates new instance of repository service
func NewService(repo Repository) Service {
	return &service{repo}
}
func (s *service) CreateRepositories(in *params.CreateRepositoriesParams) (*models.RepositoryList, error) {
	return s.repo.CreateRepositories(in.RepositoriesInput)
}

func (s *service) DeleteRepositories(in *params.DeleteRepositoriesParams) error {
	return s.repo.DeleteRepositories(in.RepositoriesInput)
}

func (s *service) ListRepositories(in *params.ListRepositoriesParams) (*models.RepositoryList, error) {
	return s.repo.ListRepositories(in)
}
