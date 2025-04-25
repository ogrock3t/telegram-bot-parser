package main

import (
	"log"
	
	"github.com/ogrock3t/telegram-bot-parser/internal/bot"
	"github.com/ogrock3t/telegram-bot-parser/internal/config"
	"github.com/ogrock3t/telegram-bot-parser/internal/configAkar4ick.com"
)

func main() {
	cfg, err := config.Load("config.json")
	if err != nil {
		log.Panic("Failed to load config: ", err)
	}

	bot.Run(cfg.TelegramToken)
}
