package cla_templates

import (
	"encoding/json"

	"github.com/ido50/sqlz"

	"github.com/communitybridge/easycla-api/gen/models"
	"github.com/jmoiron/sqlx"
)

// CLATemplatesTable is the name of events table in database
const (
	CLATemplatesTable = "cla.cla_templates"
)

// Repository interface defines methods of event repository service
type Repository interface {
	CreateCLATemplate(template *models.ClaTemplateInput) (*models.ClaTemplate, error)
}

type repository struct {
	db *sqlx.DB
}

// NewRepository creates new instance of audit event repository
func NewRepository(dbConn *sqlx.DB) Repository {
	return &repository{
		db: dbConn,
	}
}

func (r *repository) GetDB() *sqlx.DB {
	return r.db
}

func (r *repository) CreateCLATemplate(in *models.ClaTemplateInput) (*models.ClaTemplate, error) {
	var result SQLCLATemplate
	values := map[string]interface{}{
		"name":        in.Name,
		"description": in.Description,
		"version":     1,
	}
	if in.IclaHTMLBody != "" {
		values["icla_html_body"] = in.IclaHTMLBody
	}
	if in.CclaHTMLBody != "" {
		values["ccla_html_body"] = in.CclaHTMLBody
	}
	if len(in.MetaFields) != 0 {
		metaFieldJson, err := json.Marshal(in.MetaFields)
		if err != nil {
			return nil, err
		}
		values["meta_fields"] = metaFieldJson
	}
	if len(in.IclaFields) != 0 {
		iclaFieldJson, err := json.Marshal(in.IclaFields)
		if err != nil {
			return nil, err
		}
		values["icla_fields"] = iclaFieldJson
	}
	if len(in.CclaFields) != 0 {
		cclaFieldJson, err := json.Marshal(in.CclaFields)
		if err != nil {
			return nil, err
		}
		values["ccla_fields"] = cclaFieldJson
	}
	stmt := sqlz.Newx(r.GetDB()).
		InsertInto(CLATemplatesTable).
		ValueMap(values).
		Returning("id", "created_at", "updated_at")
	err := stmt.GetRow(&result)
	if err != nil {
		return nil, err
	}
	return &models.ClaTemplate{
		CclaFields:   in.CclaFields,
		CclaHTMLBody: in.CclaHTMLBody,
		CreatedAt:    result.CreatedAt.Int64,
		Description:  in.Description,
		IclaFields:   in.IclaFields,
		IclaHTMLBody: in.IclaHTMLBody,
		ID:           result.ID.String,
		MetaFields:   in.MetaFields,
		Name:         in.Name,
		UpdatedAt:    result.CreatedAt.Int64,
	}, nil
}
