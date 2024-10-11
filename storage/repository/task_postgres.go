package repository

import (
	"telegram_reminder_bot/models"

	"github.com/jmoiron/sqlx"
)

type TaskPostgres struct {
	db *sqlx.DB
}

func NewTaskPostgres(db *sqlx.DB) *TaskPostgres {
	return &TaskPostgres{db: db}
}

func (r *TaskPostgres) CreateTask(task models.Task) error {
	_, err := r.db.Exec("INSERT INTO tasks (user_id, chat_id, content, reminder_time) VALUES ($1, $2, $3, $4)", task.UserID, task.ChatID, task.Content, task.ReminderTime)
	if err != nil {
		return err
	}

	return nil
}

func (r *TaskPostgres) Tasks() ([]models.Task, error) {
	rows, err := r.db.Queryx("SELECT user_id, chat_id, content, reminder_time FROM tasks WHERE reminder_time <= NOW()")
	if err != nil {
		return nil, err
	}

	var tasks []models.Task

	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.UserID, &task.ChatID, &task.Content, &task.ReminderTime)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err = r.deleteTasks(tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *TaskPostgres) deleteTasks(tasks []models.Task) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`DELETE FROM tasks WHERE user_id = $1 AND chat_id = $2 AND content = $3 AND reminder_time = $4`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, task := range tasks {
		_, err := stmt.Exec(task.UserID, task.ChatID, task.Content, task.ReminderTime)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
