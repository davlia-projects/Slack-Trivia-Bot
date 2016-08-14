package client

import (
	"os"
	"sync"

	"github.com/nlopes/slack"
)

var (
	once   sync.Once
	client *Client
)

type Client struct {
	API *slack.Client
	RTM *slack.RTM
}

func GetClient() *Client {
	once.Do(func() {
		apiKey := os.Getenv("BOT_API_KEY")
		api := slack.New(apiKey)
		client = &Client{
			API: api,
			RTM: api.NewRTM(),
		}
	})
	return client
}
