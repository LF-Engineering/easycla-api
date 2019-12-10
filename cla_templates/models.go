package cla_templates

import (
	"database/sql"
	"encoding/json"

	"github.com/communitybridge/easycla-api/gen/models"
)

// SQLCLATemplate struct represent row of sql.events table
type SQLCLATemplate struct {
	ID           sql.NullString `db:"id"`
	CreatedAt    sql.NullInt64  `db:"created_at"`
	UpdatedAt    sql.NullInt64  `db:"updated_at"`
	Name         sql.NullString `db:"name"`
	Description  sql.NullString `db:"description"`
	Version      sql.NullInt64  `db:"version"`
	IclaHTMLBody sql.NullString `db:"icla_html_body"`
	CclaHTMLBody sql.NullString `db:"ccla_html_body"`
	MetaFields   sql.NullString `db:"meta_fields"`
	IclaFields   sql.NullString `db:"icla_fields"`
	CclaFields   sql.NullString `db:"ccla_fields"`
}

func (t *SQLCLATemplate) toClaTemplate() (*models.ClaTemplate, error) {
	template := &models.ClaTemplate{
		CclaFields:   nil,
		IclaFields:   nil,
		MetaFields:   nil,
		CclaHTMLBody: t.CclaHTMLBody.String,
		CreatedAt:    t.CreatedAt.Int64,
		Description:  t.Description.String,
		IclaHTMLBody: t.IclaHTMLBody.String,
		ID:           t.ID.String,
		Name:         t.Name.String,
		UpdatedAt:    t.UpdatedAt.Int64,
		Version:      t.Version.Int64,
	}
	if t.MetaFields.Valid {
		err := json.Unmarshal([]byte(t.MetaFields.String), &template.MetaFields)
		if err != nil {
			return nil, err
		}
	}
	if t.IclaFields.Valid {
		err := json.Unmarshal([]byte(t.IclaFields.String), &template.IclaFields)
		if err != nil {
			return nil, err
		}
	}
	if t.CclaFields.Valid {
		err := json.Unmarshal([]byte(t.CclaFields.String), &template.CclaFields)
		if err != nil {
			return nil, err
		}
	}
	return template, nil
}
