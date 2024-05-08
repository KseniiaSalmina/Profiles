package main

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"

	"github.com/KseniiaSalmina/Profiles/internal/config"
)

var cfg config.Server

func init() {
	_ = godotenv.Load(".env")
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}
}

func main() {

}
