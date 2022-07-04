package main

import (
	"time"

	"github.com/MarySmirnova/api_gateway/internal"
	"github.com/MarySmirnova/api_gateway/internal/config"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"

	log "github.com/sirupsen/logrus"
)

var cfg config.Application

func init() {
	godotenv.Load(".env")
	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}

	lvl, err := log.ParseLevel(cfg.LogLevel)
	if err != nil {
		panic(err)
	}

	log.SetLevel(lvl)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.Stamp,
	})
}

func main() {
	app := internal.NewApplication(cfg)

	app.StartServer()
}
