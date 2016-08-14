package main

import (
	"fmt"

	"github.com/dota-2-slack-bot/bot"
	"github.com/dota-2-slack-bot/client"
)

func main() {
	r := bot.NewBot()
	_ = client.GetSlackClient()
	_ = client.GetQuestionClient()
	fmt.Printf("Starting bot\n")
	r.Run()
}
