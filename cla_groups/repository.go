package cla_groups

import (
	"errors"

	"github.com/communitybridge/easycla-api/gen/models"
	"github.com/ido50/sqlz"
	"github.com/jmoiron/sqlx"
)

// CLAGroupsTable is the name of cla_groups table in database
const (
	CLAGroupsTable              = "cla.cla_groups"
	CLAGroupProjectManagerTable = "cla.cla_group_project_managers"
)

var (
	ErrClaGroupNameAlreadyExist = errors.New("Given CLA Group name already exist")
)

const (
	DuplicateCLAGroupNameError = `pq: duplicate key value violates unique constraint "cla_groups_cla_group_name_key"`
)

// Repository interface defines methods of cla_groups repository service
type Repository interface {
	CreateCLAGroup(in *models.CreateClaGroup) (*models.ClaGroup, error)
}

type repository struct {
	db *sqlx.DB
}

// NewRepository creates new instance of cla_groups repository
func NewRepository(dbConn *sqlx.DB) Repository {
	return &repository{
		db: dbConn,
	}
}

func (r *repository) GetDB() *sqlx.DB {
	return r.db
}

func (r *repository) CreateCLAGroup(in *models.CreateClaGroup) (*models.ClaGroup, error) {
	var claGroupID string
	values := make(map[string]interface{})
	values["cla_group_name"] = in.ClaGroupName
	values["foundation_id"] = in.FoundationID
	values["ccla_enabled"] = in.CclaEnabled
	values["icla_enabled"] = in.IclaEnabled
	err := sqlz.Newx(r.GetDB()).Transactional(func(tx *sqlz.Tx) error {
		err := tx.
			InsertInto(CLAGroupsTable).
			ValueMap(values).
			Returning("id").
			GetRow(&claGroupID)
		if err != nil {
			return err
		}
		for _, projectManagerID := range in.ProjectManagers {
			_, err = tx.
				InsertInto(CLAGroupProjectManagerTable).
				Columns("cla_group_id", "project_manager_id").
				Values(claGroupID, projectManagerID).
				Exec()
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		switch err.Error() {
		case DuplicateCLAGroupNameError:
			return nil, ErrClaGroupNameAlreadyExist
		default:
			return nil, err
		}
	}
	return &models.ClaGroup{
		ID:              claGroupID,
		ClaGroupName:    *in.ClaGroupName,
		FoundationID:    *in.FoundationID,
		ProjectManagers: in.ProjectManagers,
	}, nil
}
