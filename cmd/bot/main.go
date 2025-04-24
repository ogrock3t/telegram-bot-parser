package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Config struct {
	TelegramToken string `json:"telegram_token"`
}

func main() {
	configFile, err := os.Open("config.json")
	if err != nil {
		log.Panic("Error opening config.json: ", err)
	}
	defer configFile.Close()

	var config Config
	decoder := json.NewDecoder(configFile)
	if err := decoder.Decode(&config); err != nil {
		log.Panic("Error decoding config.json: ", err)
	}

	bot, err := tgbotapi.NewBotAPI(config.TelegramToken)
	if err != nil {
		log.Panic("Error creating telegram bot: ", err)
	}

	bot.Debug = true
	log.Printf("Telegram bot %s starting work", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		switch update.Message.Text {
		case "/start":
			userName := update.Message.From.FirstName
			if userName == "" {
				userName = "friend"
			}

			msg.Text = fmt.Sprintf("Welcome, %s!\n"+
				"This bot provides real-time exchange rate monitoring.\n\n"+
				"Author: github.com/ogrock3t", userName)

		case "/help":
			msg.Text = "..."

		default:
			msg.Text = "..."
		}

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}
