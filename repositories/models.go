package repositories

import (
	"database/sql"

	"github.com/communitybridge/easycla-api/gen/models"
)

// SQLRepository represents row of sql.repositories table
type SQLRepository struct {
	ID               sql.NullString `db:"id"`
	RepositoryType   sql.NullString `db:"repository_type"`
	Name             sql.NullString `db:"name"`
	ExternalID       sql.NullString `db:"external_id"`
	OrganizationName sql.NullString `db:"organization_name"`
	URL              sql.NullString `db:"url"`
	Enabled          sql.NullBool   `db:"enabled"`
	ProjectID        sql.NullString `db:"project_id"`
	ClaGroupID       sql.NullString `db:"cla_group_id"`
	CreatedAt        sql.NullInt64  `db:"created_at"`
}

func (t *SQLRepository) toRepository() *models.Repository {
	return &models.Repository{
		ClaGroupID:       t.ClaGroupID.String,
		CreatedAt:        t.CreatedAt.Int64,
		ExternalID:       t.ExternalID.String,
		Enabled:          t.Enabled.Bool,
		ID:               t.ID.String,
		Name:             t.Name.String,
		OrganizationName: t.OrganizationName.String,
		ProjectID:        t.ProjectID.String,
		RepositoryType:   t.RepositoryType.String,
		URL:              t.URL.String,
	}
}
