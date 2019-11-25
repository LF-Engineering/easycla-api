package cla_groups

import (
	"database/sql"
)

// SQLCLAGroups represent row of cla.cla_groups table
type SQLCLAGroups struct {
	ID           sql.NullString `db:"id"`
	CLAGroupName sql.NullString `db:"cla_group_name"`
	ProjectID    sql.NullString `db:"project_id"`
	CreatedAt    sql.NullInt64  `db:"created_at"`
	UpdatedAt    sql.NullInt64  `db:"updated_at"`
	CCLAEnabled  sql.NullBool   `db:"ccla_enabled"`
	ICLAEnabled  sql.NullBool   `db:"icla_enabled"`
}

// SQLCLAGroupProjectManagers represents row of cla.cla_group_project_managers table
type SQLCLAGroupProjectManagers struct {
	CLAGroupID       sql.NullString `db:"cla_group_id"`
	ProjectManagerID sql.NullString `db:"project_manager_id"`
}
