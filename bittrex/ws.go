package bittrex

import (
	"bytes"
	"compress/zlib"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"sync/atomic"
	"time"

	"github.com/alexjorgef/signalr"
	"github.com/google/uuid"
)

const (
	WS_BASE = "socket-v3.bittrex.com" // WS API endpoint
	WS_HUB  = "C3"                    // SignalR main hub

	STREAM_CANDLE          = "candle"
	STREAM_ORDERBOOK       = "orderBook"
	STREAM_TICKER          = "ticker"
	STREAM_TICKERS         = "tickers"
	STREAM_MARKETSUMMARIES = "marketSummaries"
	STREAM_MARKETSUMMARY   = "marketSummary"
	STREAM_ORDER           = "order"
	STREAM_TRADE           = "trade"
	STREAM_HEARTBEAT       = "heartbeat"
)

type Response struct {
	Success   bool        `json:"Success"`
	ErrorCode interface{} `json:"ErrorCode"`
}

// doAsyncTimeout runs f in a different goroutine
//   if f returns before timeout elapses, doAsyncTimeout returns
//     the result of f().
//     otherwise it returns "operation timeout" error, and calls tmFunc after f returns.
func doAsyncTimeout(f func() error, tmFunc func(error), timeout time.Duration) error {
	errs := make(chan error)

	go func() {
		err := f()
		select {
		case errs <- err:
		default:
			if tmFunc != nil {
				tmFunc(err)
			}
		}
	}()

	select {
	case err := <-errs:
		return err
	case <-time.After(timeout):
		return errors.New("operation timeout")
	}
}

// Some streams contain private data and require that you be authenticated prior to subscribing.
func (b *Bittrex) Authentication(c *signalr.Client) error {
	r := &Response{}

	apiTimestamp := time.Now().UnixNano() / 1000000
	UUID := uuid.New().String()

	preSign := strings.Join([]string{fmt.Sprintf("%d", apiTimestamp), UUID}, "")

	mac := hmac.New(sha512.New, []byte(b.client.apiSecret))
	_, err := mac.Write([]byte(preSign))
	if err != nil {
		return err
	}

	sig := hex.EncodeToString(mac.Sum(nil))

	auth, err := c.CallHub(WS_HUB, "Authenticate", b.client.apiKey, apiTimestamp, UUID, sig)
	if err != nil {
		return err
	}

	err = json.Unmarshal(auth, r)
	if err != nil {
		return fmt.Errorf("%s", err.Error())
	}

	if !r.Success {
		return fmt.Errorf("%s", r.ErrorCode)
	}

	return nil
}

