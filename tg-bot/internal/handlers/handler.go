package handlers

import (
	"regexp"
	"strings"
	"tg-bot/internal/models/bot"
	"tg-bot/internal/services/bot_service"
	"tg-bot/internal/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	msgs  = make(map[int64]*tgbotapi.Message)
	flags = make(map[int64]bool)
	re    = regexp.MustCompile(`@\w+ ctrl (\d+)([a-z])`)
)

type Handler interface {
	HandleUpdate(botService *bot_service.BotService, update tgbotapi.Update)
	HandleMessage(botService *bot_service.BotService, message *tgbotapi.Message)
	HandleCallbackQuery(botService *bot_service.BotService, callback *tgbotapi.CallbackQuery)
	HandleMyChatMemberUpdate(botService *bot_service.BotService, myChatMember *tgbotapi.ChatMemberUpdated)
	HandleEditedMessage(botService *bot_service.BotService, editedMessage *tgbotapi.Message)
}

type HandlerImpl struct {
	Handler
}

func NewHandler() Handler {
	return &HandlerImpl{}
}

func (h *HandlerImpl) HandleUpdate(botService *bot_service.BotService, update tgbotapi.Update) {
	switch {
	case update.Message != nil:
		h.HandleMessage(botService, update.Message)
	case update.CallbackQuery != nil:
		h.HandleCallbackQuery(botService, update.CallbackQuery)
	case update.MyChatMember != nil:
		h.HandleMyChatMemberUpdate(botService, update.MyChatMember)
	case update.EditedMessage != nil:
		h.HandleEditedMessage(botService, update.EditedMessage)
	default:
		return
	}
}

func (h *HandlerImpl) HandleMessage(botService *bot_service.BotService, message *tgbotapi.Message) {
	if message.Text == "Добавить напоминание" {
		utils.DeleteMessage(botService.Bot.(*bot.BotImpl).BotAPI, message)
		go botService.CreateReminder(message)
		flags[message.From.ID] = true
	} else if re.Match([]byte(message.Text)) {
		go botService.HandleCommand(message, msgs[message.From.ID])
	} else if !flags[message.From.ID] {
		if strings.TrimSpace(message.Text) != "" {
			msgs[message.From.ID] = message
		}
	} else {
		flags[message.From.ID] = botService.UpdateReminder(message)
	}
}

func (h *HandlerImpl) HandleCallbackQuery(botService *bot_service.BotService, callback *tgbotapi.CallbackQuery) {
	botService.HandleCallbackQuery(callback)
	flags[callback.From.ID] = false
}

func (h *HandlerImpl) HandleMyChatMemberUpdate(botService *bot_service.BotService, myChatMember *tgbotapi.ChatMemberUpdated) {
	botService.HandleMyChatMemberUpdate(myChatMember)
}

func (h *HandlerImpl) HandleEditedMessage(botService *bot_service.BotService, editedMessage *tgbotapi.Message) {
	msgs[editedMessage.From.ID] = editedMessage
}
