package bittrex

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTradeStream_SubscribeCandleUpdates(t *testing.T) {
	client := New("", "")
	ch := make(chan Candle)
	errCh := make(chan error)
	stopCh := make(chan bool)
	go func() { errCh <- client.SubscribeCandleUpdates("ADA-USD", ch, stopCh) }()
	var err error
	var candle Candle
	select {
	case candle = <-ch:
	case err = <-errCh:
	case <-time.NewTicker(3 * time.Minute).C:
		stopCh <- true
		err = errors.New("timeout")
	}
	assert.NoError(t, err)
	assert.NotEmpty(t, candle.MarketSymbol)
	assert.NotEmpty(t, candle.Interval)
	assert.NotEmpty(t, candle.StartsAt)
	assert.NotEmpty(t, candle.Open)
	open, _ := candle.Open.Float64()
	assert.Greater(t, open, float64(0))
	assert.NotEmpty(t, candle.High)
	high, _ := candle.High.Float64()
	assert.Greater(t, high, float64(0))
	assert.NotEmpty(t, candle.Low)
	low, _ := candle.Low.Float64()
	assert.Greater(t, low, float64(0))
	assert.NotEmpty(t, candle.Close)
	close, _ := candle.Close.Float64()
	assert.Greater(t, close, float64(0))
}

func TestTradeStream_SubscribeMarketSummariesUpdates(t *testing.T) {
	client := New("", "")
	ch := make(chan MarketSummary)
	errCh := make(chan error)
	stopCh := make(chan bool)
	go func() { errCh <- client.SubscribeMarketSummariesUpdates(ch, stopCh) }()
	var err error
	var marketSummary MarketSummary
	select {
	case marketSummary = <-ch:
	case err = <-errCh:
	case <-time.NewTicker(3 * time.Minute).C:
		stopCh <- true
		err = errors.New("timeout")
	}
	assert.NoError(t, err)
	assert.NotEmpty(t, marketSummary.Symbol)
	assert.NotEmpty(t, marketSummary.High)
	high, _ := marketSummary.High.Float64()
	assert.Greater(t, high, float64(0))
	assert.NotEmpty(t, marketSummary.Low)
	low, _ := marketSummary.Low.Float64()
	assert.Greater(t, low, float64(0))
	assert.NotEmpty(t, marketSummary.Volume)
	volume, _ := marketSummary.Volume.Float64()
	assert.Greater(t, volume, float64(0))
	assert.NotEmpty(t, marketSummary.QuoteVolume)
	quoteVolume, _ := marketSummary.QuoteVolume.Float64()
	assert.Greater(t, quoteVolume, float64(0))
	assert.NotEmpty(t, marketSummary.PercentChange)
	assert.NotEmpty(t, marketSummary.UpdatedAt)
}

func TestTradeStream_SubscribeMarketSummaryUpdates(t *testing.T) {
	client := New("", "")
	ch := make(chan MarketSummary)
	errCh := make(chan error)
	stopCh := make(chan bool)
	go func() { errCh <- client.SubscribeMarketSummaryUpdates("ADA-USD", ch, stopCh) }()
	var err error
	var marketSummary MarketSummary
	select {
	case marketSummary = <-ch:
	case err = <-errCh:
	case <-time.NewTicker(3 * time.Minute).C:
		stopCh <- true
		err = errors.New("timeout")
	}
	assert.NoError(t, err)
	assert.NotEmpty(t, marketSummary.Symbol)
	assert.NotEmpty(t, marketSummary.High)
	high, _ := marketSummary.High.Float64()
	assert.Greater(t, high, float64(0))
	assert.NotEmpty(t, marketSummary.Low)
	low, _ := marketSummary.Low.Float64()
	assert.Greater(t, low, float64(0))
	assert.NotEmpty(t, marketSummary.Volume)
	volume, _ := marketSummary.Volume.Float64()
	assert.Greater(t, volume, float64(0))
	assert.NotEmpty(t, marketSummary.QuoteVolume)
	quoteVolume, _ := marketSummary.QuoteVolume.Float64()
	assert.Greater(t, quoteVolume, float64(0))
	assert.NotEmpty(t, marketSummary.PercentChange)
	assert.NotEmpty(t, marketSummary.UpdatedAt)
}

