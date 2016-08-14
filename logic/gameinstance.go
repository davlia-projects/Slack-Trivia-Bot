package logic

import (
	"fmt"
	"log"
	"time"

	c "github.com/dota-2-slack-bot/client"
	"github.com/dota-2-slack-bot/config"
	"github.com/nlopes/slack"
)

var (
	slackClient    *c.SlackClient              = c.GetSlackClient()
	questionClient *c.QuestionClient           = c.GetQuestionClient()
	params         slack.PostMessageParameters = slack.NewPostMessageParameters()
)

type GameInstance struct {
	ID             string
	Game           *Game
	ContinuousMode bool
	Name           string
	HintTicker     *time.Ticker
	QuestionTimer  *time.Timer
	Config         config.Config
}

func NewGameInstance(conf config.Config, name, id string) *GameInstance {
	newGame, err := NewGame(conf)
	params.AsUser = true
	if err != nil {
		log.Printf("error: could not create new game instance (%+v)\n", err)
	}
	c := &GameInstance{
		ID:     id,
		Game:   newGame,
		Name:   name,
		Config: conf,
	}
	return c
}

func (C *GameInstance) MakeGuess(guess, pid string) {
	if C.Game.CurrentQuestion == nil || pid == "" || guess == "" {
		return
	}
	if C.Game.GetPlayerByPID(pid) == nil {
		user, err := slackClient.API.GetUserInfo(pid)
		if err != nil {
			fmt.Printf("error: could not get user info %d (%+v)\n", pid, err)
		}
		C.Game.CreatePlayer(pid, user.Name)
	}
	player := C.Game.GetPlayerByPID(pid)
	isCorrect := C.Game.MakeGuess(guess)
	if isCorrect {
		C.HintTicker.Stop()
		C.QuestionTimer.Stop()
		awardedPoints, streakChange := C.Game.Correct(pid)
		if streakChange {
			oldPlayer := C.Game.GetPlayerWithStreak()
			oldStreak := oldPlayer.Streak
			C.Game.SetNewStreak(pid)
			C.sendMessage(fmt.Sprintf("%s is correct. +%d points (total score: %d streak: %d). %s's %d win streak has been ended!", player.Name, awardedPoints, player.Score, player.Streak, oldPlayer.Name, oldStreak))
		} else {
			C.sendMessage(fmt.Sprintf("%s is correct. +%d points (total score: %d streak: %d)", player.Name, awardedPoints, player.Score, player.Streak))
		}
		if C.ContinuousMode {
			C.QuestionCommand()
		}
	} else {
		player.Guesses++
	}
}

// QuestionCommand returns a question if it exists. Otherwise it will start a new round create one.
func (C *GameInstance) QuestionCommand() {
	if C.Game.CurrentQuestion == nil {
		C.Game.Reset()
		C.Game.CurrentQuestion = questionClient.NewQuestion()
		C.HintTicker = time.NewTicker(time.Second * C.Config.HintDelay)
		C.QuestionTimer = time.NewTimer(time.Second * C.Config.QuestionTime)
		go func() {
			for _ = range C.HintTicker.C {
				C.Game.NextHint()
				C.sendMessage(C.Game.CurrentHint.Stars)
			}
		}()
		go func() {
			<-C.QuestionTimer.C
			C.sendMessage(C.Game.CurrentQuestion.Answer)
			C.HintTicker.Stop()
		}()
	}
	C.sendMessage(C.Game.CurrentQuestion.Prompt)
}

// HintCommand checks and returns a hint and true if it exists. Otherwise empty string and false will be returned.
func (C *GameInstance) HintCommand() {
	if C.Game.CurrentHint != nil {
		C.sendMessage(C.Game.CurrentHint.Stars)
	}
}

func (C *GameInstance) ContinuousModeOn() {
	C.ContinuousMode = true
}

func (C *GameInstance) ContinuousModeOff() {
	C.ContinuousMode = false
}

func (C *GameInstance) GetStatsForPlayer(pid string) {
	player := C.Game.GetPlayerByPID(pid)
	C.sendMessage(fmt.Sprintf("%s stats - Score: %d Streak: %d", player.Name, player.Score, player.Streak))
}

func (C *GameInstance) sendMessage(message string) {
	slackClient.API.PostMessage(C.ID, message, params)
	// client.RTM.SendMessage(client.RTM.NewOutgoingMessage(message, C.ID))
}
