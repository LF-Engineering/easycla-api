package cla_templates_test

import (
	"testing"
	"time"

	"github.com/communitybridge/easycla-api/gen/models"

	"github.com/communitybridge/easycla-api/cla_templates"

	params "github.com/communitybridge/easycla-api/gen/restapi/operations/cla_templates"
	"github.com/stretchr/testify/assert"
)

func TestUpdateCLATemplate(t *testing.T) {
	prepareTestDatabase()
	t.Run("update on invalid ID", func(t *testing.T) {
		got, err := claTemplatesService.UpdateCLATemplate(&params.UpdateCLATemplateParams{
			ClaTemplateID: "aaaabbbb-cccc-dddd-eeee-aabbccddeeff",
			ClaTemplate:   &models.ClaTemplateInput{},
		})
		if !assert.Equal(t, cla_templates.ErrClaTemplateNotFound, err) {
			return
		}
		if !assert.Nil(t, got) {
			return
		}
	})

	t.Run("update template1 with values like template2", func(t *testing.T) {
		timeBeforeCreation := time.Now().Unix()
		got, err := claTemplatesService.UpdateCLATemplate(&params.UpdateCLATemplateParams{
			ClaTemplateID: template1.ID,
			ClaTemplate: &models.ClaTemplateInput{
				CclaFields:   template2.CclaFields,
				CclaHTMLBody: template2.CclaHTMLBody,
				Description:  template2.Description,
				IclaFields:   template2.IclaFields,
				IclaHTMLBody: template2.IclaHTMLBody,
				MetaFields:   template2.MetaFields,
				Name:         template2.Name,
			},
		})
		if !assert.Nil(t, err, "error should be nil") {
			return
		}
		if !assert.NotNil(t, got, "response shoud not be nil") {
			return
		}
		if got.UpdatedAt < timeBeforeCreation {
			t.Log("Updated at must be greater than current_time")
			t.Fail()
			return
		}
		if !assert.Equal(t, template1.Version+1, got.Version, "template version must get incremented on update") {
			return
		}
		expected := *template2
		expected.CreatedAt = template1.CreatedAt
		expected.UpdatedAt = got.UpdatedAt
		expected.ID = template1.ID
		expected.Version = template1.Version + 1

		if !assert.Equal(t, &expected, got, "response should match with expected response") {
			return
		}
	})
}
