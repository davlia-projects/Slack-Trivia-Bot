package main

import (
	"fmt"

	"github.com/dota-2-slack-bot/bot"
)

func main() {
	r := bot.NewBot()
	fmt.Printf("Starting bot\n")
	r.Run()
}
