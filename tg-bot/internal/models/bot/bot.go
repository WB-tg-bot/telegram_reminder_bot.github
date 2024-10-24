package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
	Request(c tgbotapi.Chattable) (*tgbotapi.APIResponse, error)
	GetUpdatesChan(u tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel
	GetMe() (tgbotapi.User, error)
}

type BotImpl struct {
	*tgbotapi.BotAPI
}

func NewBot(token string) Bot {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	return &BotImpl{bot}
}

func (b *BotImpl) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	return b.BotAPI.Send(c)
}

func (b *BotImpl) Request(c tgbotapi.Chattable) (*tgbotapi.APIResponse, error) {
	return b.BotAPI.Request(c)
}

func (b *BotImpl) GetUpdatesChan(u tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel {
	return b.BotAPI.GetUpdatesChan(u)
}

func (b *BotImpl) GetMe() (tgbotapi.User, error) {
	return b.BotAPI.GetMe()
}
