package model

import "time"

type Log struct {
	ID      string    `json:"id"`
	Level   string    `json:"level"`
	Time    time.Time `json:"time"`
	Message string    `json:"message"`
	Source  string    `json:"source"`
}
