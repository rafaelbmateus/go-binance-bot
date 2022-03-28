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

var (
	Version = "no version provided"
	Commit  = "no commit hash provided"
)

func main() {
	ctx := context.Background()
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()
	log.Info().Msgf("starting with version %s and commit %s", Version, Commit)

	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatal().Msgf("error on load config file %q", err)
		return
	}

	apiKey := os.Getenv("BINANCE_API_KEY")
	apiSecret := os.Getenv("BINANCE_API_SECRET")
	webhook := os.Getenv("SLACK_WEBHOOK_URL")
	binance := binance.NewBinance(&log, &ctx, sdk.NewClient(apiKey, apiSecret))
	notify := notify.NewSlackNotify(cfg.Name, webhook)
	bot := bot.NewBot(&log, &ctx, binance, cfg, notify)
	bot.Run()

	select {} // block forever
}
