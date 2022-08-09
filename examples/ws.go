package main

import (
	"fmt"
	"time"

	"github.com/alexjorgef/go-bittrex/bittrex"
)

func realMainWs() int {
	client := bittrex.New(API_KEY, API_SECRET)

	errCh := make(chan error)
	stopCh := make(chan bool)

	// Subscribe to candle stream
	chCandle := make(chan bittrex.Candle)
	go func() { errCh <- client.SubscribeCandleUpdates("BTC-USD", chCandle, stopCh) }()
	go func() { errCh <- client.SubscribeCandleUpdates("ETH-USD", chCandle, stopCh) }()
	go func() { errCh <- client.SubscribeCandleUpdates("ADA-USD", chCandle, stopCh) }()

	fmt.Printf("Candle (SubscribeCandleUpdates):\n")
	for start := time.Now(); time.Since(start) < (5 * time.Second); {
		select {
		case candle := <-chCandle:
			fmt.Printf("\t%+v\n", candle)
		case err := <-errCh:
			fmt.Printf("\t%+v\n", err)
		}
	}

	// Subscribe to ordebook stream
	chOrderbook := make(chan bittrex.OrderBook)
	go func() { errCh <- client.SubscribeOrderbookUpdates("ADA-USD", chOrderbook, stopCh) }()

	fmt.Printf("OrderBook (SubscribeOrderbookUpdates):\n")
	for start := time.Now(); time.Since(start) < (5 * time.Second); {
		select {
		case orderbook := <-chOrderbook:
			fmt.Printf("\t%+v\n", orderbook)
		case err := <-errCh:
			fmt.Printf("\t%+v\n", err)
		}
	}

	// Subscribe to ticker stream
	chTickers := make(chan bittrex.Ticker)
	go func() { errCh <- client.SubscribeTickersUpdates(chTickers, stopCh) }()
	go func() { errCh <- client.SubscribeTickersUpdates(chTickers, stopCh) }()
	go func() { errCh <- client.SubscribeTickersUpdates(chTickers, stopCh) }()

	fmt.Printf("Ticker (SubscribeTickersUpdates):\n")
	for start := time.Now(); time.Since(start) < (5 * time.Second); {
		select {
		case ticker := <-chTickers:
			fmt.Printf("\t%+v\n", ticker)
		case err := <-errCh:
			fmt.Printf("\t%+v\n", err)
		}
	}

	// Subscribe to ticker stream
	chTicker := make(chan bittrex.Ticker)
	go func() { errCh <- client.SubscribeTickerUpdates("BTC-USD", chTicker, stopCh) }()
	go func() { errCh <- client.SubscribeTickerUpdates("ETH-USD", chTicker, stopCh) }()
	go func() { errCh <- client.SubscribeTickerUpdates("ADA-USD", chTicker, stopCh) }()

	fmt.Printf("Ticker (SubscribeTickerUpdates):\n")
	for start := time.Now(); time.Since(start) < (5 * time.Second); {
		select {
		case ticker := <-chTicker:
			fmt.Printf("\t%+v\n", ticker)
		case err := <-errCh:
			fmt.Printf("\t%+v\n", err)
		}
	}

	// Subscribe to trade stream
	chTrade := make(chan bittrex.Trade)
	go func() { errCh <- client.SubscribeTradeUpdates("BTC-USD", chTrade, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("BTC-USDT", chTrade, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("ETH-USD", chTrade, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("ETH-EUR", chTrade, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("ETH-USDC", chTrade, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("ETH-USDT", chTrade, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("ETH-BTC", chTrade, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("ADA-BTC", chTrade, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("ADA-USD", chTrade, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("ADA-USDT", chTrade, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("ADA-EUR", chTrade, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("ADA-ETH", chTrade, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("LINK-USD", chTrade, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("AAVE-USD", chTrade, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("DOT-USD", chTrade, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("DOT-USDT", chTrade, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("DOT-EUR", chTrade, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("DOT-BTC", chTrade, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("DOT-ETH", chTrade, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("DOGE-USDT", chTrade, stopCh) }()

	fmt.Printf("Trade (SubscribeTradeUpdates):\n")
	for start := time.Now(); time.Since(start) < (5 * time.Second); {
		select {
		case trade := <-chTrade:
			fmt.Printf("\t%+v\n", trade)
		case err := <-errCh:
			fmt.Printf("\t%+v\n", err)
		}
	}

	return 0
}
