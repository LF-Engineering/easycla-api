package cla_templates_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/communitybridge/easycla-api/cla_templates"
	params "github.com/communitybridge/easycla-api/gen/restapi/operations/cla_templates"
)

func TestDeleteCLATemplate(t *testing.T) {
	prepareTestDatabase()
	tests := []struct {
		name          string
		claTemplateID string
		wantErr       error
	}{
		{
			name:          "template is present",
			claTemplateID: template1.ID,
			wantErr:       nil,
		},
		{
			name:          "template not present",
			claTemplateID: "aaaabbbb-cccc-dddd-eeee-aabbccddeeff",
			wantErr:       cla_templates.ErrClaTemplateNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			previousCount := countOfTemplatesInDB()
			err := claTemplatesService.DeleteCLATemplate(&params.DeleteCLATemplateParams{
				ClaTemplateID: tt.claTemplateID,
			})
			if err != tt.wantErr {
				t.Errorf("DeleteCLATemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				currentCount := countOfTemplatesInDB()
				if !assert.Equal(t, previousCount-1, currentCount, "cla_template should be deleted") {
					return
				}
			}
		})
	}
}
