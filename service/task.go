package service

import "telegram_reminder_bot/repository"

type TaskService struct {
	repo repository.Task
}

func NewTaskService(repo repository.Task) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) Create(chatID int64, task string, interval int, unit string, username string) error {
	return s.repo.Create(chatID, task, interval, unit, username)
}
