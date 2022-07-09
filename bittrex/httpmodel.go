package bittrex

import (
	"time"

	"github.com/shopspring/decimal"
)

// Currencies

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

// Markets

type Market struct {
	Symbol              string          `json:"symbol"`
	BaseCurrencySymbol  string          `json:"baseCurrencySymbol"`
	QuoteCurrencySymbol string          `json:"quoteCurrencySymbol"`
	MinTradeSize        decimal.Decimal `json:"minTradeSize"`
	Precision           int             `json:"precision"`
	Status              string          `json:"status"`
	CreatedAt           time.Time       `json:"createdAt"`
	Notice              string          `json:"notice"`
	ProhibitedIn        []string        `json:"prohibitedIn"`
}

type Ticker struct {
	Symbol        string          `json:"symbol"`
	LastTradeRate decimal.Decimal `json:"lastTradeRate"`
	BidRate       decimal.Decimal `json:"bidRate"`
	AskRate       decimal.Decimal `json:"askRate"`
}
