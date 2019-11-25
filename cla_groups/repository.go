package cla_groups

import (
	"database/sql"
	"errors"
	"time"

	"github.com/communitybridge/easycla-api/gen/models"
	"github.com/communitybridge/easycla-api/gen/restapi/operations/cla_groups"
	"github.com/ido50/sqlz"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

const (
	// CLAGroupsTable is the name of cla_groups table in database
	CLAGroupsTable = "cla.cla_groups"
	// CLAGroupProjectManagerTable is the name of cla_group_project_managers table in database
	CLAGroupProjectManagerTable = "cla.cla_group_project_managers"
)

var (
	// ErrClaGroupNameAlreadyExist is the error which happens when we create CLA Group with
	// name which already exist in system
	ErrClaGroupNameAlreadyExist = errors.New("given cla group name already exist")
	// ErrClaGroupNotFound is the error that CLA Group does not exist in system
	ErrClaGroupNotFound = errors.New("cla group does not exist")
)

const (
	// DuplicateCLAGroupNameError is an error string
	// returned by postgres when unique contraint of cla_group_name fails
	DuplicateCLAGroupNameError = `pq: duplicate key value violates unique constraint "cla_groups_project_id_cla_group_name_key"`
)

// Repository interface defines methods of cla_groups repository service
type Repository interface {
	CreateCLAGroup(in *models.CreateClaGroup) (*models.ClaGroup, error)
	DeleteCLAGroup(claGroupID string) error
	UpdateCLAGroup(claGroupID string, in *models.UpdateClaGroup) error
	ListCLAGroups(params *cla_groups.ListCLAGroupsParams) (*models.ClaGroupList, error)
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
	values["project_id"] = in.ProjectID
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
		ProjectID:       *in.ProjectID,
		ProjectManagers: in.ProjectManagers,
	}, nil
}

func (r *repository) DeleteCLAGroup(claGroupID string) error {
	err := sqlz.Newx(r.GetDB()).Transactional(func(tx *sqlz.Tx) error {
		var err error
		var res sql.Result
		var rowsAffected int64
		_, err = tx.
			DeleteFrom(CLAGroupProjectManagerTable).
			Where(sqlz.Eq("cla_group_id", claGroupID)).
			Exec()
		if err != nil {
			return err
		}
		res, err = tx.
			DeleteFrom(CLAGroupsTable).
			Where(sqlz.Eq("id", claGroupID)).
			Exec()
		if err != nil {
			return err
		}
		rowsAffected, err = res.RowsAffected()
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			err = ErrClaGroupNotFound
		}
		return err
	})
	return err
}

func (r *repository) UpdateCLAGroup(claGroupID string, in *models.UpdateClaGroup) error {
	values := make(map[string]interface{})
	values["cla_group_name"] = in.ClaGroupName
	values["ccla_enabled"] = in.CclaEnabled
	values["icla_enabled"] = in.IclaEnabled
	values["updated_at"] = time.Now().Unix()
	res, err := sqlz.Newx(r.GetDB()).
		Update(CLAGroupsTable).
		SetMap(values).
		Where(sqlz.Eq("id", claGroupID)).
		Exec()
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrClaGroupNotFound
	}
	return nil
}

func createListCLAGroupQuery(db *sqlx.DB, in *cla_groups.ListCLAGroupsParams) *sqlz.SelectStmt {
	stmt := sqlz.Newx(db).
		Select(`
			cg.id,
			cg.project_id,
			cg.cla_group_name,
			cg.icla_enabled,
			cg.ccla_enabled,
			cg.created_at,
			cg.updated_at,
			array_agg(cgm.project_manager_id) as project_managers`).
		From(CLAGroupsTable+" cg").
		LeftJoin(CLAGroupProjectManagerTable+" cgm", sqlz.Eq("cg.id", sqlz.Indirect("cgm.cla_group_id")))

	var conditions []sqlz.WhereCondition
	if in.ProjectID != nil {
		conditions = append(conditions,
			sqlz.Eq("cg.project_id", *in.ProjectID))
	}

	if in.ProjectManagerID != nil {
		conditions = append(conditions,
			sqlz.SQLCond("cg.id IN (select cla_group_id from cla.cla_group_project_managers where project_manager_id = ?)", *in.ProjectManagerID))
	}

	if len(conditions) != 0 {
		stmt = stmt.Where(conditions...)
	}

	if in.Offset != nil {
		stmt = stmt.Offset(*in.Offset)
	}

	if in.PageSize != nil {
		stmt = stmt.Limit(*in.PageSize)
	}

	stmt = stmt.GroupBy("cg.id")

	orderBy := "cla_group_name"
	if in.OrderBy != nil {
		orderBy = *in.OrderBy
	}

	if in.SortOrder != nil && *in.SortOrder == "desc" {
		stmt = stmt.OrderBy(sqlz.Desc(orderBy))
	} else {
		stmt = stmt.OrderBy(sqlz.Asc(orderBy))
	}
	return stmt
}

func (r *repository) ListCLAGroups(in *cla_groups.ListCLAGroupsParams) (*models.ClaGroupList, error) {
	stmt := createListCLAGroupQuery(r.GetDB(), in)
	var out models.ClaGroupList
	out.ClaGroups = make([]*models.ClaGroup, 0)
	rows, err := stmt.GetAllAsRows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var cg models.ClaGroup
		err := rows.Scan(&cg.ID,
			&cg.ProjectID,
			&cg.ClaGroupName,
			&cg.IclaEnabled,
			&cg.CclaEnabled,
			&cg.CreatedAt,
			&cg.UpdatedAt,
			pq.Array(&cg.ProjectManagers))
		if err != nil {
			return nil, err
		}
		out.ClaGroups = append(out.ClaGroups, &cg)
	}
	return &out, nil
}
