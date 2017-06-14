package main

import (
	"fmt"
	"log"
	"net/http"

	"gopkg.in/telegram-bot-api.v4"
)

// HandlerFunc is a function definition that will be processing bot messages
type HandlerFunc func(message string) string

// Bot type aggregate all bot logic
type Bot struct {
	api *tgbotapi.BotAPI

	// updates is a channel that's used to communicate with telegram server
	updates <-chan tgbotapi.Update

	// handler is a function that's creating response for the message
	handler HandlerFunc
}

// CreateBot creates and initializes new bot
func CreateBot(config *Config, handlerFunc HandlerFunc) *Bot {
	bot := Bot{handler: handlerFunc}
	var err error

	if bot.api, err = tgbotapi.NewBotAPI(config.BotID); err != nil {
		log.Panic(err)
	}

	bot.api.Debug = true

	if len(config.WebhookURL) == 0 {
		log.Println("WebhookURL environment variable is not set. Using polling.")

		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60

		bot.updates, err = bot.api.GetUpdatesChan(u)

	} else {
		log.Printf("WebhookURL is set to '%s'. Using webhooks", config.WebhookURL)
		webHookURL := fmt.Sprintf("%s/bot%s", config.WebhookURL, bot.api.Token)
		webHook := tgbotapi.NewWebhook(webHookURL)
		if _, err = bot.api.SetWebhook(webHook); err != nil {
			log.Panic(err)
		}

		bot.updates = bot.api.ListenForWebhook("/bot" + bot.api.Token)
		go http.ListenAndServe("0.0.0.0:8080", nil)
	}

	return &bot
}

// Listen start listening on message updates and calling provided handler for processing incoming messages
func (bot *Bot) Listen() {
	for update := range bot.updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		messageText := bot.handler(update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, messageText)

		msg.ParseMode = tgbotapi.ModeMarkdown

		_, err := bot.api.Send(msg)

		if err != nil {
			log.Println(err)
		}
	}
}
