package bittrex

import (
	"time"

	"github.com/shopspring/decimal"
)

type TradeSlice struct {
	Deltas []struct {
		ID         string          `json:"id"`
		ExecutedAt time.Time       `json:"executedAt"`
		Quantity   decimal.Decimal `json:"quantity"`
		Rate       decimal.Decimal `json:"rate"`
		TakerSide  string          `json:"takerSide"`
	} `json:"deltas"`
	Sequence     int    `json:"sequence"`
	MarketSymbol string `json:"marketSymbol"`
}

type TickerSlice struct {
	Sequence int `json:"sequence"`
	Deltas   []struct {
		Symbol        string          `json:"symbol"`
		LastTradeRate decimal.Decimal `json:"lastTradeRate"`
		BidRate       decimal.Decimal `json:"bidRate"`
		AskRate       decimal.Decimal `json:"askRate"`
	}
}

type OrderBookSlice struct {
	MarketSymbol string  `json:"marketSymbol"`
	Depth        int     `json:"depth"`
	Sequence     int     `json:"sequence"`
	BidDeltas    []Order `json:"bidDeltas"`
	AskDeltas    []Order `json:"askDeltas"`
}
