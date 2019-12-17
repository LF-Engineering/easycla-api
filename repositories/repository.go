package repositories

import (
	"errors"

	"github.com/ido50/sqlz"

	"github.com/communitybridge/easycla-api/gen/models"
	params "github.com/communitybridge/easycla-api/gen/restapi/operations/repositories"
	"github.com/jmoiron/sqlx"
)

// RepositoryType is the type of repository
type RepositoryType string

// valid values of RepositoryType
const (
	RepositoryTypeGithub RepositoryType = "github"
	RepositoryTypeGitlab RepositoryType = "gitlab"
	RepositoryTypeGerrit RepositoryType = "gerrit"
)

const (
	// CLARepositoryTable is name of repositories table in database
	CLARepositoryTable = "cla.repositories"
	// CLAGroupsTable is name of cla_groups table in database
	CLAGroupsTable = "cla.cla_groups"
)

// Errors
var (
	ErrInvalidClgGroupAndProjectID = errors.New("given cla_group is not present for given project")
)

// Repository interface defines methods of repository service
type Repository interface {
	CreateRepositories(in *models.CreateRepositoriesInput) (*models.RepositoryList, error)
	DeleteRepositories(in *models.DeleteRepositoriesInput) error
	ListRepositories(in *params.ListRepositoriesParams) (*models.RepositoryList, error)
}

type repository struct {
	db *sqlx.DB
}

// NewRepository creates new instance of audit event repository
func NewRepository(dbConn *sqlx.DB) Repository {
	return &repository{
		db: dbConn,
	}
}

func (r *repository) GetDB() *sqlx.DB {
	return r.db
}

func (r *repository) CreateRepositories(in *models.CreateRepositoriesInput) (*models.RepositoryList, error) {
	var err error
	var ids []string
	type returnResult struct {
		ID string `db:"id"`
	}
	terr := sqlz.Newx(r.GetDB()).Transactional(func(tx *sqlz.Tx) error {
		var count int64
		var id returnResult
		count, err = tx.Select("*").From(CLAGroupsTable).
			Where(sqlz.Eq("id", in.ClaGroupID)).
			Where(sqlz.Eq("project_id", in.ProjectID)).GetCount()
		if err != nil {
			return err
		}
		if count == 0 {
			err = ErrInvalidClgGroupAndProjectID
			return err
		}
		for _, repo := range in.Repositories {
			values := map[string]interface{}{
				"repository_type":   repo.RepositoryType,
				"name":              repo.Name,
				"organization_name": repo.OrganizationName,
				"url":               repo.URL,
				"project_id":        in.ProjectID,
				"cla_group_id":      in.ClaGroupID,
				"external_id":       repo.ExternalID,
			}
			if repo.Enabled != nil {
				values["enabled"] = *repo.Enabled
			}
			stmt := tx.
				InsertInto(CLARepositoryTable).
				ValueMap(values).Returning("id")
			err = stmt.GetRow(&id)
			if err != nil {
				return err
			}
			ids = append(ids, id.ID)
		}
		return nil
	})
	if terr != nil {
		return nil, err
	}
	return r.getRepositories(*in.ProjectID, *in.ClaGroupID, ids)
}

func (r *repository) DeleteRepositories(in *models.DeleteRepositoriesInput) error {
	ids := make([]interface{}, len(in.RepositoryIds))
	for i, v := range in.RepositoryIds {
		ids[i] = v
	}
	_, err := sqlz.Newx(r.GetDB()).
		DeleteFrom(CLARepositoryTable).
		Where(sqlz.Eq("cla_group_id", in.ClaGroupID)).
		Where(sqlz.Eq("project_id", in.ProjectID)).
		Where(sqlz.In("id", ids...)).
		Exec()
	return err
}

func (r *repository) ListRepositories(in *params.ListRepositoriesParams) (*models.RepositoryList, error) {
	stmt := sqlz.Newx(r.GetDB()).
		Select("*").
		From(CLARepositoryTable)

	var conditions []sqlz.WhereCondition
	if in.ClaGroupID != nil {
		conditions = append(conditions, sqlz.Eq("cla_group_id", *in.ClaGroupID))
	}
	conditions = append(conditions, sqlz.Eq("project_id", in.ProjectID))
	if in.RepositoryType != nil {
		conditions = append(conditions, sqlz.Eq("repository_type", *in.RepositoryType))
	}
	if in.RepositoryOrganization != nil {
		conditions = append(conditions, sqlz.Eq("organization_name", *in.RepositoryOrganization))
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
	orderBy := "name"
	if in.OrderBy != nil {
		orderBy = *in.OrderBy
	}
	if in.SortOrder != nil && *in.SortOrder == "desc" {
		stmt = stmt.OrderBy(sqlz.Desc(orderBy))
	} else {
		stmt = stmt.OrderBy(sqlz.Asc(orderBy))
	}
	rows, err := stmt.GetAllAsRows()
	if err != nil {
		return nil, err
	}

	var result models.RepositoryList
	result.Repositories = make([]*models.Repository, 0)
	defer rows.Close()
	for rows.Next() {
		var t SQLRepository
		err = rows.StructScan(&t)
		if err != nil {
			return nil, err
		}
		result.Repositories = append(result.Repositories, t.toRepository())
	}
	return &result, nil
}

func (r *repository) getRepositories(projectID string, claGroupID string, repositoryIDs []string) (*models.RepositoryList, error) {
	ids := make([]interface{}, len(repositoryIDs))
	for i, v := range repositoryIDs {
		ids[i] = v
	}
	rows, err := sqlz.Newx(r.GetDB()).
		Select("*").From(CLARepositoryTable).
		Where(sqlz.Eq("cla_group_id", claGroupID)).
		Where(sqlz.Eq("project_id", projectID)).
		Where(sqlz.In("id", ids...)).GetAllAsRows()
	if err != nil {
		return nil, err
	}
	var result models.RepositoryList
	result.Repositories = make([]*models.Repository, 0)
	defer rows.Close()
	for rows.Next() {
		var t SQLRepository
		err = rows.StructScan(&t)
		if err != nil {
			return nil, err
		}
		result.Repositories = append(result.Repositories, t.toRepository())
	}
	return &result, nil

}
