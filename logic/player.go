package logic

type Player struct {
	ID      string
	Name    string
	Score   int
	Streak  int
	Guesses int
}

func NewPlayer(id, name string) *Player {
	p := &Player{
		ID:   id,
		Name: name,
	}
	return p
}
