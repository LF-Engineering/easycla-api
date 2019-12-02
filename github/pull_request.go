package github

import (
	"context"

	"github.com/google/go-github/github"
)

// Commits contains commits for particular pull request
type Commits struct {
	commits []*github.RepositoryCommit
}

// CommitAuthor is information about author of particular commit
type CommitAuthor struct {
	SHA            string
	AuthorID       int64
	AuthorUsername string
	AuthorEmail    string
}

// ListCommits returns list of commits on pull request
func (gc *Client) ListCommits(owner string, repo string, pullReqNumber int) (*Commits, error) {
	commits, _, err := gc.client.PullRequests.ListCommits(context.TODO(), owner, repo, pullReqNumber, nil)
	if err != nil {
		return nil, err
	}
	return &Commits{commits: commits}, nil
}

// GetCommitAuthors function return list of authors of pull request
func (c *Commits) GetCommitAuthors() []*CommitAuthor {
	var authors []*CommitAuthor
	for _, commit := range c.commits {
		var ca CommitAuthor
		ca.SHA = commit.GetSHA()
		if commit.GetAuthor() != nil {
			// github named user
			namedAuthor := commit.GetAuthor()
			if namedAuthor.GetName() != "" {
				ca.AuthorUsername = namedAuthor.GetName()
			} else if namedAuthor.GetLogin() != "" {
				ca.AuthorUsername = namedAuthor.GetLogin()
			}
			ca.AuthorID = namedAuthor.GetID()
			ca.AuthorEmail = namedAuthor.GetEmail()
		} else if commit.GetCommit().GetAuthor() != nil {
			// git user
			commitAuther := commit.GetCommit().GetAuthor()
			ca.AuthorUsername = commitAuther.GetName()
			ca.AuthorEmail = commitAuther.GetEmail()
		}
		authors = append(authors, &ca)
	}
	return authors
}
