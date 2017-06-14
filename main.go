package main

import (
	"log"

	"github.com/slovnik/slovnik"
)

type Env struct {
	cfg      *Config
	template *Template
}

func main() {
	env := Env{
		cfg:      Setup(),
		template: CreateTemplate(),
	}

	b := CreateBot(env.cfg, env.handleIt)
	b.Listen()
}

func (e *Env) handleIt(message string) string {
	c, err := slovnik.NewClient(e.cfg.SlovnikURL)
	if err != nil {
		log.Println(err)
	}

	words, err := c.Translate(message)

	if err != nil {
		log.Println(err)
		return e.template.Execute(nil)
	}

	return e.template.Execute(words)
}
