package cla_templates_test

import (
	"testing"

	"github.com/communitybridge/easycla-api/gen/models"

	"github.com/stretchr/testify/assert"

	params "github.com/communitybridge/easycla-api/gen/restapi/operations/cla_templates"
)

func TestListCLATemplates(t *testing.T) {
	prepareTestDatabase()
	// when templates are present
	res, err := claTemplatesService.ListCLATemplate(&params.ListCLATemplatesParams{})
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, res) {
		return
	}
	if !assert.Equal(t, 2, len(res.ClaTemplates)) {
		return
	}
	if !assert.Equal(t, &models.ClaTemplateList{ClaTemplates: []*models.ClaTemplate{template1, template2}}, res) {
		return
	}
	// when there are no templates
	clearTemplatesTable()
	res, err = claTemplatesService.ListCLATemplate(&params.ListCLATemplatesParams{})
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, res) {
		return
	}
	if !assert.Equal(t, 0, len(res.ClaTemplates)) {
		return
	}
}