// Provides regular updates of the current market summary data for all markets.
//   Market summary data is different from candles in that it is a rolling 24-hour number as opposed to data for a fixed interval like candles.
func (b *Bittrex) SubscribeMarketSummariesUpdates(marketSummaries chan<- MarketSummary, stop <-chan bool) error {
	const timeout = 5 * time.Second
	client := signalr.NewWebsocketClient()

	var updTime int64

	client.OnClientMethod = func(hub string, method string, messages []json.RawMessage) {
		if hub != WS_HUB {
			return
		}

		switch method {
		case STREAM_HEARTBEAT, STREAM_MARKETSUMMARIES:
			atomic.StoreInt64(&updTime, time.Now().Unix())
		default:
			fmt.Printf("unsupported message type: %s %v\n", method, messages)
		}

		for _, msg := range messages {
			dbuf, err := base64.StdEncoding.DecodeString(strings.Trim(string(msg), `"`))
			if err != nil {
				fmt.Printf("DecodeString error: %s %s\n", err.Error(), string(msg))
				continue
			}

			r, err := zlib.NewReader(bytes.NewReader(append([]byte{120, 156}, dbuf...)))
			if err != nil {
				fmt.Printf("unzip error %s %s \n", err.Error(), string(msg))
				continue
			}
			defer r.Close()

			var out bytes.Buffer
			written, _ := io.Copy(&out, r)

			if written > 0 {
				marketSummarySlice := MarketSummarySlice{}
				err = json.Unmarshal(out.Bytes(), &marketSummarySlice)
				if err != nil {
					fmt.Printf("unmarshal error %s\n", err.Error())
					continue
				}

				for _, delta := range marketSummarySlice.Deltas {
					marketSummary := MarketSummary{}
					marketSummary.Symbol = delta.Symbol
					marketSummary.High = delta.High
					marketSummary.Low = delta.Low
					marketSummary.Volume = delta.Volume
					marketSummary.QuoteVolume = delta.QuoteVolume
					marketSummary.PercentChange = delta.PercentChange
					marketSummary.UpdatedAt = delta.UpdatedAt
					select {
					case marketSummaries <- marketSummary:
					default:
						if b.client.debug {
							log.Printf("marketSummaries send err: %d\n", len(marketSummaries))
						}
					}
				}
			}
		}
	}

	client.OnMessageError = func(err error) {
		fmt.Printf("ERROR OCCURRED: %s\n", err.Error())
	}

	err := doAsyncTimeout(
		func() error {
			return client.Connect("https", WS_BASE, []string{WS_HUB})
		}, func(err error) {
			if err == nil {
				client.Close()
			}
		}, timeout)
	if err != nil {
		return err
	}

	defer client.Close()

	_, err = client.CallHub(WS_HUB, "Subscribe", []interface{}{"heartbeat", "market_summaries"})
	if err != nil {
		return err
	}

	tick := time.NewTicker(1 * time.Minute)

	for {
		select {
		case signal := <-stop:
			if signal {
				return errors.New("client.stop")
			}
		case <-client.DisconnectedChannel:
			return errors.New("client.DisconnectedChannel")
		case <-tick.C:
			if time.Now().Unix()-atomic.LoadInt64(&updTime) > 60 {
				return errors.New("marketSummaries messages timeout")
			}
		}
	}
}

// Provides regular updates of the current market summary data for a given market.
//   Market summary data is different from candles in that it is a rolling 24-hour number as opposed to data for a fixed interval like candles.
func (b *Bittrex) SubscribeMarketSummaryUpdates(market string, marketSummaries chan<- MarketSummary, stop <-chan bool) error {
	const timeout = 5 * time.Second
	client := signalr.NewWebsocketClient()

	var updTime int64

	client.OnClientMethod = func(hub string, method string, messages []json.RawMessage) {
		if hub != WS_HUB {
			return
		}

		switch method {
		case STREAM_HEARTBEAT, STREAM_MARKETSUMMARY:
			atomic.StoreInt64(&updTime, time.Now().Unix())
		default:
			fmt.Printf("unsupported message type: %s %v\n", method, messages)
		}

		for _, msg := range messages {
			dbuf, err := base64.StdEncoding.DecodeString(strings.Trim(string(msg), `"`))
			if err != nil {
				fmt.Printf("DecodeString error: %s %s\n", err.Error(), string(msg))
				continue
			}

			r, err := zlib.NewReader(bytes.NewReader(append([]byte{120, 156}, dbuf...)))
			if err != nil {
				fmt.Printf("unzip error %s %s \n", err.Error(), string(msg))
				continue
			}
			defer r.Close()

			var out bytes.Buffer
			written, _ := io.Copy(&out, r)

			if written > 0 {
				marketSummary := MarketSummary{}
				err = json.Unmarshal(out.Bytes(), &marketSummary)
				if err != nil {
					fmt.Printf("unmarshal error %s\n", err.Error())
					continue
				}
				select {
				case marketSummaries <- marketSummary:
				default:
					if b.client.debug {
						log.Printf("marketSummary send err: %s %d\n", market, len(marketSummaries))
					}
				}
			}
		}
	}

	client.OnMessageError = func(err error) {
		fmt.Printf("ERROR OCCURRED: %s\n", err.Error())
	}

	err := doAsyncTimeout(
		func() error {
			return client.Connect("https", WS_BASE, []string{WS_HUB})
		}, func(err error) {
			if err == nil {
				client.Close()
			}
		}, timeout)
	if err != nil {
		return err
	}

	defer client.Close()

	_, err = client.CallHub(WS_HUB, "Subscribe", []interface{}{"heartbeat", "market_summary_" + market})
	if err != nil {
		return err
	}

	tick := time.NewTicker(1 * time.Minute)

	for {
		select {
		case signal := <-stop:
			if signal {
				return errors.New("client.stop")
			}
		case <-client.DisconnectedChannel:
			return errors.New("client.DisconnectedChannel")
		case <-tick.C:
			if time.Now().Unix()-atomic.LoadInt64(&updTime) > 60 {
				return errors.New("marketSummary messages timeout")
			}
		}
	}
}

