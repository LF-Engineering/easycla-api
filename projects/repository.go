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
	compactProject := `{"projectID":"1234566....","projectName":"CFF Migration","templates":[{"type":"icla","legalEntityName":"","author":"some user lfid","url":"https://....","version":"v1","created":"iso datetime","updated":"iso datetime"},{"type":"ccla","legalEntityName":"","author":"some user lfid","url":"https://....","version":"v1","created":"iso datetime","updated":"iso datetime"}],"cclaEnabled":true,"iclaEnabled":true,"iclaSignatureRequired":true,"cclaCompaniesSigned":["companyID1","companyID2","companyID3","companyID4","companyID5","companyID6","companyID7","companyID8","companyID9","companyID10","companyID11","companyID12"],"whitelist":{"users":[{"userID":"<user id 1>","signed":true,"signed_on":"iso datetime","last_activity":"iso datetime"},{"userID":"<user id 2>","signed":true,"signed_on":"iso datetime","last_activity":"iso datetime"},{"userID":"<user id 3>","signed":false,"signed_on":null,"last_activity":null}],"domains":[{"domain":"<domain string>"}],"emails":[{"email":"<email string/regex>"}],"githubUsers":[{"githubUser":"<github user id"}],"githubOrgs":[{"githubOrg":"<github org name/id"}]}}`
	err := json.Unmarshal([]byte(compactProject), &p)
	return &p, err
}
