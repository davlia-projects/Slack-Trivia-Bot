package models

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
