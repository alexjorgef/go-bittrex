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

	// Subscribe to ordebook stream
	chOrderbook := make(chan bittrex.OrderBook)
	go func() { errCh <- client.SubscribeOrderbookUpdates("ADA-USD", chOrderbook, stopCh) }()

	fmt.Printf("OrderBook (Symbol, Depth, Bid, Ask):\n")
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
	go func() { errCh <- client.SubscribeTickerUpdates("BTC-USD", chTickers, stopCh) }()
	go func() { errCh <- client.SubscribeTickerUpdates("ETH-USD", chTickers, stopCh) }()
	go func() { errCh <- client.SubscribeTickerUpdates("ADA-USD", chTickers, stopCh) }()

	fmt.Printf("Ticker (Symbol, LastTradeRate, BitRate, AskRate):\n")
	for start := time.Now(); time.Since(start) < (5 * time.Second); {
		select {
		case ticker := <-chTickers:
			// fmt.Printf("\t%s %s %s %s\n", ticker.Symbol, ticker.AskRate.String(), ticker.BidRate.String(), ticker.LastTradeRate.String())
			fmt.Printf("\t%+v\n", ticker)
		case err := <-errCh:
			fmt.Printf("\t%+v\n", err)
		}
	}

	// Subscribe to trade stream
	chTrades := make(chan bittrex.Trade)
	go func() { errCh <- client.SubscribeTradeUpdates("BTC-USD", chTrades, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("BTC-USDT", chTrades, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("ETH-USD", chTrades, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("ETH-EUR", chTrades, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("ETH-USDC", chTrades, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("ETH-USDT", chTrades, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("ETH-BTC", chTrades, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("ADA-BTC", chTrades, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("ADA-USD", chTrades, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("ADA-USDT", chTrades, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("ADA-EUR", chTrades, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("ADA-ETH", chTrades, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("LINK-USD", chTrades, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("AAVE-USD", chTrades, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("DOT-USD", chTrades, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("DOT-USDT", chTrades, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("DOT-EUR", chTrades, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("DOT-BTC", chTrades, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("DOT-ETH", chTrades, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("DOGE-USDT", chTrades, stopCh) }()

	fmt.Printf("Trade (Symbol, ID, ExecutedAt, Quantity, Rate, TakerSide):\n")
	for start := time.Now(); time.Since(start) < (5 * time.Second); {
		select {
		case trade := <-chTrades:
			fmt.Printf("\t%+v\n", trade)
		case err := <-errCh:
			fmt.Printf("\t%+v\n", err)
		}
	}

	return 0
}
