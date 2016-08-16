package bot

import (
	"fmt"
	"log"

	c "github.com/dota-2-slack-bot/client"
	"github.com/dota-2-slack-bot/config"
	"github.com/dota-2-slack-bot/logic"

	"github.com/nlopes/slack"
)

type Bot struct {
	Channels map[string]*logic.GameInstance
}

var (
	defaultConfig config.Config = config.Config{
		MaxPoints:    5,
		MaxHintCount: 3,
		HintDelay:    5,
		QuestionTime: 30,
	}
	slackClient    *c.SlackClient    = c.GetSlackClient()
	questionClient *c.QuestionClient = c.GetQuestionClient()
)

func NewBot() *Bot {
	s := &Bot{
		Channels: map[string]*logic.GameInstance{},
	}
	return s
}

func (B *Bot) onStart() {
	channels, err := slackClient.API.GetChannels(true)
	if err != nil {
		log.Fatalf("error: could not retrieve list of channels (%+v)\n", err)
	}
	for _, channel := range channels {
		if channel.IsMember {
			B.Channels[channel.ID] = logic.NewGameInstance(defaultConfig, channel.Name, channel.ID)
		}
	}
}

func (B *Bot) HandleMessageEvent(ev *slack.MessageEvent) {
	fmt.Printf("%+v\n", ev)
	channel := B.Channels[ev.Channel]
	if ev.BotID != "" { // ignore bots
		return
	}
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
	default:
		if ev.SubType == "" {
			channel.MakeGuess(ev.Text, ev.User)
		}
	}
}

func (B *Bot) HandleChannelJoinedEvent(ev *slack.ChannelJoinedEvent) {
	fmt.Printf("Joined channel %s\n", ev.Channel.Name)
	B.Channels[ev.Channel.ID] = logic.NewGameInstance(defaultConfig, ev.Channel.Name, ev.Channel.ID)

}

func (B *Bot) Run() {
	B.onStart()
	go slackClient.RTM.ManageConnection()
Loop:
	for {
		select {
		case msg := <-slackClient.RTM.IncomingEvents:
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
