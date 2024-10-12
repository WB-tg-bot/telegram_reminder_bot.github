package main

import (
	"log"
	"tg-bot/models/bot"

	"os"
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
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

	re := regexp.MustCompile(`@\w+ ctrl (\d+)([a-z])`)

	msgs := make(map[int64]string)

	tgBot.RestoreTasks()

	for update := range updates {
		if update.MyChatMember != nil {
			tgBot.HandleMyChatMemberUpdate(update.MyChatMember)
			continue
		}

		if re.Match([]byte(update.Message.Text)) {
			go tgBot.HandleCommand(update.Message, msgs[update.Message.From.ID])
		} else {
			msgs[update.Message.From.ID] = update.Message.Text
		}

	}

	select {}

}
