package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UserTest struct {
	ChatID  int64
	Current int
	Answers []string
}

var tests = make(map[int64]*UserTest)
var userTokens = make(map[int64]string)

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

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Добро пожаловать на тест! Войдите, чтобы начать.")
				msg.ReplyMarkup = keyboard
				bot.Send(msg)

			default:
				if test, exists := tests[update.Message.Chat.ID]; exists {
					switch test.Current {
					case 0:
						test.Answers = append(test.Answers, update.Message.Text)
						test.Current++
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите ваш пароль:")
						bot.Send(msg)
					case 1:
						test.Answers = append(test.Answers, update.Message.Text)
						userLogin(test.Answers[0], test.Answers[1], bot, update.Message.Chat.ID)
					}
				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Я не понимаю. Используйте /help для получения информации.")
					bot.Send(msg)
				}
			}
		} else if update.CallbackQuery != nil {
			switch update.CallbackQuery.Data {
			case "login":
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Введите вашу почту:")
				bot.Send(msg)
				tests[update.CallbackQuery.Message.Chat.ID] = &UserTest{
					ChatID:  update.CallbackQuery.Message.Chat.ID,
					Current: 0,
					Answers: make([]string, 0),
				}
			}
		}
	}
}

func userLogin(email string, password string, bot *tgbotapi.BotAPI, chatID int64) {
	loginData := map[string]string{
		"email":    email,
		"password": password,
	}

	jsonData, err := json.Marshal(loginData)
	if err != nil {
		log.Printf("Ошибка маршализации данных: %v", err)
		return
	}

	req, err := http.NewRequest("POST", "https://your-api-url.com/api/users/login", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Ошибка создания запроса: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Ошибка выполнения запроса: %v", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		var result map[string]interface{}
		err := json.Unmarshal(body, &result)
		if err != nil {
			log.Printf("Ошибка декодирования ответа: %v, ответ: %s", err, string(body))
			return
		}

		log.Printf("Ответ от API: %v", result)

		token, ok := result["token"].(string)
		if !ok {
			log.Printf("Токен отсутствует или неверного формата: %v", result)
			msg := tgbotapi.NewMessage(chatID, "Ошибка входа: токен не найден.")
			bot.Send(msg)
			return
		}

		msg := tgbotapi.NewMessage(chatID, "Вход выполнен успешно!")
		bot.Send(msg)

		saveToken(chatID, token)

		getUserProfile(chatID, bot)
	} else {
		log.Printf("Ошибка входа, статус: %d, тело ответа: %s", resp.StatusCode, string(body))
		msg := tgbotapi.NewMessage(chatID, "Ошибка входа. Проверьте данные.")
		bot.Send(msg)
	}
}


func saveToken(chatID int64, token string) {
	userTokens[chatID] = token
}

func makeAuthenticatedRequest(chatID int64, endpoint string) (*http.Response, error) {
	token, exists := userTokens[chatID]
	if !exists {
		return nil, errors.New("пользователь не авторизован")
	}

	req, err := http.NewRequest("GET", "https://your-api-url.com"+endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	return client.Do(req)
}

func getUserProfile(chatID int64, bot *tgbotapi.BotAPI) {
	resp, err := makeAuthenticatedRequest(chatID, "/api/users/profile")
	if err != nil {
		log.Printf("Ошибка при выполнении запроса профиля: %v", err)
		msg := tgbotapi.NewMessage(chatID, "Ошибка получения профиля.")
		bot.Send(msg)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Ваш профиль: %s", string(body)))
		bot.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(chatID, "Не удалось получить профиль.")
		bot.Send(msg)
	}
}