package service

import (
	"telegram_reminder_bot/models"
	"telegram_reminder_bot/repository"
)

type Task interface {
	CreateTask(models.Task) error
	Tasks() ([]models.Task, error)
}

type Service struct {
	Task
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		NewTaskService(repo.Task),
	}
}
