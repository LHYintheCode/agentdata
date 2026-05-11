package model

import "time"

type Session struct {
	ID       string    `json:"id"`
	Source   string    `json:"source"`
	Project  string    `json:"project"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role      string    `json:"role"`
	Text      string    `json:"text"`
	Timestamp time.Time `json:"timestamp"`
}
