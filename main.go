package main

import (
	"log"
	"net/http"

	"os"

	"github.com/slovnik/slovnik"

	"gopkg.in/telegram-bot-api.v4"
)

const (
	// EnvBotID contains environment variable name with Bot ID received from Botfather
	EnvBotID = "SLOVNIK_BOT_ID"

	// EnvAPIURL contains
	EnvAPIURL = "SLOVNIK_API_URL"

	// EnvWebhookURL contains
	EnvWebhookURL = "SLOVNIK_WEBHOOK_URL"
)

func main() {

	botID, ok := os.LookupEnv(EnvBotID)
	if !ok {
		log.Panic(EnvBotID + " is not set!")
	}

	slovnikURL, ok := os.LookupEnv(EnvAPIURL)

	if !ok {
		log.Panic(EnvAPIURL + " is not set!")
	}

	var bot *tgbotapi.BotAPI
	var err error

	if bot, err = tgbotapi.NewBotAPI(botID); err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	webhookURL, ok := os.LookupEnv(EnvWebhookURL)

	var updates <-chan tgbotapi.Update

	if !ok {
		log.Printf("%s environment variable is not set. Using polling.", EnvWebhookURL)

		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60

		updates, err = bot.GetUpdatesChan(u)

	} else {
		log.Printf("%s is set to '%s'. Using webhooks", EnvWebhookURL, webhookURL)

		webHook := tgbotapi.NewWebhook(webhookURL + "/bot" + bot.Token)
		if _, err = bot.SetWebhook(webHook); err != nil {
			log.Panic(err)
		}

		updates = bot.ListenForWebhook("/bot" + bot.Token)
		go http.ListenAndServe("0.0.0.0:8080", nil)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		messageText := ""

		c, err := slovnik.NewClient(slovnikURL)
		if err != nil {
			log.Fatalln(err)
		}

		w, err := c.Translate(update.Message.Text)

		if err != nil || len(w.Word) <= 0 {
			messageText = "Specified word not found :("
		} else {
			messageText = w.String()
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, messageText)

		msg.ParseMode = tgbotapi.ModeMarkdown

		bot.Send(msg)
	}
}
