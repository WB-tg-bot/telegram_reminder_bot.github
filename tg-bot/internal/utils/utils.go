package utils

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func DeleteMessage(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	deleteConfig := tgbotapi.NewDeleteMessage(msg.Chat.ID, msg.MessageID)
	_, err := bot.Request(deleteConfig)
	if err != nil {
		log.Printf("Failed to delete message: %v", err)
	} else {
		log.Printf("Successfully deleted message %d in chat %d", msg.MessageID, msg.Chat.ID)
	}
}
