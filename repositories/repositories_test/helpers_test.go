package repositories_test

import (
	"log"

	"github.com/communitybridge/easycla-api/cla_groups"
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

func isCLAGroupPresent(claGroupID string) bool {
	count, err := sqlz.Newx(sqlxDB).
		Select("*").
		From(cla_groups.CLAGroupsTable).
		Where(sqlz.Eq("id", claGroupID)).GetCount()
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

func getCLAGroup(id string) (*cla_groups.SQLCLAGroups, error) {
	var res cla_groups.SQLCLAGroups
	err := sqlz.Newx(sqlxDB).
		Select("*").From(cla_groups.CLAGroupsTable).
		Where(sqlz.Eq("id", id)).
		GetRow(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func numberOfProjectManagers(claGroupID string) int64 {
	count, err := sqlz.Newx(sqlxDB).
		Select("*").
		From(cla_groups.CLAGroupProjectManagerTable).
		Where(sqlz.Eq("cla_group_id", claGroupID)).GetCount()
	if err != nil {
		log.Fatal(err)
	}
	return count
}
