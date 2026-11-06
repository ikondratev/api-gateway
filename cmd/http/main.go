package main

import (
	"os"

	"github.com/ikondratev/api-gateway/internal/application"
)

func main() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	app, err := application.New(env)
	if err != nil {
		panic(err)
	}

	if err := app.Run(); err != nil {
		panic(err)
	}
}