func TestTradeStream_SubscribeOrderbookUpdates(t *testing.T) {
	client := New("", "")
	ch := make(chan OrderBook)
	errCh := make(chan error)
	stopCh := make(chan bool)
	go func() { errCh <- client.SubscribeOrderbookUpdates("ADA-USD", ch, stopCh) }()
	var err error
	var orderbook OrderBook
	select {
	case orderbook = <-ch:
	case err = <-errCh:
	case <-time.NewTicker(3 * time.Minute).C:
		stopCh <- true
		err = errors.New("timeout")
	}
	assert.NoError(t, err)
	assert.NotEmpty(t, orderbook.Symbol)
	assert.NotEmpty(t, orderbook.Depth)
	assert.True(t, (len(orderbook.Bid) > 0 || len(orderbook.Ask) > 0))
	if len(orderbook.Bid) > 0 {
		assert.NotEmpty(t, orderbook.Bid)
		bidRate, _ := orderbook.Bid[0].Rate.Float64()
		assert.Greater(t, bidRate, float64(0))
	}
	if len(orderbook.Ask) > 0 {
		assert.NotEmpty(t, orderbook.Ask)
		askRate, _ := orderbook.Ask[0].Rate.Float64()
		assert.Greater(t, askRate, float64(0))
	}
}

func TestTradeStream_SubscribeTickersUpdates(t *testing.T) {
	client := New("", "")
	ch := make(chan Ticker)
	errCh := make(chan error)
	stopCh := make(chan bool)
	go func() { errCh <- client.SubscribeTickersUpdates(ch, stopCh) }()
	var err error
	var ticker Ticker
	select {
	case ticker = <-ch:
	case err = <-errCh:
	case <-time.NewTicker(3 * time.Minute).C:
		stopCh <- true
		err = errors.New("timeout")
	}
	assert.NoError(t, err)
	assert.NotEmpty(t, ticker.Symbol)
	assert.NotEmpty(t, ticker.AskRate)
	assert.NotEmpty(t, ticker.BidRate)
	assert.NotEmpty(t, ticker.LastTradeRate)
	rate, _ := ticker.LastTradeRate.Float64()
	assert.Greater(t, rate, float64(0))
}

func TestTradeStream_SubscribeTickerUpdates(t *testing.T) {
	client := New("", "")
	ch := make(chan Ticker)
	errCh := make(chan error)
	stopCh := make(chan bool)
	go func() { errCh <- client.SubscribeTickerUpdates("BTC-USD", ch, stopCh) }()
	go func() { errCh <- client.SubscribeTickerUpdates("ETH-USD", ch, stopCh) }()
	go func() { errCh <- client.SubscribeTickerUpdates("ADA-USD", ch, stopCh) }()
	var err error
	var ticker Ticker
	select {
	case ticker = <-ch:
	case err = <-errCh:
	case <-time.NewTicker(3 * time.Minute).C:
		stopCh <- true
		err = errors.New("timeout")
	}
	assert.NoError(t, err)
	assert.NotEmpty(t, ticker.Symbol)
	assert.NotEmpty(t, ticker.AskRate)
	assert.NotEmpty(t, ticker.BidRate)
	assert.NotEmpty(t, ticker.LastTradeRate)
	rate, _ := ticker.LastTradeRate.Float64()
	assert.Greater(t, rate, float64(0))
}

func TestTradeStream_SubscribeTradeUpdates(t *testing.T) {
	client := New("", "")
	ch := make(chan Trade)
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
	var err error
	var trade Trade
	select {
	case trade = <-ch:
	case err = <-errCh:
	case <-time.NewTicker(3 * time.Minute).C:
		stopCh <- true
		err = errors.New("timeout")
	}
	assert.NoError(t, err)
	assert.NotEmpty(t, trade.Symbol)
	assert.NotEmpty(t, trade.ID)
	assert.NotEmpty(t, trade.TakerSide)
	assert.NotEmpty(t, trade.ExecutedAt)
	assert.NotEmpty(t, trade.Rate)
	assert.NotEmpty(t, trade.Quantity)
	rate, _ := trade.Rate.Float64()
	assert.Greater(t, rate, float64(0))
}
