package model

import "time"

type Article struct {
	ID       string    `json:"id"`
	AuthorID string    `json:"author_id"`
	Time     time.Time `json:"time"`
	Read     int       `json:"read"`
	Like     int       `json:"like"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	Status   string    `json:"status"`
}
