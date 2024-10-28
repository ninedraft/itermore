package itermore_test

import (
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/ninedraft/itermore"
)

func TestTick(t *testing.T) {
	defer assertGoroutineLeak(t)()

	i := 0
	for _ = range itermore.Tick(time.Millisecond) {
		i++
		if i > 10 {
			break
		}
	}
}

func TestTickCtx(t *testing.T) {
	defer assertGoroutineLeak(t)()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	i := 0
	for _ = range itermore.TickCtx(ctx, time.Millisecond) {
		i++
		if i > 10 {
			cancel()
		}
	}
}

func TestTimer(t *testing.T) {
	t.Run("no reset", func(t *testing.T) {
		defer assertGoroutineLeak(t)()

		for range itermore.Timer(time.Millisecond) {
		}
	})

	t.Run("several iterations", func(t *testing.T) {
		defer assertGoroutineLeak(t)()

		prev, i := time.Now(), 0
		const N = 4

		for tp, reset := range itermore.Timer(time.Millisecond) {
			if !tp.After(prev) {
				t.Fatalf("")
			}
			if i < N {
				reset(time.Millisecond)
				i++
			}
		}

		if i != N {
			t.Fatalf("want 4 iterations, got %d", N)
		}
	})
}

func assertGoroutineLeak(t *testing.T) func() {
	t.Helper()

	start := runtime.NumGoroutine()
	return func() {
		end := runtime.NumGoroutine()
		if start < end {
			t.Fatalf("goroutines leak: %d", end-start)
		}
	}
}
