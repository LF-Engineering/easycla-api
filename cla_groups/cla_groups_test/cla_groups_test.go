package cla_groups_test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/communitybridge/easycla-api/cla_groups"
	"github.com/communitybridge/easycla-api/events"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"gopkg.in/testfixtures.v2"
)

var (
	db       *sql.DB
	fixtures *testfixtures.Context
)

var claGroupsService cla_groups.Service
var eventsService events.Service

func TestMain(m *testing.M) {
	var err error

	testfixtures.SkipDatabaseNameCheck(true)
	db, err = sql.Open("postgres", "dbname=cla password=prasanna user=prasanna port=5432 sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	/*
		err = testfixtures.GenerateFixtures(db, &testfixtures.PostgreSQL{}, "testdata/fixtures")
		if err != nil {
			log.Fatalf("Error generating fixtures: %v", err)
		}
	*/
	fixtures, err = testfixtures.NewFolder(db, &testfixtures.PostgreSQL{}, "testdata/fixtures")
	if err != nil {
		log.Fatal(err)
	}

	sqlxDB := sqlx.NewDb(db, "postgres")
	eventsRepo := events.NewRepository(sqlxDB)
	eventsService = events.NewService(eventsRepo)

	claGroupsRepo := cla_groups.NewRepository(sqlxDB)
	claGroupsService = cla_groups.NewService(claGroupsRepo, eventsService)

	os.Exit(m.Run())
}

func prepareTestDatabase() {
	if err := fixtures.Load(); err != nil {
		log.Fatal(err)
	}
}
