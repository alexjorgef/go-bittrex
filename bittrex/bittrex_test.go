package bittrex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Currencies

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
	assert.Equal(t, "1ECO-BTC", markets[0].Symbol)
}
