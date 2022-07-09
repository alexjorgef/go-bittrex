// Package bittrex is an implementation of the Biitrex API in Golang.
package bittrex

import (
	"encoding/json"
	"net/http"
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

// Retrieve the ticker for a specific market.
func (b *Bittrex) GetTicker(market string) (ticker Ticker, err error) {
	r, err := b.client.do("GET", "markets/"+strings.ToUpper(market)+"/ticker", "", false)
	if err != nil {
		return
	}

	err = json.Unmarshal(r, &ticker)
	return
}