// TODO: Add a default and a function with parameters for candleInterval
// Sends a message at the start of each candle (based on the subscribed interval) and when trades have occurred on the market.
//   Note that this means on an active market you will receive many updates over the course of each candle interval as trades occur.
//   You will always recieve an update at the start of each interval.
//   If no trades occurred yet, this update will be a 0-volume placeholder that carries forward the Close of the previous interval as the current interval's OHLC values.
func (b *Bittrex) SubscribeCandleUpdates(market string, candles chan<- Candle, stop <-chan bool) error {
	const timeout = 5 * time.Second
	client := signalr.NewWebsocketClient()

	var updTime int64

	client.OnClientMethod = func(hub string, method string, messages []json.RawMessage) {
		if hub != WS_HUB {
			return
		}

		switch method {
		case STREAM_HEARTBEAT, STREAM_CANDLE:
			atomic.StoreInt64(&updTime, time.Now().Unix())
		default:
			fmt.Printf("unsupported message type: %s %v\n", method, messages)
		}

		for _, msg := range messages {
			dbuf, err := base64.StdEncoding.DecodeString(strings.Trim(string(msg), `"`))
			if err != nil {
				fmt.Printf("DecodeString error: %s %s\n", err.Error(), string(msg))
				continue
			}

			r, err := zlib.NewReader(bytes.NewReader(append([]byte{120, 156}, dbuf...)))
			if err != nil {
				fmt.Printf("unzip error %s %s \n", err.Error(), string(msg))
				continue
			}
			defer r.Close()

			var out bytes.Buffer
			written, _ := io.Copy(&out, r)

			if written > 0 {
				candleSlice := CandleSlice{}
				err = json.Unmarshal(out.Bytes(), &candleSlice)
				if err != nil {
					fmt.Printf("unmarshal error %s\n", err.Error())
					continue
				}

				candle := Candle{
					MarketSymbol: candleSlice.MarketSymbol,
					Interval:     candleSlice.Interval,
					StartsAt:     candleSlice.Delta.StartsAt,
					Open:         candleSlice.Delta.Open,
					High:         candleSlice.Delta.High,
					Low:          candleSlice.Delta.Low,
					Close:        candleSlice.Delta.Close,
					Volume:       candleSlice.Delta.Volume,
					QuoteVolume:  candleSlice.Delta.QuoteVolume,
				}
				select {
				case candles <- candle:
				default:
					if b.client.debug {
						log.Printf("candle send err: %s %d\n", market, len(candles))
					}
				}

			}
		}
	}

	client.OnMessageError = func(err error) {
		fmt.Printf("ERROR OCCURRED: %s\n", err.Error())
	}

	err := doAsyncTimeout(
		func() error {
			return client.Connect("https", WS_BASE, []string{WS_HUB})
		}, func(err error) {
			if err == nil {
				client.Close()
			}
		}, timeout)
	if err != nil {
		return err
	}

	defer client.Close()

	_, err = client.CallHub(WS_HUB, "Subscribe", []interface{}{"heartbeat", "candle_" + market + "_MINUTE_1"})
	if err != nil {
		return err
	}

	tick := time.NewTicker(1 * time.Minute)

	for {
		select {
		case signal := <-stop:
			if signal {
				return errors.New("client.stop")
			}
		case <-client.DisconnectedChannel:
			return errors.New("client.DisconnectedChannel")
		case <-tick.C:
			if time.Now().Unix()-atomic.LoadInt64(&updTime) > 60 {
				return errors.New("candle messages timeout")
			}
		}
	}
}

