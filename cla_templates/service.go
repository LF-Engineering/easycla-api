package cla_templates

import (
	"github.com/communitybridge/easycla-api/gen/models"
	params "github.com/communitybridge/easycla-api/gen/restapi/operations/cla_templates"
)

// Service interface defines methods of cla_template service
type Service interface {
	CreateCLATemplate(in *params.CreateCLATemplateParams) (*models.ClaTemplate, error)
	UpdateCLATemplate(in *params.UpdateCLATemplateParams) (*models.ClaTemplate, error)
	DeleteCLATemplate(in *params.DeleteCLATemplateParams) error
	GetCLATemplate(in *params.GetCLATemplateParams) (*models.ClaTemplate, error)
	ListCLATemplate(in *params.ListCLATemplatesParams) (*models.ClaTemplateList, error)
}

type service struct {
	repo Repository
}

// NewService creates new instance of cla_template service
func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) CreateCLATemplate(in *params.CreateCLATemplateParams) (*models.ClaTemplate, error) {
	return s.repo.CreateCLATemplate(in.ClaTemplate)
}

func (s *service) GetCLATemplate(in *params.GetCLATemplateParams) (*models.ClaTemplate, error) {
	return s.repo.GetCLATemplate(in.ClaTemplateID)
}

func (s *service) UpdateCLATemplate(in *params.UpdateCLATemplateParams) (*models.ClaTemplate, error) {
	return s.repo.UpdateCLATemplate(in.ClaTemplateID, in.ClaTemplate)
}

func (s *service) DeleteCLATemplate(in *params.DeleteCLATemplateParams) error {
	return s.repo.DeleteCLATemplate(in.ClaTemplateID)
}

func (s *service) ListCLATemplate(in *params.ListCLATemplatesParams) (*models.ClaTemplateList, error) {
	return s.repo.ListCLATemplates()
}
