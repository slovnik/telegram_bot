package main

type Env struct {
	config   *Config
	template *Template
}

func main() {
	env := Env{
		config:   Setup(),
		template: CreateTemplate(),
	}

	b := CreateBot(&env)
	b.Listen()
}
