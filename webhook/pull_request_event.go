package webhook

import (
	ghb "github.com/communitybridge/easycla-api/github"
	log "github.com/communitybridge/easycla-api/logging"
	"github.com/google/go-github/github"
)

// ProcessPullRequestEvent functions processes the pull request event
func ProcessPullRequestEvent(event *github.PullRequestEvent) {
	log.Debugf("received pull request event of action %s\n", event.GetAction())
	installationID := event.GetInstallation().GetID()
	pullReqNumber := event.GetPullRequest().GetNumber()
	owner := event.GetRepo().GetOwner().GetLogin()
	repo := event.GetRepo().GetName()

	client, err := ghb.NewGithubAppClient(installationID)
	if err != nil {
		log.Debugf("Error in ProcessPullRequestEvent: InstallationID=%d, PullRequestNumber=%d, owner=%s, repo =%s, error=%s\n",
			installationID, pullReqNumber, owner, repo, err.Error())
		return
	}
	commits, err := client.ListCommits(owner, repo, pullReqNumber)
	if err != nil {
		log.Debugf("Error in ProcessPullRequestEvent: InstallationID=%d, PullRequestNumber=%d, owner=%s, repo =%s, error=%s\n",
			installationID, pullReqNumber, owner, repo, err.Error())
		return
	}
	log.Debugf("list of commit and authors : \n")
	for _, ca := range commits.GetCommitAuthors() {
		log.Debugf("SHA : %s, auther_id : %d, auther_username : %s, auther_email : %s\n", ca.SHA, ca.AuthorID, ca.AuthorUsername, ca.AuthorEmail)
	}
}
