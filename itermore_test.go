package itermore_test

import (
	"iter"
	"slices"
	"testing"

	"github.com/ninedraft/itermore"
)

type pair struct {
	a int
	b string
}

func TestNone(t *testing.T) {
	t.Parallel()

	assertBreak(t, itermore.None[int])

	t.Run("iter", func(t *testing.T) {
		t.Parallel()

		for _ = range itermore.None[int] {
			t.Fatal("no iterations are expected")
		}
	})
}

func TestNone2(t *testing.T) {
	t.Parallel()

	assertBreak2(t, itermore.None2[int, string])

	t.Run("iter", func(t *testing.T) {
		t.Parallel()

		for _ = range itermore.None[int] {
			t.Fatal("no iterations are expected")
		}
	})
}

func TestOne(t *testing.T) {
	t.Parallel()

	assertBreak(t, itermore.One(1))

	t.Run("iter", func(t *testing.T) {
		t.Parallel()

		got := []int{}
		for x := range itermore.One(1) {
			got = append(got, x)
		}

		want := []int{1}
		if !slices.Equal(got, want) {
			t.Errorf("got:  %v", got)
			t.Errorf("want: %v", want)
		}
	})
}

func TestOne2(t *testing.T) {
	t.Parallel()

	assertBreak2(t, itermore.One2(1, "a"))

	t.Run("iter", func(t *testing.T) {
		t.Parallel()

		got := []pair{}
		for a, b := range itermore.One2(1, "a") {
			got = append(got, pair{a, b})
		}

		want := []pair{{1, "a"}}
		if !slices.Equal(got, want) {
			t.Errorf("got:  %v", got)
			t.Errorf("want: %v", want)
		}
	})
}

func TestForever(t *testing.T) {
	t.Parallel()

	assertBreak(t, itermore.Forever(1))

	t.Run("iter", func(t *testing.T) {
		t.Parallel()

		got := []int{}
		for x := range itermore.Forever(1) {
			got = append(got, x)
			if len(got) >= 10 {
				break
			}
		}

		want := []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1}

		if !slices.Equal(got, want) {
			t.Errorf("got:  %v", got)
			t.Errorf("want: %v", want)
		}
	})
}

func TestForeverFn(t *testing.T) {
	t.Parallel()

	assertBreak(t, itermore.ForeverFn(func() int { return 1 }))

	t.Run("iter", func(t *testing.T) {
		t.Parallel()

		fn := func() int { return 1 }
		got := []int{}
		for x := range itermore.ForeverFn(fn) {
			got = append(got, x)
			if len(got) >= 10 {
				break
			}
		}

		want := []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1}

		if !slices.Equal(got, want) {
			t.Errorf("got:  %v", got)
			t.Errorf("want: %v", want)
		}
	})
}

func TestEnumerate(t *testing.T) {
	t.Parallel()

	{
		seq := itermore.Slice([]int{1, 2, 3})
		assertBreak2(t, itermore.Enumerate(seq))
	}

	t.Run("iter", func(t *testing.T) {
		t.Parallel()

		input := itermore.Slice([]string{
			"a", "b", "c", "d", "e",
		})

		got := []pair{}
		for i, v := range itermore.Enumerate(input) {
			got = append(got, pair{i, v})
		}

		want := []pair{
			{0, "a"}, {1, "b"}, {2, "c"}, {3, "d"}, {4, "e"},
		}

		if !slices.Equal(got, want) {
			t.Errorf("got:  %v", got)
			t.Errorf("want: %v", want)
		}
	})
}

func TestNext(t *testing.T) {
	t.Parallel()

	{
		next := func() (int, bool) { return 1, true }
		assertBreak(t, itermore.Next(next))
	}

	t.Run("iter", func(t *testing.T) {
		t.Parallel()

		i := 10
		next := func() (int, bool) {
			i--
			if i < 0 {
				return 0, false
			}
			return i, true
		}

		got := []int{}
		for x := range itermore.Next(next) {
			got = append(got, x)
		}

		want := []int{9, 8, 7, 6, 5, 4, 3, 2, 1, 0}
		if !slices.Equal(got, want) {
			t.Errorf("got:  %v", got)
			t.Errorf("want: %v", want)
		}
	})
}

