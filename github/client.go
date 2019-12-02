package github

import (
	"log"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/github"
)

// Client is container for go-github Client
type Client struct {
	client *github.Client
}

// NewGithubAppClient function creates github client for installation
func NewGithubAppClient(installationID int64) (*Client, error) {
	itr, err := ghinstallation.New(http.DefaultTransport, getGithubAppID(), int(installationID), []byte(getGitubAppPrivateKey()))
	if err != nil {
		log.Println("Unable to create github installation client", err)
		return nil, err
	}
	client := github.NewClient(&http.Client{Transport: itr})
	return &Client{
		client: client,
	}, nil
}
