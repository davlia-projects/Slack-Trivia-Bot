package bot

import (
	"fmt"
	"log"

	c "github.com/dota-2-slack-bot/client"
	"github.com/dota-2-slack-bot/config"

	"github.com/nlopes/slack"
)

type Bot struct {
	Channels map[string]*Channel
}

var (
	defaultConfig config.Config = config.Config{
		MaxPoints:    5,
		MaxHintCount: 3,
		HintDelay:    8,
		QuestionTime: 30,
	}
	client *c.Client                   = c.GetClient()
	params slack.PostMessageParameters = slack.NewPostMessageParameters()
)

func NewBot() *Bot {
	s := &Bot{
		Channels: map[string]*Channel{},
	}
	return s
}

func (B *Bot) onStart() {
	params.AsUser = true
	channels, err := client.API.GetChannels(true)
	if err != nil {
		log.Fatalf("error: could not retrieve list of channels (%+v)\n", err)
	}
	for _, channel := range channels {
		if channel.IsMember {
			B.Channels[channel.ID] = NewChannel(defaultConfig, channel.Name, channel.ID)
		}
	}
}

func (B *Bot) HandleMessageEvent(ev *slack.MessageEvent) {
	fmt.Printf("%+v\n", ev)
	channel := B.Channels[ev.Channel]
	switch ev.Text {
	case "!q":
		channel.QuestionCommand()
	case "!h":
		channel.HintCommand()
	case "!c":
		channel.ContinuousModeOn()
	case "!o":
		channel.ContinuousModeOff()
	case "!s":
		channel.GetStatsForPlayer(ev.User)
	case "!debug":
		client.API.PostMessage(ev.Channel, "debug", params)
	default:
		channel.MakeGuess(ev.Text, ev.User)
	}
}

func (B *Bot) HandleChannelJoinedEvent(ev *slack.ChannelJoinedEvent) {
	fmt.Printf("Joined channel %s\n", ev.Channel.Name)
	B.Channels[ev.Channel.ID] = &Channel{}

}

func (B *Bot) Run() {
	B.onStart()
	go client.RTM.ManageConnection()
Loop:
	for {
		select {
		case msg := <-client.RTM.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				fmt.Printf("Connected\n")

			case *slack.MessageEvent:
				B.HandleMessageEvent(ev)

			case *slack.ChannelJoinedEvent:
				B.HandleChannelJoinedEvent(ev)

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break Loop

			default:
				// fmt.Printf("Unhandled Message Type %s: %+v\n", msg.Type, msg.Data)
			}
		}
	}
}
