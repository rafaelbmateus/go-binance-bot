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

	me.Notify.SendMessage(notify.NewMessage(me.welcomeMessage(*me.Config)))
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

// create an order to buy.
func (me *Bot) buy(trade config.Trade, price float64) error {
	wallet, err := me.Binance.SymbolBalance(trade.GetSymbol())
	if err != nil {
		return err
	}

	if math.Floor(wallet) > 0 {
		me.Log.Debug().Msgf("[%s] already bought this coin", trade.Symbol)
		return nil
	}

	wallet, err = me.Binance.SymbolBalance(trade.BuyWith())
	if err != nil {
		return err
	}

	if wallet == 0 || wallet < trade.Limit {
		me.Log.Debug().Msgf("[%s] don't have %s enough on wallet to buy",
			trade.Symbol, trade.BuyWith())
		return nil
	}

	quantity := math.Floor(trade.Limit / price)
	if quantity == 0 {
		me.Log.Debug().Msgf("[%s] don't have quantity enough to buy",
			trade.Symbol)
		return nil
	}

	order, err := me.Binance.CreateOrder(trade.GetSymbol(),
		trade.BuyWith(), "BUY", quantity, price)
	if err != nil {
		return err
	}

	me.Log.Info().Msgf("[%s] order to buy for %d", trade.Symbol, price)
	me.Notify.SendMessage(notify.NewMessage(
		fmt.Sprintf("[%s] Order to buy created for %s, quantity: %.2f, wallet: %.2f",
			trade.Symbol, order.Price, quantity, wallet)))

	return nil
}

// create an order to sell.
func (me *Bot) sell(trade config.Trade, price float64) error {
	wallet, err := me.Binance.SymbolBalance(trade.GetSymbol())
	if err != nil {
		return err
	}

	if wallet == 0 {
		me.Log.Debug().Msgf("[%s] don't have enough on wallet to sell",
			trade.Symbol)
		return nil
	}

	quantity := math.Floor(wallet / price)
	if quantity == 0 {
		me.Log.Debug().Msgf("[%s] don't have quantity enough to sell",
			trade.Symbol)
		return nil
	}

	order, err := me.Binance.CreateOrder(trade.GetSymbol(),
		trade.BuyWith(), "SELL", quantity, price)
	if err != nil {
		return err
	}

	me.Log.Debug().Msgf("[%s] sell for %d", trade.Symbol, price)
	me.Notify.SendMessage(notify.NewMessage(
		fmt.Sprintf("[%s] Order to buy sell for %s, quantity: %.2f, wallet: %.2f",
			trade.Symbol, order.Price, quantity, wallet)))

	return nil
}

// format welcome message notification.
func (me *Bot) welcomeMessage(config config.Config) string {
	msg := fmt.Sprintf("%s started! :money_mouth_face:\n\n", config.Name)
	for i, trade := range config.Trades {
		currentPrice, _ := me.Binance.SymbolPrice(trade.GetSymbol(), trade.BuyWith(), "1m")
		msg += fmt.Sprintf("> *%d: %s*\n", i+1, trade.Symbol)
		msg += fmt.Sprintf("> *Interval:* %s\n", trade.Interval)
		msg += fmt.Sprintf("> *Current price:* %.2f\n", currentPrice)
		msg += fmt.Sprintf("> *Buy price:* %.2f\n", trade.BuyPrice)
		msg += fmt.Sprintf("> *Sell Price:* %.2f\n", trade.SellPrice)
		msg += fmt.Sprintf("> *Profit:* %.2f%%\n\n", (trade.SellPrice-trade.BuyPrice)/(trade.SellPrice)*100)
	}

	return msg
}
