package tasks

import "time"

type Task struct {
	ChatID       int64     `json:"chat_id"`
	UserName     string    `json:"username"`
	Content      string    `json:"content"`
	ReminderTime time.Time `json:"reminder_time"`
}
