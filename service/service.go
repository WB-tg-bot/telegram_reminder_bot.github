package service

import "telegram_reminder_bot/repository"

type Task interface {
	Create(chatID int64, task string, interval int, unit string, username string) error
}

type Service struct {
	Task
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		NewTaskService(repo.Task),
	}
}
