# go-bittrex

[![Test](https://github.com/alexjorgef/go-bittrex/workflows/Test/badge.svg)](https://github.com/alexjorgef/go-bittrex/actions?query=workflow%3ATest)
[![Lint](https://github.com/alexjorgef/go-bittrex/workflows/Lint/badge.svg)](https://github.com/alexjorgef/go-bittrex/actions?query=workflow%3ALint)
[![codecov](https://codecov.io/gh/alexjorgef/go-bittrex/branch/main/graph/badge.svg)](https://codecov.io/gh/alexjorgef/go-bittrex)
[![Go Report Card](https://goreportcard.com/badge/github.com/alexjorgef/go-bittrex)](https://goreportcard.com/report/github.com/alexjorgef/go-bittrex)
[![GoDoc](https://godoc.org/github.com/alexjorgef/go-bittrex?status.svg)](https://godoc.org/github.com/alexjorgef/go-bittrex)

go-bittrex is a Go client library for accessing the [Bittrex API](https://bittrex.github.io/api).

# Install

```console
go get github.com/alexjorgef/go-bittrex
```

## Quick Start

```go
package main

import (
	"fmt"
	"log"

	"github.com/alexjorgef/go-bittrex/bittrex"
)

func main() {
	client := bittrex.New("", "")
	currency, err := client.GetCurrency("ETH")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%+v\n", currency)
}
```

## Examples

### REST API

```go
package main

import (
	"fmt"
	"log"

	"github.com/alexjorgef/go-bittrex/bittrex"
)

func main() {
	client := bittrex.New("", "")
	currency, err := client.GetCurrency("ETH")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%+v\n", currency)
}
```

### Websocket

```go
package main

import (
	"fmt"
	"log"

	"github.com/alexjorgef/go-bittrex/bittrex"
)

func main() {
	client := bittrex.New("", "")
	ch := make(chan bittrex.Trade)
	errCh := make(chan error)
	stopCh := make(chan bool)
	go func() {
		errCh <- bt.SubscribeTradeUpdates("BTC-USD", ch, stopCh)
	}()
	select {
	case trade := <-ch:
		fmt.Printf("%+v\n", trade)
	case err := <-errCh:
		fmt.Printf("%+v\n", err)
	}
}
```

## Todos

- [ ] REST API
    - [ ] Account
		- [ ] GET /account
		- [ ] GET /account/fees/fiat
		- [ ] GET /account/fees/fiat/{currencySymbol}
		- [ ] GET /account/fees/trading
		- [ ] GET /account/fees/trading/{marketSymbol}
		- [ ] GET /account/volume
		- [ ] GET /account/permissions/markets
		- [ ] GET /account/permissions/markets/{marketSymbol}
		- [ ] GET /account/permissions/currencies
		- [ ] GET /account/permissions/cu
    - [ ] Addresses
    	- [ ] GET /addresses
    	- [ ] POST /addresses
    	- [ ] GET /addresses/{currencySymb
    - [ ] Balances
    	- [ ] GET /balances
    	- [ ] HEAD /balances
    	- [ ] GET /balances/{currencySymbo
    - [ ] Batch
		- [ ] POST /batch
    - [ ] ConditionalOrders
		- [ ] GET /conditional-orders/{conditionalOrderId}
		- [ ] DELETE /conditional-orders/{conditionalOrderId}
		- [ ] GET /conditional-orders/closed
		- [ ] GET /conditional-orders/open
		- [ ] HEAD /conditional-orders/open
		- [ ] POST /conditional-orders
    - [ ] Currencies
		- [X] GET /currencies
		- [X] GET /currencies/{symbol}
    - [ ] Deposits
		- [ ] GET /deposits/open
		- [ ] HEAD /deposits/open
		- [ ] GET /deposits/closed
		- [ ] GET /deposits/ByTxId/{txId}
		- [ ] GET /deposits/{depositId}
    - [ ] Executions
		- [ ] GET /executions/{executionId}
		- [ ] GET /executions
		- [ ] GET /executions/last-id
		- [ ] HEAD /executions/last-id
    - [ ] FundsTransferMethods
		- [ ] GET /funds-transfer-methods/{fundsTransferMethodId}
    - [ ] Markets
		- [X] GET /markets
		- [X] GET /markets/summaries
		- [ ] HEAD /markets/summaries
		- [X] GET /markets/tickers
		- [ ] HEAD /markets/tickers
		- [X] GET /markets/{marketSymbol}/ticker
		- [ ] GET /markets/{marketSymbol}
		- [X] GET /markets/{marketSymbol}/summary
		- [X] GET /markets/{marketSymbol}/orderbook
		- [ ] HEAD /markets/{marketSymbol}/orderbook
		- [X] GET /markets/{marketSymbol}/trades
		- [ ] HEAD /markets/{marketSymbol}/trade
		- [ ] GET /markets/{marketSymbol}/candles/{candleType}/{candleInterval}/recent
		- [ ] HEAD /markets/{marketSymbol}/candles/{candleType}/{candleInterval}/recent
		- [ ] GET /markets/{marketSymbol}/candles/{candleType}/{candleInterval}/historical/{year}/{month}/{day}
    - [ ] Orders
		- [ ] GET /orders/closed
		- [ ] GET /orders/open
		- [ ] DELETE /orders/open
		- [ ] HEAD /orders/open
		- [ ] GET /orders/{orderId}
		- [ ] DELETE /orders/{orderId}
		- [ ] GET /orders/{orderId}/executions
		- [ ] POST /orders
	- [ ] Ping
		- [ ] GET /ping
    - [ ] Subaccounts
		- [ ] GET /subaccounts
		- [ ] POST /subaccounts
		- [ ] GET /subaccounts/{subaccountId}
		- [ ] GET /subaccounts/withdrawals/open
		- [ ] GET /subaccounts/withdrawals/closed
		- [ ] GET /subaccounts/deposits/open
		- [ ] GET /subaccounts/deposits/cl
    - [ ] Transfers
		- [ ] GET /transfers/sent
		- [ ] GET /transfers/received
		- [ ] GET /transfers/{transferId}
		- [ ] POST /transfers
    - [ ] Withdrawals
		- [ ] GET /withdrawals/open
		- [ ] GET /withdrawals/closed
		- [ ] GET /withdrawals/ByTxId/{txId}
		- [ ] GET /withdrawals/{withdrawalId}
		- [ ] DELETE /withdrawals/{withdrawalId}
		- [ ] POST /withdrawals
		- [ ] GET /withdrawals/allowed-addr
- [ ] Websocket API
    - [ ] Authenticate
    - [ ] IsAuthenticated
    - [ ] Subscribe
    - [ ] Unsubscribe
	- [ ] Streams
		- [ ] Balance
		- [ ] Candle
		- [ ] Conditional Order
		- [ ] Deposit
		- [ ] Execution
		- [ ] Heartbeat
		- [ ] Market Summaries
		- [ ] Market Summary
		- [ ] Order
		- [ ] Orderbook
		- [ ] Tickers
		- [ ] Ticker
		- [x] Trade

## References

This repository is a cleaned & updated version of [toorop/go-bittrex](https://github.com/toorop/go-bittrex) repo (inspired from [alexeykaravan/go-bittrex-v3](https://github.com/alexeykaravan/go-bittrex-v3) fork.