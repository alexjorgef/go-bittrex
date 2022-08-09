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

	STREAM_ORDERBOOK = "orderBook"
	STREAM_TICKER    = "ticker"
	STREAM_ORDER     = "order"
	STREAM_TRADE     = "trade"
	STREAM_HEARTBEAT = "heartbeat"
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
						log.Printf("trade send err: %s %d \n", market, len(tickers))
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
				return errors.New("trade messages timeout")
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
							log.Printf("trade send err: %s %d \n", market, len(trades))
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
