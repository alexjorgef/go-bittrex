package main

import (
	"fmt"
	"time"

	"github.com/alexjorgef/go-bittrex/bittrex"
)

func realMainWs() int {
	// Bittrex client
	client := bittrex.New(API_KEY, API_SECRET)

	// Open channels and start a websocket connection
	chTrades := make(chan bittrex.Trade)
	errCh := make(chan error)
	stopCh := make(chan bool)
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

	// Read from channels and stop after X time
	fmt.Printf("StreamTrade (TakerSide, Symbol, Quantity, Rate):\n")
	for start := time.Now(); time.Since(start) < (10 * time.Second); {
		select {
		case trade := <-chTrades:
			fmt.Printf("\t%s\t%s %s at %s\n", trade.TakerSide, trade.Symbol, trade.Quantity.String(), trade.Rate.String())
		case err := <-errCh:
			fmt.Printf("\t%+v\n", err)
		}
	}

	// Open channels and start a websocket connection
	chTickers := make(chan bittrex.Ticker)
	go func() { errCh <- client.SubscribeTickerUpdates("BTC-USD", chTickers, stopCh) }()
	go func() { errCh <- client.SubscribeTickerUpdates("BTC-USDT", chTickers, stopCh) }()
	go func() { errCh <- client.SubscribeTickerUpdates("ETH-USD", chTickers, stopCh) }()
	go func() { errCh <- client.SubscribeTickerUpdates("ETH-EUR", chTickers, stopCh) }()
	go func() { errCh <- client.SubscribeTickerUpdates("ETH-USDC", chTickers, stopCh) }()
	go func() { errCh <- client.SubscribeTickerUpdates("ETH-USDT", chTickers, stopCh) }()
	go func() { errCh <- client.SubscribeTickerUpdates("ETH-BTC", chTickers, stopCh) }()
	go func() { errCh <- client.SubscribeTickerUpdates("ADA-BTC", chTickers, stopCh) }()
	go func() { errCh <- client.SubscribeTickerUpdates("ADA-USD", chTickers, stopCh) }()
	go func() { errCh <- client.SubscribeTickerUpdates("ADA-USDT", chTickers, stopCh) }()
	go func() { errCh <- client.SubscribeTickerUpdates("ADA-EUR", chTickers, stopCh) }()
	go func() { errCh <- client.SubscribeTickerUpdates("ADA-ETH", chTickers, stopCh) }()
	go func() { errCh <- client.SubscribeTickerUpdates("LINK-USD", chTickers, stopCh) }()
	go func() { errCh <- client.SubscribeTickerUpdates("AAVE-USD", chTickers, stopCh) }()
	go func() { errCh <- client.SubscribeTickerUpdates("DOT-USD", chTickers, stopCh) }()
	go func() { errCh <- client.SubscribeTickerUpdates("DOT-USDT", chTickers, stopCh) }()
	go func() { errCh <- client.SubscribeTickerUpdates("DOT-EUR", chTickers, stopCh) }()
	go func() { errCh <- client.SubscribeTickerUpdates("DOT-BTC", chTickers, stopCh) }()
	go func() { errCh <- client.SubscribeTickerUpdates("DOT-ETH", chTickers, stopCh) }()
	go func() { errCh <- client.SubscribeTickerUpdates("DOGE-USDT", chTickers, stopCh) }()

	// Read from channels and stop after X time
	fmt.Printf("StreamTicker (Symbol, Ask, Bid, LastTradeRate):\n")
	for start := time.Now(); time.Since(start) < (35 * time.Second); {
		select {
		case ticker := <-chTickers:
			fmt.Printf("\t%s %s %s %s\n", ticker.Symbol, ticker.AskRate.String(), ticker.BidRate.String(), ticker.LastTradeRate.String())
		case err := <-errCh:
			fmt.Printf("\t%+v\n", err)
		}
	}

	return 0
}
