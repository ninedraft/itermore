package itermore

import (
	"context"
	"iter"
	"sync/atomic"
	"time"
)

// Tick creates a sequence that yields time.Time values each dt interval.
// It will panic if dt is non-positive.
// It will prevent goroutine leak if the sequence is not fully consumed.
func Tick(dt time.Duration) iter.Seq[time.Time] {
	return func(yield func(time.Time) bool) {
		ticker := time.NewTicker(dt)
		defer ticker.Stop()

		YieldFrom(yield, Chan(ticker.C))
	}
}

// TickCtx creates a sequence that yields time.Time values each dt interval.
// It will panic if dt is non-positive.
// It will prevent goroutine leak if the sequence is not fully consumed.
// It will stop the sequence when the given context is canceled.
func TickCtx(ctx context.Context, dt time.Duration) iter.Seq[time.Time] {
	return func(yield func(time.Time) bool) {
		ticker := time.NewTicker(dt)
		defer ticker.Stop()

		YieldFrom(yield, ChanCtx(ctx, ticker.C))
	}
}

// Timer creates and immediately starts a timer.
// Returned seq emits timestamps and a reset function, which can be used to set next timer timeout.
// If reset function is not called, then seq is stopped.
// This function cleanup timer after seq is consumed or stopped.
func Timer(dt time.Duration) iter.Seq2[time.Time, func(time.Duration)] {
	return func(yield func(time.Time, func(time.Duration)) bool) {
		timer := time.NewTimer(dt)
		isResetted := &atomic.Bool{}
		reset := func(dt time.Duration) {
			isResetted.Store(true)
			timer.Reset(dt)
		}
		if timer != nil {
			defer timer.Stop()
		}

		for tick := range timer.C {
			isResetted.Store(false)

			if !yield(tick, reset) {
				break
			}
			if !isResetted.Load() {
				break
			}
		}
	}
}
