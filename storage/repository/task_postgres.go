package repository

import (
	"log"
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
	_, err := r.db.Exec("INSERT INTO tasks (username, chat_id, content, reminder_time) VALUES ($1, $2, $3, $4)", task.UserName, task.ChatID, task.Content, task.ReminderTime)
	if err != nil {
		return err
	}

	log.Printf("Created task with username %s and chatID %d", task.UserName, task.ChatID)
	return nil
}

func (r *TaskPostgres) Tasks() ([]models.Task, error) {
	rows, err := r.db.Queryx("SELECT username, chat_id, content, reminder_time FROM tasks WHERE reminder_time <= NOW()")
	if err != nil {
		return nil, err
	}

	var tasks []models.Task

	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.UserName, &task.ChatID, &task.Content, &task.ReminderTime)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	log.Printf("Found %d tasks", len(tasks))

	if err = r.deleteTasks(tasks); err != nil {
		return nil, err
	}

	log.Printf("Deleted %d tasks", len(tasks))
	return tasks, nil
}

func (r *TaskPostgres) deleteTasks(tasks []models.Task) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`DELETE FROM tasks WHERE username = $1 AND chat_id = $2 AND content = $3 AND reminder_time = $4`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, task := range tasks {
		_, err := stmt.Exec(task.UserName, task.ChatID, task.Content, task.ReminderTime)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
