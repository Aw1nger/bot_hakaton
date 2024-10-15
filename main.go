package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UserTest struct {
	ChatID  int64
	Current int
	Answers []string
}

var tests = make(map[int64]*UserTest)

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
				if test, exists := tests[update.Message.Chat.ID]; exists {
					handleAnswer(test, update.Message.Text, bot)
				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Я не понимаю. Используйте /help для получения информации.")
					bot.Send(msg)
				}
			}
		} else if update.CallbackQuery != nil {
			switch update.CallbackQuery.Data {
			case "login":
				keyboard := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Начать тест", "testStart"),
					),
				)

				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Нажмите чтобы начать тест:")
				msg.ReplyMarkup = keyboard
				bot.Send(msg)

			case "testStart":
				if _, exists := tests[update.CallbackQuery.Message.Chat.ID]; exists {
					msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Вы уже начали тест. Пожалуйста, ответьте на текущий вопрос.")
					bot.Send(msg)
				} else {
					handleLogin(update.CallbackQuery.Message.Chat.ID, bot)
				}
			}
		}
	}
}

func handleLogin(chatID int64, bot *tgbotapi.BotAPI) {
	test := &UserTest{
		ChatID:  chatID,
		Current: 0,
		Answers: make([]string, 0),
	}
	tests[chatID] = test
	startTest(test, bot)
}

func startTest(test *UserTest, bot *tgbotapi.BotAPI) {
	questions := []string{
		"Вопрос 1:",
		"Вопрос 2:",
		"Вопрос 3:",
		"Вопрос 4:",
		"Вопрос 5:",
		"Вопрос 6:",
		"Вопрос 7:",
		"Вопрос 8:",
		"Вопрос 9:",
		"Вопрос 10:",
	}

	if test.Current < len(questions) {
		msg := tgbotapi.NewMessage(test.ChatID, questions[test.Current])
		bot.Send(msg)
	} else {
		finishTest(test, bot)
	}
}

func handleAnswer(test *UserTest, answer string, bot *tgbotapi.BotAPI) {
	test.Answers = append(test.Answers, answer)
	test.Current++
	startTest(test, bot)
}

func finishTest(test *UserTest, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(test.ChatID, "Тест успешно пройден! Возвращайтесь на сайт: [ссылка](https://example.com)")
	msg.ParseMode = "Markdown"
	bot.Send(msg)

	log.Printf("Ответы пользователя %d: %v", test.ChatID, test.Answers)

	delete(tests, test.ChatID)
}