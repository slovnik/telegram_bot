package main

import (
	"fmt"
	"log"
	"strings"

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

	words, err := c.Translate(message)

	if err != nil {
		response = "Указанное слово не найдено :("
	} else {
		if len(words) > 1 {
			for _, w := range words {
				response += fmt.Sprintln(shortTranslation(w))
			}
		} else if len(words) == 1 {
			response = fullTranslation(words[0])
		} else {
			response = "Указанное слово не найдено :("
		}

	}

	return
}

func fullTranslation(w slovnik.Word) string {
	out := fmt.Sprintf("*%s* - %s\n\n", w.Word, strings.Join(w.Translations, ", "))
	out += fmt.Sprintf("*%s*\n", w.WordType)

	if len(w.Synonyms) > 0 {
		out += fmt.Sprintln("\n*Synonyms:*")
		out += fmt.Sprintln(strings.Join(w.Synonyms, ", "))
	}
	if len(w.Antonyms) > 0 {
		out += fmt.Sprintln("\n*Antonyms:*")
		out += fmt.Sprintln(strings.Join(w.Antonyms, ", "))
	}

	if len(w.DerivedWords) > 0 {
		out += fmt.Sprintln("\n*Derived words:*")
		out += fmt.Sprintln(strings.Join(w.DerivedWords, ", "))
	}
	return out
}

func shortTranslation(w slovnik.Word) string {
	return fmt.Sprintf("*%s* - %s", w.Word, strings.Join(w.Translations, ", "))
}
