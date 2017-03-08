package main

import (
	"log"

	"os"

	"github.com/slovnik/slovnik"

	"net/http"

	"gopkg.in/telegram-bot-api.v4"
)

func main() {

	botID := os.Getenv("SLOVNIK_BOT_ID")
	slovnikURL := os.Getenv("SLOVNIK_API_URL")
	webhookURL := os.Getenv("SLOVNIK_WEBHOOK_URL")

	var bot *tgbotapi.BotAPI
	var err error

	if bot, err = tgbotapi.NewBotAPI(botID); err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	webHook := tgbotapi.NewWebhook(webhookURL + "/bot" + bot.Token)
	if _, err = bot.SetWebhook(webHook); err != nil {
		log.Panic(err)
	}

	updates := bot.ListenForWebhook("/bot" + bot.Token)
	go http.ListenAndServe("0.0.0.0:8080", nil)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		messageText := ""

		c, _ := slovnik.NewClient(slovnikURL)
		w, err := c.Translate(update.Message.Text)

		if err != nil || len(w.Word) <= 0 {
			messageText = "Specified word not found :("
		} else {
			messageText = w.String()
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, messageText)

		msg.ParseMode = tgbotapi.ModeMarkdown
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}
