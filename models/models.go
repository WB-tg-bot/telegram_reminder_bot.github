package models

type Chats struct {
	Id   int    `json:"id" bd:"id"`
	Name string `json:"name" bd:"name"`
}
type Jobs struct {
	Id           int    `json:"id" bd:"id"`
	Task         string `json:"task" bd:"task"`
	ReminderTime string `json:"reminder_time" bd:"reminder_time"`
	Done         bool   `json:"done" db:"done"`
}

type JobsChats struct {
	Id     int
	JobID  int
	ChatID int
}

type Members struct {
	Id   int    `json:"id" bd:"id"`
	Name string `json:"name" bd:"name"`
}

type JabsMembers struct {
	Id       int
	JobID    int
	MemberID int
}

// пока придержем, еще разбираемся со структурой
var Json struct {
	ChatID   int64  `json:"chat_id"`  // ID чата
	Task     string `json:"task"`     // Описание задачи
	Interval int    `json:"interval"` // Интервал времени
	Unit     string `json:"unit"`     // Единица зимерения (час, день и т.д)
	Username string `json:"username"` // Имя пользователя
}
