package bot

import (
	"context"
	"fmt"
	"time"

	"github.com/rafaelbmateus/binance-bot/binance"
	"github.com/rafaelbmateus/binance-bot/config"
	"github.com/rafaelbmateus/binance-bot/notify"
	"github.com/rafaelbmateus/binance-bot/strategy"
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
			me.Log.Error().Msgf("[%s] watchdog sleep error %q", trade.Symbol, err)
		}

		time.Sleep(trade.Interval)
	}
}

// Monitor symbol to buy or sell.
func (me *Bot) Monitor(trade config.Trade) error {
	price, err := me.Binance.SymbolPrice(trade.GetAsset(), trade.RSIInterval)
	if err != nil {
		return err
	}

	wallet, err := me.Binance.SymbolBalance(trade.GetSymbol())
	if err != nil {
		return err
	}

	rsi, err := me.calculateRSI(trade)
	if err != nil {
		return err
	}

	me.Log.Info().Msgf("[%s] wallet: %.5f amount: %.5f price: %.5f rsi: %.2f",
		trade.Symbol, wallet, trade.Amount, price, rsi)

	if rsi <= trade.RSIBuy {
		if trade.Amount == 0 {
			return nil
		}

		err := me.buy(trade, wallet, price, rsi)
		if err != nil {
			return err
		}

		return nil
	}

	if rsi >= trade.RSISell {
		if trade.Amount == 0 {
			return nil
		}

		err := me.sell(trade, wallet, price, rsi)
		if err != nil {
			return err
		}

		return nil
	}

	return nil
}

func (me *Bot) calculateRSI(trade config.Trade) (float64, error) {
	klines, err := me.Binance.Client.NewKlinesService().Symbol(
		fmt.Sprintf("%s%s", trade.GetSymbol(), trade.BuyWith())).
		Interval(trade.RSIInterval).Limit(trade.RSILimit).Do(*me.Context)
	if err != nil {
		return 0, err
	}

	rsi, err := strategy.CalculateRSI(klines)
	if err != nil {
		return 0, err
	}

	return rsi, nil
}

// create an order to buy.
func (me *Bot) buy(trade config.Trade, wallet, price, rsi float64) error {
	if wallet >= trade.Amount {
		me.Log.Debug().Msgf("[%s] already bought, %.5f on wallet", trade.Symbol, wallet)
		return nil
	}

	order, err := me.Binance.CreateOrder(trade.GetSymbol(),
		trade.BuyWith(), "BUY", trade.Amount, price)
	if err != nil {
		return err
	}

	me.Notify.SendMessage(notify.NewMessage(
		fmt.Sprintf("*[%s] Order to buy created!*\nPrice: %s\nWallet: %.5f\nAmount: %.5f\nRSI: %.2f",
			trade.Symbol, order.Price, wallet, trade.Amount, rsi)))

	return nil
}

// create an order to sell.
func (me *Bot) sell(trade config.Trade, wallet, price, rsi float64) error {
	if wallet < trade.Amount {
		me.Log.Debug().Msgf("[%s] don't have enough to sell, %.5f on wallet", trade.Symbol, wallet)
		return nil
	}

	order, err := me.Binance.CreateOrder(trade.GetSymbol(),
		trade.BuyWith(), "SELL", trade.Amount, price)
	if err != nil {
		return err
	}

	me.Notify.SendMessage(notify.NewMessage(
		fmt.Sprintf("*[%s] Order to sell created!*\nPrice: %s\nWallet: %.5f\nAmount: %.5f\nRSI: %.2f",
			trade.Symbol, order.Price, wallet, trade.Amount, rsi)))

	return nil
}

// check rules before buy.
func (me *Bot) rules(trade *config.Trade) error {
	exchange, err := me.Binance.SymbolExchangeInfo(trade.Symbol)
	if err != nil {
		return err
	}

	minNotional := exchange.Symbols[0].MinNotionalFilter().MinNotional
	minQuantity := exchange.Symbols[0].LotSizeFilter().MinQuantity
	me.Log.Debug().Msgf("[%s] minNotional=%.2f minQuantity=%.2f",
		trade.Symbol, minNotional, minQuantity)

	return nil
}

// format welcome message notification.
func (me *Bot) welcomeMessage(config config.Config) string {
	msg := fmt.Sprintf("%s started! :money_mouth_face:\n\n", config.Name)
	for i, trade := range config.Trades {
		msg += fmt.Sprintf("> *%d: %s*\n", i+1, trade.Symbol)
		msg += fmt.Sprintf("> *Interval:* %s\n", trade.Interval)
		msg += fmt.Sprintf("> *Amount:* %.5f\n", trade.Amount)
		msg += fmt.Sprintf("> *Buy when RSI is below:* %.2f\n", trade.RSIBuy)
		msg += fmt.Sprintf("> *Sell when RSI is upper:* %.2f\n", trade.RSISell)
		msg += fmt.Sprintf("> *Interval RSI is:* %s\n", trade.RSIInterval)
		msg += fmt.Sprintf("> *History Limit RSI is:* %d\n", trade.RSILimit)
	}

	return msg
}
