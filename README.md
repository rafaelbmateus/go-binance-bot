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
    interval: "10s"
    buyPrice: 34322.0
    sellPrice: 50640.0
    limit: 100
  - symbol: "LOKA/USDT"
    interval: "1m"
    buyPrice: 1.51
    sellPrice: 3.08
    limit: 10
```

* symbol: Symbol name to trade - `string`
* interval: Interval to the next trade - `time.Duration`
* buyPrice: Create a buy order when the price is below - `float64`
* sellPrice: Create a sell order when the price is high - `float64`
* limit: Limit of USDT that will negotiate - `float64`

# How to run?

To run the bot on your computer you need to have docker and compose installed.

If you are running in the first time you need create `.env` and `config.yaml` files.
Use this command to generate the files in the correct location:

```console
make config
```

Then change your environment vars in `.env`
with your binance api credentials.

Another file you should customize is the `config.yaml`
with the cryptocurrencies you want to trade.

Finally, to run the project, exec:

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