// TODO: Add a default and a function with parameters for depth
// Sends a message when there are changes to the order book within the subscribed depth.
func (b *Bittrex) SubscribeOrderbookUpdates(market string, orderbooks chan<- OrderBook, stop <-chan bool) error {
	const timeout = 5 * time.Second
	client := signalr.NewWebsocketClient()

	var updTime int64

	client.OnClientMethod = func(hub string, method string, messages []json.RawMessage) {
		if hub != WS_HUB {
			return
		}

		switch method {
		case STREAM_HEARTBEAT, STREAM_ORDERBOOK:
			atomic.StoreInt64(&updTime, time.Now().Unix())
		default:
			fmt.Printf("unsupported message type: %s %v\n", method, messages)
		}

		for _, msg := range messages {
			dbuf, err := base64.StdEncoding.DecodeString(strings.Trim(string(msg), `"`))
			if err != nil {
				fmt.Printf("DecodeString error: %s %s\n", err.Error(), string(msg))
				continue
			}

			r, err := zlib.NewReader(bytes.NewReader(append([]byte{120, 156}, dbuf...)))
			if err != nil {
				fmt.Printf("unzip error %s %s \n", err.Error(), string(msg))
				continue
			}
			defer r.Close()

			var out bytes.Buffer
			written, _ := io.Copy(&out, r)

			if written > 0 {
				orderbookSlice := OrderBookSlice{}
				err = json.Unmarshal(out.Bytes(), &orderbookSlice)
				if err != nil {
					fmt.Printf("unmarshal error %s\n", err.Error())
					continue
				}

				orderbook := OrderBook{Symbol: orderbookSlice.MarketSymbol, Depth: orderbookSlice.Depth}
				for _, delta := range orderbookSlice.AskDeltas {
					orderbook.Ask = append(orderbook.Ask, Order{Quantity: delta.Quantity, Rate: delta.Rate})
				}
				for _, delta := range orderbookSlice.BidDeltas {
					orderbook.Bid = append(orderbook.Bid, Order{Quantity: delta.Quantity, Rate: delta.Rate})
				}
				select {
				case orderbooks <- orderbook:
				default:
					if b.client.debug {
						log.Printf("orderbook send err: %s %d\n", market, len(orderbooks))
					}
				}

			}
		}
	}

	client.OnMessageError = func(err error) {
		fmt.Printf("ERROR OCCURRED: %s\n", err.Error())
	}

	err := doAsyncTimeout(
		func() error {
			return client.Connect("https", WS_BASE, []string{WS_HUB})
		}, func(err error) {
			if err == nil {
				client.Close()
			}
		}, timeout)
	if err != nil {
		return err
	}

	defer client.Close()

	_, err = client.CallHub(WS_HUB, "Subscribe", []interface{}{"heartbeat", "orderbook_" + market + "_25"})
	if err != nil {
		return err
	}

	tick := time.NewTicker(1 * time.Minute)

	for {
		select {
		case signal := <-stop:
			if signal {
				return errors.New("client.stop")
			}
		case <-client.DisconnectedChannel:
			return errors.New("client.DisconnectedChannel")
		case <-tick.C:
			if time.Now().Unix()-atomic.LoadInt64(&updTime) > 60 {
				return errors.New("orderook messages timeout")
			}
		}

	}
}

