package itermore_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/ninedraft/itermore"
)

func ExampleZip() {
	xx := itermore.Items(10, 20, 30)
	yy := itermore.Items("a", "b", "c")

	for x, y := range itermore.Zip(xx, yy) {
		fmt.Println(x, y)
	}
	// Output: 10 a
	// 20 b
	// 30 c
}

func TestZip(t *testing.T) {
	t.Parallel()

	{
		a, b := itermore.One(1), itermore.One("a")
		seq := itermore.Zip(a, b)

		assertBreak2(t, seq)
	}

	t.Run("none", func(t *testing.T) {
		t.Parallel()
		a, b := itermore.None[int], itermore.None[string]
		for _, _ = range itermore.Zip(a, b) {
			t.Fatal("must not iterate over empty seq")
		}
	})

	t.Run("none-one", func(t *testing.T) {
		t.Parallel()
		a, b := itermore.One(1), itermore.None[string]
		for _, _ = range itermore.Zip(a, b) {
			t.Fatal("must not iterate over empty seq")
		}
	})

	t.Run("many-many", func(t *testing.T) {
		t.Parallel()
		a := itermore.Slice([]int{1, 2, 3})
		b := itermore.Slice([]string{"a", "b", "c", "d"})

		got := []pair{}
		expect := []pair{
			{1, "a"}, {2, "b"}, {3, "c"},
		}

		for av, bv := range itermore.Zip(a, b) {
			got = append(got, pair{av, bv})
		}

		if !slices.Equal(got, expect) {
			t.Errorf("got:  %v", got)
			t.Errorf("want: %v", expect)
		}
	})
}

func TestZipLongest(t *testing.T) {
	t.Parallel()

	{
		a, b := itermore.One(1), itermore.One("a")
		seq := itermore.ZipLongest(a, b)

		assertBreak2(t, seq)
	}

	t.Run("none", func(t *testing.T) {
		t.Parallel()
		a, b := itermore.None[int], itermore.None[string]
		for _, _ = range itermore.ZipLongest(a, b) {
			t.Fatal("must not iterate over empty seq")
		}
	})

	t.Run("none-one", func(t *testing.T) {
		t.Parallel()
		a, b := itermore.One(1), itermore.None[string]

		got := []pair{}
		expect := []pair{{1, ""}}

		for av, bv := range itermore.ZipLongest(a, b) {
			got = append(got, pair{av, bv})
		}

		if !slices.Equal(got, expect) {
			t.Errorf("got:  %v", got)
			t.Errorf("want: %v", expect)
		}
	})

	t.Run("many-many", func(t *testing.T) {
		t.Parallel()
		a := itermore.Slice([]int{1, 2, 3})
		b := itermore.Slice([]string{"a", "b", "c", "d"})

		got := []pair{}
		expect := []pair{
			{1, "a"}, {2, "b"}, {3, "c"}, {0, "d"},
		}

		for av, bv := range itermore.ZipLongest(a, b) {
			got = append(got, pair{av, bv})
		}

		if !slices.Equal(got, expect) {
			t.Errorf("got:  %v", got)
			t.Errorf("want: %v", expect)
		}
	})
}
