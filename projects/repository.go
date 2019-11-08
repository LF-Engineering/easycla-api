package projects

import (
	"context"
	"encoding/json"

	"github.com/communitybridge/easycla-api/gen/models"
)

// Repository interface defines methods of project repository service
type Repository interface {
	GetProject(ctx context.Context, projectID string) (*models.Project, error)
}

type repository struct {
}

// NewRepository creates new instance of project repository
func NewRepository() Repository {
	return &repository{}
}

func (r *repository) GetProject(ctx context.Context, projectID string) (*models.Project, error) {
	var p models.Project
	compactProject := `{"project_id":"1234566....","project_name":"CFF Migration","templates":[{"type":"icla","legal_entity_name":"","author":"some user lfid","url":"https://....","version":"v1","created":"iso datetime","updated":"iso datetime"},{"type":"ccla","legal_entity_name":"","author":"some user lfid","url":"https://....","version":"v1","created":"iso datetime","updated":"iso datetime"}],"ccla_enabled":true,"icla_enabled":true,"icla_signature_required":true,"ccla_companies_signed":["companyID1","companyID2","companyID3","companyID4","companyID5","companyID6","companyID7","companyID8","companyID9","companyID10","companyID11","companyID12"],"whitelist":{"users":[{"user_id":"<user id 1>","signed":true,"signed_on":"iso datetime","last_activity":"iso datetime"},{"user_id":"<user id 2>","signed":true,"signed_on":"iso datetime","last_activity":"iso datetime"},{"user_id":"<user id 3>","signed":false,"signed_on":null,"last_activity":null}],"domains":[{"domain":"<domain string>"}],"emails":[{"email":"<email string/regex>"}],"github_users":[{"github_user":"<github user id"}],"github_orgs":[{"github_org":"<github org name/id"}]}}`
	err := json.Unmarshal([]byte(compactProject), &p)
	return &p, err
}
