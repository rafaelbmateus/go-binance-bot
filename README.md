# go-binance-bot

Go Bincance Bot run trades automatically.
Put the buy or sell price and let the bot work for you.

Simple like:
> Buy when it's cheap or sell when it's expensive.

## Setup Environment

Create a env file in root folder with your
[binance api key](https://www.binance.com/en/support/faq/360002502072)
like this:

```bash
BINANCE_API_KEY=<YOUR_API_KEY_HERE>
BINANCE_API_SECRET=<YOUR_API_SECRET_HERE>
SLACK_WEBHOOK_URL=<YOUR_SLACK_WEBHOOK_URL_HERE>
```

## Bot Config

The bot configuration is provided by yaml file that
has the following trade parameters.
In this file, put the symbol price you want to buy and sell.

```yaml
name: "Binance Bot"
trades:
  - symbol: "BTC/USDT"
    interval: "1m"
    buyPrice: 44630.0
    sellPrice: 44640.0
  - symbol: "LOKA/USDT"
    interval: "1m"
    buyPrice: 2.295
    sellPrice: 2.34
```

## How to use?

To up the containers run this make command:

```bash
make up
```
