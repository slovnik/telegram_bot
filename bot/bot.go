package bot

import (
	"log"
	"net/http"

	"github.com/slovnik/telegram_bot/config"
	"gopkg.in/telegram-bot-api.v4"
)

// HandlerFunc is a function definition that will be processing bot messages
type HandlerFunc func(message string) string

var bot *tgbotapi.BotAPI

// updates is a channel that's used to communicate with telegram server
var updates <-chan tgbotapi.Update

// handler is a function that's creating response for the message
var handler HandlerFunc

// Create new bot
func Create(config *config.Config, handlerFunc HandlerFunc) {

	handler = handlerFunc

	var err error

	if bot, err = tgbotapi.NewBotAPI(config.BotID); err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	if len(config.WebhookURL) == 0 {
		log.Println("WebhookURL environment variable is not set. Using polling.")

		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60

		updates, err = bot.GetUpdatesChan(u)

	} else {
		log.Printf("WebhookURL is set to '%s'. Using webhooks", config.WebhookURL)

		webHook := tgbotapi.NewWebhook(config.WebhookURL + "/bot" + bot.Token)
		if _, err = bot.SetWebhook(webHook); err != nil {
			log.Panic(err)
		}

		updates = bot.ListenForWebhook("/bot" + bot.Token)
		go http.ListenAndServe("0.0.0.0:8080", nil)
	}
}

// Listen start listening on message updates and calling provided handler for processing incoming messages
func Listen() {
	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		messageText := handler(update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, messageText)

		msg.ParseMode = tgbotapi.ModeMarkdown

		bot.Send(msg)
	}
}
