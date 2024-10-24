package main

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"tg-bot/internal/handlers"
	"tg-bot/internal/models/bot"
	"tg-bot/internal/services/bot_service"
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
	botService := bot_service.NewBotService(tgBot)

	tgBot.(*bot.BotImpl).Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	handler := handlers.NewHandler()

	updates := tgBot.GetUpdatesChan(u)

	ctx, cancel := initContext()

	botService.RestoreTasks()

	go receiveUpdates(ctx, botService, handler, updates)

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

func receiveUpdates(ctx context.Context, botService *bot_service.BotService, handler handlers.Handler, updates tgbotapi.UpdatesChannel) {
	for {
		select {
		case <-ctx.Done():
			return
		case update := <-updates:
			go handler.HandleUpdate(botService, update)
		}
	}
}
