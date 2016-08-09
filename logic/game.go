package logic

import (
	"log"
	"math"
	"math/rand"
	"regexp"
	"strings"
)

type Config struct {
	MaxPoints    int
	MaxHintCount int
}

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
	Players          map[string]*Player
	Config           Config
}

func NewGame(conf Config) (*Game, error) {
	g := &Game{
		Config: conf,
	}
	return g, nil
}

func (G *Game) GetNewQuestion() error {
	G.CurrentQuestion = &Question{
		Category: "Testing",
		Prompt:   "Prompt",
		Answer:   "Answer",
	}
	return nil
}

func (G *Game) GetNextHint() error {
	question := G.CurrentQuestion
	hint := G.CurrentHint
	length := len(question.Answer)
	if hint.Count >= G.Config.MaxHintCount ||
		hint.Revealed > length/2 {
		return nil
	}
	if hint == nil {
		pat, err := regexp.Compile(`\w`)
		if err != nil {
			log.Fatalf("regex did not compile\n")
		}
		stars := string(pat.ReplaceAll([]byte(question.Prompt), []byte("*")))
		hint := &Hint{
			Stars: stars,
			Count: 1,
		}
		G.CurrentHint = hint
	} else {
		newHint := []string{}
		tokens := strings.Split(hint.Stars, " ")
		offset := 0
		for _, token := range tokens {
			chars := strings.Split(token, "")
			for {
				index := rand.Intn(len(chars))
				if chars[index] != "*" {
					chars[index] = "*"
					break
				}
			}
			newHint = append(newHint, strings.Join(chars, ""))
			offset += len(token) + 1
			hint.Revealed++
		}
		hint.Count++
		hint.Stars = strings.Join(newHint, " ")
	}
	return nil
}

func (G *Game) ResetGuesses() {
	for _, player := range G.Players {
		player.Guesses = 0
	}
}

func (G *Game) StartRound() error {
	G.PastQuestions = append(G.PastQuestions, G.CurrentQuestion)
	err := G.GetNewQuestion()
	if err != nil {
		log.Fatalf("error: could not get new question (%+v)\n", err)
	}
	G.ResetGuesses()
	return err
}

func (G *Game) MakeGuess(guess string, pid string) (bool, bool) {
	if rawString(guess) == rawString(G.CurrentQuestion.Answer) {
		player := G.Players[pid]
		awardedPoints := math.Max(float64(G.Config.MaxPoints-player.Guesses), 0)
		player.Score += int(awardedPoints)
		player.Streak++
		streakChange := G.PlayerWithStreak == pid
		G.PlayerWithStreak = pid
		return true, streakChange
	} else {
		player := G.Players[pid]
		player.Guesses++
		return false, false
	}
}

func rawString(str string) string {
	return strings.Replace(strings.TrimSpace(strings.ToLower(str)), "'", "", -1)
}
