package repository

import (
	"github.com/jmoiron/sqlx"
)

type Task interface {
	Create(chatID int64, task string, interval int, unit string, username string) error
	//GetAll() []Tasks
}

type Repository struct {
	Task
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Task: NewTaskPostrgres(db),
	}
}
