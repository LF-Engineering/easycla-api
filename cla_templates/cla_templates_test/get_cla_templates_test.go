package cla_templates_test

import (
	"testing"

	"github.com/communitybridge/easycla-api/cla_templates"

	"github.com/stretchr/testify/assert"

	"github.com/communitybridge/easycla-api/gen/models"
	params "github.com/communitybridge/easycla-api/gen/restapi/operations/cla_templates"
)

func TestGetCLATemplate(t *testing.T) {
	prepareTestDatabase()
	tests := []struct {
		name          string
		claTemplateID string
		want          *models.ClaTemplate
		wantErr       error
	}{
		{
			name:          "template is present",
			claTemplateID: template1.ID,
			want:          template1,
			wantErr:       nil,
		},
		{
			name:          "template not present",
			claTemplateID: "aaaabbbb-cccc-dddd-eeee-aabbccddeeff",
			want:          nil,
			wantErr:       cla_templates.ErrClaTemplateNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := claTemplatesService.GetCLATemplate(&params.GetCLATemplateParams{
				ClaTemplateID: tt.claTemplateID,
			})
			if err != tt.wantErr {
				t.Errorf("GetCLATemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got, "actual response should match with expected response") {
				return
			}
		})
	}
}
