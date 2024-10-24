package bot_service

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/robfig/cron/v3"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tg-bot/internal/models/bot"
	"tg-bot/internal/models/reminder"
	"tg-bot/internal/models/tasks"
	"tg-bot/internal/services"
	"tg-bot/internal/transport"
	"tg-bot/internal/utils"
)

var (
	reminders  = make(map[int64]reminder.Reminder)
	botMessage = make(map[int64]tgbotapi.Message)
)

type BotService struct {
	Bot bot.Bot
}

func NewBotService(bot bot.Bot) *BotService {
	return &BotService{bot}
}

func (b *BotService) CreateReminder(msg *tgbotapi.Message) {
	rmdr, exists := reminders[msg.From.ID]
	if !exists {
		reminders[msg.From.ID] = reminder.NewReminder(msg.From.ID)
	} else {
		b.deleteBotMessage(rmdr.GetUserID())
	}

	botMessageConfig := tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("@%s, введите текст вашего напоминания", msg.From.UserName))
	botMessage[msg.From.ID], _ = b.Bot.Send(botMessageConfig)
}

func (b *BotService) UpdateReminder(msg *tgbotapi.Message) bool {
	rmdr, exists := reminders[msg.From.ID]
	if !exists {
		return false
	}

	if msg.From.ID != rmdr.GetUserID() {
		return false
	}

	if rmdr.GetTask() == nil || strings.TrimSpace(rmdr.GetTask().Text) == "" {
		b.deleteBotMessage(rmdr.GetUserID())

		rmdr.SetTask(msg)

		if strings.TrimSpace(rmdr.GetTask().Text) == "" {
			botMessageConfig := tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("@%s, что-то пошло не так, пожалуйста попробуйте еще раз.\nвведите текст вашего напоминания", msg.From.UserName))
			botMessage[rmdr.GetUserID()], _ = b.Bot.Send(botMessageConfig)

			utils.DeleteMessage(b.Bot.(*bot.BotImpl).BotAPI, msg)
			return true
		}

		botMessageConfig := tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("@%s, введите интервал ожидания\n(целое число)", msg.From.UserName))
		reminders[msg.From.ID] = rmdr
		botMessage[rmdr.GetUserID()], _ = b.Bot.Send(botMessageConfig)

		utils.DeleteMessage(b.Bot.(*bot.BotImpl).BotAPI, msg)
		return true

	} else if rmdr.GetInterval() == "" {
		b.deleteBotMessage(rmdr.GetUserID())

		rmdr.SetInterval(msg.Text)

		num, err := strconv.Atoi(rmdr.GetInterval())
		if err != nil || num < 0 {
			botMessageConfig := tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("@%s, что-то пошло не так, пожалуйста попробуйте еще раз.\nвведите интервал ожидания\n(целое число)", msg.From.UserName))
			reminders[msg.From.ID] = rmdr
			botMessage[rmdr.GetUserID()], err = b.Bot.Send(botMessageConfig)
			if err != nil {
				log.Printf("Failed to send message: %v", err)
			}
			if msg.MessageID != 0 {
				utils.DeleteMessage(b.Bot.(*bot.BotImpl).BotAPI, msg)
			}
			rmdr.SetInterval("")
			reminders[msg.From.ID] = rmdr
			return true
		}

		botMessageConfig := tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("@%s, выберите продолжительность:", msg.From.UserName))
		botMessageConfig.ReplyMarkup = services.TimeKeyboard
		botMessage[rmdr.GetUserID()], err = b.Bot.Send(botMessageConfig)
		if err != nil {
			log.Printf("Failed to send message: %v", err)
		}
		if msg.MessageID != 0 {
			utils.DeleteMessage(b.Bot.(*bot.BotImpl).BotAPI, msg)
		}
		reminders[msg.From.ID] = rmdr
	}
	return true
}

func (b *BotService) HandleCallbackQuery(callback *tgbotapi.CallbackQuery) {
	rmdr, exists := reminders[callback.From.ID]
	if !exists || rmdr.GetTask() == nil {
		return
	}

	b.deleteBotMessage(rmdr.GetUserID())

	rmdr.SetDuration(callback.Data)
	command := fmt.Sprintf("@%s ctrl %s%s", b.Bot.(*bot.BotImpl).Self.UserName, rmdr.GetInterval(), rmdr.GetDuration())
	input := tgbotapi.Message{
		Chat: callback.Message.Chat,
		From: callback.From,
		Text: command,
	}

	b.HandleCommand(&input, rmdr.GetTask())
	reminders[callback.From.ID] = reminder.NewReminder(callback.From.ID)
}

func (b *BotService) HandleCommand(message *tgbotapi.Message, task *tgbotapi.Message) {
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
		_, err := b.Bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
		return
	}

	if task == nil || task.Chat.ID != message.Chat.ID {
		utils.DeleteMessage(b.Bot.(*bot.BotImpl).BotAPI, message)

		msg := tgbotapi.NewMessage(message.Chat.ID, "Нет текста задачи!")
		_, err := b.Bot.Send(msg)
		if err != nil {
			log.Println(err)
		}

		return
	}

	inputTask := tasks.NewTask(message.Chat.ID, message.From.UserName, task.Text, time.Now().Add(durationTime))

	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	resp, err := transport.PostJSON(viper.GetString("url_create_task"), inputTask)
	if err != nil {
		log.Println("Error sending user to server: ", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Error saving user: ", resp.Status)
		return
	}

	utils.DeleteMessage(b.Bot.(*bot.BotImpl).BotAPI, message)
	utils.DeleteMessage(b.Bot.(*bot.BotImpl).BotAPI, task)

	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("@%s, #Задача# принята. Напомню о ней через %d%s", message.From.UserName, value, duration))
	_, err = b.Bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func (b *BotService) deleteBotMessage(userID int64) {
	if idBotMsg, ok := botMessage[userID]; ok {
		utils.DeleteMessage(b.Bot.(*bot.BotImpl).BotAPI, &idBotMsg)
	}
}

func (b *BotService) RestoreTasks() {
	c := cron.New()

	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	_, err := c.AddFunc("@every 1s", func() {
		resp, err := transport.GetJSON(viper.GetString("url_get_task"))
		if err != nil {
			log.Println("Error getting tasks from server: ", err)
			return
		}
		defer resp.Body.Close()

		var outputTasks []tasks.TaskImpl
		err = json.NewDecoder(resp.Body).Decode(&outputTasks)
		if err != nil {
			log.Println("Error decoding tasks: ", err)
			return
		}

		for _, task := range outputTasks {
			_, err = b.Bot.Send(tgbotapi.NewMessage(task.GetChatID(), fmt.Sprintf("Напоминание для @%s:\n\n%s", task.GetUserName(), task.GetContent())))
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

func (b *BotService) HandleMyChatMemberUpdate(myChatMember *tgbotapi.ChatMemberUpdated) {
	if myChatMember.NewChatMember.User.ID == b.Bot.(*bot.BotImpl).Self.ID {
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
				b.Bot.(*bot.BotImpl).Self.UserName, b.Bot.(*bot.BotImpl).Self.UserName)

			msg := tgbotapi.NewMessage(myChatMember.Chat.ID, messageText)
			_, err := b.Bot.Send(msg)
			if err != nil {
				log.Println(err)
			}

			msg = tgbotapi.NewMessage(myChatMember.Chat.ID, "Пожалуйста, выберите опцию:")
			msg.ReplyMarkup = services.Menu
			_, err = b.Bot.Send(msg)
			if err != nil {
				log.Println(err)
			}

		default:
			return
		}
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
