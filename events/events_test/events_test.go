package events_test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/communitybridge/easycla-api/gen/models"

	"github.com/spf13/viper"

	"github.com/communitybridge/easycla-api/events"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"gopkg.in/testfixtures.v2"
)

var (
	sqlxDB   *sqlx.DB
	fixtures *testfixtures.Context
)

var eventsService events.Service

func TestMain(m *testing.M) {
	var err error

	viper.SetDefault("TEST_DATABASE_DSN", "dbname=cla-test password=test user=test port=5432 sslmode=disable")
	err = viper.BindEnv("TEST_DATABASE_DSN")
	if err != nil {
		log.Fatal(err)
	}
	db, err := sql.Open("postgres", viper.GetString("TEST_DATABASE_DSN"))
	if err != nil {
		log.Fatal(err)
	}

	fixtures, err = testfixtures.NewFolder(db, &testfixtures.PostgreSQL{}, "testdata/fixtures")
	if err != nil {
		log.Fatal(err)
	}

	sqlxDB = sqlx.NewDb(db, "postgres")
	eventsRepo := events.NewRepository(sqlxDB)
	eventsService = events.NewService(eventsRepo)

	os.Exit(m.Run())
}

func prepareTestDatabase() {
	if err := fixtures.Load(); err != nil {
		log.Fatal(err)
	}
}

func newInt64(v int64) *int64 {
	return &v
}

func newString(v string) *string {
	return &v
}

var (
	Event1 = &testEventList[0]
	Event2 = &testEventList[1]
	Event3 = &testEventList[2]
	Event4 = &testEventList[3]
	Event5 = &testEventList[4]
	Event6 = &testEventList[5]
)

const (
	// TotalNumberOfEvents = 6
	// UserPrasanna      = "prasannamahajan"
	// UserDdeal         = "ddeal"
	UserMani = "mani"
	// Company1          = "company-1"
	// Company2          = "company-2"
	// Company3          = "company-3"
	// EventUserAdded    = "UserAdded"
	EventUserUpdated = "UserUpdated"
	// EventUserDeleted  = "UserDeleted"
	// EventEmpty        = "EmptyEventData"
	ProjectKubernetes = "Kubernetes"
	ProjectPrometheus = "Prometheus"
	// ProjectGrpc       = "gRPC"
)

var testEventList = []models.Event{
	{
		CompanyID: "company-3",
		EventData: "{\"user\": \"mani\", \"event\": \"UserDeleted\"}",
		EventTime: 1,
		EventType: "UserDeleted",
		ID:        "ac3980ce-0374-4bbe-88a0-abb59ad10392",
		ProjectID: "gRPC",
		UserID:    "mani",
	},
	{
		CompanyID: "company-2",
		EventData: "{\"user\": \"prasannamahajan\", \"event\": \"UserUpdated\"}",
		EventTime: 2,
		EventType: "UserUpdated",
		ID:        "7c493c31-3dd7-475f-84b3-fc8f876eceb5",
		ProjectID: "Prometheus",
		UserID:    "prasannamahajan",
	},
	{
		CompanyID: "company-1",
		EventData: "",
		EventTime: 3,
		EventType: "EmptyEventData",
		ID:        "3529ad99-8be4-408d-92b1-fa4de6e677e1",
		ProjectID: "Kubernetes",
		UserID:    "mani",
	},
	{
		CompanyID: "company-3",
		EventData: "{\"user\": \"ddeal\", \"event\": \"UserAdded\"}",
		EventTime: 4,
		EventType: "UserAdded",
		ID:        "7e3df980-8288-46ce-8e8a-40e50d8222fe",
		ProjectID: "",
		UserID:    "ddeal",
	},
	{
		CompanyID: "company-1",
		EventData: "{\"user\": \"mani\", \"event\": \"UserDeleted\"}",
		EventTime: 5,
		EventType: "UserDeleted",
		ID:        "296abfb7-07da-4184-b905-63d8423ba2a2",
		ProjectID: "Prometheus",
		UserID:    "mani",
	},
	{
		CompanyID: "company-2",
		EventData: "{\"user\": \"ddeal\", \"event\": \"UserUpdated\"}",
		EventTime: 6,
		EventType: "UserUpdated",
		ID:        "3f7f4ab3-3214-4be8-8ad5-ef42f9f223f0",
		ProjectID: "",
		UserID:    "ddeal",
	},
}
