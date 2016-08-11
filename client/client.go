package client

import (
	"sync"

	"github.com/nlopes/slack"
)

const (
	APIKey = "xoxb-57834688131-S4MhbAfABG2iURPN0HhzwGYb"
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
		api := slack.New(APIKey)
		client = &Client{
			API: api,
			RTM: api.NewRTM(),
		}
	})
	return client
}
