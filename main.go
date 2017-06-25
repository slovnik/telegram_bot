package main

import (
	"log"
)

type Env struct {
	config   *Config
	template *Template
}

func main() {
	env := Env{
		config:   Setup(),
		template: CreateTemplate(),
	}

	b, err := CreateBot(&env)
	if err != nil {
		log.Fatalln(err)
	}
	b.Listen()
}
