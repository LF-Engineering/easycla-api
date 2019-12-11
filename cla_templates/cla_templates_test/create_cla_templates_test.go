package cla_templates_test

import (
	"testing"
	"time"

	"github.com/go-openapi/strfmt"

	"github.com/stretchr/testify/assert"

	"github.com/communitybridge/easycla-api/gen/models"
	params "github.com/communitybridge/easycla-api/gen/restapi/operations/cla_templates"
)

func TestCreateCLATemplate(t *testing.T) {
	prepareTestDatabase()
	clearTemplatesTable()
	t.Run("valid input", func(t *testing.T) {
		timeBeforeCreation := time.Now().Unix()
		got, err := claTemplatesService.CreateCLATemplate(&params.CreateCLATemplateParams{
			ClaTemplate: &models.ClaTemplateInput{
				CclaFields:   template1.CclaFields,
				CclaHTMLBody: template1.CclaHTMLBody,
				Description:  template1.Description,
				IclaFields:   template1.IclaFields,
				IclaHTMLBody: template1.IclaHTMLBody,
				MetaFields:   template1.MetaFields,
				Name:         template1.Name,
			},
		})
		if !assert.Nil(t, err, "err should be nil") {
			return
		}
		if !assert.NotNil(t, got, "response should not be nil") {
			return
		}
		if !assert.Equal(t, int64(1), countOfTemplatesInDB(), "template should be created in database") {
			return
		}
		if !strfmt.IsUUID4(got.ID) {
			t.Log("id should be of type UUID")
			t.Fail()
			return
		}
		if got.CreatedAt < timeBeforeCreation {
			t.Log("invalid created_at time")
			t.Fail()
			return
		}
		if got.UpdatedAt < timeBeforeCreation {
			t.Log("invalid updated_at time")
			t.Fail()
			return
		}
		if !assert.Equal(t, int64(1), got.Version, "template version should be 1") {
			return
		}
		expected := *template1
		expected.CreatedAt = got.CreatedAt
		expected.UpdatedAt = got.UpdatedAt
		expected.ID = got.ID
		expected.Version = 1

		if !assert.Equal(t, &expected, got, "actual response should match with expected response") {
			return
		}

		gotTemplate, err := claTemplatesService.GetCLATemplate(&params.GetCLATemplateParams{ClaTemplateID: got.ID})
		if !assert.Nil(t, err, "err should be nil") {
			return
		}
		if !assert.NotNil(t, gotTemplate, "response should not be nil") {
			return
		}
		if !assert.Equal(t, gotTemplate, got, "values in database must equal to actual response") {
			return
		}
	})
}
