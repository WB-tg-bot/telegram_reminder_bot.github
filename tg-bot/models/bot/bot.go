package bot

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

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

	if taskText == "" {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Нет текста задачи!")
		_, err := b.Send(msg)
		if err != nil {
			log.Println(err)
		}
		return
	}

	task := fmt.Sprintf("@"+message.From.UserName+" #Задача#: %s", taskText)
	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("#Задача# принята. Напомню о ней через %d%s", value, duration))
	_, err = b.Send(msg)
	if err != nil {
		log.Println(err)
	}

	time.AfterFunc(durationTime, func() {
		reminder := tgbotapi.NewMessage(message.Chat.ID, task)
		b.Send(reminder)
	})
}
