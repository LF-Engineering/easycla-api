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
		if !assert.Nil(t, err) {
			return
		}
		if !assert.NotNil(t, got) {
			return
		}
		if !assert.Equal(t, int64(1), countOfTemplatesInDB()) {
			return
		}
		if !strfmt.IsUUID4(got.ID) {
			t.Fail()
			return
		}
		if got.CreatedAt < timeBeforeCreation {
			t.Fail()
			return
		}
		if got.UpdatedAt < timeBeforeCreation {
			t.Fail()
			return
		}
		if !assert.Equal(t, int64(1), got.Version) {
			return
		}
		expected := *template1
		expected.CreatedAt = got.CreatedAt
		expected.UpdatedAt = got.UpdatedAt
		expected.ID = got.ID
		expected.Version = 1

		if !assert.Equal(t, &expected, got) {
			return
		}

		gotTemplate, err := claTemplatesService.GetCLATemplate(&params.GetCLATemplateParams{ClaTemplateID: got.ID})
		if !assert.Nil(t, err) {
			return
		}
		if !assert.NotNil(t, gotTemplate) {
			return
		}
		if !assert.Equal(t, gotTemplate, got) {
			return
		}
	})
}
