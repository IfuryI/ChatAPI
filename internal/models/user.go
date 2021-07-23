package models

import "time"

// User структура пользователя
type User struct {
	ID        string    `json:"user_id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
