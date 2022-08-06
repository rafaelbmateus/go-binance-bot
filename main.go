package main

import (
	"context"
	"os"

	sdk "github.com/adshao/go-binance/v2"
	"github.com/rafaelbmateus/binance-bot/binance"
	"github.com/rafaelbmateus/binance-bot/bot"
	"github.com/rafaelbmateus/binance-bot/config"
	"github.com/rafaelbmateus/binance-bot/notify"
	"github.com/rs/zerolog"
)

func main() {
	ctx := context.Background()

	log := zerolog.New(os.Stdout).With().Timestamp().Logger()
	log.Level(zerolog.InfoLevel)
	log.Debug().Msgf("service starging...")

	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatal().Msgf("error on load config file %q", err)
		return
	}

	binance := binance.NewBinance(&log, &ctx,
		sdk.NewClient(cfg.Binance.ApiKey, cfg.Binance.ApiSecret))
	notify := notify.NewSlackNotify(cfg.Name, cfg.Notify.WebhookURL)
	bot := bot.NewBot(&log, &ctx, binance, cfg, notify)
	bot.Run()

	select {} // block forever
}
