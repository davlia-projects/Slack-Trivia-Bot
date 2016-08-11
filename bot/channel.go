package bot

import (
	"fmt"
	"log"
	"time"

	"github.com/dota-2-slack-bot/config"
	"github.com/dota-2-slack-bot/logic"
)

// var client *c.Client = c.GetClient()

type Channel struct {
	ID             string
	GameInstance   *logic.Game
	ContinuousMode bool
	Name           string
	HintTicker     *time.Ticker
	QuestionTimer  *time.Timer
	Config         config.Config
}

func NewChannel(conf config.Config, name, id string) *Channel {
	newGame, err := logic.NewGame(conf)
	if err != nil {
		log.Printf("error: could not create new game instance (%+v)\n", err)
	}
	c := &Channel{
		ID:           id,
		GameInstance: newGame,
		Name:         name,
		Config:       conf,
	}
	return c
}

func (C *Channel) MakeGuess(guess, pid string) {
	if C.GameInstance.CurrentQuestion == nil || pid == "" || guess == "" {
		return
	}
	if C.GameInstance.GetPlayerByPID(pid) == nil {
		user, err := client.API.GetUserInfo(pid)
		if err != nil {
			fmt.Printf("error: could not get user info %d (%+v)\n", pid, err)
		}
		C.GameInstance.CreatePlayer(pid, user.Name)
	}
	correct, streakChange := C.GameInstance.MakeGuess(guess, pid)
	if correct {
		C.HintTicker.Stop()
		C.QuestionTimer.Stop()
		if streakChange {
			C.GameInstance.ClearStreak()
			C.sendMessage("Correct answer with streak change")
		} else {
			C.sendMessage("Correct answer without streak change")
		}
		if C.ContinuousMode {
			C.QuestionCommand()
		}
	}
}

// QuestionCommand returns a question if it exists. Otherwise it will start a new round create one.
func (C *Channel) QuestionCommand() {
	if C.GameInstance.CurrentQuestion == nil {
		C.GameInstance.NewRound()
		C.HintTicker = time.NewTicker(time.Second * C.Config.HintDelay)
		C.QuestionTimer = time.NewTimer(time.Second * C.Config.QuestionTime)
		go func() {
			for _ = range C.HintTicker.C {
				C.GameInstance.NextHint()
				C.sendMessage(C.GameInstance.CurrentHint.Stars)
			}
		}()
		go func() {
			<-C.QuestionTimer.C
			C.sendMessage(C.GameInstance.CurrentQuestion.Answer)
			C.HintTicker.Stop()
		}()
	}
	C.sendMessage(C.GameInstance.CurrentQuestion.Prompt)
}

// HintCommand checks and returns a hint and true if it exists. Otherwise empty string and false will be returned.
func (C *Channel) HintCommand() {
	if C.GameInstance.CurrentHint != nil {
		C.sendMessage(C.GameInstance.CurrentHint.Stars)
	}
}

func (C *Channel) ContinuousModeOn() {
	C.ContinuousMode = true
}

func (C *Channel) ContinuousModeOff() {
	C.ContinuousMode = false
}

func (C *Channel) GetPlayer(pid string) *logic.Player {
	return C.GameInstance.GetPlayerByPID(pid)
}

func (C *Channel) GetStatsForPlayer(pid string) {
	player := C.GameInstance.GetPlayerByPID(pid)
	C.sendMessage(fmt.Sprintf("%s stats - Score: %d Streak: %d", player.Name, player.Score, player.Streak))
}

func (C *Channel) sendMessage(message string) {
	client.API.PostMessage(C.ID, message, params)
	// client.RTM.SendMessage(client.RTM.NewOutgoingMessage(message, C.ID))
}
