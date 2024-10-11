package service

import (
	"telegram_reminder_bot/models"
	"telegram_reminder_bot/repository"
)

type TaskService struct {
	repo repository.Task
}

func NewTaskService(repo repository.Task) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(task models.Task) error {
	return s.repo.CreateTask(task)
}

func (s *TaskService) Tasks() ([]models.Task, error) {
	return s.repo.Tasks()
}
