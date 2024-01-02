package itermore_test

import (
	"iter"
	"slices"
	"testing"

	"github.com/ninedraft/itermore"
)

func TestSlice(t *testing.T) {
	t.Parallel()

	assertBreak(t, itermore.Slice([]int{1, 2, 3}))

	t.Run("iter", func(t *testing.T) {
		t.Parallel()

		got := []int{}
		want := []int{1, 2, 3}

		for x := range itermore.Slice(want) {
			got = append(got, x)
		}

		if !slices.Equal(got, want) {
			t.Errorf("got:  %v", got)
			t.Errorf("want: %v", want)
		}
	})

	t.Run("nil", func(t *testing.T) {
		t.Parallel()

		for value := range itermore.Slice[int](nil) {
			t.Fatalf("must not iterate over nil seq, got: %v", value)
		}
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		for value := range itermore.Slice([]int{}) {
			t.Fatalf("must not iterate over empty seq, got: %v", value)
		}
	})
}

func TestLoop(t *testing.T) {
	t.Parallel()

	assertBreak(t, itermore.Loop([]int{1, 2, 3}))

	t.Run("iter", func(t *testing.T) {
		t.Parallel()

		want := []int{1, 2, 3, 1, 2, 3}
		input := want[:3]

		next, stop := iter.Pull(itermore.Loop(input))
		defer stop()

		got := []int{}
		for range len(want) {
			x, ok := next()
			if !ok {
				t.Fatalf("must not stop iteration")
			}

			got = append(got, x)
		}

		if !slices.Equal(got, want) {
			t.Errorf("got:  %v", got)
			t.Errorf("want: %v", want)
		}
	})
}

func TestCollect(t *testing.T) {
	t.Parallel()

	t.Run("slice", func(t *testing.T) {
		t.Parallel()

		want := []int{1, 2, 3}
		got := itermore.Collect([]int{}, itermore.Slice(want))

		if !slices.Equal(got, want) {
			t.Errorf("got:  %v", got)
			t.Errorf("want: %v", want)
		}
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		want := []int{}
		got := itermore.Collect([]int{}, itermore.Slice(want))

		if !slices.Equal(got, want) {
			t.Errorf("got:  %v", got)
			t.Errorf("want: %v", want)
		}
	})
}
