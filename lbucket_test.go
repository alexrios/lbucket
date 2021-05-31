package lbucket

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"
)

func TestNewTickLeakyBucket(t *testing.T) {
	bucket := NewTickLeakyBucket(0, 1*time.Second, WithCustomTicker(NewFakeTicker()))
	defer bucket.Fix()

	size := bucket.Size()
	if size != 0 {
		t.Fatalf("bucket size should be zero but it was %d", size)
	}
}

func TestRefill(t *testing.T) {
	t.Run("should refill the bucket just until the capacity", func(t *testing.T) {
		fTicker := NewFakeTicker()
		bucket := NewTickLeakyBucket(1, 1*time.Second, WithCustomTicker(fTicker))
		defer bucket.Fix()

		// try to refill
		err := bucket.Refill()
		if err != nil {
			log.Fatalln(err)
		}
		size := bucket.Size()
		if size != 1 {
			t.Fatalf("bucket size should be 1 but it was %d", size)
		}
		// Should fail
		err = bucket.Refill()
		if !errors.Is(err, ErrBucketReachedCap) {
			t.Fatalf("bucket size should be exceeded")
		}
	})
	t.Run("should be able to refill in the same pace it leaks", func(t *testing.T) {
		fTicker := NewFakeTicker()
		bucket := NewTickLeakyBucket(1, 1*time.Second, WithCustomTicker(fTicker))
		defer bucket.Fix()

		// try to refill
		err := bucket.Refill()
		if err != nil {
			log.Fatalln(err)
		}
		// Leak once
		fTicker.NextTick()

		// refilling again
		err = bucket.Refill()
		if err != nil {
			log.Fatalln(err)
		}

		// Should fail
		err = bucket.Refill()
		if !errors.Is(err, ErrBucketReachedCap) {
			t.Fatalf("bucket size should be exceeded")
		}
	})
}

func TestFix(t *testing.T) {
	t.Run("should never leak after fix", func(t *testing.T) {
		fTicker := NewFakeTicker()
		bucket := NewTickLeakyBucket(3, 1*time.Second, WithCustomTicker(fTicker))

		bucket.Fix()

		if fTicker.IsStopped() {
			t.Fatal("this ticker should be stopped")
		}
	})
}

// FakeTicker is solely used on test suites.
type FakeTicker struct {
	isStopped bool
	once      sync.Once
	TimeC     chan time.Time
	done      chan struct{}
}

func NewFakeTicker() *FakeTicker {
	return &FakeTicker{
		isStopped: false,
		once:      sync.Once{},
		TimeC:     make(chan time.Time),
		done:      make(chan struct{}),
	}
}

func (f *FakeTicker) Stop() {
	f.once.Do(func() {
		f.done <- struct{}{}
	})
}

func (f *FakeTicker) IsStopped() bool {
	for {
		select {
		case <-f.done:
			return false
		default:
			return true
		}
	}
}

func (f *FakeTicker) Receiver() <-chan time.Time {
	return f.TimeC
}

func (f *FakeTicker) NextTick() {
	// since it's not important for testing, any time will work
	f.TimeC <- time.Now()
}

func Example() {
	bucket := NewTickLeakyBucket(3, time.Second)

	for i := 0; i < 10; i++ {
		time.Sleep(500 * time.Millisecond)
		err := bucket.Refill()
		if errors.Is(err, ErrBucketReachedCap) {
			fmt.Println("FULL CAP! cannot refill")
			continue
		}
		fmt.Println(bucket.Size())
	}
}
