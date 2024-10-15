package main

import (
    "log"
    "os"
    "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

    log.Printf("Авторизован как %s", bot.Self.UserName)

    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60

    updates := bot.GetUpdatesChan(u)

    for update := range updates {
        if update.Message != nil {
            log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
            msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Ты написал: "+update.Message.Text)
            bot.Send(msg)
        }
    }
}
