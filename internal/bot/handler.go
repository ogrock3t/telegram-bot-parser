package bot

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Run(token string) {
	botAPI, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic("Error creating bot: ", err)
	}

	botAPI.Debug = true
	log.Printf("Bot %s started", botAPI.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := botAPI.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		handleMessage(botAPI, update.Message)
	}
}

func handleMessage(botAPI *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	log.Printf("[%s] %s", msg.From.UserName, msg.Text)

	reply := tgbotapi.NewMessage(msg.Chat.ID, "")
	switch msg.Text {
	case "/start":
		reply.Text = generateWelcomeMessage(msg.From.FirstName)
	case "/help":
		reply.Text = "Select the ruble exchange rate to:"
		reply.ReplyMarkup = generateHelpKeyboard()
	default:
		reply.Text = "..."
	}

	if _, err := botAPI.Send(reply); err != nil {
		log.Println("Error sending reply: ", err)
	}
}

func generateWelcomeMessage(userName string) string {
	if userName == "" {
		userName = "friend"
	}

	return fmt.Sprintf("Welcome, %s!\nThis bot provides real-time exchange rate monitoring.\n\nAuthor: github.com/ogrock3t", userName)
}

func generateHelpKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("USD", "currency_usd"),
			tgbotapi.NewInlineKeyboardButtonData("EUR", "currency_eur"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("BTC", "currency_btc"),
		),
	)
}
