package main

import (
	"log"
	"strings"

	"github.com/slovnik/slovnik"

	"text/template"

	"bytes"

	"github.com/slovnik/telegram_bot/bot"
	"github.com/slovnik/telegram_bot/config"
)

var cfg *config.Config

var tmpl *template.Template

var templateFiles = []string{
	"./tmpl/full-word.gotmpl",
	"./tmpl/short-word.gotmpl",
	"./tmpl/translation.gotmpl",
}

func main() {

	cfg = config.Setup()

	funcs := template.FuncMap{
		"join": strings.Join,
	}

	var err error

	tmpl, err = template.New("").Funcs(funcs).ParseFiles(templateFiles...)

	if err != nil {
		log.Println(err)
	}

	bot.Create(cfg, handleIt)
	bot.Listen()
}

func handleIt(message string) string {
	c, err := slovnik.NewClient(cfg.SlovnikURL)
	if err != nil {
		log.Fatalln(err)
	}

	words, err := c.Translate(message)

	if err != nil {
		log.Println(err)
	}

	var buf bytes.Buffer

	tmpl.ExecuteTemplate(&buf, "translation", words)

	return buf.String()
}
