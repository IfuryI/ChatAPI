package models

import "time"

// Message структура сообщения
type Message struct {
	ID        string    `json:"message_id,omitempty"`
	ChatID    string    `json:"chat"`
	UserID    string    `json:"author"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
