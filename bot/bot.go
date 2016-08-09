package bot

import (
	"fmt"
	"log"

	"github.com/dota-2-slack-bot/logic"
	"github.com/nlopes/slack"
)

type Channel struct {
	GameInstance *logic.Game
}

type Bot struct {
	API      *slack.Client
	RTM      *slack.RTM
	Channels map[string]*Channel
}

var defaultConfig logic.Config = logic.Config{
	MaxPoints:    5,
	MaxHintCount: 5,
}

func NewBot() *Bot {
	client := slack.New("xoxb-57834688131-S4MhbAfABG2iURPN0HhzwGYb")
	s := &Bot{
		API:      client,
		RTM:      client.NewRTM(),
		Channels: map[string]*Channel{},
	}
	return s
}

func (B *Bot) onStart() {
	channels, err := B.API.GetChannels(true)
	if err != nil {
		log.Fatalf("error: could not retrieve list of channels (%+v)\n", err)
	}
	for _, channel := range channels {
		if channel.IsMember {
			newGame, err := logic.NewGame(defaultConfig)
			if err != nil {
				log.Printf("error: could not create new game instance (%+v)\n", err)
				continue
			}
			B.Channels[channel.ID] = &Channel{
				GameInstance: newGame,
			}
		}
	}
}

func (B *Bot) HandleMessageEvent(ev *slack.MessageEvent) {
	fmt.Printf("Message: %+v\n", ev)
	switch ev.Text {
	case "!q":

	case "!h":

	case "!c":

	case "!o":

	case "!s":

	default:

	}
}

func (B *Bot) HandleChannelJoinedEvent(ev *slack.ChannelJoinedEvent) {
	fmt.Printf("Joined channel %s\n", ev.Channel.Name)
	B.Channels[ev.Channel.ID] = &Channel{}

}

func (B *Bot) Run() {
	B.onStart()
	go B.RTM.ManageConnection()
Loop:
	for {
		select {
		case msg := <-B.RTM.IncomingEvents:
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
