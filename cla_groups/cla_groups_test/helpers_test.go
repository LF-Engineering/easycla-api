package cla_groups_test

import (
	"log"

	"github.com/communitybridge/easycla-api/cla_groups"
	"github.com/ido50/sqlz"
)

func newString(s string) *string {
	return &s
}

func newBool(b bool) *bool {
	return &b
}

func isCLAGroupPresent(claGroupId string) bool {
	count, err := sqlz.Newx(sqlxDB).
		Select("*").
		From(cla_groups.CLAGroupsTable).
		Where(sqlz.Eq("id", claGroupId)).GetCount()
	if err != nil {
		log.Fatal(err)
	}
	return count == 1
}

func numberOfCLAGroups() int64 {
	count, err := sqlz.Newx(sqlxDB).
		Select("*").
		From(cla_groups.CLAGroupsTable).
		GetCount()
	if err != nil {
		log.Fatal(err)
	}
	return count
}

func numberOfProjectManagers(claGroupId string) int64 {
	count, err := sqlz.Newx(sqlxDB).
		Select("*").
		From(cla_groups.CLAGroupProjectManagerTable).
		Where(sqlz.Eq("cla_group_id", claGroupId)).GetCount()
	if err != nil {
		log.Fatal(err)
	}
	return count
}
