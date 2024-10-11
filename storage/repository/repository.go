package repository

import (
	"telegram_reminder_bot/models"

	"github.com/jmoiron/sqlx"
)

type Task interface {
	CreateTask(models.Task) error
	Tasks() ([]models.Task, error)
}

type Repository struct {
	Task
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Task: NewTaskPostgres(db),
	}
}
