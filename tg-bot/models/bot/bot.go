package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"net/http"
	"strconv"
	"strings"
	"tg-bot/models/tasks"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

/*
var (
	menu = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Добавить напоминание")))

	timeKeyboard = tgbotapi.NewInlineKeyboardMarkup(
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
*/

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

/*
	func (b *Bot) MarkupHandler(chat tgbotapi.Chat, user tgbotapi.User, updates tgbotapi.UpdatesChannel) {
		msg := tgbotapi.NewMessage(chat.ID, fmt.Sprintf("@%s, введите текст вашего напоминания", user.UserName))
		_, err := b.Send(msg)
		if err != nil {
			log.Println(err)
		}

		var taskText string
		var time string
		var quantity string

		re := regexp.MustCompile(`[shdwm]`)
		reQuantity := regexp.MustCompile(`\d+`)

		for update := range updates {
			if update.Message.Text != "" && update.Message.From.ID == user.ID {
				if taskText == "" {
					taskText = update.Message.Text
					msg = tgbotapi.NewMessage(chat.ID, fmt.Sprintf("@%s, выберите единицу времени:", user.UserName))
					msg.ReplyMarkup = timeKeyboard
					_, err = b.Send(msg)
					if err != nil {
						log.Println(err)
					}
				} else if time == "" {
					if re.MatchString(update.Message.Text) {
						time = update.Message.Text
						msg = tgbotapi.NewMessage(chat.ID, fmt.Sprintf("@%s, введите продолжительность (целым числом)", user.UserName))
						_, err = b.Send(msg)
						if err != nil {
							log.Println(err)
						}
					}
				} else if quantity == "" {
					if reQuantity.MatchString(update.Message.Text) {
						quantity = update.Message.Text
						break
					}
				}
			}
		}

		inputText := fmt.Sprintf("@%s ctrl %s%s", b.Self.UserName, quantity, time)
		input := tgbotapi.Message{
			Chat: &chat,
			From: &user,
			Text: inputText,
		}

		b.HandleCommand(&input, taskText)
	}
*/
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

	inputTask := tasks.Task{
		ChatID:       message.Chat.ID,
		UserName:     message.From.UserName,
		Content:      taskText,
		ReminderTime: time.Now().Add(durationTime),
	}

	taskJSON, err := json.Marshal(inputTask)
	if err != nil {
		log.Println("Error marshaling user: ", err)
		return
	}

	resp, err := http.Post("http://telegram-reminder-bot:8000/create-task", "application/json", bytes.NewBuffer(taskJSON))
	// resp, err := http.Post("http://localhost:8000/create-task", "application/json", bytes.NewBuffer(taskJSON))
	if err != nil {
		log.Println("Error sending user to server: ", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Error saving user: ", resp.Status)
		return
	}

	deleteConfig := tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID)
	_, err = b.Request(deleteConfig)
	if err != nil {
		log.Printf("Failed to delete message: %v", err)
	} else {
		log.Printf("Successfully deleted message %d in chat %d", message.MessageID, message.Chat.ID)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("#Задача# принята. Напомню о ней через %d%s", value, duration))
	_, err = b.Send(msg)
	if err != nil {
		log.Println(err)
	}

}

func (b *Bot) RestoreTasks() {
	c := cron.New()

	_, err := c.AddFunc("@every 1s", func() {

		resp, err := http.Get("http://telegram-reminder-bot:8000/tasks")
		// resp, err := http.Get("http://localhost:8000/tasks")
		if err != nil {
			log.Println("Error getting tasks from server: ", err)
			return
		}
		defer resp.Body.Close()

		var outputTasks []tasks.Task
		err = json.NewDecoder(resp.Body).Decode(&outputTasks)
		if err != nil {
			log.Println("Error decoding tasks: ", err)
			return
		}

		for _, task := range outputTasks {
			_, err = b.Send(tgbotapi.NewMessage(task.ChatID, fmt.Sprintf("Напоминание для @%s:\n\n%s", task.UserName, task.Content)))
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

func (b *Bot) HandleMyChatMemberUpdate(myChatMember *tgbotapi.ChatMemberUpdated) {
	if myChatMember.NewChatMember.User.UserName == b.Self.UserName {
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
				"\n• s - секунды, \n• h -часы, \n• d - дни, \n• w - недели, \n• m - месяцы",
				b.Self.UserName, b.Self.UserName)
			msg := tgbotapi.NewMessage(myChatMember.Chat.ID, messageText)
			_, err := b.Send(msg)
			if err != nil {
				log.Println(err)
			}
			/*
				msg = tgbotapi.NewMessage(myChatMember.Chat.ID, "Пожалуйста, выберите опцию:")
				msg.ReplyMarkup = menu
				_, err = b.Send(msg)
				if err != nil {
					log.Println(err)
				}
			*/
		default:

		}
	}
}
