package main

import (
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"

	_ "github.com/KseniiaSalmina/Profiles/docs"
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

// @title Profiles managment API
// @version 1.0.0
// @description service to managment users profiles

// @host localhost:8080
// @BasePath /

// @securityDefinitions.basic BasicAuth
// @in header
// @name Authorization
func main() {
	application, err := app.NewApplication(cfg)
	if err != nil {
		log.Fatal(err)
	}
	application.Run()
}