// Sends a message with the best bid price, best ask price, and last trade price for all markets as there are changes to the order book or trades.
func (b *Bittrex) SubscribeTickersUpdates(tickers chan<- Ticker, stop <-chan bool) error {
	const timeout = 15 * time.Second
	client := signalr.NewWebsocketClient()

	var updTime int64

	client.OnClientMethod = func(hub string, method string, messages []json.RawMessage) {
		if hub != WS_HUB {
			return
		}

		switch method {
		case STREAM_HEARTBEAT, STREAM_TICKERS:
			atomic.StoreInt64(&updTime, time.Now().Unix())
		default:
			fmt.Printf("unsupported message type: %s\n", method)
		}

		for _, msg := range messages {
			dbuf, err := base64.StdEncoding.DecodeString(strings.Trim(string(msg), `"`))
			if err != nil {
				fmt.Printf("DecodeString error: %s %s\n", err.Error(), string(msg))
				continue
			}

			r, err := zlib.NewReader(bytes.NewReader(append([]byte{120, 156}, dbuf...)))
			if err != nil {
				fmt.Printf("unzip error %s %s\n", err.Error(), string(msg))
				continue
			}
			defer r.Close()

			var out bytes.Buffer
			written, _ := io.Copy(&out, r)

			if written > 0 {
				tickerSlice := TickerSlice{}
				err = json.Unmarshal(out.Bytes(), &tickerSlice)
				if err != nil {
					fmt.Printf("unmarshal error %s\n", err.Error())
					continue
				}

				for _, delta := range tickerSlice.Deltas {
					ticker := Ticker{}
					ticker.Symbol = delta.Symbol
					ticker.LastTradeRate = delta.LastTradeRate
					ticker.BidRate = delta.BidRate
					ticker.AskRate = delta.AskRate
					select {
					case tickers <- ticker:
					default:
						if b.client.debug {
							log.Printf("tickers send err: %d\n", len(tickers))
						}
					}
				}
			}
		}
	}

	client.OnMessageError = func(err error) {
		fmt.Printf("ERROR OCCURRED: %s\n", err.Error())
	}

	err := doAsyncTimeout(
		func() error {
			return client.Connect("https", WS_BASE, []string{WS_HUB})
		}, func(err error) {
			if err == nil {
				client.Close()
			}
		}, timeout)
	if err != nil {
		return err
	}

	defer client.Close()

	_, err = client.CallHub(WS_HUB, "Subscribe", []interface{}{"heartbeat", "tickers"})
	if err != nil {
		return err
	}

	tick := time.NewTicker(1 * time.Minute)

	// Blocking loop
	for {
		select {
		case signal := <-stop:
			if signal {
				return errors.New("client.stop")
			}
		case <-client.DisconnectedChannel:
			return errors.New("client.DisconnectedChannel")
		case <-tick.C:
			if time.Now().Unix()-atomic.LoadInt64(&updTime) > 60 {
				return errors.New("tickers messages timeout")
			}
		}
	}
}

// Sends a message with the best bid and ask price for the given market as well as the last trade price whenever there is a relevant change to the order book or a trade.
func (b *Bittrex) SubscribeTickerUpdates(market string, tickers chan<- Ticker, stop <-chan bool) error {
	const timeout = 15 * time.Second
	client := signalr.NewWebsocketClient()

	var updTime int64

	client.OnClientMethod = func(hub string, method string, messages []json.RawMessage) {
		if hub != WS_HUB {
			return
		}

		switch method {
		case STREAM_HEARTBEAT, STREAM_TICKER:
			atomic.StoreInt64(&updTime, time.Now().Unix())
		default:
			fmt.Printf("unsupported message type: %s\n", method)
		}

		for _, msg := range messages {
			dbuf, err := base64.StdEncoding.DecodeString(strings.Trim(string(msg), `"`))
			if err != nil {
				fmt.Printf("DecodeString error: %s %s\n", err.Error(), string(msg))
				continue
			}

			r, err := zlib.NewReader(bytes.NewReader(append([]byte{120, 156}, dbuf...)))
			if err != nil {
				fmt.Printf("unzip error %s %s\n", err.Error(), string(msg))
				continue
			}
			defer r.Close()

			var out bytes.Buffer
			written, _ := io.Copy(&out, r)

			if written > 0 {
				ticker := Ticker{}
				err = json.Unmarshal(out.Bytes(), &ticker)
				if err != nil {
					fmt.Printf("unmarshal error %s\n", err.Error())
					continue
				}

				select {
				case tickers <- ticker:
				default:
					if b.client.debug {
						log.Printf("ticker send err: %s %d\n", market, len(tickers))
					}
				}
			}
		}
	}

	client.OnMessageError = func(err error) {
		fmt.Printf("ERROR OCCURRED: %s\n", err.Error())
	}

	err := doAsyncTimeout(
		func() error {
			return client.Connect("https", WS_BASE, []string{WS_HUB})
		}, func(err error) {
			if err == nil {
				client.Close()
			}
		}, timeout)
	if err != nil {
		return err
	}

	defer client.Close()

	_, err = client.CallHub(WS_HUB, "Subscribe", []interface{}{"heartbeat", "ticker_" + market})
	if err != nil {
		return err
	}

	tick := time.NewTicker(1 * time.Minute)

	// Blocking loop
	for {
		select {
		case signal := <-stop:
			if signal {
				return errors.New("client.stop")
			}
		case <-client.DisconnectedChannel:
			return errors.New("client.DisconnectedChannel")
		case <-tick.C:
			if time.Now().Unix()-atomic.LoadInt64(&updTime) > 60 {
				return errors.New("ticker messages timeout")
			}
		}
	}
}

