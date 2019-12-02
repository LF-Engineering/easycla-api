package github

import (
	"errors"

	"github.com/communitybridge/easycla-api/config"

	"github.com/google/go-github/github"
)

var gs *service

type service struct {
	client *github.Client
}

var githubAppPrivateKey string
var githubAppID int

// Init function initialize github service.
// service is used internally only by this package.
// TODO() : add authorization in client.
func Init(conf config.Config) {
	gs = &service{
		client: github.NewClient(nil),
	}
	githubAppPrivateKey = conf.GithubAppPrivateKey
	githubAppID = conf.GithubAppID
}

func getService() (*service, error) {
	if gs == nil {
		return nil, errors.New("github service not initialized")
	}
	return gs, nil
}

func getGitubAppPrivateKey() string {
	return githubAppPrivateKey
}

func getGithubAppID() int {
	return githubAppID
}
