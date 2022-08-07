// Package bittrex is an implementation of the Biitrex API in Golang.
package bittrex

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	CANDLETYPE_TRADE    = "TRADE"
	CANDLETYPE_MIDPOINT = "MIDPOINT"

	INTERVAL_DAY1    = "DAY_1"
	INTERVAL_HOUR1   = "HOUR_1"
	INTERVAL_MINUTE5 = "MINUTE_5"
	INTERVAL_MINUTE1 = "MINUTE_1"
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

// List currencies.
func (b *Bittrex) GetCurrencies() (currencies []Currency, err error) {
	r, err := b.client.do("GET", "currencies", "", false)
	if err != nil {
		return
	}

	err = json.Unmarshal(r, &currencies)
	return
}

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

// Retrieve the ticker for a specific market.
func (b *Bittrex) GetMarket(marketSymbol string) (market Market, err error) {
	r, err := b.client.do("GET", "markets/"+strings.ToUpper(marketSymbol), "", false)
	if err != nil {
		return
	}

	err = json.Unmarshal(r, &market)
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

// Retrieve the recent trades for a specific market.
func (b *Bittrex) GetTrades(marketSymbol string) (trades []Trade, err error) {
	r, err := b.client.do("GET", "markets/"+strings.ToUpper(marketSymbol)+"/trades", "", false)
	if err != nil {
		return
	}

	err = json.Unmarshal(r, &trades)
	return
}

type GetCandlesOpts struct {
	CandleType string
}

// Retrieve recent candles for a specific market and candle interval.
//   The maximum age of the returned candles depends on the interval as follows:
//   (MINUTE_1: 1 day, MINUTE_5: 1 day, HOUR_1: 31 days, DAY_1: 366 days).
//   Candles for intervals without any trading activity will match the previous close and volume will be zero.
func (b *Bittrex) GetCandles(marketSymbol string, candleInterval string) (candles []Candle, err error) {
	return b.GetCandlesWithOpts(marketSymbol, candleInterval, &GetCandlesOpts{})
}

// Retrieve recent candles for a specific market and candle interval.
//   The maximum age of the returned candles depends on the interval as follows:
//   (MINUTE_1: 1 day, MINUTE_5: 1 day, HOUR_1: 31 days, DAY_1: 366 days).
//   Candles for intervals without any trading activity will match the previous close and volume will be zero.
func (b *Bittrex) GetCandlesWithOpts(marketSymbol string, candleInterval string, opts *GetCandlesOpts) (candles []Candle, err error) {
	v := reflect.ValueOf(opts)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return candles, errors.New("invalid opts pointer")
	}

	endpoint := "markets/" + strings.ToUpper(marketSymbol) + "/candles/" + strings.ToUpper(candleInterval) + "/recent"
	if !reflect.DeepEqual(opts, &GetCandlesOpts{}) {
		if len(v.Elem().Field(0).Interface().(string)) > 0 {
			endpoint = "markets/" + strings.ToUpper(marketSymbol) + "/candles/" + strings.ToUpper(opts.CandleType) + "/" + strings.ToUpper(candleInterval) + "/recent"
		}
	}

	r, err := b.client.do("GET", endpoint, "", false)
	if err != nil {
		return
	}

	err = json.Unmarshal(r, &candles)
	return
}

type GetCandlesHistoryOpts struct {
	CandleType   string
	HistoryMonth int
	HistoryDay   int
}

// Retrieve recent candles for a specific market and candle interval.
//   The date range of returned candles depends on the interval as follows:
//   (MINUTE_1: 1 day, MINUTE_5: 1 day, HOUR_1: 31 days, DAY_1: 366 days).
//   Candles for intervals without any trading activity will match the previous close and volume will be zero.
func (b *Bittrex) GetCandlesHistory(marketSymbol string, candleInterval string, year int) (candles []Candle, err error) {
	return b.GetCandlesHistoryWithOpts(marketSymbol, candleInterval, year, &GetCandlesHistoryOpts{})
}

// Retrieve recent candles for a specific market and candle interval.
//   The date range of returned candles depends on the interval as follows:
//   (MINUTE_1: 1 day, MINUTE_5: 1 day, HOUR_1: 31 days, DAY_1: 366 days).
//   Candles for intervals without any trading activity will match the previous close and volume will be zero.
func (b *Bittrex) GetCandlesHistoryWithOpts(marketSymbol string, candleInterval string, year int, opts *GetCandlesHistoryOpts) (candles []Candle, err error) {
	v := reflect.ValueOf(opts)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return candles, errors.New("invalid opts pointer")
	}

	endpoint := "markets/" + strings.ToUpper(marketSymbol) + "/candles/" + strings.ToUpper(candleInterval) + "/historical/" + strconv.Itoa(year)
	if !reflect.DeepEqual(opts, &GetCandlesHistoryOpts{}) {
		var endpointHistPart string
		if candleInterval == INTERVAL_DAY1 {
			endpointHistPart = "/historical/" + strconv.Itoa(year)
		}
		if candleInterval == INTERVAL_HOUR1 {
			if v.Elem().Field(1).Interface().(int) != 0 {
				endpointHistPart = "/historical/" + strconv.Itoa(year) + "/" + strconv.Itoa(opts.HistoryMonth)
			} else {
				return candles, errors.New("invalid HistoryMonth option")
			}
		}
		if candleInterval == INTERVAL_MINUTE5 || candleInterval == INTERVAL_MINUTE1 {
			if v.Elem().Field(1).Interface().(int) != 0 && v.Elem().Field(2).Interface().(int) != 0 {
				endpointHistPart = "/historical/" + strconv.Itoa(year) + "/" + strconv.Itoa(opts.HistoryMonth) + "/" + strconv.Itoa(opts.HistoryDay)
			} else {
				return candles, errors.New("invalid HistoryDay option")
			}
		}
		if len(v.Elem().Field(0).Interface().(string)) > 0 {
			endpoint = "markets/" + strings.ToUpper(marketSymbol) + "/candles/" + strings.ToUpper(opts.CandleType) + "/" + strings.ToUpper(candleInterval) + endpointHistPart
		} else {
			endpoint = "markets/" + strings.ToUpper(marketSymbol) + "/candles/" + strings.ToUpper(candleInterval) + endpointHistPart
		}
	}

	r, err := b.client.do("GET", endpoint, "", false)
	if err != nil {
		return
	}

	err = json.Unmarshal(r, &candles)
	return
}

// Ping

// Pings the service
func (b *Bittrex) Ping() (serverTime int64, err error) {
	r, err := b.client.do("GET", "ping", "", false)
	if err != nil {
		return
	}
	var pingR Ping

	err = json.Unmarshal(r, &pingR)
	serverTime = pingR.ServerTime
	return
}
