package bot

import (
	"strconv"
	"github.com/YotoHana/tgbot-golang/config"
	db "github.com/YotoHana/tgbot-golang/internal/database"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var userStates = make(map[int64]string)

func Run() error{
	cfg, err := config.ReadCfg()
	if err != nil {
		return err
	}


	bot, err := tgbotapi.NewBotAPI(cfg.ApiKey)
	if err != nil {
		return err
	}

	bot.Debug = false

	u := tgbotapi.NewUpdate(0)

	u.Timeout = 30

	updates := bot.GetUpdatesChan(u)

	for update := range updates {

		if update.CallbackQuery == nil && update.Message == nil {
			continue
		}

		if update.CallbackQuery != nil {
			task := update.CallbackQuery.Data
			sql, err := db.ConnectToDb()
			if err != nil {
				return err
			}
			err = db.DeleteData(sql, update.CallbackQuery.Message.Chat.ID, task)
			if err != nil {
				return err
			}
			message := "Задача " + task + " удалена!"
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, message)
			bot.Send(msg)
		} else if update.Message != nil {
			if update.Message.IsCommand() {
				chatID := update.Message.Chat.ID
				messageID := update.Message.MessageID
				command := update.Message.Command()
				args := update.Message.CommandArguments()

				switch command {
				case "start":
					message := "Привет! Я создан для того, чтобы запоминать твои задачи! Пока я работаю в тестовом режиме, о всех недочетах сообщайте в /report. Для помощи напиши /help"
					msg := tgbotapi.NewMessage(chatID, message)
					bot.Send(msg)
				case "help":
					message := "/add чтобы добавить новую задачу \n /list чтобы увидеть список задач \n /remove чтобы удалить задачу"
					msg := tgbotapi.NewMessage(chatID, message)
					msg.ReplyToMessageID = messageID
					bot.Send(msg)
				case "add":
					if args == "" {
						message := "Введите задачу"
						userStates[chatID] = "awaiting_task"
						msg := tgbotapi.NewMessage(chatID, message)
						msg.ReplyToMessageID = messageID
						bot.Send(msg)

					} else {
						sql, err := db.ConnectToDb()
						if err != nil {
							return err
						}
						db.CreateTable(sql)
						err = db.InsertData(sql, chatID, args)
						if err != nil {
							return err
						}
						message := "Задача добавлена"
						msg := tgbotapi.NewMessage(chatID, message)
						msg.ReplyToMessageID = messageID
						bot.Send(msg)
					}
				case "list":
					sql, err := db.ConnectToDb()
					if err != nil {
						return err
					}
					var title []db.Title
					title, err = db.QueryData(sql, chatID)
					if err != nil {
						return err
					}
					if title == nil {
						msg := tgbotapi.NewMessage(chatID, "У вас нет задач!")
						bot.Send(msg)
					}
					var message string
					for i, v := range title {
						message += strconv.Itoa(i+1) + ". " + v.Title + "\n"
					}
					
					msg := tgbotapi.NewMessage(chatID, message)
					bot.Send(msg)

				case "time":
					msg := tgbotapi.NewMessage(chatID, "Через какое время вас оповестить? (В форме чч:мм)")
					userStates[chatID] = "awaiting_time"
					bot.Send(msg) 
				case "remove":
					sql, err := db.ConnectToDb()
					if err != nil {
						return err
					}
					msg := tgbotapi.NewMessage(chatID, "Какую задачу хотите удалить?")
					
					var title []db.Title
					title, err = db.QueryData(sql, chatID)
					if err != nil {
						return err
					}

					keyboard := tgbotapi.InlineKeyboardMarkup{}
					for _, v := range title {
						var row []tgbotapi.InlineKeyboardButton
						btn := tgbotapi.NewInlineKeyboardButtonData(v.Title, v.Title)
						row = append(row, btn)
						keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
					}

					msg.ReplyMarkup = keyboard
					bot.Send(msg)

				case "report":
					msg := tgbotapi.NewMessage(chatID, "Введите сообщение репорта:")
					userStates[chatID] = "awaiting_report"
					bot.Send(msg)
				default:
					msg := tgbotapi.NewMessage(chatID, "Команда не найдена")
					msg.ReplyToMessageID = messageID
					bot.Send(msg)
				}
			} else {
				if userStates[update.Message.Chat.ID] == "awaiting_task" {
					task := update.Message.Text
					sql, err := db.ConnectToDb()
					if err != nil {
						return err
					}
					err = db.InsertData(sql, update.Message.Chat.ID, task)
					if err != nil {
						return err
					}
					delete(userStates, update.Message.Chat.ID)
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Задача добавлена: " + task)
					bot.Send(msg)
				} else if userStates[update.Message.Chat.ID] == "awaiting_report" {
					report := update.Message.Text
					message := "Вам пришел репорт! \n - "
					msg := tgbotapi.NewMessage(cfg.AdminChatID, message + report)
					bot.Send(msg)
					delete(userStates, update.Message.Chat.ID)
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Репорт успешно отправлен!")
					bot.Send(msg)
				} else if userStates[update.Message.Chat.ID] == "awaiting_time" {
					timeText := update.Message.Text
					delete(userStates, update.Message.Chat.ID)
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, timeText)
					bot.Send(msg)
				}
			}
		}
	}
	return nil
}