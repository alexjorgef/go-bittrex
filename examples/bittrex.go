package main

import (
	"fmt"
	"os"

	"github.com/alexjorgef/go-bittrex/bittrex"
)

const (
	API_KEY    = ""
	API_SECRET = ""
)

func main() {
	os.Exit(realMainHttp())
}

func realMainHttp() int {

	// Bittrex client
	client := bittrex.New(API_KEY, API_SECRET)

	// Currencies

	currencies, err := client.GetCurrencies()
	if err != nil {
		return 1
	}
	fmt.Printf("Currencies (Symbol, Name):\n")
	for i := 0; i < 3; i++ {
		fmt.Printf("\t%s\t%s\n", currencies[i].Symbol, currencies[i].Name)
	}

	currency, err := client.GetCurrency("XLM")
	if err != nil {
		return 1
	}
	fmt.Printf("Currency XLM:\n")
	fmt.Printf("\tSymbol:      \t%s\n", currency.Symbol)
	fmt.Printf("\tName:        \t%s\n", currency.Name)
	fmt.Printf("\tBaseAddress: \t%s\n", currency.BaseAddress)
	fmt.Printf("\tCoinType:    \t%s\n", currency.CoinType)
	fmt.Printf("\tLogoURL:     \t%s\n", currency.LogoURL)
	fmt.Printf("\tNotice:      \t%s\n", currency.Notice)
	fmt.Printf("\tProhibitedIn:\t%+q\n", currency.ProhibitedIn)

	// Markets

	markets, err := client.GetMarkets()
	if err != nil {
		return 1
	}
	fmt.Printf("Markets (Symbol, Status):\n")
	for i := 0; i < 3; i++ {
		fmt.Printf("\t%s\t%s\n", markets[i].Symbol, markets[i].Status)
	}

	marketsSummaries, err := client.GetMarketsSummaries()
	if err != nil {
		return 1
	}
	fmt.Printf("Summaries (Symbol, Volume):\n")
	for i := 0; i < 3; i++ {
		fmt.Printf("\t%s\t%s\n", marketsSummaries[i].Symbol, marketsSummaries[i].Volume.String())
	}

	marketsTickers, err := client.GetMarketsTickers()
	if err != nil {
		return 1
	}
	fmt.Printf("Tickers (Symbol, LastTradeRate):\n")
	for i := 0; i < 3; i++ {
		fmt.Printf("\t%s\t%s\n", marketsTickers[i].Symbol, marketsTickers[i].LastTradeRate.String())
	}

	ticker, err := client.GetTicker("ETH-USD")
	if err != nil {
		return 1
	}
	fmt.Printf("Ticker ETH-USD:\n")
	fmt.Printf("\tSymbol:       \t%s\n", ticker.Symbol)
	fmt.Printf("\tLastTradeRate:\t%s\n", ticker.LastTradeRate.String())
	fmt.Printf("\tAskRate:      \t%s\n", ticker.AskRate.String())
	fmt.Printf("\tBidRate:      \t%s\n", ticker.BidRate.String())

	market, err := client.GetMarket("ETH-USD")
	if err != nil {
		return 1
	}
	fmt.Printf("Market ETH-USD:\n")
	fmt.Printf("\tSymbol:             \t%s\n", market.Symbol)
	fmt.Printf("\tBaseCurrencySymbol: \t%s\n", market.BaseCurrencySymbol)
	fmt.Printf("\tQuoteCurrencySymbol:\t%s\n", market.QuoteCurrencySymbol)
	fmt.Printf("\tMinTradeSize:       \t%s\n", market.MinTradeSize.String())
	fmt.Printf("\tPrecision:          \t%d\n", market.Precision)
	fmt.Printf("\tStatus:             \t%s\n", market.Status)
	fmt.Printf("\tCreatedAt:          \t%s\n", market.CreatedAt)
	fmt.Printf("\tProhibitedIn:       \t%+q\n", market.ProhibitedIn)

	summary, err := client.GetSummary("ETH-USD")
	if err != nil {
		return 1
	}
	fmt.Printf("Summary ETH-USD:\n")
	fmt.Printf("\tSymbol:       \t%s\n", summary.Symbol)
	fmt.Printf("\tHigh:         \t%s\n", summary.High.String())
	fmt.Printf("\tLow:          \t%s\n", summary.Low.String())
	fmt.Printf("\tVolume:       \t%s\n", summary.Volume.String())
	fmt.Printf("\tQuoteVolume:  \t%s\n", summary.QuoteVolume.String())
	fmt.Printf("\tPercentChange:\t%s\n", summary.PercentChange.String())
	fmt.Printf("\tUpdatedAt:    \t%s\n", summary.UpdatedAt)

	orderBooks, err := client.GetOrderBook("ETH-USD", 0)
	if err != nil {
		return 1
	}
	fmt.Printf("OrderBooks ETH-USD (BID/ASK, Quantity, Rate):\n")
	for i := 0; i < 5; i++ {
		fmt.Printf("\tBID\t%s at %s\n", orderBooks.Bid[i].Quantity.String(), orderBooks.Bid[i].Rate.String())
		fmt.Printf("\tASK\t%s at %s\n", orderBooks.Ask[i].Quantity.String(), orderBooks.Ask[i].Rate.String())
	}

	trades, err := client.GetTrades("ETH-USD")
	if err != nil {
		return 1
	}
	fmt.Printf("Trades ETH-USD (TakerSide, Quantity, Rate):\n")
	for i := 0; i < 7; i++ {
		fmt.Printf("\t%s\t%s at %s\n", trades[i].TakerSide, trades[i].Quantity.String(), trades[i].Rate.String())
	}

	candlesOpts := &bittrex.GetCandlesOpts{
		CandleType: bittrex.CANDLETYPE_MIDPOINT,
	}
	candles, err := client.GetCandlesWithOpts("ETH-USD", bittrex.INTERVAL_DAY1, candlesOpts)
	if err != nil {
		return 1
	}
	fmt.Printf("Candles ETH-USD (StartsAt, Open, Close):\n")
	for i := 0; i < 3; i++ {
		fmt.Printf("\t%s\t%s %s\n", candles[i].StartsAt, candles[i].Open.String(), candles[i].Close.String())
	}

	candlesHistoryOpts := &bittrex.GetCandlesHistoryOpts{
		HistoryMonth: 10,
		HistoryDay:   3,
		CandleType:   bittrex.CANDLETYPE_MIDPOINT,
	}
	candlesHistory, err := client.GetCandlesHistoryWithOpts("ETH-USD", bittrex.INTERVAL_DAY1, 2021, candlesHistoryOpts)
	if err != nil {
		return 1
	}
	fmt.Printf("Historical Candles ETH-USD (StartsAt, Open, Close):\n")
	for i := 0; i < 3; i++ {
		fmt.Printf("\t%s\t%s %s\n", candlesHistory[i].StartsAt, candlesHistory[i].Open.String(), candlesHistory[i].Close.String())
	}

	// Ping

	ping, err := client.Ping()
	if err != nil {
		return 1
	}
	fmt.Printf("OK. Server time: %d\n", ping)

	return 0
}
