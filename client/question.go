package client

import (
	"sync"

	. "github.com/dota-2-slack-bot/models"
)

var (
	questionOnce   sync.Once
	questionClient *QuestionClient
)

type QuestionClient struct {
}

func GetQuestionClient() *QuestionClient {
	questionOnce.Do(func() {
		questionClient = &QuestionClient{}
	})
	return questionClient
}

func (Q *QuestionClient) NewQuestion() *Question {
	q := &Question{
		Prompt:   "Prompt",
		Answer:   "Answer",
		Category: "Category",
	}
	return q
}
