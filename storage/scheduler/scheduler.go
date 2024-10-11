package scheduler

import (
	"errors"
	"fmt"
	"time"

	"github.com/robfig/cron/v3" // Используем эту библиотеку для планировщиков задач
)

// Определяем структуру Job дл хранения информации о задаче
type Job struct {
	ChatID   int64  // ID чата, куда будет отправляться напоминание
	Task     string // Описание задачи
	Username string // Имя пользователя, которому принадлежит задача

}

// Глобальная переменная для планировщика
var cronScheduler *cron.Cron

// Инициализация планировщика
func InitScheduler() {
	// Создаем новый экземпляр планировщика
	cronScheduler = cron.New(cron.WithLocation(time.UTC))
	// Запускаем планировщик
	cronScheduler.Start()
}

// Фунция для планирования задачи
func ScheduleTask(chatID int64, task string, interval int, unit string, username string) {
	var cronExpr string // Переменная для хранения выражения cron

	// Определяем выражение cron в зависимости от единицы измерения
	switch unit {
	case "h": // Если указан час
		cronExpr = "@every" + time.Duration(interval).String() + "h"
	case "d": // Если указан день
		cronExpr = "@every" + time.Duration(interval).String() + "d"
	case "w": // Если указан неделя
		cronExpr = "@every" + time.Duration(interval).String() + "w"
	case "mo": // Если указан месяц
		cronExpr = "@every" + time.Duration(interval*30).String() + "d"
	default:
		//fmt.Errorf("Incorrect value", cronExpr) - нужно добавить логику обработки ошибки, я написал как пример, можете развить дальше))
		return // Если единица не распознана, выходим из функции

	}
	// Добавляем задачу в планировщик
	cronScheduler.AddFunc(cronExpr, func() {
		SendReminder(chatID, task, username) // При срабатывании задачи отправляем напоминание
	})

}

// Функция для отправки напоминания в Telegram
func SendReminder(chatID int64, task string, username string) {
	message := "Напоминание: @" + username + " " + task // Формируем сообщение
	// Здесь нужно реализовать логику отправки сообщения через Telegram API
	fmt.Println(message) // Это сделанно для того, что бы убрать ошибку))).
}

func StopScheduler() {
	// Тут нужно описать логику остановки планировщика.
}

// Функция для удаления задачи по ID
func DeleteTask(taskID string) error {
	// Логика для удаления задачи
	// Если задача не найдена, можно вернуть ошибку
	if taskID == "" {
		return errors.New("task ID cannot be empty")
	}
	// Здесь нужно добавить код для фактического удаления задачи.
	return nil // Возрващаем ошибку если удаление не удалось
}
