package bittrex

import (
	"time"

	"github.com/shopspring/decimal"
)

type Currency struct {
	Symbol                   string        `json:"symbol"`
	Name                     string        `json:"name"`
	CoinType                 string        `json:"coinType"`
	Status                   string        `json:"status"`
	MinConfirmations         int           `json:"minConfirmations"`
	Notice                   string        `json:"notice"`
	TxFee                    string        `json:"txFee"`
	LogoURL                  string        `json:"logoUrl"`
	ProhibitedIn             []interface{} `json:"prohibitedIn"`
	BaseAddress              string        `json:"baseAddress"`
	AssociatedTermsOfService []interface{} `json:"associatedTermsOfService"`
	Tags                     []interface{} `json:"tags"`
}

type Market struct {
	Symbol                   string          `json:"symbol"`
	BaseCurrencySymbol       string          `json:"baseCurrencySymbol"`
	QuoteCurrencySymbol      string          `json:"quoteCurrencySymbol"`
	MinTradeSize             decimal.Decimal `json:"minTradeSize"`
	Precision                int             `json:"precision"`
	Status                   string          `json:"status"`
	CreatedAt                time.Time       `json:"createdAt"`
	ProhibitedIn             []string        `json:"prohibitedIn"`
	AssociatedTermsOfService []string        `json:"associatedTermsOfService"`
	Tags                     []string        `json:"tags"`
}

type MarketSummary struct {
	Symbol        string          `json:"symbol"`
	High          decimal.Decimal `json:"high"`
	Low           decimal.Decimal `json:"low"`
	Volume        decimal.Decimal `json:"volume"`
	QuoteVolume   decimal.Decimal `json:"quoteVolume"`
	PercentChange decimal.Decimal `json:"percentChange"`
	UpdatedAt     time.Time       `json:"updatedAt"`
}

type Ticker struct {
	Symbol        string          `json:"symbol"`
	LastTradeRate decimal.Decimal `json:"lastTradeRate"`
	BidRate       decimal.Decimal `json:"bidRate"`
	AskRate       decimal.Decimal `json:"askRate"`
}

type Order struct {
	Quantity decimal.Decimal `json:"quantity"`
	Rate     decimal.Decimal `json:"rate"`
}

type OrderBook struct {
	Symbol string
	Depth  int
	Bid    []Order `json:"bid"`
	Ask    []Order `json:"ask"`
}

type Trade struct {
	Symbol     string
	ID         string          `json:"id"`
	ExecutedAt time.Time       `json:"executedAt"`
	Quantity   decimal.Decimal `json:"quantity"`
	Rate       decimal.Decimal `json:"rate"`
	TakerSide  string          `json:"takerSide"`
}

type Ping struct {
	ServerTime int64 `json:"serverTime"`
}

type Candle struct {
	MarketSymbol string
	Interval     string
	StartsAt     time.Time       `json:"startsAt"`
	Open         decimal.Decimal `json:"open"`
	High         decimal.Decimal `json:"high"`
	Low          decimal.Decimal `json:"low"`
	Close        decimal.Decimal `json:"close"`
	Volume       decimal.Decimal `json:"volume"`
	QuoteVolume  decimal.Decimal `json:"quoteVolume"`
}
