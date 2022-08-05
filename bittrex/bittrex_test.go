package bittrex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Currencies

func TestCurrenciesService_GetCurrencies(t *testing.T) {
	bt := New("", "")
	currencies, err := bt.GetCurrencies()
	assert.NoError(t, err)
	assert.NotEmpty(t, currencies[0].Name)
}

func TestCurrenciesService_GetCurrency(t *testing.T) {
	bt := New("", "")
	currency, err := bt.GetCurrency("BTC")
	assert.NoError(t, err)
	assert.Equal(t, "Bitcoin", currency.Name)
}

// Markets

func TestMarketsService_GetMarkets(t *testing.T) {
	bt := New("", "")
	markets, err := bt.GetMarkets()
	assert.NoError(t, err)
	assert.NotEmpty(t, markets[0].Symbol)
}

func TestMarketsService_GetMarketsSummaries(t *testing.T) {
	bt := New("", "")
	marketSummaries, err := bt.GetMarketsSummaries()
	assert.NoError(t, err)
	assert.NotEmpty(t, marketSummaries[0].Volume)
}

func TestMarketsService_GetMarketsTickers(t *testing.T) {
	bt := New("", "")
	marketTickers, err := bt.GetMarketsTickers()
	assert.NoError(t, err)
	assert.NotEmpty(t, marketTickers[0].Symbol)
}

func TestMarketsService_GetTicker(t *testing.T) {
	bt := New("", "")
	ticker, err := bt.GetTicker("BTC-USD")
	assert.NoError(t, err)
	assert.NotEmpty(t, ticker.Symbol)
}

func TestMarketsService_GetSummary(t *testing.T) {
	bt := New("", "")
	summary, err := bt.GetSummary("BTC-USD")
	assert.NoError(t, err)
	assert.NotEmpty(t, summary.Volume)
}

func TestMarketsService_GetOrderBook(t *testing.T) {
	bt := New("", "")
	orderBook, err := bt.GetOrderBook("BTC-USD", 0)
	assert.NoError(t, err)
	assert.Len(t, orderBook.Ask, 25)
	assert.Len(t, orderBook.Bid, 25)
	orderBook, err = bt.GetOrderBook("BTC-USD", 2)
	assert.Error(t, err)
	orderBook, err = bt.GetOrderBook("BTC-USD", 1)
	assert.NoError(t, err)
	assert.Len(t, orderBook.Ask, 1)
	assert.Len(t, orderBook.Bid, 1)
	assert.NotEmpty(t, orderBook.Ask[0].Quantity)
	assert.NotEmpty(t, orderBook.Bid[0].Quantity)
}
