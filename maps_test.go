package itermore_test

import (
	"maps"
	"testing"

	"github.com/ninedraft/itermore"
)

func TestMap(t *testing.T) {
	t.Parallel()

	want := map[int]string{
		1: "a", 2: "b", 3: "c",
	}

	assertBreak2(t, itermore.Map(want))

	t.Run("iter", func(t *testing.T) {
		t.Parallel()

		got := map[int]string{}

		for a, b := range itermore.Map(want) {
			got[a] = b
		}

		if !maps.Equal(got, want) {
			t.Errorf("got:  %v", got)
			t.Errorf("want: %v", want)
		}
	})

	t.Run("nil", func(t *testing.T) {
		t.Parallel()

		for a, b := range itermore.Map[int, string](nil) {
			t.Fatalf("must not iterate over nil seq, got: %v, %v", a, b)
		}
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		for a, b := range itermore.Map(map[int]string{}) {
			t.Fatalf("must not iterate over empty seq, got: %v, %v", a, b)
		}
	})
}

func TestCollectMap(t *testing.T) {
	t.Parallel()

	t.Run("iter", func(t *testing.T) {
		t.Parallel()

		want := map[int]string{
			1: "a", 2: "b", 3: "c",
		}
		got := map[int]string{}

		itermore.CollectMap(got, itermore.Map(want))

		if !maps.Equal(got, want) {
			t.Errorf("got:  %v", got)
			t.Errorf("want: %v", want)
		}
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		got := map[int]string{}
		itermore.CollectMap(got, itermore.None2[int, string]())

		if !maps.Equal(got, map[int]string{}) {
			t.Errorf("got:  %v", got)
			t.Errorf("want: %v", map[int]string{})
		}
	})
}

func TestCollectMapKeys(t *testing.T) {
	t.Parallel()

	t.Run("iter", func(t *testing.T) {
		t.Parallel()

		want := map[int]struct{}{
			1: {}, 2: {}, 3: {},
		}
		got := map[int]struct{}{}

		itermore.CollectKeys(got, struct{}{}, itermore.Items(1, 2, 3))

		if !maps.Equal(got, want) {
			t.Errorf("got:  %v", got)
			t.Errorf("want: %v", want)
		}
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		got := map[int]struct{}{}
		itermore.CollectKeys(got, struct{}{}, itermore.None[int]())

		if !maps.Equal(got, map[int]struct{}{}) {
			t.Errorf("got:  %v", got)
			t.Errorf("want: %v", map[int]struct{}{})
		}
	})
}
