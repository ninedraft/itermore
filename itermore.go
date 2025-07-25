package itermore

import (
	"cmp"
	"iter"

	"golang.org/x/exp/constraints"
)

// None creates an empty sequence.
func None[E any](yield func(E) bool) {}

// None2 creates an empty sequence of pairs.
func None2[A, B any](yield func(A, B) bool) {}

// One creates a sequence that yields a single value.
func One[E any](value E) iter.Seq[E] {
	return func(yield func(E) bool) {
		yield(value)
	}
}

// One2 creates a sequence that yields a single pair.
func One2[A, B any](a A, b B) iter.Seq2[A, B] {
	return func(yield func(A, B) bool) {
		yield(a, b)
	}
}

// Forever creates an infinite sequence that yields a single value.
func Forever[E any](value E) iter.Seq[E] {
	return func(yield func(E) bool) {
		for {
			if !yield(value) {
				return
			}
		}
	}
}

// ForeverFn creates an infinite sequence that yields a single value from the given function.
func ForeverFn[E any](fn func() E) iter.Seq[E] {
	return func(yield func(E) bool) {
		for {
			value := fn()
			if !yield(value) {
				return
			}
		}
	}
}

// Enumerate creates a sequence that yields values from the given sequence with their indexes.
// If number of values in the sequence is greater then math.MaxInt, counter will wrap around.
func Enumerate[E any](seq iter.Seq[E]) iter.Seq2[int, E] {
	return func(yield func(int, E) bool) {
		var i int
		for value := range seq {
			if !yield(i, value) {
				return
			}
			i++
		}
	}
}

// Next creates sequence that yields values from the given function.
func Next[E any](next func() (E, bool)) iter.Seq[E] {
	return func(yield func(E) bool) {
		for {
			value, ok := next()
			if !ok {
				return
			}
			if !yield(value) {
				return
			}
		}
	}
}

// YieldFrom pulls all values from the given sequence and yields them.
func YieldFrom[E any](yield func(E) bool, seq iter.Seq[E]) bool {
	for value := range seq {
		if !yield(value) {
			return false
		}
	}

	return true
}

// YieldFrom2 pulls all pairs from the given sequence and yields them.
func YieldFrom2[A, B any](yield func(A, B) bool, seq iter.Seq2[A, B]) bool {
	for a, b := range seq {
		if !yield(a, b) {
			return false
		}
	}

	return true
}

// Chain creates a sequence that yields values from all provided sequences.
// Values are yielded in the order they appear in the arguments.
func Chain[E any](seqs ...iter.Seq[E]) iter.Seq[E] {
	return func(yield func(E) bool) {
		for _, seq := range seqs {
			if !YieldFrom(yield, seq) {
				return
			}
		}
	}
}

// SkipN skips first n values from the given sequence.
// If n is greater then number of values in the sequence, SkipN will return an empty sequence.
// If n is negative or zero, SkipN will return an empty sequence.
func SkipN[E any](n int, seq iter.Seq[E]) iter.Seq[E] {
	return func(yield func(E) bool) {
		next, cancel := iter.Pull(seq)
		defer cancel()

		for i := 0; i < n; i++ {
			_, ok := next()
			if !ok {
				return
			}
		}

		for {
			value, ok := next()
			if !ok {
				return
			}
			if !yield(value) {
				return
			}
		}
	}
}

// TakeN yields first n values from the given sequence.
// If n is greater then number of values in the sequence, TakeN will yield all values from the sequence.
// If n is negative or zero, TakeN will return an empty sequence.
func TakeN[E any](n int, seq iter.Seq[E]) iter.Seq[E] {
	return func(yield func(E) bool) {
		next, cancel := iter.Pull(seq)
		defer cancel()

		for i := 0; i < n; i++ {
			value, ok := next()
			if !ok {
				return
			}
			if !yield(value) {
				return
			}
		}
	}
}

// Number is a type, which can be added and compared.
type Number interface {
	constraints.Integer | constraints.Float
}

// For behaves like a for loop with a counter.
// It can be though as following code:
//
//	start <= to: for i := start; i < to; i += step
//	start > to:  for i := start; i > to; i += step
func For[N Number](start, to, step N) iter.Seq[N] {
	if start <= to {
		return func(yield func(N) bool) {
			forIncr(start, to, step, yield)
		}
	}

	return func(yield func(N) bool) {
		forDecr(start, to, step, yield)
	}
}

func forIncr[N Number](start, to, step N, yield func(N) bool) {
	for i := start; i < to; i += step {
		if !yield(i) {
			return
		}
	}
}

func forDecr[N Number](start, to, step N, yield func(N) bool) {
	for i := start; i > to; i += step {
		if !yield(i) {
			return
		}
	}
}

// Then returns a sequence, which yields values from seq, then always calls then function.
// It calls then function even if seq is empty.
func Then[E any](seq iter.Seq[E], then func()) iter.Seq[E] {
	return func(yield func(E) bool) {
		defer then()
		YieldFrom(yield, seq)
	}
}

// Compact returns a sequence, which yields values from seq, but skips consecutive duplicates.
// It roughly equal to slices.Compact.
func Compact[E comparable](seq iter.Seq[E]) iter.Seq[E] {
	return func(yield func(E) bool) {
		var prev E
		ok := false
		for value := range seq {
			if ok && prev == value {
				continue
			}
			prev, ok = value, true
			if !yield(value) {
				return
			}
		}
	}
}

// Max returns the largest element in the sequence.
// If sequence is empty, Max returns false.
// If there is a single item in the sequence, Max returns it.
// It roughly equal to slices.Max.
func Max[E cmp.Ordered](seq iter.Seq[E]) (E, bool) {
	var x E
	ok := false

	for value := range seq {
		if !ok {
			x, ok = value, true
		}
		x = max(x, value)
	}

	return x, ok
}

// Min returns the smallest element in the sequence.
// If sequence is empty, Min returns false.
// If there is a single item in the sequence, Min returns it.
// It roughly equal to slices.Min.
func Min[E cmp.Ordered](seq iter.Seq[E]) (E, bool) {
	var x E
	ok := false

	for value := range seq {
		if !ok {
			x, ok = value, true
		}
		x = min(x, value)
	}

	return x, ok
}

// PairsPadded creates a sequence that yields pairs of values from the given sequence.
// If number of values in the sequence is odd, Pairs will pad the last pair with the given value.
func PairsPadded[E any](seq iter.Seq[E], pad E) iter.Seq2[E, E] {
	return func(yield func(E, E) bool) {
		next, cancel := iter.Pull(seq)
		defer cancel()

		for {
			a, ok := next()
			if !ok {
				return
			}

			b, ok := next()

			if !ok {
				b = pad
			}

			if !yield(a, b) {
				return
			}
		}
	}
}

// Pairs creates a sequence that yields pairs of values from the given sequence.
// If number of values in the sequence is odd, Pairs will skip the last value.
func Pairs[E any](seq iter.Seq[E]) iter.Seq2[E, E] {
	return func(yield func(E, E) bool) {
		next, cancel := iter.Pull(seq)
		defer cancel()

		for {
			a, ok := next()
			if !ok {
				return
			}

			b, ok := next()
			if !ok {
				return
			}

			if !yield(a, b) {
				return
			}
		}
	}
}
