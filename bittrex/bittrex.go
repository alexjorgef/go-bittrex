// Package bittrex is an implementation of the Biitrex API in Golang.
package bittrex

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	APIBASE    = "https://api.bittrex.com/" // HTTP API endpoint
	WSBASE     = "socket-v3.bittrex.com"    // WS API endpoint
	APIVERSION = "v3"                       // API version
	WSHUB      = "C3"                       // SignalR main hub

	CHANNEL_ORDERBOOK = "orderBook"
	CHANNEL_TICKER    = "ticker"
	CHANNEL_ORDER     = "order"
	CHANNEL_TRADE     = "trade"
	CHANNEL_HEARTBEAT = "heartbeat"
)

type Bittrex struct {
	client *Client
}

// New returns an instantiated bittrex struct
func New(apiKey, apiSecret string) *Bittrex {
	client := NewClient(apiKey, apiSecret)
	return &Bittrex{client}
}

// NewWithCustomHTTPClient returns an instantiated bittrex struct with custom http client
func NewWithCustomHTTPClient(apiKey, apiSecret string, httpClient *http.Client) *Bittrex {
	client := NewClientWithCustomHTTPConfig(apiKey, apiSecret, httpClient)
	return &Bittrex{client}
}

// NewWithCustomTimeout returns an instantiated bittrex struct with custom timeout
func NewWithCustomTimeout(apiKey, apiSecret string, timeout time.Duration) *Bittrex {
	client := NewClientWithCustomTimeout(apiKey, apiSecret, timeout)
	return &Bittrex{client}
}

// SetDebug set enable/disable http request/response dump
func (b *Bittrex) SetDebug(enable bool) {
	b.client.debug = enable
}

// Currencies

// Retrieve info on a specified currency.
func (b *Bittrex) GetCurrency(symbol string) (currency Currency, err error) {
	r, err := b.client.do("GET", "currencies/"+symbol, "", false)
	if err != nil {
		return
	}

	err = json.Unmarshal(r, &currency)
	return
}

// Markets

// List markets.
func (b *Bittrex) GetMarkets() (markets []Market, err error) {
	r, err := b.client.do("GET", "markets", "", false)
	if err != nil {
		return
	}

	err = json.Unmarshal(r, &markets)
	return
}

// List summaries of the last 24 hours of activity for all markets.
func (b *Bittrex) GetMarketsSummaries() (marketSummaries []MarketSummary, err error) {
	r, err := b.client.do("GET", "markets/summaries", "", false)
	if err != nil {
		return
	}

	err = json.Unmarshal(r, &marketSummaries)
	return
}

// List tickers for all markets.
func (b *Bittrex) GetMarketsTickers() (marketTickers []Ticker, err error) {
	r, err := b.client.do("GET", "markets/tickers", "", false)
	if err != nil {
		return
	}

	err = json.Unmarshal(r, &marketTickers)
	return
}

// Retrieve the ticker for a specific market.
func (b *Bittrex) GetTicker(marketSymbol string) (ticker Ticker, err error) {
	r, err := b.client.do("GET", "markets/"+strings.ToUpper(marketSymbol)+"/ticker", "", false)
	if err != nil {
		return
	}

	err = json.Unmarshal(r, &ticker)
	return
}

// Retrieve summary of the last 24 hours of activity for a specific market.
func (b *Bittrex) GetSummary(marketSymbol string) (marketSummary MarketSummary, err error) {
	r, err := b.client.do("GET", "markets/"+strings.ToUpper(marketSymbol)+"/summary", "", false)
	if err != nil {
		return
	}

	err = json.Unmarshal(r, &marketSummary)
	return
}

// Retrieve the order book for a specific market.
func (b *Bittrex) GetOrderBook(marketSymbol string, depth int) (orderBook OrderBook, err error) {
	if depth != 1 && depth != 25 && depth != 500 && depth != 0 {
		return orderBook, errors.New("invalid depth")
	}
	if depth == 0 {
		depth = 25
	}

	r, err := b.client.do("GET", "markets/"+strings.ToUpper(marketSymbol)+"/orderbook?depth="+strconv.Itoa(depth), "", false)
	if err != nil {
		return
	}

	err = json.Unmarshal(r, &orderBook)
	return
}
