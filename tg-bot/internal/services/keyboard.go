package services

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var (
	Menu = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Добавить напоминание")))

	TimeKeyboard = tgbotapi.NewInlineKeyboardMarkup(
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
