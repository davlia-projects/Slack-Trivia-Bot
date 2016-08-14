package client

import (
	"os"
	"sync"

	"github.com/nlopes/slack"
)

var (
	slackOnce   sync.Once
	slackClient *SlackClient
)

type SlackClient struct {
	API *slack.Client
	RTM *slack.RTM
}

func GetSlackClient() *SlackClient {
	slackOnce.Do(func() {
		apiKey := os.Getenv("BOT_API_KEY")
		api := slack.New(apiKey)
		slackClient = &SlackClient{
			API: api,
			RTM: api.NewRTM(),
		}
	})
	return slackClient
}
