package config

import "time"

type Config struct {
	MaxPoints    int
	MaxHintCount int
	HintDelay    time.Duration
	QuestionTime time.Duration
}
