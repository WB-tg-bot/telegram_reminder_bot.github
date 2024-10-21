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

var botMessage tgbotapi.Message

var (
	menu = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Добавить напоминание")))

	timeKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("секунды", "s"),
			tgbotapi.NewInlineKeyboardButtonData("часы", "h"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("дни", "d"),
			tgbotapi.NewInlineKeyboardButtonData("недели", "w"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("месяцы", "m"),
		),
	)
)

type Bot struct {
	*tgbotapi.BotAPI
}

type Reminder struct {
	Task     *tgbotapi.Message
	Interval string
	Duration string
}

func NewReminder() *Reminder {
	return &Reminder{}
}

func NewBot(token string) *Bot {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	return &Bot{bot}
}

var reminders = make(map[int64]Reminder)

func (b *Bot) CreateReminder(msg *tgbotapi.Message) {
	reminders[msg.From.ID] = *NewReminder()
	botMessageConfig := tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("@%s, введите текст вашего напоминания", msg.From.UserName))

	botMessage, _ = b.Send(botMessageConfig)
}

func (b *Bot) UpdateReminder(msg *tgbotapi.Message) bool {
	reminder, exists := reminders[msg.From.ID]
	if !exists {
		return false
	}

	if reminder.Task == nil {
		b.DeleteMessage(&botMessage)

		reminder.Task = msg
		botMessageConfig := tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("@%s, введите интервал ожидания\n(целое число)", msg.From.UserName))
		reminders[msg.From.ID] = reminder

		botMessage, _ = b.Send(botMessageConfig)

		b.DeleteMessage(msg)

		return true

	} else if reminder.Interval == "" {
		b.DeleteMessage(&botMessage)

		reminder.Interval = msg.Text
		botMessageConfig := tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("@%s, выберите продолжительность:", msg.From.UserName))
		botMessageConfig.ReplyMarkup = timeKeyboard

		botMessage, _ = b.Send(botMessageConfig)
		b.DeleteMessage(msg)

		reminders[msg.From.ID] = reminder
	}
	return true
}

func (b *Bot) HandleCallbackQuery(callback *tgbotapi.CallbackQuery) {
	b.DeleteMessage(&botMessage)

	reminder, exists := reminders[callback.From.ID]
	if !exists {
		return
	}

	reminder.Duration = callback.Data
	command := fmt.Sprintf("@%s ctrl %s%s", b.Self.UserName, reminder.Interval, reminder.Duration)
	input := tgbotapi.Message{
		Chat: callback.Message.Chat,
		From: callback.From,
		Text: command,
	}

	b.HandleCommand(&input, reminder.Task)
	reminders[callback.From.ID] = Reminder{
		Task:     nil,
		Interval: "",
		Duration: "",
	}
}

func (b *Bot) HandleCommand(message *tgbotapi.Message, task *tgbotapi.Message) {
	args := strings.Split(message.Text, " ")
	if len(args) < 3 || args[1] != "ctrl" {
		return
	}

	interval := args[2]
	duration := interval[len(interval)-1:]
	valueStr := interval[:len(interval)-1]

	value, err := strconv.Atoi(valueStr)
	if err != nil {
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
			log.Println(err)
		}
		return
	}

	if task == nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Нет текста задачи!")
		_, err := b.Send(msg)
		if err != nil {
			log.Println(err)
		}
		return
	}

	inputTask := tasks.Task{
		ChatID:       message.Chat.ID,
		UserName:     message.From.UserName,
		Content:      task.Text,
		ReminderTime: time.Now().Add(durationTime),
	}

	taskJSON, err := json.Marshal(inputTask)
	if err != nil {
		log.Println("Error marshaling user: ", err)
		return
	}

	resp, err := http.Post("http://telegram-reminder-bot:8000/create-task", "application/json", bytes.NewBuffer(taskJSON))
	//resp, err := http.Post("http://localhost:8000/create-task", "application/json", bytes.NewBuffer(taskJSON))
	if err != nil {
		log.Println("Error sending user to server: ", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Error saving user: ", resp.Status)
		return
	}

	b.DeleteMessage(message)
	b.DeleteMessage(task)

	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("@%s, #Задача# принята. Напомню о ней через %d%s", message.From.UserName, value, duration))
	_, err = b.Send(msg)
	if err != nil {
		log.Println(err)
	}

}

func (b *Bot) DeleteMessage(msg *tgbotapi.Message) {
	deleteConfig := tgbotapi.NewDeleteMessage(msg.Chat.ID, msg.MessageID)
	_, err := b.Request(deleteConfig)
	if err != nil {
		log.Printf("Failed to delete message: %v", err)
	} else {
		log.Printf("Successfully deleted message %d in chat %d", msg.MessageID, msg.Chat.ID)
	}
}

func (b *Bot) RestoreTasks() {
	c := cron.New()

	_, err := c.AddFunc("@every 1s", func() {

		resp, err := http.Get("http://telegram-reminder-bot:8000/tasks")
		//resp, err := http.Get("http://localhost:8000/tasks")
		if err != nil {
			log.Println("Error getting tasks from server: ", err)
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
				log.Println(err)
			}
		}
	})
	if err != nil {
		log.Println(err)
	}

	c.Start()

}

func (b *Bot) HandleMyChatMemberUpdate(myChatMember *tgbotapi.ChatMemberUpdated) {
	if myChatMember.NewChatMember.User.ID == b.Self.ID {
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
				"\n• s - секунды, \n• h - часы, \n• d - дни, \n• w - недели, \n• m - месяцы",
				b.Self.UserName, b.Self.UserName)

			msg := tgbotapi.NewMessage(myChatMember.Chat.ID, messageText)
			_, err := b.Send(msg)
			if err != nil {
				log.Println(err)
			}

			msg = tgbotapi.NewMessage(myChatMember.Chat.ID, "Пожалуйста, выберите опцию:")
			msg.ReplyMarkup = menu
			_, err = b.Send(msg)
			if err != nil {
				log.Println(err)
			}

		default:
			return
		}
	}
}
