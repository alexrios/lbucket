package lbucket

import (
	"time"
)

// BucketTicker is a the way bucket interact with tickers
type BucketTicker interface {
	Stop()
	Receiver() <-chan time.Time
}

// RealTicker is the default ticker laying over the stdlib ticker.
type realTicker struct {
	*time.Ticker
}

func (f realTicker) Receiver() <-chan time.Time {
	return f.C
}
