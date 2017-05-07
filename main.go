package main

import (
	"log"

	"github.com/slovnik/slovnik"

	"github.com/slovnik/telegram_bot/bot"
	"github.com/slovnik/telegram_bot/config"
)

var cfg *config.Config

func main() {

	cfg = config.Setup()

	bot.Create(cfg, handleIt)
	bot.Listen()
}

func handleIt(message string) (response string) {
	c, err := slovnik.NewClient(cfg.SlovnikURL)
	if err != nil {
		log.Fatalln(err)
	}

	w, err := c.Translate(message)

	if err != nil || len(w.Word) <= 0 {
		response = "Specified word not found :("
	} else {
		response = w.String()
	}

	return
}
