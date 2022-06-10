# Go Binance Bot

This project runs trades automatically and send notifications.
Basically you run the container and let the bot working for you.

Simple like:
> Buy when it's cheap! And sell when it's expensive.

## Binance API

Fist, you need a binance api credentials.

If you don't have a binance api read
[How to create binance api](https://www.binance.com/pt-BR/support/faq/360002502072).

Save your secrets in a safe place then we will put it in the `.env`.

# How to run?

To run the bot on your computer you need to have
[docker](https://docker.com) and [compose](https://docs.docker.com/compose) installed.

Make sure you created the `.env` and `config.yaml` files:

```console
cp config-example.yaml config.yaml
cp example.env .env
```

## Config file

The bot configuration is provided by yaml file that
has the following trade parameters `config.yaml`.
In this file, put the symbol price you want to buy and sell.

```yaml
name: "go-binance-bot"

trades:
  - symbol: "BTC/USDT"
    interval: "10s"
    amount: 0.0004
    rsi_buy: 30
    rsi_sell: 70
    rsi_limit: 14
    rsi_interval: "15m"
  - symbol: "BNB/USDT"
    interval: "10s"
    amount: 0.05
    rsi_buy: 30
    rsi_sell: 70
    rsi_limit: 14
    rsi_interval: "15m"
  - symbol: "LOKA/USDT"
    interval: "10s"
    amount: 20
    rsi_buy: 30
    rsi_sell: 70
    rsi_limit: 14
    rsi_interval: "15m"
  - symbol: "ETH/USDT"
    interval: "10s"
    amount: 0 # set zero to test the logic.
    rsi_buy: 30
    rsi_sell: 70
    rsi_limit: 14
    rsi_interval: "15m"
```

* symbol: Symbol name to trade - `string`
* interval: Interval to the next trade - `time.Duration`
* buyPrice: Create a buy order when the price is below - `float64`
* sellPrice: Create a sell order when the price is high - `float64`
* limit: Limit of USDT that will negotiate - `float64`

## Env file

Change your environment vars in `.env` file
with your binance api credentials.

```
BINANCE_API_KEY=xXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxX
BINANCE_API_SECRET=xXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxXxX
SLACK_WEBHOOK_URL=https://hooks.slack.com/services/ABABABABA/BABABABABAB/BLABLABLABLABLABLABLABLA
```

## Running the bot

Finally start the containers to monitor
the coins and open buy or sell order:

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

To stop the bot and remove container, execute:

```console
make clean
```

## 📫 Contributing

To contribute to the project, follow these steps:

1. Clone the repository: `git clone git@github.com:rafaelbmateus/go-binance-bot.git`
2. Create a feature branch: `git switch -c feature-a`
3. Make changes and confirm (try using [conventional commits](https://www.conventionalcommits.org)): `git commit -m 'feat: new bot feature'`
4. Push the feature branch: `git push origin feature-a`
5. Create a [pull request](https://help.github.com/en/github/collaborating-with-issues-and-pull-requests/creating-a-pull-request)
6. Get reviews from other users
7. Merge to `main` branch (we encourage using commit squash)
8. Remove the branch merged.
