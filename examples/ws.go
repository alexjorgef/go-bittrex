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
	ch := make(chan bittrex.StreamTrade)
	errCh := make(chan error)
	stopCh := make(chan bool)
	go func() { errCh <- client.SubscribeTradeUpdates("BTC-USD", ch, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("BTC-USDT", ch, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("ETH-USD", ch, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("ETH-EUR", ch, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("ETH-USDC", ch, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("ETH-USDT", ch, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("ETH-BTC", ch, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("ADA-BTC", ch, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("ADA-USD", ch, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("ADA-USDT", ch, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("ADA-EUR", ch, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("ADA-ETH", ch, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("LINK-USD", ch, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("AAVE-USD", ch, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("DOT-USD", ch, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("DOT-USDT", ch, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("DOT-EUR", ch, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("DOT-BTC", ch, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("DOT-ETH", ch, stopCh) }()
	go func() { errCh <- client.SubscribeTradeUpdates("DOGE-USDT", ch, stopCh) }()

	// Read from channels and stop after X time
	fmt.Printf("StreamTrade (TakerSide, Symbol, Quantity, Rate):\n")
	for start := time.Now(); time.Since(start) < (35 * time.Second); {
		select {
		case trade := <-ch:
			fmt.Printf("\t%s\t%s %s at %s\n", trade.TakerSide, trade.Symbol, trade.Quantity.String(), trade.Rate.String())
		case err := <-errCh:
			fmt.Printf("\t%+v\n", err)
		}
	}

	return 0
}
