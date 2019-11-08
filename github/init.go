package github

import (
	"errors"

	"github.com/google/go-github/github"
)

var gs *service

type service struct {
	client *github.Client
}

// Init function initialize github service.
// service is used internally only by this package.
// TODO() : add authorization in client.
func Init() {
	gs = &service{
		client: github.NewClient(nil),
	}
}

func getService() (*service, error) {
	if gs == nil {
		return nil, errors.New("github service not initialized")
	}
	return gs, nil
}
