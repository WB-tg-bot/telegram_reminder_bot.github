package models

import "time"

type Task struct {
	UserID       int64     `json:"user_id"`
	ChatID       int64     `json:"chat_id"`
	Content      string    `json:"content"`
	ReminderTime time.Time `json:"reminder_time"`
}

type User struct {
	ID       int64  `json:"user_id"`
	ChatID   int64  `json:"chat_id"`
	Username string `json:"username"`
}

type Chat struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
