package itermore

import (
	"context"
	"iter"
)

// Chan returns a new sequence that iterates over values from the given channel.
// If channel is nil, Chan returns an empty sequence.
// Closing channel will stop iteration.
func Chan[E any](ch <-chan E) iter.Seq[E] {
	return func(yield func(E) bool) {
		if ch == nil {
			return
		}

		for value := range ch {
			if !yield(value) {
				return
			}
		}
	}
}

// ChanCtx returns a new sequence that iterates over values from the given channel.
// If channel is nil, Chan returns an empty sequence.
// Closing channel will stop iteration.
// If ctx is canceled, ChanCtx will stop iteration.
func ChanCtx[E any](ctx context.Context, ch <-chan E) iter.Seq[E] {
	return func(yield func(E) bool) {
		if ch == nil {
			return
		}

		for {
			select {
			case <-ctx.Done():
				return
			case value, ok := <-ch:
				if !ok {
					return
				}
				if !yield(value) {
					return
				}
			}
		}
	}
}

func CollectChan[E any](ch chan<- E, seq iter.Seq[E]) {
	for value := range seq {
		ch <- value
	}
}

func CollectChanCtx[E any](ctx context.Context, ch chan<- E, seq iter.Seq[E]) {
	for value := range seq {
		select {
		case <-ctx.Done():
			return
		case ch <- value:
			// pass
		}
	}
}
