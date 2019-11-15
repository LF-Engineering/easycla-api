// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package health

import (
	"database/sql"

	log "github.com/communitybridge/easycla-api/logging"
	"github.com/ido50/sqlz"
	"github.com/jmoiron/sqlx"
)

// RepositoryService is the health repository service interface definition
type RepositoryService interface {
	IsAlive() bool
}

type repository struct {
	db *sqlx.DB
}

// NewRepository creates a new health repository service instance
func NewRepository(db *sqlx.DB) RepositoryService {
	return repository{
		db,
	}
}

// IsAlive returns true if we're able to connect to the repository, false otherwise.
func (r repository) IsAlive() bool {
	// Ask about the DB uptime
	dbQuery := sqlz.Newx(r.db).
		Select("count(*) as count").
		From("pg_stat_activity")

	//Select("date_trunc('second', current_timestamp - pg_postmaster_start_time()) as uptime")

	// DEBUG
	rawSQL, bindings := dbQuery.ToSQL(false)
	log.Debugf("HealthCheck Query, SQL: %s with params: %v", rawSQL, bindings)
	// END DEBUG

	// Quick response model to capture the SQL response - use sql.NullString since we might get nulls back
	type ResponseModel struct {
		Count sql.NullInt32 `json:"count"`
	}
	var responses []*ResponseModel
	err := dbQuery.GetAll(&responses)
	if err != nil {
		log.Warnf("Error querying database - error: %v", err)
		return false
	}

	if len(responses) == 0 {
		log.Warn("Unable to query database for uptime.")
		return false
	} else if len(responses) > 1 {
		log.Warn("Multiple query results from uptime query.")
		return false
	}

	if !responses[0].Count.Valid {
		log.Warn("Invalid query response from uptime query.")
		return false
	}
	log.Debugf("Database connections: %d", responses[0].Count.Int32)

	return true
}
