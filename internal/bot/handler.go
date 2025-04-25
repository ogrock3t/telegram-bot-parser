package bot

import (
	"fmt"
	"log"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	lastUserMsgID int
	lastBotMsgID  int
	msgMutex      sync.Mutex
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
		reply.Text = "I don't know the command you wrote :(\n"
		reply.Text += "I know commands:\n"
		reply.ReplyMarkup = generateHelpKeyboardAllCommands()
	}

	if _, err := botAPI.Send(reply); err != nil {
		log.Println("Error sending reply: ", err)
	}
}

func handleCallback(botAPI *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	msgMutex.Lock()
	defer msgMutex.Unlock()

	chatID := callback.Message.Chat.ID

	if lastBotMsgID != 0 {
		if _, err := botAPI.Send(tgbotapi.NewDeleteMessage(chatID, lastBotMsgID)); err != nil {
			log.Println("Error deleting previous bot message:", err)
		}
	}

	if _, err := botAPI.Send(tgbotapi.NewDeleteMessage(chatID, callback.Message.MessageID)); err != nil {
		log.Println("Error deleting button message:", err)
	}

	msg := tgbotapi.NewMessage(chatID, "")

	switch callback.Data {
	case "game_truth_or_dare":
		msg.Text = "You selected Truth or Dare!\n\nWhat's your choice?"
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Truth", "truth"),
				tgbotapi.NewInlineKeyboardButtonData("Dare", "dare"),
			),
		)
	case "start":
		msg.Text = generateWelcomeMessage(callback.From.FirstName)
	case "game", "play":
		msg.Text = "Choose a game:"
		msg.ReplyMarkup = generateHelpKeyboardForTruthOrDate()
	default:
		msg.Text = "Unknown selection"
	}

	sentMsg, err := botAPI.Send(msg)
	if err != nil {
		log.Println("Error sending reply:", err)
	} else {
		lastBotMsgID = sentMsg.MessageID
	}

	if callback.Message.ReplyToMessage != nil {
		lastUserMsgID = callback.Message.ReplyToMessage.MessageID
	}

	if _, err := botAPI.Request(tgbotapi.NewCallback(callback.ID, "")); err != nil {
		log.Println("Error answering callback:", err)
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

func generateHelpKeyboardAllCommands() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("/start", "start"),
		), tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("/game", "game"),
		),
	)
}
