package logic

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"regexp"
	"strings"

	"github.com/dota-2-slack-bot/config"
)

type Player struct {
	ID      string
	Name    string
	Score   int
	Streak  int
	Guesses int
}

type Question struct {
	Category string
	Prompt   string
	Answer   string
}

type Hint struct {
	Stars    string
	Count    int
	Revealed int
}

type Game struct {
	CurrentQuestion  *Question
	CurrentHint      *Hint
	PastQuestions    []*Question
	PlayerWithStreak string
	Players          map[string]*Player // maps pid to player
	Config           config.Config
}

func NewPlayer(id, name string) *Player {
	p := &Player{
		ID:   id,
		Name: name,
	}
	return p
}

func NewGame(conf config.Config) (*Game, error) {
	g := &Game{
		Config:  conf,
		Players: map[string]*Player{},
	}
	return g, nil
}

func (G *Game) NewQuestion() error {
	G.CurrentQuestion = &Question{
		Category: "Testing",
		Prompt:   "Prompt",
		Answer:   "Answer",
	}
	return nil
}

func (G *Game) NextHint() error {
	question := G.CurrentQuestion
	hint := G.CurrentHint
	length := len(question.Answer)
	if hint == nil {
		pat, err := regexp.Compile(`\w`)
		if err != nil {
			log.Fatalf("regex did not compile\n")
		}
		stars := string(pat.ReplaceAll([]byte(question.Prompt), []byte("*")))
		hint := &Hint{
			Stars:    stars,
			Count:    1,
			Revealed: 0,
		}
		G.CurrentHint = hint
	} else if hint.Count >= G.Config.MaxHintCount || hint.Revealed > length/2 {
		return nil
	} else {
		newHint := []string{}
		hintTokens := strings.Split(hint.Stars, " ")
		ansTokens := strings.Split(G.CurrentQuestion.Answer, " ")
		offset := 0
		for t := range hintTokens {
			hintChars := strings.Split(hintTokens[t], "")
			ansChars := strings.Split(ansTokens[t], "")
			for {
				index := rand.Intn(len(hintChars))
				if hintChars[index] == "*" {
					hintChars[index] = ansChars[index]
					break
				}
			}
			newHint = append(newHint, strings.Join(hintChars, ""))
			offset += len(hintTokens[t]) + 1
			hint.Revealed++
		}
		hint.Count++
		hint.Stars = strings.Join(newHint, " ")
		fmt.Printf("%s\n", strings.Join(newHint, " "))
	}
	return nil
}

func (G *Game) ResetGuesses() {
	for _, player := range G.Players {
		player.Guesses = 0
	}
}

func (G *Game) NewRound() error {
	if G.CurrentQuestion != nil {
		G.PastQuestions = append(G.PastQuestions, G.CurrentQuestion)
	}
	G.CurrentQuestion = nil
	G.CurrentHint = nil
	G.ResetGuesses()
	err := G.NewQuestion()
	if err != nil {
		log.Fatalf("error: could not get new question (%+v)\n", err)
	}
	return err
}

func (G *Game) MakeGuess(guess string) (isCorrect bool) {
	return rawString(guess) == rawString(G.CurrentQuestion.Answer)
}

func (G *Game) Correct(pid string) (awardedPoints int, streakChange bool) {
	player := G.Players[pid]
	awardedPoints = int(math.Max(float64(G.Config.MaxPoints-player.Guesses), 1))
	player.Score += awardedPoints
	player.Streak++
	streakChange = G.PlayerWithStreak == pid
	return
}

func (G *Game) SetNewStreak(pid string) {
	G.PlayerWithStreak = pid
	for _, player := range G.Players {
		if player.ID != pid {
			player.Streak = 0
		}
	}
}

func (G *Game) GetPlayerByPID(pid string) *Player {
	player, ok := G.Players[pid]
	if !ok {
		return nil
	}
	return player
}

func (G *Game) CreatePlayer(pid, name string) {
	G.Players[pid] = NewPlayer(pid, name)
}

func rawString(str string) string {
	return strings.Replace(strings.TrimSpace(strings.ToLower(str)), "'", "", -1)
}
