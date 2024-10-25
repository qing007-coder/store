package model

import "time"

type TokenStore struct {
	ClientID            string        `json:"ClientID"`
	UserID              string        `json:"UserID"`
	RedirectURI         string        `json:"RedirectURI"`
	Scope               string        `json:"Scope"`
	Code                string        `json:"Code"`
	CodeChallenge       string        `json:"CodeChallenge"`
	CodeChallengeMethod string        `json:"CodeChallengeMethod"`
	CodeCreateAt        time.Time     `json:"CodeCreateAt"`
	CodeExpiresIn       time.Duration `json:"CodeExpiresIn"`
	Access              string        `json:"Access"`
	AccessCreateAt      time.Time     `json:"AccessCreateAt"`
	AccessExpiresIn     time.Duration `json:"AccessExpiresIn"`
	Refresh             string        `json:"Refresh"`
	RefreshCreateAt     time.Time     `json:"RefreshCreateAt"`
	RefreshExpiresIn    time.Duration `json:"RefreshExpiresIn"`
}
