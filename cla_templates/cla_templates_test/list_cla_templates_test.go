package cla_templates_test

import (
	"testing"

	"github.com/communitybridge/easycla-api/gen/models"

	"github.com/stretchr/testify/assert"

	params "github.com/communitybridge/easycla-api/gen/restapi/operations/cla_templates"
)

func TestListCLATemplates(t *testing.T) {
	prepareTestDatabase()
	t.Run("templates are present", func(t *testing.T) {
		res, err := claTemplatesService.ListCLATemplate(&params.ListCLATemplatesParams{})
		if !assert.Nil(t, err, "err should be nil") {
			return
		}
		if !assert.NotNil(t, res, "response should not be nil") {
			return
		}
		if !assert.Equal(t, 2, len(res.ClaTemplates), "number of templates returned should be equal to total number of templates in database") {
			return
		}
		if !assert.Equal(t, &models.ClaTemplateList{ClaTemplates: []*models.ClaTemplate{template1, template2}}, res, "response should match with expected response") {
			return
		}
	})
	t.Run("templates are not present", func(t *testing.T) {
		clearTemplatesTable()
		res, err := claTemplatesService.ListCLATemplate(&params.ListCLATemplatesParams{})
		if !assert.Nil(t, err, "error should be nil") {
			return
		}
		if !assert.NotNil(t, res, "response should not be nil") {
			return
		}
		if !assert.Equal(t, 0, len(res.ClaTemplates), "number of templates received should be zero") {
			return
		}
	})
}
