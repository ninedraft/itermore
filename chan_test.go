package itermore_test

import (
	"context"
	"slices"
	"testing"

	"github.com/ninedraft/itermore"
)

func TestChan(t *testing.T) {
	t.Parallel()

	{
		ch := make(chan int, 1)
		ch <- 10
		assertBreak(t, itermore.Chan(ch))
	}

	t.Run("iter", func(t *testing.T) {
		t.Parallel()

		want := []int{1, 2, 3}
		input := make(chan int)
		go func() {
			defer close(input)
			for _, x := range want {
				input <- x
			}
		}()

		got := []int{}
		for x := range itermore.Chan(input) {
			got = append(got, x)
		}

		if !slices.Equal(got, want) {
			t.Errorf("got:  %v", got)
			t.Errorf("want: %v", want)
		}
	})

	t.Run("closed", func(t *testing.T) {
		t.Parallel()

		input := make(chan int)
		close(input)

		for x := range itermore.Chan(input) {
			t.Fatalf("must not iterate over closed chan, got: %v", x)
		}
	})

	t.Run("nil", func(t *testing.T) {
		t.Parallel()

		var ch chan int
		for x := range itermore.Chan(ch) {
			t.Fatalf("must not iterate over closed chan, got: %v", x)
		}
	})
}

func TestChanCtx(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	{
		ch := make(chan int, 1)
		ch <- 10
		assertBreak(t, itermore.ChanCtx(ctx, ch))
	}

	t.Run("iter", func(t *testing.T) {
		t.Parallel()

		want := []int{1, 2, 3}
		input := make(chan int)
		go func() {
			defer close(input)
			for _, x := range want {
				input <- x
			}
		}()

		got := []int{}
		for x := range itermore.ChanCtx(ctx, input) {
			got = append(got, x)
		}

		if !slices.Equal(got, want) {
			t.Errorf("got:  %v", got)
			t.Errorf("want: %v", want)
		}
	})

	t.Run("closed", func(t *testing.T) {
		t.Parallel()

		input := make(chan int)
		close(input)

		for x := range itermore.ChanCtx(ctx, input) {
			t.Fatalf("must not iterate over closed chan, got: %v", x)
		}
	})

	t.Run("nil", func(t *testing.T) {
		t.Parallel()

		var ch chan int
		for x := range itermore.ChanCtx(ctx, ch) {
			t.Fatalf("must not iterate over closed chan, got: %v", x)
		}
	})

	t.Run("ctx-cancel", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(ctx)
		cancel()

		input := make(chan int)
		for x := range itermore.ChanCtx(ctx, input) {
			t.Fatalf("must not iterate over closed chan, got: %v", x)
		}
	})
}

func TestCollectChan(t *testing.T) {
	t.Parallel()

	t.Run("iter", func(t *testing.T) {
		t.Parallel()

		want := []int{1, 2, 3}
		seq := itermore.Slice(want)

		ch := make(chan int, len(want))
		itermore.CollectChan(ch, seq)
		close(ch)

		got := []int{}
		for x := range ch {
			got = append(got, x)
		}

		if !slices.Equal(got, want) {
			t.Errorf("got:  %v", got)
			t.Errorf("want: %v", want)
		}
	})
}

func TestCollectChanCtx(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("iter", func(t *testing.T) {
		t.Parallel()

		want := []int{1, 2, 3}
		seq := itermore.Slice(want)

		ch := make(chan int, len(want))
		itermore.CollectChanCtx(ctx, ch, seq)
		close(ch)

		got := []int{}
		for x := range ch {
			got = append(got, x)
		}

		if !slices.Equal(got, want) {
			t.Errorf("got:  %v", got)
			t.Errorf("want: %v", want)
		}
	})

	t.Run("ctx-cancel", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(ctx)
		cancel()

		ch := make(chan int)
		itermore.CollectChanCtx(ctx, ch, itermore.Slice([]int{1, 2, 3}))
		close(ch)

		for x := range ch {
			t.Fatalf("must not iterate over closed chan, got: %v", x)
		}
	})
}
