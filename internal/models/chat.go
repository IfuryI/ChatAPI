package models

import "time"

// Chat структура чата
type Chat struct {
	ID        string    `json:"chat_id,omitempty"`
	Name      string    `json:"name"`
	ChatUsers []string  `json:"users"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
