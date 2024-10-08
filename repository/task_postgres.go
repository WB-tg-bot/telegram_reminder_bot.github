package repository

import "github.com/jmoiron/sqlx"

type TaskPostgres struct {
	db *sqlx.DB
}

func NewTaskPostrgres(db *sqlx.DB) *TaskPostgres {
	return &TaskPostgres{db: db}
}

func (r *TaskPostgres) Create(chatID int64, task string, interval int, unit string, username string) error {
	return nil
}

func (r *TaskPostgres) Get(chatID int64, task string) (int64, error) {
	return 0, nil
}
