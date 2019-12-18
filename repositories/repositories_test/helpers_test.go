package repositories_test

import (
	"log"

	"github.com/communitybridge/easycla-api/repositories"
	"github.com/ido50/sqlz"
)

func newString(s string) *string {
	return &s
}

func newInt64(i int64) *int64 {
	return &i
}

func newBool(b bool) *bool {
	return &b
}

func isRepositoryPresent(repoID string) bool {
	count, err := sqlz.Newx(sqlxDB).
		Select("*").
		From(repositories.CLARepositoryTable).
		Where(sqlz.Eq("id", repoID)).GetCount()
	if err != nil {
		log.Fatal(err)
	}
	return count == 1
}

func deleteAllRepositories() {
	_, err := sqlz.Newx(sqlxDB).
		DeleteFrom(repositories.CLARepositoryTable).
		Exec()
	if err != nil {
		log.Fatal(err)
	}
}

func numberOfRepositories() int64 {
	count, err := sqlz.Newx(sqlxDB).
		Select("*").
		From(repositories.CLARepositoryTable).
		GetCount()
	if err != nil {
		log.Fatal(err)
	}
	return count
}
