package github

import (
	"context"
)

// GetUsernameFromID function returns github username of provided github user id.
func GetUsernameFromID(id int64) (string, error) {
	gs, err := getService()
	if err != nil {
		return "", err
	}
	user, _, err := gs.client.Users.GetByID(context.TODO(), id)
	if err != nil {
		return "", err
	}
	return user.GetLogin(), nil
}

// GetIDFromUsername function returns github id of the user with provided github username
func GetIDFromUsername(username string) (int64, error) {
	gs, err := getService()
	if err != nil {
		return 0, err
	}
	user, _, err := gs.client.Users.Get(context.TODO(), username)
	if err != nil {
		return 0, err
	}
	return user.GetID(), nil
}

// GetUserGithubOrganizations returns list of organizations of which given user is part of.
// It only returns organizations of the user where organization membership info is public.
// It takes github username as input.
func GetUserGithubOrganizations(username string) ([]string, error) {
	gs, err := getService()
	if err != nil {
		return []string{}, err
	}
	organizations, _, err := gs.client.Organizations.List(context.TODO(), username, nil)
	if err != nil {
		return []string{}, err
	}
	orgs := make([]string, len(organizations))
	for _, org := range organizations {
		orgs = append(orgs, org.GetLogin())
	}
	return orgs, nil
}
