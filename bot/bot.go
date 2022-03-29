package bot

import (
	"context"
	"fmt"
	"math"
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
	currentPrice, err := me.Binance.SymbolPrice(trade.GetSymbol(), trade.BuyWith(), "1m")
	if err != nil {
		return err
	}

	if currentPrice <= trade.BuyPrice {
		me.Log.Debug().Msgf("[%s] current price=%.2f - nice to buy!",
			trade.Symbol, currentPrice, trade.BuyPrice)
		err := me.buy(trade, currentPrice)
		if err != nil {
			return err
		}

		return nil
	}

	if currentPrice >= trade.SellPrice {
		me.Log.Debug().Msgf("[%s] current price=%.2f - nice to sell!",
			trade.Symbol, currentPrice, trade.SellPrice)
		err := me.sell(trade, currentPrice)
		if err != nil {
			return err
		}

		return nil
	}

	me.Log.Debug().Msgf("[%s] current price=%.2f buy=%.2f sell=%.2f",
		trade.Symbol, currentPrice, trade.BuyPrice, trade.SellPrice)
	return nil
}

// trade create an order to buy or sell.
func (me *Bot) buy(trade config.Trade, price float64) error {
	wallet, err := me.Binance.SymbolBalance(trade.BuyWith())
	if err != nil {
		return err
	}

	if wallet == 0 {
		return nil
	}

	quantity := math.Floor(wallet / price)
	if quantity == 0 {
		return nil
	}

	order, err := me.Binance.CreateOrder(trade.GetSymbol(), trade.BuyWith(),
		"BUY", quantity, trade.BuyPrice)
	if err != nil {
		return err
	}

	me.Log.Debug().Msgf("[%s] buy for %d", trade.Symbol, price)
	me.Notify.SendMessage(notify.NewMessage(fmt.Sprintf("Order to buy created %v", order)))

	return nil
}

// trade create an order to buy or sell.
func (me *Bot) sell(trade config.Trade, price float64) error {
	wallet, err := me.Binance.SymbolBalance(trade.GetSymbol())
	if err != nil {
		return err
	}

	if wallet == 0 {
		return nil
	}

	quantity := math.Floor(wallet / price)
	if quantity == 0 {
		return nil
	}

	order, err := me.Binance.CreateOrder(trade.GetSymbol(), trade.BuyWith(),
		"SELL", quantity, price)
	if err != nil {
		return err
	}

	me.Log.Debug().Msgf("[%s] sell for %d", trade.Symbol, price)
	me.Notify.SendMessage(notify.NewMessage(fmt.Sprintf("Order to sell created %v", order)))

	return nil
}

func welcomeMessage(config config.Config) string {
	msg := fmt.Sprintf("%s started! :money_mouth_face:\n\n", config.Name)
	for i, trade := range config.Trades {
		msg += fmt.Sprintf("> **%d:** %s\n", i+1, trade.Symbol)
		msg += fmt.Sprintf("> **Interval:** %s\n", trade.Interval)
		msg += fmt.Sprintf("> **Purchase price:** %.2f\n", trade.BuyPrice)
		msg += fmt.Sprintf("> **Sale Price:** %.2f\n", trade.SellPrice)
	}

	return msg
}
