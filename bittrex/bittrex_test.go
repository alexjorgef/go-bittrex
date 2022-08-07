package bittrex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Currencies

func TestCurrenciesService_GetCurrencies(t *testing.T) {
	bt := New("", "")
	currencies, err := bt.GetCurrencies()
	assert.NoError(t, err)
	assert.NotEmpty(t, currencies[0].Name)
}

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
	assert.NotEmpty(t, markets[0].Symbol)
}

func TestMarketsService_GetMarketsSummaries(t *testing.T) {
	bt := New("", "")
	marketSummaries, err := bt.GetMarketsSummaries()
	assert.NoError(t, err)
	assert.NotEmpty(t, marketSummaries[0].Volume)
}

func TestMarketsService_GetMarketsTickers(t *testing.T) {
	bt := New("", "")
	marketTickers, err := bt.GetMarketsTickers()
	assert.NoError(t, err)
	assert.NotEmpty(t, marketTickers[0].Symbol)
}

func TestMarketsService_GetTicker(t *testing.T) {
	bt := New("", "")
	ticker, err := bt.GetTicker("ETH-USD")
	assert.NoError(t, err)
	assert.NotEmpty(t, ticker.Symbol)
}

func TestMarketsService_GetMarket(t *testing.T) {
	bt := New("", "")
	market, err := bt.GetMarket("ETH-USD")
	assert.NoError(t, err)
	assert.NotEmpty(t, market.Symbol)
}

func TestMarketsService_GetSummary(t *testing.T) {
	bt := New("", "")
	summary, err := bt.GetSummary("ETH-USD")
	assert.NoError(t, err)
	assert.NotEmpty(t, summary.Volume)
}

func TestMarketsService_GetOrderBook(t *testing.T) {
	bt := New("", "")
	orderBook, err := bt.GetOrderBook("ETH-USD", 0)
	assert.NoError(t, err)
	assert.Len(t, orderBook.Ask, 25)
	assert.Len(t, orderBook.Bid, 25)
	orderBook, err = bt.GetOrderBook("ETH-USD", 2)
	assert.Error(t, err)
	orderBook, err = bt.GetOrderBook("ETH-USD", 1)
	assert.NoError(t, err)
	assert.Len(t, orderBook.Ask, 1)
	assert.Len(t, orderBook.Bid, 1)
	assert.NotEmpty(t, orderBook.Ask[0].Quantity)
	assert.NotEmpty(t, orderBook.Bid[0].Quantity)
}

func TestMarketsService_GetTrades(t *testing.T) {
	bt := New("", "")
	trades, err := bt.GetTrades("ETH-USD")
	assert.NoError(t, err)
	assert.NotEmpty(t, trades[0].Quantity)
	assert.NotEmpty(t, trades[0].Rate)
	assert.NotEmpty(t, trades[0].TakerSide)
}

func TestMarketsService_GetCandles(t *testing.T) {
	bt := New("", "")
	candles, err := bt.GetCandles("ETH-USD", INTERVAL_DAY1)
	assert.NoError(t, err)
	assert.NotEmpty(t, candles[0].StartsAt)
	assert.GreaterOrEqual(t, candles[0].High.IntPart(), candles[0].Low.IntPart())
	assert.NotEmpty(t, candles[0].Open)
	assert.NotEmpty(t, candles[0].Close)
	assert.NotEmpty(t, candles[0].Volume)
	assert.NotEmpty(t, candles[0].QuoteVolume)
	candles, err = bt.GetCandlesWithOpts("ETH-USD", INTERVAL_DAY1, &GetCandlesOpts{CandleType: CANDLETYPE_TRADE})
	assert.NoError(t, err)
	assert.NotEmpty(t, candles[0].StartsAt)
	assert.NotEmpty(t, candles[0].Volume)
	assert.NotEmpty(t, candles[0].QuoteVolume)
	candles, err = bt.GetCandlesWithOpts("ETH-USD", INTERVAL_DAY1, &GetCandlesOpts{CandleType: CANDLETYPE_MIDPOINT})
	assert.NoError(t, err)
	assert.NotEmpty(t, candles[0].StartsAt)
	assert.Empty(t, candles[0].Volume)
	assert.Empty(t, candles[0].QuoteVolume)
}

func TestMarketsService_GetCandlesHistory(t *testing.T) {
	bt := New("", "")
	candles, err := bt.GetCandlesHistory("ETH-USD", INTERVAL_DAY1, 2021)
	assert.NoError(t, err)
	assert.NotEmpty(t, candles[0].StartsAt)
	assert.GreaterOrEqual(t, candles[0].High.IntPart(), candles[0].Low.IntPart())
	assert.NotEmpty(t, candles[0].Open)
	assert.NotEmpty(t, candles[0].Close)
	assert.NotEmpty(t, candles[0].Volume)
	assert.NotEmpty(t, candles[0].QuoteVolume)
	candles, err = bt.GetCandlesHistoryWithOpts("ETH-USD", INTERVAL_DAY1, 2021, &GetCandlesHistoryOpts{CandleType: CANDLETYPE_TRADE, HistoryMonth: 10, HistoryDay: 1})
	assert.NoError(t, err)
	assert.NotEmpty(t, candles[0].StartsAt)
	assert.NotEmpty(t, candles[0].Volume)
	assert.NotEmpty(t, candles[0].QuoteVolume)
	candles, err = bt.GetCandlesHistoryWithOpts("ETH-USD", INTERVAL_HOUR1, 2021, &GetCandlesHistoryOpts{CandleType: CANDLETYPE_MIDPOINT, HistoryMonth: 10, HistoryDay: 1})
	assert.NoError(t, err)
	assert.NotEmpty(t, candles[0].StartsAt)
	assert.Empty(t, candles[0].Volume)
	assert.Empty(t, candles[0].QuoteVolume)
	candles, err = bt.GetCandlesHistoryWithOpts("ETH-USD", INTERVAL_MINUTE5, 2021, &GetCandlesHistoryOpts{HistoryMonth: 10, HistoryDay: 1})
	assert.NoError(t, err)
	assert.NotEmpty(t, candles[0].StartsAt)
	assert.NotEmpty(t, candles[0].Volume)
	assert.NotEmpty(t, candles[0].QuoteVolume)
	candles, err = bt.GetCandlesHistoryWithOpts("ETH-USD", INTERVAL_MINUTE5, 2021, &GetCandlesHistoryOpts{CandleType: CANDLETYPE_MIDPOINT, HistoryMonth: 10, HistoryDay: 1})
	assert.NoError(t, err)
	assert.NotEmpty(t, candles[0].StartsAt)
	assert.Empty(t, candles[0].Volume)
	assert.Empty(t, candles[0].QuoteVolume)
	candles, err = bt.GetCandlesHistoryWithOpts("ETH-USD", INTERVAL_HOUR1, 2021, &GetCandlesHistoryOpts{})
	assert.Error(t, err)
	assert.Empty(t, candles)
	candles, err = bt.GetCandlesHistoryWithOpts("ETH-USD", INTERVAL_MINUTE5, 2021, &GetCandlesHistoryOpts{HistoryMonth: 10})
	assert.Error(t, err)
	assert.Empty(t, candles)
}

func TestPingService_Ping(t *testing.T) {
	bt := New("", "")
	ping, err := bt.Ping()
	assert.NoError(t, err)
	assert.NotEmpty(t, ping)
}
