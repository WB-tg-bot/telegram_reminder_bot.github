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

	ctx, cancel := initContext()

	tgBot.RestoreTasks()

	go receiveUpdates(ctx, tgBot, updates)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop
	cancel()
	log.Println("Bot stopped")
}

func initContext() (context.Context, context.CancelFunc) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	return ctx, cancel
}

func receiveUpdates(ctx context.Context, tgBot *bot.Bot, updates tgbotapi.UpdatesChannel) {
	for {
		select {
		case <-ctx.Done():
			return
		case update := <-updates:
			go handleUpdate(tgBot, update)
		}
	}
}

func handleUpdate(tgBot *bot.Bot, update tgbotapi.Update) {
	switch {
	case update.Message != nil:
		handleMessage(tgBot, update.Message)
	case update.CallbackQuery != nil:
		handleCallbackQuery(tgBot, update.CallbackQuery)
	case update.MyChatMember != nil:
		handleMyChatMemberUpdate(tgBot, update.MyChatMember)
	case update.EditedMessage != nil:
		handleEditedMessage(tgBot, update.EditedMessage)
	default:
		return
	}
}

func handleMessage(tgBot *bot.Bot, message *tgbotapi.Message) {
	if message.Text == "Добавить напоминание" {
		tgBot.DeleteMessage(message)
		go tgBot.CreateReminder(message)
		flags[message.From.ID] = true
	} else if re.Match([]byte(message.Text)) {
		go tgBot.HandleCommand(message, msgs[message.From.ID])
	} else if !flags[message.From.ID] {
		msgs[message.From.ID] = message
	} else {
		flags[message.From.ID] = tgBot.UpdateReminder(message)
	}
}

func handleCallbackQuery(tgBot *bot.Bot, callback *tgbotapi.CallbackQuery) {
	tgBot.HandleCallbackQuery(callback)
	flags[callback.From.ID] = false
}

func handleMyChatMemberUpdate(tgBot *bot.Bot, myChatMember *tgbotapi.ChatMemberUpdated) {
	tgBot.HandleMyChatMemberUpdate(myChatMember)
}

func handleEditedMessage(tgBot *bot.Bot, editedMessage *tgbotapi.Message) {
	msgs[editedMessage.From.ID] = editedMessage
}
