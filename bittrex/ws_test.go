package bittrex

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTradeStream_SubscribeTradeUpdates(t *testing.T) {
	bt := New("", "")
	ch := make(chan Trade)
	errCh := make(chan error)
	stopCh := make(chan bool)
	go func() {
		errCh <- bt.SubscribeTradeUpdates("BTC-USD", ch, stopCh)
	}()
	var err error
	var trade Trade
	select {
	case trade = <-ch:
	case err = <-errCh:
	case <-time.NewTicker(3 * time.Minute).C:
		stopCh <- true
		err = errors.New("timeout")
	}
	assert.NoError(t, err)
	assert.Greater(t, trade.Rate.IntPart(), int64(0))
}
