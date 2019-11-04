package orgs

import (
	"context"

	"github.com/communitybridge/easycla-api/gen/models"
	"github.com/communitybridge/easycla-api/gen/restapi/operations/organization"
)

// Service interface defines methods of organization service
type Service interface {
	GetOrgFoundations(ctx context.Context, in *organization.GetOrgFoundationsParams) (models.Organization, error)
}

type service struct {
	repo Repository
}

// NewService creates new instance of organization service
func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) GetOrgFoundations(ctx context.Context, in *organization.GetOrgFoundationsParams) (models.Organization, error) {
	return s.repo.GetOrgFoundations(ctx, in.SalesForceOrganizationID)
}
