package model

type Log struct {
	ID     string `json:"id"`
	Level  string `json:"level"`
	Time   int64  `json:"time"`
	Msg    string `json:"msg"`
	Source string `json:"source"`
}
