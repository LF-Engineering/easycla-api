package cla_templates

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/ido50/sqlz"

	"github.com/communitybridge/easycla-api/gen/models"
	"github.com/jmoiron/sqlx"
)

const (
	// CLATemplatesTable is the name of events table in database
	CLATemplatesTable = "cla.cla_templates"
	// NoResultErrorString in error return by sql when it does not get any result
	NoResultErrorString = "sql: no rows in result set"
)

var (
	// ErrClaTemplateNotFound is error returned when requested cla template does not exist in system
	ErrClaTemplateNotFound = errors.New("cla template does not exist")
)

// Repository interface defines methods of cla_templates repository service
type Repository interface {
	CreateCLATemplate(template *models.ClaTemplateInput) (*models.ClaTemplate, error)
	GetCLATemplate(claTemplateID string) (*models.ClaTemplate, error)
	UpdateCLATemplate(claTemplateID string, template *models.ClaTemplateInput) (*models.ClaTemplate, error)
	DeleteCLATemplate(claTemplateID string) error
	ListCLATemplates() (*models.ClaTemplateList, error)
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
		metaFieldJSON, err := json.Marshal(in.MetaFields)
		if err != nil {
			return nil, err
		}
		values["meta_fields"] = metaFieldJSON
	}
	if len(in.IclaFields) != 0 {
		iclaFieldJSON, err := json.Marshal(in.IclaFields)
		if err != nil {
			return nil, err
		}
		values["icla_fields"] = iclaFieldJSON
	}
	if len(in.CclaFields) != 0 {
		cclaFieldJSON, err := json.Marshal(in.CclaFields)
		if err != nil {
			return nil, err
		}
		values["ccla_fields"] = cclaFieldJSON
	}
	stmt := sqlz.Newx(r.GetDB()).
		InsertInto(CLATemplatesTable).
		ValueMap(values).
		Returning("id", "created_at", "updated_at")
	err := stmt.GetRow(&result)
	if err != nil {
		return nil, err
	}
	return r.GetCLATemplate(result.ID.String)
}

func (r *repository) DeleteCLATemplate(claTemplateID string) error {
	res, err := sqlz.Newx(r.GetDB()).
		DeleteFrom(CLATemplatesTable).
		Where(sqlz.Eq("id", claTemplateID)).
		Exec()
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrClaTemplateNotFound
	}
	return nil
}

func (r *repository) GetCLATemplate(claTemplateID string) (*models.ClaTemplate, error) {
	var template SQLCLATemplate
	err := sqlz.Newx(r.GetDB()).
		Select("*").
		From(CLATemplatesTable).
		Where(sqlz.Eq("id", claTemplateID)).GetRow(&template)
	if err != nil {
		if err.Error() == NoResultErrorString {
			return nil, ErrClaTemplateNotFound
		}
		return nil, err
	}
	return template.toClaTemplate()
}

func (r *repository) UpdateCLATemplate(claTemplateID string, in *models.ClaTemplateInput) (*models.ClaTemplate, error) {
	values := map[string]interface{}{
		"name":        in.Name,
		"description": in.Description,
		"updated_at":  time.Now().Unix(),
	}
	if in.IclaHTMLBody != "" {
		values["icla_html_body"] = in.IclaHTMLBody
	}
	if in.CclaHTMLBody != "" {
		values["ccla_html_body"] = in.CclaHTMLBody
	}
	if len(in.MetaFields) != 0 {
		metaFieldJSON, err := json.Marshal(in.MetaFields)
		if err != nil {
			return nil, err
		}
		values["meta_fields"] = metaFieldJSON
	} else {
		values["meta_fields"] = nil
	}
	if len(in.IclaFields) != 0 {
		iclaFieldJSON, err := json.Marshal(in.IclaFields)
		if err != nil {
			return nil, err
		}
		values["icla_fields"] = iclaFieldJSON
	} else {
		values["icla_fields"] = nil
	}
	if len(in.CclaFields) != 0 {
		cclaFieldJSON, err := json.Marshal(in.CclaFields)
		if err != nil {
			return nil, err
		}
		values["ccla_fields"] = cclaFieldJSON
	} else {
		values["ccla_fields"] = nil
	}
	res, err := sqlz.Newx(r.GetDB()).
		Update(CLATemplatesTable).
		SetMap(values).
		Set("version", sqlz.Indirect("version + 1")).
		Where(sqlz.Eq("id", claTemplateID)).
		Exec()
	if err != nil {
		return nil, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, ErrClaTemplateNotFound
	}
	return r.GetCLATemplate(claTemplateID)
}

func (r *repository) ListCLATemplates() (*models.ClaTemplateList, error) {
	rows, err := sqlz.Newx(r.GetDB()).
		Select("*").
		From(CLATemplatesTable).
		OrderBy(sqlz.Asc("name")).
		GetAllAsRows()
	if err != nil {
		return nil, err
	}
	var result models.ClaTemplateList
	result.ClaTemplates = make([]*models.ClaTemplate, 0)
	defer rows.Close()
	for rows.Next() {
		var t SQLCLATemplate
		err = rows.StructScan(&t)
		if err != nil {
			return nil, err
		}
		template, err := t.toClaTemplate()
		if err != nil {
			return nil, err
		}
		result.ClaTemplates = append(result.ClaTemplates, template)
	}
	return &result, nil
}
