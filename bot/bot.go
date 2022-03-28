package bot

import (
	"context"
	"fmt"
	"time"

	"github.com/rafaelbmateus/binance-bot/binance"
	"github.com/rafaelbmateus/binance-bot/config"
	"github.com/rafaelbmateus/binance-bot/notify"
	"github.com/rs/zerolog"
)

// Bot is a representation of this package.
type Bot struct {
	Log     *zerolog.Logger
	Context *context.Context
	Binance *binance.Binance
	Config  *config.Config
	Notify  *notify.SlackNotify
}

// NewBot create a new bot instance.
func NewBot(log *zerolog.Logger, ctx *context.Context,
	binance *binance.Binance, cfg *config.Config, notify *notify.SlackNotify) *Bot {
	return &Bot{
		Log:     log,
		Context: ctx,
		Binance: binance,
		Config:  cfg,
		Notify:  notify,
	}
}

// Run bot to each trade with go routine.
func (me *Bot) Run() {
	me.Log.Info().Msg("bot is running")

	for _, trade := range me.Config.Trades {
		go me.Watchdog(trade)
	}

	me.Notify.SendMessage(notify.NewMessage(welcomeMessage(*me.Config)))
}

// Watchdog the trade forever in a time interval.
func (me *Bot) Watchdog(trade config.Trade) {
	for {
		err := me.Monitor(trade)
		if err != nil {
			me.Log.Error().Msgf("watchdog sleep error %q", err)
		}

		time.Sleep(trade.Interval)
	}
}

// Monitor symbol to buy or sell.
func (me *Bot) Monitor(trade config.Trade) error {
	me.Log.Debug().Msgf("monitor %v", trade)

	currentPrice, err := me.Binance.SymbolPrice(trade.GetSymbol(), trade.BuyWith(), "1m")
	me.Log.Debug().Msgf("current price of %s is %v", trade.Symbol, currentPrice)
	if err != nil {
		return err
	}

	if currentPrice < trade.BuyPrice {
		err := me.trade("BUY", trade, currentPrice)
		if err != nil {
			return err
		}
	}

	if currentPrice > trade.SellPrice {
		err := me.trade("SELL", trade, currentPrice)
		if err != nil {
			return err
		}
	}

	me.Log.Debug().Msgf("nothing to do %v", trade)

	return nil
}

// trade create an order to buy or sell.
func (me *Bot) trade(side string, trade config.Trade, price float64) error {
	me.Log.Debug().Msgf("time to %s price %v", side, price)
	wallet, err := me.Binance.SymbolBalance(trade.GetSymbol())
	if err != nil {
		return err
	}

	if wallet == 0 {
		return nil
	}

	me.Log.Debug().Msgf("available to %s in wallet %f", side, wallet)
	order, err := me.Binance.CreateOrder(trade.GetSymbol(), trade.BuyWith(),
		side, wallet, trade.SellPrice)

	if err != nil {
		return err
	}

	me.Log.Debug().Msgf("order to %s created %v", side, trade)
	me.Notify.SendMessage(notify.NewMessage(fmt.Sprintf("order to %s created %v", side, order)))

	return nil
}

func welcomeMessage(config config.Config) string {
	msg := fmt.Sprintf("%s started! :money_mouth_face:\n\n", config.Name)
	for i, trade := range config.Trades {
		msg += fmt.Sprintf("> %d: %s\n", i+1, trade.Symbol)
		msg += fmt.Sprintf("> Interval: %s\n", trade.Interval)
		msg += fmt.Sprintf("> Purchase price: %f\n", trade.BuyPrice)
		msg += fmt.Sprintf("> Sale Price: %f\n", trade.SellPrice)
	}

	return msg
}