func TestYieldFrom(t *testing.T) {
	t.Parallel()

	newSeq := func() iter.Seq[int] {
		a := itermore.Slice([]int{1, 2, 3})
		b := itermore.Slice([]int{4, 5, 6})

		return func(yield func(int) bool) {
			if !itermore.YieldFrom(yield, a) {
				return
			}

			if !itermore.YieldFrom(yield, b) {
				return
			}
		}
	}

	assertBreak(t, newSeq())

	t.Run("iter", func(t *testing.T) {
		t.Parallel()

		got := []int{}
		for x := range newSeq() {
			got = append(got, x)
		}

		want := []int{1, 2, 3, 4, 5, 6}
		if !slices.Equal(got, want) {
			t.Errorf("got:  %v", got)
			t.Errorf("want: %v", want)
		}
	})
}

func TestYieldFrom2(t *testing.T) {
	t.Parallel()

	newSeq := func() iter.Seq2[int, string] {
		a := itermore.Slice([]string{"a", "b", "c"})
		b := itermore.Slice([]string{"d", "e", "f"})

		return func(yield func(int, string) bool) {
			if !itermore.YieldFrom2(yield, itermore.Enumerate(a)) {
				return
			}

			if !itermore.YieldFrom2(yield, itermore.Enumerate(b)) {
				return
			}
		}
	}

	assertBreak2(t, newSeq())

	t.Run("iter", func(t *testing.T) {
		t.Parallel()

		got := []pair{}
		for a, b := range newSeq() {
			got = append(got, pair{a, b})
		}

		want := []pair{
			{0, "a"}, {1, "b"}, {2, "c"},
			{0, "d"}, {1, "e"}, {2, "f"},
		}

		if !slices.Equal(got, want) {
			t.Errorf("got:  %v", got)
			t.Errorf("want: %v", want)
		}
	})
}

func TestChain(t *testing.T) {
	t.Parallel()

	newSeq := func() iter.Seq[int] {
		a := itermore.Slice([]int{1, 2, 3})
		b := itermore.Slice([]int{4, 5, 6})

		return itermore.Chain(a, b)
	}

	assertBreak(t, newSeq())

	t.Run("iter", func(t *testing.T) {
		t.Parallel()

		got := []int{}
		for x := range newSeq() {
			got = append(got, x)
		}

		want := []int{1, 2, 3, 4, 5, 6}
		if !slices.Equal(got, want) {
			t.Errorf("got:  %v", got)
			t.Errorf("want: %v", want)
		}
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		for _ = range itermore.Chain[int]() {
			t.Fatalf("must not iterate over empty seq")
		}
	})
}

func TestSkipN(t *testing.T) {
	t.Parallel()

	newSeq := func() iter.Seq[int] {
		a := itermore.Slice([]int{1, 2, 3, 4, 5, 6})
		return itermore.SkipN(2, a)
	}

	assertBreak(t, newSeq())

	t.Run("iter", func(t *testing.T) {
		t.Parallel()

		got := []int{}
		for x := range newSeq() {
			got = append(got, x)
		}

		want := []int{3, 4, 5, 6}
		if !slices.Equal(got, want) {
			t.Errorf("got:  %v", got)
			t.Errorf("want: %v", want)
		}
	})

	t.Run("short", func(t *testing.T) {
		t.Parallel()

		seq := itermore.Slice([]int{1, 2})

		for value := range itermore.SkipN(10, seq) {
			t.Errorf("must not iterate over short seq, got: %v", value)
		}
	})
}

func TestTakeN(t *testing.T) {
	t.Parallel()

	newSeq := func() iter.Seq[int] {
		a := itermore.Slice([]int{1, 2, 3, 4, 5, 6})
		return itermore.TakeN(2, a)
	}

	assertBreak(t, newSeq())

	t.Run("iter", func(t *testing.T) {
		t.Parallel()

		got := []int{}
		for x := range newSeq() {
			got = append(got, x)
		}

		want := []int{1, 2}
		if !slices.Equal(got, want) {
			t.Errorf("got:  %v", got)
			t.Errorf("want: %v", want)
		}
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		a := itermore.Slice([]int{})
		for _ = range itermore.TakeN(2, a) {
			t.Fatalf("must not iterate over empty seq")
		}
	})

	t.Run("zero", func(t *testing.T) {
		t.Parallel()

		a := itermore.Slice([]int{1, 2, 3})
		for _ = range itermore.TakeN(0, a) {
			t.Fatalf("must not iterate over empty seq")
		}
	})

	t.Run("negative", func(t *testing.T) {
		t.Parallel()

		want := []int{1, 2, 3}
		a := itermore.Slice(want)

		for _ = range itermore.TakeN(-1, a) {
			t.Fatalf("must not iterate over empty seq")
		}
	})
}

