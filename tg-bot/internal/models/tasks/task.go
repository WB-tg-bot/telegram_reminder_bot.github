package tasks

import "time"

type Task interface {
	GetChatID() int64
	GetUserName() string
	GetContent() string
	GetReminderTime() time.Time
}

type TaskImpl struct {
	ChatID       int64     `json:"chat_id"`
	UserName     string    `json:"username"`
	Content      string    `json:"content"`
	ReminderTime time.Time `json:"reminder_time"`
}

func NewTask(chatID int64, userName string, content string, reminderTime time.Time) *TaskImpl {
	return &TaskImpl{
		ChatID:       chatID,
		UserName:     userName,
		Content:      content,
		ReminderTime: reminderTime,
	}
}

func (t *TaskImpl) GetChatID() int64 {
	return t.ChatID
}

func (t *TaskImpl) GetUserName() string {
	return t.UserName
}

func (t *TaskImpl) GetContent() string {
	return t.Content
}

func (t *TaskImpl) GetReminderTime() time.Time {
	return t.ReminderTime
}
