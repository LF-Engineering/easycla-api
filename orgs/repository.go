package orgs

import (
	"context"
	"encoding/json"

	"github.com/communitybridge/easycla-api/gen/models"
)

// Repository interface defines methods of organization repository service
type Repository interface {
	GetOrgFoundations(ctx context.Context, sfOrgID string) (models.Organization, error)
}

type repository struct {
}

// NewRepository creates new instance of organization repository
func NewRepository() Repository {
	return &repository{}
}

func (r *repository) GetOrgFoundations(ctx context.Context, sfOrgID string) (models.Organization, error) {
	var result models.Organization
	compactResult := `{"<salesForceOrganizationId>":[{"foundationId":"<foundationId>","claGroups":[{"Cloud Native Computing Foundation (CNCF)":{"claGroupId":"9320840-980202983-8901489085","organizationSigned":false,"waitingForSignature":false,"employeesSigned":0,"employeesPending":0,"projectsCovered":["< projectId1 >","< projectId2 >","< projectId3 >","< projectId4 >","< projectId5 >"]}}]},{"foundationId":"<foundationId>","claGroups":[{"O-RAN Software Community (ORAN-SC)":{"claGroupId":"32049890-980202983-8901489085","organizationSigned":true,"waitingForSignature":false,"employeesSigned":25,"employeesPending":10,"projectsCovered":["< projectId1 >","< projectId2 >"]}},{"O-RAN Specification Code Project (ORAN-SCP)":{"claGroupId":"932-092450-0840-980202983-8901489085","organizationSigned":true,"waitingForSignature":false,"employeesSigned":12,"employeesPending":34,"projectsCovered":["< projectId1 >","< projectId2 >"]}}]},{"foundationId":"<foundationId>","claGroups":[{"OpenEXR":{"claGroupId":"01847938-980202983-8901489085","organizationSigned":true,"waitingForSignature":false,"employeesSigned":13,"employeesPending":14,"projectsCovered":["< projectId1 >","< projectId2 >","< projectId3 >"]}},{"OpenVDB":{"claGroupId":"1290-980202983-8901489085","organizationSigned":false,"waitingForSignature":true,"employeesSigned":0,"employeesPending":0,"projectsCovered":["< projectId1 >","< projectId2 >","< projectId3 >"]}},{"OpenColorIO":{"claGroupId":"0871971-980202983-8901489085","organizationSigned":true,"waitingForSignature":false,"employeesSigned":37,"employeesPending":8,"projectsCovered":["< projectId1 >","< projectId2 >","< projectId3 >"]}}]}]}`
	err := json.Unmarshal([]byte(compactResult), &result)
	return result, err
}