func TestFor(t *testing.T) {
	t.Parallel()

	{
		seq := itermore.For(0, 10, 1)
		assertBreak(t, seq)
	}

	{
		seq := itermore.For(10, 0, -1)
		assertBreak(t, seq)
	}

	t.Run("iter", func(t *testing.T) {
		t.Parallel()

		seq := itermore.For(0, 10, 1)

		got := itermore.Collect[[]int](nil, seq)

		want := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

		if !slices.Equal(got, want) {
			t.Errorf("got:  %v", got)
			t.Errorf("want: %v", want)
		}
	})

	t.Run("iter-reverse", func(t *testing.T) {
		t.Parallel()

		seq := itermore.For(10, 0, -1)

		got := itermore.Collect[[]int](nil, seq)

		want := []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}

		if !slices.Equal(got, want) {
			t.Errorf("got:  %v", got)
			t.Errorf("want: %v", want)
		}
	})
}

func TestThen(t *testing.T) {
	t.Parallel()

	t.Run("non-empty", func(t *testing.T) {
		t.Parallel()

		ok := false
		then := func() {
			ok = true
		}

		seq := itermore.Then(itermore.One(1), then)

		for _ = range seq {
			// pass
		}

		if !ok {
			t.Errorf("then must be called")
		}
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		ok := false
		then := func() {
			ok = true
		}

		seq := itermore.Then(itermore.None[int], then)

		for _ = range seq {
			t.Fatalf("must not iterate over empty seq")
		}

		if !ok {
			t.Errorf("then must be called")
		}
	})
}

func TestCompact(t *testing.T) {
	t.Parallel()

	{
		seq := itermore.Items(1, 1, 2, 2, 3, 3, 4, 5)
		assertBreak(t, itermore.Compact(seq))
	}

	t.Run("non-empty", func(t *testing.T) {
		input := []int{1, 1, 2, 2, 3, 3, 4, 5}
		seq := itermore.Slice(input)
		compacted := itermore.Compact(seq)

		got := itermore.Collect[[]int](nil, compacted)

		want := slices.Compact(slices.Clone(input))

		if !slices.Equal(got, want) {
			t.Errorf("got:  %v", got)
			t.Errorf("want: %v", want)
		}
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		seq := itermore.Compact(itermore.None[int])

		for _ = range seq {
			t.Fatalf("must not iterate over empty seq")
		}
	})
}

func TestMax(t *testing.T) {
	t.Parallel()

	t.Run("non-empty", func(t *testing.T) {
		t.Parallel()

		seq := itermore.Items(2, 3, 1, 5, 4, 5)
		got, ok := itermore.Max(seq)

		want := 5
		if got != want {
			t.Errorf("got:  %v", got)
			t.Errorf("want: %v", want)
		}

		if !ok {
			t.Errorf("ok must be true")
		}
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		x, ok := itermore.Max(itermore.None[int])

		if x != 0 {
			t.Errorf("got:  %v", x)
			t.Errorf("want: %v", 0)
		}

		if ok {
			t.Errorf("ok must be false")
		}
	})
}

func TestMin(t *testing.T) {
	t.Parallel()

	t.Run("non-empty", func(t *testing.T) {
		t.Parallel()

		seq := itermore.Items(2, 3, 1, 4, 1, 5)
		got, ok := itermore.Min(seq)

		want := 1
		if got != want {
			t.Errorf("got:  %v", got)
			t.Errorf("want: %v", want)
		}

		if !ok {
			t.Errorf("ok must be true")
		}
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		x, ok := itermore.Min(itermore.None[int])

		if x != 0 {
			t.Errorf("got:  %v", x)
			t.Errorf("want: %v", 0)
		}

		if ok {
			t.Errorf("ok must be false")
		}
	})
}

func assertBreak[E any](t *testing.T, seq iter.Seq[E]) {
	t.Helper()

	t.Run("break", func(t *testing.T) {
		t.Parallel()

		for _ = range seq {
			break
		}
	})
}

func assertBreak2[A, B any](t *testing.T, seq iter.Seq2[A, B]) {
	t.Helper()

	t.Run("break", func(t *testing.T) {
		t.Parallel()

		for _, _ = range seq {
			break
		}
	})
}
