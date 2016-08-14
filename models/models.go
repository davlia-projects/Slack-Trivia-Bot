package models

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
