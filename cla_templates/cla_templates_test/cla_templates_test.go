package cla_templates_test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/ido50/sqlz"

	"github.com/communitybridge/easycla-api/cla_templates"

	"github.com/communitybridge/easycla-api/gen/models"

	"github.com/spf13/viper"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"gopkg.in/testfixtures.v2"
)

var (
	sqlxDB   *sqlx.DB
	fixtures *testfixtures.Context
)

var claTemplatesService cla_templates.Service

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
	claTemplatesRepo := cla_templates.NewRepository(sqlxDB)
	claTemplatesService = cla_templates.NewService(claTemplatesRepo)
	os.Exit(m.Run())
}

func prepareTestDatabase() {
	if err := fixtures.Load(); err != nil {
		log.Fatal(err)
	}
}

func clearTemplatesTable() {
	_, err := sqlz.Newx(sqlxDB).DeleteFrom(cla_templates.CLATemplatesTable).Exec()
	if err != nil {
		log.Fatal("Unable to delete all entries in cla_templates table")
	}
}

func countOfTemplatesInDB() int64 {
	count, err := sqlz.Newx(sqlxDB).
		Select("*").
		From(cla_templates.CLATemplatesTable).
		GetCount()
	if err != nil {
		log.Fatal(err)
	}
	return count
}

func newInt64(v int64) *int64 {
	return &v
}

func newString(v string) *string {
	return &v
}

var (
	template1 = testTemplateList.ClaTemplates[0]
	template2 = testTemplateList.ClaTemplates[1]
)

var testTemplateList = &models.ClaTemplateList{
	ClaTemplates: []*models.ClaTemplate{
		{
			CclaFields: []*models.Field{
				{
					AnchorString: "Please sign:",
					FieldType:    "sign",
					Height:       0,
					ID:           "sign",
					IsEditable:   false,
					IsOptional:   false,
					Name:         "Please Sign",
					Offsetx:      100,
					Offsety:      -6,
					Width:        0,
				},
			},
			CclaHTMLBody: "<html><body><p>CCLA Project Name: {{ PROJECT_NAME }}</p></body></html>",
			CreatedAt:    1575971771,
			Description:  "For use of projects under the Apache style of CLA.",
			IclaFields: []*models.Field{
				{
					AnchorString: "Full name:",
					FieldType:    "text_unlocked",
					Height:       20,
					ID:           "full_name",
					IsEditable:   false,
					IsOptional:   false,
					Name:         "Full Name",
					Offsetx:      65,
					Offsety:      -8,
					Width:        340,
				},
			},
			IclaHTMLBody: "<html><body><p>ICLA Project Name: {{ PROJECT_NAME }}</br></p></body></html>",
			ID:           "132c332e-66ff-4361-9833-594344a072bf",
			MetaFields: []*models.MetaField{
				{
					Description:      "Project's Full Name.",
					Name:             "Project Name",
					TemplateVariable: "PROJECT_NAME",
					Value:            "",
				},
			},
			Name:      "Apache Style",
			UpdatedAt: 1575971771,
			Version:   1,
		},
		{
			CclaFields: []*models.Field{
				{
					AnchorString: "Date:",
					FieldType:    "date",
					Height:       0,
					ID:           "date",
					IsEditable:   false,
					IsOptional:   false,
					Name:         "Date",
					Offsetx:      40,
					Offsety:      -7,
					Width:        0,
				},
			},
			CclaHTMLBody: "<html><body><p>CCLA Project Name: {{ PROJECT_NAME }}</p></body></html>",
			CreatedAt:    1575972053,
			Description:  "For use of projects under the LF style of CLA.",
			IclaFields: []*models.Field{
				{
					AnchorString: "Mailing Address:",
					FieldType:    "text_unlocked",
					Height:       20,
					ID:           "mailing_address1",
					IsEditable:   false,
					IsOptional:   false,
					Name:         "Mailing Address",
					Offsetx:      105,
					Offsety:      -7,
					Width:        300,
				},
			},
			IclaHTMLBody: "<html><body><p>ICLA Project Name: {{ PROJECT_NAME }}</br></p></body></html>",
			ID:           "03405605-6b1b-4ece-a280-f97e4d91fea8",
			MetaFields: []*models.MetaField{
				{
					Description:      "Project's Full Name.",
					Name:             "Project Name",
					TemplateVariable: "PROJECT_NAME",
					Value:            "",
				},
			},
			Name:      "LF Style",
			UpdatedAt: 1575972094,
			Version:   2,
		},
	},
}
