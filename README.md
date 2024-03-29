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
	"time"

	"github.com/alexjorgef/go-bittrex/bittrex"
)

func main() {
	// Bittrex client
	client := bittrex.New("", "")

	// Open channels and start a websocket connection that write to it
	ch := make(chan bittrex.Ticker)
	errCh := make(chan error)
	stopCh := make(chan bool)
	go func() { errCh <- client.SubscribeTickerUpdates("BTC-USD", ch, stopCh) }()
	go func() { errCh <- client.SubscribeTickerUpdates("ETH-USD", ch, stopCh) }()
	go func() { errCh <- client.SubscribeTickerUpdates("ADA-USD", ch, stopCh) }()

	// Read from ticker/error channels only 1 time
	select {
	case ticker := <-ch:
		fmt.Printf("%+v\n", ticker)
	case err := <-errCh:
		fmt.Printf("%+v\n", err)
	}

	// Read from channels and stop after 35 seconds
	for start := time.Now(); time.Since(start) < (35 * time.Second); {
		select {
		case ticker := <-ch:
			fmt.Printf("%+v\n", ticker)
		case err := <-errCh:
			fmt.Printf("%+v\n", err)
		}
	}

	// Read from channels infinitely
	for {
		select {
		case ticker := <-ch:
			fmt.Printf("%+v\n", ticker)
		case err := <-errCh:
			fmt.Printf("%+v\n", err)
		}
	}
}
```

## References

This repository is a cleaned & updated version of [toorop/go-bittrex](https://github.com/toorop/go-bittrex) repo (inspired from [alexeykaravan/go-bittrex-v3](https://github.com/alexeykaravan/go-bittrex-v3) fork).

[i13]: https://github.com/alexjorgef/go-bittrex/issues/13