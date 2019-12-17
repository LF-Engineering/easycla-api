package repositories_test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/communitybridge/easycla-api/gen/models"

	"github.com/communitybridge/easycla-api/events"
	"github.com/communitybridge/easycla-api/repositories"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"gopkg.in/testfixtures.v2"
)

var (
	sqlxDB   *sqlx.DB
	fixtures *testfixtures.Context
)

var repositoriesService repositories.Service
var eventsService events.Service

func TestMain(m *testing.M) {
	var err error
	var db *sql.DB

	viper.SetDefault("TEST_DATABASE_DSN", "dbname=cla-test password=test user=test port=5432 sslmode=disable")
	err = viper.BindEnv("TEST_DATABASE_DSN")
	if err != nil {
		log.Fatal(err)
	}
	db, err = sql.Open("postgres", viper.GetString("TEST_DATABASE_DSN"))
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

	repositoriesRepo := repositories.NewRepository(sqlxDB)
	repositoriesService = repositories.NewService(repositoriesRepo, eventsService)

	os.Exit(m.Run())
}

func prepareTestDatabase() {
	if err := fixtures.Load(); err != nil {
		log.Fatal(err)
	}
}

var (
	// Prometheus cla-group repo
	ClientJava       = testRepositories[0]
	HAProxyExportter = testRepositories[1]
	NodeExporter     = testRepositories[2]
	Prometheus       = testRepositories[3]
	// Kuberenetes cla_group repo
	Minikube = testRepositories[4]
)

var (
	PrometheusClaGroupID  = "de1e355b-7b2c-45e2-86e0-ee6f995af44b"
	Prometheus2ClaGroupID = "de1e355b-7b2c-45e2-86e0-ee6f995af44a"
	PrometheusProjectID   = "0014100000Te1ERAAA"
	KubernetesClaGroupID  = "4b65422e-b47f-4631-ab7c-47ecf6f2200c"
	KubernetesProjectID   = "0014100000Te1ERAAZ"
)

var testRepositories = []*models.Repository{
	{
		ClaGroupID:       "de1e355b-7b2c-45e2-86e0-ee6f995af44b",
		CreatedAt:        1576582373,
		Enabled:          true,
		ExternalID:       "7997879",
		ID:               "93c404c4-c188-4cbe-b3e0-3a3cbd6279f0",
		Name:             "prometheus/client_java",
		OrganizationName: "prometheus",
		ProjectID:        "0014100000Te1ERAAA",
		RepositoryType:   "github",
		URL:              "https://github.com/prometheus/client_java",
	},
	{
		ClaGroupID:       "de1e355b-7b2c-45e2-86e0-ee6f995af44b",
		CreatedAt:        1576582373,
		Enabled:          true,
		ExternalID:       "7939398",
		ID:               "6115dbd0-a7fc-4e49-910e-fdc6504fa8b5",
		Name:             "prometheus/haproxy_exporter",
		OrganizationName: "prometheus",
		ProjectID:        "0014100000Te1ERAAA",
		RepositoryType:   "github",
		URL:              "https://github.com/prometheus/haproxy_exporter",
	},
	{
		ClaGroupID:       "de1e355b-7b2c-45e2-86e0-ee6f995af44a",
		CreatedAt:        1576582373,
		Enabled:          true,
		ExternalID:       "9524057",
		ID:               "f24bbb69-1187-42cd-9d78-054570aded71",
		Name:             "prometheut/node_exporter",
		OrganizationName: "prometheut",
		ProjectID:        "0014100000Te1ERAAA",
		RepositoryType:   "gerrit",
		URL:              "https://github.com/prometheus/node_exporter",
	},
	{
		ClaGroupID:       "de1e355b-7b2c-45e2-86e0-ee6f995af44b",
		CreatedAt:        1576582373,
		Enabled:          true,
		ExternalID:       "6838921",
		ID:               "ae85c1d2-1e93-4e4d-b6d5-216b1a41ff17",
		Name:             "prometheus/prometheus",
		OrganizationName: "prometheus",
		ProjectID:        "0014100000Te1ERAAA",
		RepositoryType:   "github",
		URL:              "https://github.com/prometheus/prometheus",
	},
	{
		ClaGroupID:       "4b65422e-b47f-4631-ab7c-47ecf6f2200c",
		CreatedAt:        1576582589,
		Enabled:          true,
		ExternalID:       "56353740",
		ID:               "d5ed115b-b697-4b0e-998a-c88ce9a1844e",
		Name:             "kubernetes/minikube",
		OrganizationName: "kubernetes",
		ProjectID:        "0014100000Te1ERAAZ",
		RepositoryType:   "github",
		URL:              "https://github.com/kubernetes/minikube",
	},
}
