package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"tg-bot/models/tasks"
	"time"

	"github.com/robfig/cron/v3"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	*tgbotapi.BotAPI
}

func NewBot(token string) *Bot {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	return &Bot{bot}
}

func (b *Bot) HandleCommand(message *tgbotapi.Message, taskText string) {
	args := strings.Split(message.Text, " ")
	if len(args) < 3 || args[1] != "ctrl" {
		return
	}

	interval := args[2]
	duration := interval[len(interval)-1:]
	valueStr := interval[:len(interval)-1]

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Error parsing interval value: %v", err)
		return
	}

	var durationTime time.Duration
	switch duration {
	case "s":
		durationTime = time.Second * time.Duration(value)
	case "h":
		durationTime = time.Hour * time.Duration(value)
	case "d":
		durationTime = time.Hour * 24 * time.Duration(value)
	case "w":
		durationTime = time.Hour * 24 * 7 * time.Duration(value)
	case "m":
		durationTime = time.Hour * 24 * 30 * time.Duration(value)
	default:
		msg := tgbotapi.NewMessage(message.Chat.ID, "Неверный формат времени! Используйте только s, h, d, w, m")
		_, err := b.Send(msg)
		if err != nil {
			log.Println("Error sending message:", err)
		}
		return
	}

	if taskText == "" {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Нет текста задачи!")
		_, err := b.Send(msg)
		if err != nil {
			log.Println("Error sending message:", err)
		}
		return
	}

	inputTask := tasks.Task{
		ChatID:       message.Chat.ID,
		UserName:     message.From.UserName,
		Content:      taskText,
		ReminderTime: time.Now().Add(durationTime),
	}

	taskJSON, err := json.Marshal(inputTask)
	if err != nil {
		log.Println("Error marshaling task: ", err)
		return
	}

	resp, err := http.Post("http://telegram-reminder-bot:8000/create-task", "application/json", bytes.NewBuffer(taskJSON))
	if err != nil {
		log.Println("Error sending task to server: ", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Server error response: ", resp.Status)
		return
	}

	deleteConfig := tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID)
	_, err = b.Request(deleteConfig)
	if err != nil {
		log.Printf("Failed to delete message: %v", err)
	} else {
		log.Printf("Successfully deleted message %d in chat %d", message.MessageID, message.Chat.ID)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("#Задача# принята. Напомню о ней через %d%s", value, duration))
	_, err = b.Send(msg)
	if err != nil {
		log.Println("Error sending confirmation message:", err)
	}
}

func (b *Bot) RestoreTasks() {
	c := cron.New()

	_, err := c.AddFunc("@every 1s", func() {
		resp, err := http.Get("http://telegram-reminder-bot:8000/tasks")
		if err != nil {
			log.Println("Error fetching tasks from server: ", err)
			return
		}
		defer resp.Body.Close()

		var outputTasks []tasks.Task
		err = json.NewDecoder(resp.Body).Decode(&outputTasks)
		if err != nil {
			log.Println("Error decoding tasks: ", err)
			return
		}

		for _, task := range outputTasks {
			_, err = b.Send(tgbotapi.NewMessage(task.ChatID, fmt.Sprintf("Напоминание для @%s:\n\n%s", task.UserName, task.Content)))
			if err != nil {
				log.Println("Error sending reminder:", err)
			}
		}
	})
	if err != nil {
		log.Println("Error scheduling cron job:", err)
	}

	c.Start()
}

func (b *Bot) HandleMyChatMemberUpdate(myChatMember *tgbotapi.ChatMemberUpdated) {
	if myChatMember.NewChatMember.User.UserName == b.Self.UserName {
		switch myChatMember.NewChatMember.Status {
		case "member":
			messageText := fmt.Sprintf("Привет!"+
				"\nЯ @%s — бот планировщик для ваших задач."+
				"\n\n"+
				"Что я могу?\n\nКоманда @%s ctrl [число][время]:"+
				"\n\nВаше предыдущее сообщение сохраняется как #Задача#."+
				"\nЯ напомню вам о ней через указанное время."+
				"\n\n[число] - интервал (целое число)"+
				"\n[время] - продолжительность \n"+
				"\n• s - секунды, \n• h -часы, \n• d - дни, \n• w - недели, \n• m - месяцы",
				b.Self.UserName, b.Self.UserName)
			msg := tgbotapi.NewMessage(myChatMember.Chat.ID, messageText)
			_, err := b.Send(msg)
			if err != nil {
				log.Println("Error sending welcome message:", err)
			}
		default:
			// Нужно что-то тут))
		}
	}
}
