package projects

import (
	"context"

	"github.com/communitybridge/easycla-api/gen/models"
	"github.com/communitybridge/easycla-api/gen/restapi/operations/project"
)

// Service interface defines methods of project service
type Service interface {
	GetProject(ctx context.Context, in *project.GetProjectParams) (*models.Project, error)
}

type service struct {
	repo Repository
}

// NewService creates new instance of project service
func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) GetProject(ctx context.Context, in *project.GetProjectParams) (*models.Project, error) {
	return s.repo.GetProject(ctx, in.ProjectID)
}
