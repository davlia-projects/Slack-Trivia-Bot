package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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
	url := &url.URL{
		Scheme: "http",
		Host:   "localhost",
		Path:   "question",
	}
	resp, err := http.Get(url.String())
	if err != nil {
		fmt.Errorf("error: couldn't get question (%+v)\n", err)
		return nil
	}
	if resp.StatusCode >= 500 {
		fmt.Errorf("error: gateway error (%d)\n", resp.StatusCode)
		return nil
	}
	var question *Question
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(question); err != nil {
		fmt.Errorf("error: couldn't unmarshall question (%+v)\n", err)
		return nil
	}
	return question
}
