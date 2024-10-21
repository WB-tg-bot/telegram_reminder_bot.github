package main

import (
	"context"
	"log"
	"os/signal"
	"tg-bot/models/bot"

	"os"
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

var (
	msgs  = make(map[int64]*tgbotapi.Message)
	flags = make(map[int64]bool)
	re    = regexp.MustCompile(`@\w+ ctrl (\d+)([a-z])`)
)

func main() {
	// Загрузка переменных окружения из файла .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Получение токена бота из переменной окружения
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not set in .env file")
	}

	tgBot := bot.NewBot(botToken)

	tgBot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := tgBot.GetUpdatesChan(u)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	tgBot.RestoreTasks()

	go receiveUpdates(ctx, tgBot, updates)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop
	cancel()
	log.Println("Bot stopped")
}

func receiveUpdates(ctx context.Context, tgBot *bot.Bot, updates tgbotapi.UpdatesChannel) {
	for {
		select {
		case <-ctx.Done():
			return
		case update := <-updates:
			handleUpdate(tgBot, update)
		}
	}
}

func handleUpdate(tgBot *bot.Bot, update tgbotapi.Update) {
	switch {
	case update.Message != nil:
		if update.Message.Text == "Добавить напоминание" {
			tgBot.DeleteMessage(update.Message)
			go tgBot.CreateReminder(update.Message)
			flags[update.Message.From.ID] = true

		} else if re.Match([]byte(update.Message.Text)) {
			go tgBot.HandleCommand(update.Message, msgs[update.Message.From.ID])

		} else if !flags[update.Message.From.ID] {
			msgs[update.Message.From.ID] = update.Message

		} else {
			flags[update.Message.From.ID] = tgBot.UpdateReminder(update.Message)
		}

	case update.CallbackQuery != nil:
		tgBot.HandleCallbackQuery(update.CallbackQuery)
		flags[update.CallbackQuery.From.ID] = false

	case update.MyChatMember != nil:
		tgBot.HandleMyChatMemberUpdate(update.MyChatMember)
		return

	default:
		return
	}
}
