package main

import (
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"

	app "github.com/KseniiaSalmina/Profiles/internal"
	"github.com/KseniiaSalmina/Profiles/internal/config"
)

var cfg config.Application

func init() {
	_ = godotenv.Load(".env")
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}
}

func main() {
	application, err := app.NewApplication(cfg)
	if err != nil {
		log.Fatal(err)
	}
	application.Run()
}
