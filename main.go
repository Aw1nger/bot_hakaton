package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	token, err := os.ReadFile("token.txt")
	if err != nil {
		log.Fatalf("Ошибка чтения файла с токеном: %v", err)
	}

	bot, err := tgbotapi.NewBotAPI(string(token))
	if err != nil {
		log.Fatalf("Ошибка создания бота: %v", err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			switch update.Message.Text {
			case "/start":

				keyboard := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Войти", "login"),
					),
				)

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Добро пожаловать на тест! Войдите чтобы начать тест")

				msg.ReplyMarkup = keyboard
				bot.Send(msg)

			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Я не понимаю. Используйте /help для получения информации.")
				bot.Send(msg)
			}
		} else if update.CallbackQuery != nil {
			switch update.CallbackQuery.Data {
			case "login":
				handleLogin(update.CallbackQuery.Message.Chat.ID, bot)
			}
		}
	}
}

func handleLogin(chatID int64, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(chatID, "Вы нажали 'Войти'. Здесь будет логика входа.")
	bot.Send(msg)
}

//dima loh
