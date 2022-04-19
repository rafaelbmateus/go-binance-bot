# Go Binance Bot

This project runs trades automatically and send notifications.
Basically you run the container and let the bot working for you.

Simple like:
> Buy when it's cheap! And sell when it's expensive.

## Binance API

Fist, you need a binance api credentials.

If you don't have a binance api read
(https://www.binance.com/pt-BR/support/faq/360002502072)[How to create binance api]

Save your secrets in a safe place then we will put it in the `.env`

## Bot Config

The bot configuration is provided by yaml file that
has the following trade parameters.
In this file, put the symbol price you want to buy and sell.

```yaml
name: "My binance bot"
trades:
  - symbol: "BTC/USDT"
    interval: "1m"
    buyPrice: 34322.0
    sellPrice: 50640.0
  - symbol: "LOKA/USDT"
    interval: "30s"
    buyPrice: 1.51
    sellPrice: 3.08
```

* symbol: Symbol name to trade - `string`
* interval: Interval to the next trade - `time.Duration`
* buyPrice: Create a buy order when the price is below - `float64`
* sellPrice: Create a sell order when the price is high - `float64`

# How to run?

To run the bot on your computer, you need to have docker and compose installed.

First, create a `.env` file in root folder with your binance api:

```console
BINANCE_API_KEY=<YOUR_API_KEY_HERE>
BINANCE_API_SECRET=<YOUR_API_SECRET_HERE>
SLACK_WEBHOOK_URL=<YOUR_SLACK_WEBHOOK_URL_HERE>
```

See the project structure to know where to put the `.env` file with your keys:

```console
├── binance/
├── bot/
├── build/
├── config/
├── notify/
├── .env  <-----  # put your env file here!
├── config.yaml
├── go.mod
├── go.sum
├── main.go
├── Makefile
└── README.md
```

Finally to run the project, exec:

```console
make up logs

...
app_1  | {"level":"info","time":"2022-03-28T13:19:31Z","message":"bot is running"}
app_1  | {"level":"debug","time":"2022-03-28T13:19:31Z","message":"monitor {BTC/USDT 1m0s 44630 44640}"}
app_1  | {"level":"debug","time":"2022-03-28T13:19:31Z","message":"monitor {LOKA/USDT 1m0s 2.295 2.34}"}
app_1  | {"level":"debug","time":"2022-03-28T13:19:35Z","message":"current price of BTC/USDT is 47282"}
app_1  | {"level":"debug","time":"2022-03-28T13:19:35Z","message":"time to SELL price 47282"}
...
```

To stop the bot and remove container, exec:

```console
make clean
```
