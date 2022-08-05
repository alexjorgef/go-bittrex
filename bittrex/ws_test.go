package bittrex

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTradeStream_SubscribeTradeUpdates(t *testing.T) {
	bt := New("", "")
	ch := make(chan StreamTrade)
	errCh := make(chan error)
	stopCh := make(chan bool)
	go func() {
		errCh <- bt.SubscribeTradeUpdates("BTC-USD", ch, stopCh)
	}()
	var err error
	var trade StreamTrade
	select {
	case trade = <-ch:
	case err = <-errCh:
	case <-time.NewTicker(3 * time.Minute).C:
		stopCh <- true
		err = errors.New("timeout")
	}
	assert.NoError(t, err)
	assert.NotEmpty(t, trade.Symbol)
	assert.NotEmpty(t, trade.ID)
	assert.NotEmpty(t, trade.TakerSide)
	assert.NotEmpty(t, trade.ExecutedAt)
	assert.NotEmpty(t, trade.Rate)
	assert.NotEmpty(t, trade.Quantity)
	assert.Greater(t, trade.Rate.IntPart(), int64(0))
}
