package bittrex

import (
	"time"

	"github.com/shopspring/decimal"
)

type Trade struct {
	Symbol     string
	ID         string
	ExecutedAt time.Time
	Quantity   decimal.Decimal
	Rate       decimal.Decimal
	TakerSide  string
}

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
