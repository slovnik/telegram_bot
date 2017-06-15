package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/slovnik/slovnik"

	"strings"

	"gopkg.in/telegram-bot-api.v4"
)

// Bot type aggregate all bot logic
type Bot struct {
	api *tgbotapi.BotAPI

	// updates is a channel that's used to communicate with telegram server
	updates <-chan tgbotapi.Update

	// env contains environvent for the bot
	env *Env

	// client for the slovnik API
	slovnikAPIClient *slovnik.Client
}

// CreateBot creates and initializes new bot
func CreateBot(env *Env) *Bot {
	var err error

	bot := Bot{env: env}
	bot.slovnikAPIClient, err = slovnik.NewClient(bot.env.config.SlovnikURL)

	if bot.api, err = tgbotapi.NewBotAPI(env.config.BotID); err != nil {
		log.Panic(err)
	}

	if len(env.config.WebhookURL) == 0 {
		log.Println("WebhookURL environment variable is not set. Using polling.")

		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60

		bot.updates, err = bot.api.GetUpdatesChan(u)

	} else {
		log.Printf("WebhookURL is set to '%s'. Using webhooks\n", env.config.WebhookURL)
		webHookURL := fmt.Sprintf("%s/bot%s", env.config.WebhookURL, bot.api.Token)
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
		if update.Message != nil {
			bot.handleMessage(&update)
		} else if update.CallbackQuery != nil {
			bot.handleCallbackQuery(&update)
		}

	}
}

func (bot *Bot) handleMessage(update *tgbotapi.Update) {
	words, err := bot.slovnikAPIClient.Translate(update.Message.Text)
	if err != nil {
		log.Println(err)
	}
	messageText := bot.env.template.Execute(words)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, messageText)
	msg.ParseMode = tgbotapi.ModeMarkdown
	msg.ReplyMarkup = bot.addMessageKeyboard(words)
	_, err = bot.api.Send(msg)

	if err != nil {
		log.Println(err)
	}
}

func (bot *Bot) handleCallbackQuery(update *tgbotapi.Update) {
	callbackData := update.CallbackQuery.Data

	if strings.HasPrefix(callbackData, "phrases:") {
		w := strings.TrimPrefix(callbackData, "phrases:")
		words, err := bot.slovnikAPIClient.Translate(w)
		messageText := bot.env.template.Phrases(words)
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, messageText)
		msg.ParseMode = tgbotapi.ModeMarkdown

		_, err = bot.api.Send(msg)
		if err != nil {
			log.Println(err)
		}
	}
}

func (bot *Bot) addMessageKeyboard(words []slovnik.Word) *tgbotapi.InlineKeyboardMarkup {
	if words == nil || len(words) > 1 || len(words[0].Samples) <= 0 {
		return nil
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Фразы", "phrases:"+words[0].Word),
		),
	)

	return &keyboard
}
