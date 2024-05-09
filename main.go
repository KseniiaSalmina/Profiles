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

// @title Profiles management API
// @version 1.0.0
// @description API to manage users profiles
// @host localhost:8080
// @BasePath /
func main() {
	application, err := app.NewApplication(cfg)
	if err != nil {
		log.Fatal(err)
	}
	application.Run()
}
