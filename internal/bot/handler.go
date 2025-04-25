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
	u.Timeout = 20
	updates := botAPI.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			handleMessage(botAPI, update.Message)
		} else if update.CallbackQuery != nil {
			handleCallback(botAPI, update.CallbackQuery)
		}

	}
}

func handleMessage(botAPI *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	log.Printf("[%s] %s", msg.From.UserName, msg.Text)

	reply := tgbotapi.NewMessage(msg.Chat.ID, "")
	switch msg.Text {
	case "/start":
		reply.Text = generateWelcomeMessage(msg.From.FirstName)
	case "/help":
		reply.Text = "..."
	case "/play", "/game":
		reply.Text = "Choose any games:"
		reply.ReplyMarkup = generateHelpKeyboardForTruthOrDate()
	default:
		reply.Text = "..."
	}

	if _, err := botAPI.Send(reply); err != nil {
		log.Println("Error sending reply: ", err)
	}
}

func handleCallback(botAPI *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "")

	switch callback.Data {
	case "game_truth_or_dare":
		msg.Text = "You selected Truth or Dare!\n\n"
	default:
		msg.Text = "Unknown game selection"
	}

	if _, err := botAPI.Send(msg); err != nil {
		log.Println("Error sending callback reply: ", err)
	}

	callbackConfig := tgbotapi.NewCallback(callback.ID, "")
	if _, err := botAPI.Request(callbackConfig); err != nil {
		log.Println("Error answering callback: ", err)
	}
}

func generateWelcomeMessage(userName string) string {
	if userName == "" {
		userName = "friend"
	}

	return fmt.Sprintf("Welcome, %s!\nThis bot provides real-time exchange rate monitoring.\n\nAuthor: github.com/ogrock3t", userName)
}

func generateHelpKeyboardForTruthOrDate() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Truth or Dare", "game_truth_or_dare"),
		),
	)
}
