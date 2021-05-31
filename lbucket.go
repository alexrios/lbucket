package lbucket

import (
	"errors"
	"sync"
	"time"
)

var ErrBucketReachedCap = errors.New("bucket reached its max capacity")

// TickLeakyBucket represents a leaky bucket backed by time.ticker
type TickLeakyBucket struct {
	currentVolume uint
	fullBucket    uint
	mu            sync.Mutex
	// guarantees that the Fix method will be effectively called only once
	fixOnce sync.Once
	// signaling the ticker to stop
	done chan struct{}
}

// NewTickLeakyBucket creates a new leaky bucket with a predetermined capacity and leak itself based on freq
// Call fix after create a new bucket to ensure that the bucket will not leak forever
func NewTickLeakyBucket(capacity uint, freq time.Duration, opts ...tickerOption) *TickLeakyBucket {
	bucket := TickLeakyBucket{
		currentVolume: 0,
		fullBucket:    capacity,
		mu:            sync.Mutex{},
		done:          make(chan struct{}),
		fixOnce:       sync.Once{},
	}

	// apply options
	to := tickerOptions{}
	for _, opt := range opts {
		opt(&to)
	}

	// Define the ticker
	var ticker BucketTicker
	switch {
	case to.bt != nil:
		ticker = to.bt
	default:
		ticker = realTicker{Ticker: time.NewTicker(freq)}
	}

	go func() {
		defer ticker.Stop()

		rcv := ticker.Receiver()
		for {
			select {
			case <-bucket.done:
				return
			case _ = <-rcv:
				bucket.leak()
			}
		}
	}()

	return &bucket
}

func (c *TickLeakyBucket) leak() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.currentVolume > 0 {
		c.currentVolume--
	}
}

// Size returns the current volume inside the bucket
func (c *TickLeakyBucket) Size() uint {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.currentVolume
}

// Refill tries to fill the bucket again.
// returns ErrBucketReachedCap when it reaches total capacity.
func (c *TickLeakyBucket) Refill() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.currentVolume < c.fullBucket {
		c.currentVolume++
		return nil
	}
	return ErrBucketReachedCap
}

// Fix stop a bucket from leak ever again
func (c *TickLeakyBucket) Fix() {
	c.fixOnce.Do(func() {
		c.done <- struct{}{}
	})
}

type tickerOptions struct {
	bt BucketTicker
}

type tickerOption func(opt *tickerOptions)

// WithCustomTicker can change the behavior with a custom ticker according to BucketTicker interface
func WithCustomTicker(ticker BucketTicker) func(opt *tickerOptions) {
	return func(opt *tickerOptions) {
		opt.bt = ticker
		return
	}
}
