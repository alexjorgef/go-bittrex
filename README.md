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

Check more advanced examples [here](examples/).

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

	ch := make(chan bittrex.StreamTrade)
	errCh := make(chan error)
	stopCh := make(chan bool)

	go func() { errCh <- client.SubscribeTradeUpdates("BTC-USD", ch, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("ETH-USD", ch, stopCh) }()

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
    - [ ] Account ([#13][i13])
    - [ ] Addresses ([#13][i13])
    - [ ] Balances ([#13][i13])
    - [ ] Batch ([#13][i13])
    - [ ] ConditionalOrders
    - [X] Currencies
    - [ ] Deposits ([#13][i13])
    - [ ] Executions ([#13][i13])
    - [ ] FundsTransferMethods ([#13][i13])
    - [X] Markets
    - [ ] Orders ([#13][i13])
	- [X] Ping
    - [ ] Subaccounts ([#13][i13])
    - [ ] Transfers ([#13][i13])
    - [ ] Withdrawals ([#13][i13])
- [ ] Websocket API
    - [ ] Authenticate ([#13][i13])
    - [ ] IsAuthenticated ([#13][i13])
    - [X] Subscribe
    - [ ] Unsubscribe
	- [ ] Streams
		- [ ] Balance ([#13][i13])
		- [ ] Candle
		- [ ] Conditional Order ([#13][i13])
		- [ ] Deposit ([#13][i13])
		- [ ] Execution ([#13][i13])
		- [X] Heartbeat
		- [ ] Market Summaries
		- [ ] Market Summary
		- [ ] Order ([#13][i13])
		- [ ] Orderbook
		- [ ] Tickers
		- [ ] Ticker
		- [X] Trade

## References

This repository is a cleaned & updated version of [toorop/go-bittrex](https://github.com/toorop/go-bittrex) repo (inspired from [alexeykaravan/go-bittrex-v3](https://github.com/alexeykaravan/go-bittrex-v3) fork.

[i13]: https://github.com/alexjorgef/go-bittrex/issues/13