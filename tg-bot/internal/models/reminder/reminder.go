package reminder

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Reminder interface {
	GetUserID() int64
	GetTask() *tgbotapi.Message
	GetInterval() string
	GetDuration() string
	SetTask(task *tgbotapi.Message)
	SetInterval(interval string)
	SetDuration(duration string)
}

type ReminderImpl struct {
	UserID   int64
	Task     *tgbotapi.Message
	Interval string
	Duration string
}

func NewReminder(id int64) Reminder {
	return &ReminderImpl{
		UserID: id,
	}
}

func (r *ReminderImpl) GetUserID() int64 {
	return r.UserID
}

func (r *ReminderImpl) GetTask() *tgbotapi.Message {
	return r.Task
}

func (r *ReminderImpl) GetInterval() string {
	return r.Interval
}

func (r *ReminderImpl) GetDuration() string {
	return r.Duration
}

func (r *ReminderImpl) SetTask(task *tgbotapi.Message) {
	r.Task = task
}

func (r *ReminderImpl) SetInterval(interval string) {
	r.Interval = interval
}

func (r *ReminderImpl) SetDuration(duration string) {
	r.Duration = duration
}
