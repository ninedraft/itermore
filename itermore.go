package itermore

import (
	"cmp"
	"iter"
	"sync/atomic"

	"golang.org/x/exp/constraints"
)

// None creates an empty sequence.
func None[E any](yield func(E) bool) {}

// None2 creates an empty sequence of pairs.
func None2[A, B any](yield func(A, B) bool) {}

// One creates a sequence that yields a single value.
func One[E any](value E) iter.Seq[E] {
	isDrained := &atomic.Bool{}

	return func(yield func(E) bool) {
		if isDrained.Load() {
			return
		}
		isDrained.Store(true)

		yield(value)
	}
}

// One2 creates a sequence that yields a single pair.
func One2[A, B any](a A, b B) iter.Seq2[A, B] {
	isDrained := &atomic.Bool{}
	return func(yield func(A, B) bool) {
		if isDrained.Load() {
			return
		}
		isDrained.Store(true)

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
	isDrained := &atomic.Bool{}

	return func(yield func(E) bool) {
		for !isDrained.Load() {
			value, ok := next()
			if !ok {
				isDrained.Store(true)
				return
			}

			if !yield(value) {
				return
			}
		}
	}
}

// Next2 creates sequence that yields pairs of values from the given function.
// It is similar to Next, but yields two values at once.
func Next2[A, B any](next func() (A, B, bool)) iter.Seq2[A, B] {
	isDrained := &atomic.Bool{}

	return func(yield func(A, B) bool) {
		for !isDrained.Load() {
			a, b, ok := next()
			if !ok {
				isDrained.Store(true)
				return
			}

			if !yield(a, b) {
				return
			}
		}
	}
}

// YieldFrom pulls all values from the given sequence and yields them.
// Returns false if the yield function returns false.
func YieldFrom[E any](yield func(E) bool, seq iter.Seq[E]) bool {
	for value := range seq {
		if !yield(value) {
			return false
		}
	}

	return true
}

// YieldFrom2 pulls all pairs from the given sequence and yields them.
// Returns false if the yield function returns false.
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
		for i, seq := range seqs {
			if seq == nil {
				// skip nil and drained
				continue
			}

			if !YieldFrom(yield, seq) {
				return
			}
			seqs[i] = nil // mark as drained
		}
	}
}

// Chain2 creates a sequence that yields pairs from all provided sequences.
// Pairs are yielded in the order they appear in the arguments.
func Chain2[A, B any](seqs ...iter.Seq2[A, B]) iter.Seq2[A, B] {
	return func(yield func(A, B) bool) {
		for i, seq := range seqs {
			if seq == nil {
				// skip nil and drained
				continue
			}
			if !YieldFrom2(yield, seq) {
				return
			}
			seqs[i] = nil // mark as drained
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
	return func(yield func(N) bool) {
		forImpl(start, to, step, yield)
	}
}

func forImpl[N Number](start, to, step N, yield func(N) bool) {
	if step == 0 {
		panic("step cannot be zero")
	}

	if step > 0 && start > to {
		return
	}
	if step < 0 && start < to {
		return
	}

	stop := func(v N) bool {
		return v >= to
	}

	if step < 0 {
		stop = func(v N) bool {
			return v <= to
		}
	}

	for v := start; ; {
		if stop(v) {
			return
		}

		if !yield(v) {
			return
		}

		next := v + step

		// Overflow detection: for a positive step, overflow wraps around ↓; for
		// a negative step, it wraps around ↑.  The `<`/`>` tests catch both
		// signed and unsigned kinds.
		if (step > 0 && next < v) || (step < 0 && next > v) {
			return
		}

		v = next
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

// KeysOf creates a sequence that yields keys from the given sequence of pairs.
// It is useful to extract keys from a sequence of pairs.
func KeysOf[K, V any](seq iter.Seq2[K, V]) iter.Seq[K] {
	return func(yield func(K) bool) {
		for k := range seq {
			if !yield(k) {
				return
			}
		}
	}
}

// ValuesOf creates a sequence that yields values from the given sequence of pairs.
func ValuesOf[K, V any](seq iter.Seq2[K, V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range seq {
			if !yield(v) {
				return
			}
		}
	}
}

// Drain input sequence into oblivion.
// It is useful to clear sequences that are not needed anymore.
func Drain[E any](seq iter.Seq[E]) {
	for range seq {
	}
}
