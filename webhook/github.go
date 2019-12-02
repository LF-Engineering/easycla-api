package webhook

import (
	"log"

	"github.com/communitybridge/easycla-api/config"
	"github.com/communitybridge/easycla-api/gen/restapi/operations"
	"github.com/communitybridge/easycla-api/gen/restapi/operations/webhook"
	"github.com/go-openapi/runtime/middleware"
	"github.com/google/go-github/github"
)

// Configure github webhook
func Configure(api *operations.ClaAPI, config config.Config) {
	api.WebhookGithubWebhookHandler = webhook.GithubWebhookHandlerFunc(
		func(params webhook.GithubWebhookParams) middleware.Responder {
			handleGithubEvent(params, []byte(config.GithubWebhookSecret))
			return webhook.NewGithubWebhookOK()
		})
}

func handleGithubEvent(params webhook.GithubWebhookParams, webhookSecretKey []byte) {
	r := params.HTTPRequest
	payload, err := github.ValidatePayload(r, webhookSecretKey)
	if err != nil {
		log.Println("error", err)
		return
	}
	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		log.Println("error", err)
		return
	}
	switch event := event.(type) {
	case *github.InstallationEvent:
		log.Printf(`got installation event 
			installation id : %d,
			app_id : %d,
	`, event.GetInstallation().GetID(),
			event.GetInstallation().GetAppID())
	case *github.PullRequestEvent:
		ProcessPullRequestEvent(event)
	default:
		log.Println(event)
	}
}