// Sends a message with the quantity and rate of trades on a market as they occur.
func (b *Bittrex) SubscribeTradeUpdates(market string, trades chan<- Trade, stop <-chan bool) error {
	const timeout = 15 * time.Second
	client := signalr.NewWebsocketClient()

	var updTime int64

	client.OnClientMethod = func(hub string, method string, messages []json.RawMessage) {

		if hub != WS_HUB {
			return
		}

		switch method {
		case STREAM_HEARTBEAT, STREAM_TRADE:
			atomic.StoreInt64(&updTime, time.Now().Unix())

		default:
			fmt.Printf("unsupported message type: %s\n", method)
		}

		for _, msg := range messages {
			dbuf, err := base64.StdEncoding.DecodeString(strings.Trim(string(msg), `"`))
			if err != nil {
				fmt.Printf("DecodeString error: %s %s\n", err.Error(), string(msg))
				continue
			}

			r, err := zlib.NewReader(bytes.NewReader(append([]byte{120, 156}, dbuf...)))
			if err != nil {
				fmt.Printf("unzip error %s %s\n", err.Error(), string(msg))
				continue
			}
			defer r.Close()

			var out bytes.Buffer
			written, _ := io.Copy(&out, r)

			if written > 0 {
				tradeSlice := TradeSlice{}
				err = json.Unmarshal(out.Bytes(), &tradeSlice)
				if err != nil {
					fmt.Printf("unmarshal error %s\n", err.Error())
					continue
				}

				trade := Trade{Symbol: tradeSlice.MarketSymbol}

				for _, delta := range tradeSlice.Deltas {
					trade.ID = delta.ID
					trade.ExecutedAt = delta.ExecutedAt
					trade.Quantity = delta.Quantity
					trade.Rate = delta.Rate
					trade.TakerSide = delta.TakerSide
					select {
					case trades <- trade:
					default:
						if b.client.debug {
							log.Printf("trade send err: %s %d\n", market, len(trades))
						}
					}
				}
			}
		}
	}

	client.OnMessageError = func(err error) {
		fmt.Printf("ERROR OCCURRED: %s\n", err.Error())
	}

	err := doAsyncTimeout(
		func() error {
			return client.Connect("https", WS_BASE, []string{WS_HUB})
		}, func(err error) {
			if err == nil {
				client.Close()
			}
		}, timeout)
	if err != nil {
		return err
	}

	defer client.Close()

	_, err = client.CallHub(WS_HUB, "Subscribe", []interface{}{"heartbeat", "trade_" + market})
	if err != nil {
		return err
	}

	tick := time.NewTicker(1 * time.Minute)

	// Blocking loop
	for {
		select {
		case signal := <-stop:
			if signal {
				return errors.New("client.stop")
			}
		case <-client.DisconnectedChannel:
			return errors.New("client.DisconnectedChannel")
		case <-tick.C:
			if time.Now().Unix()-atomic.LoadInt64(&updTime) > 60 {
				return errors.New("trade messages timeout")
			}
		}
	}
}
