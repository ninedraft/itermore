package itermore

import (
	"iter"
)

// GroupByFn creates a sequence of groups from the given sequence.
// Each group is a sequence of values that share the same key, which is determined by the
// provided function `toKey`. The keys are yielded in the order they first appear in the sequence.
// The function `toKey` should return a comparable key for each value in the sequence.
//
// This operation is similar to UNIX `groupby` command. It groups items from input sequence,
// emitting a new group each time the key changes whule iterating over input sequence.
// This is different from SQL `GROUP BY` operation, which groups all items from dataset at once.
//
// Each output sequence shares the same memory with the input sequence, so it state is invalidated
// on next iteration. Collect group to slice for later usage or use it immediately.
func GroupByFn[E any, K comparable](seq iter.Seq[E], toKey func(E) K) iter.Seq2[K, iter.Seq[E]] {
	return func(yield func(K, iter.Seq[E]) bool) {

		next, nextStop := iter.Pull(seq)
		defer nextStop()

		start, ok := next()
		if !ok {
			return
		}

		key := toKey(start)

		stopGroup := false
		inputDrained := false
		group := func(yield func(E) bool) {
			if stopGroup {
				return
			}

			if !yield(start) {
				return
			}

			for {
				value, ok := next()
				if !ok {
					inputDrained = true
					return
				}

				key = toKey(value)
				if key != toKey(start) {
					start = value
					key = toKey(value)
					stopGroup = true
					return
				}

				if !yield(value) {
					return
				}
			}
		}

		for !inputDrained {
			stopGroup = false
			ok = yield(key, group)
			if !ok {
				return
			}
		}
	}
}

// Group creates a sequence of groups from the given sequence.
// Each group consists of equal values from the sequence.
//
// The values are yielded in the order they first appear in the sequence.
// This operation is similar to UNIX `groupby` command. It groups items from input sequence,
// emitting a new group each time the value changes while iterating over input sequence.
//
// This is different from SQL `GROUP BY` operation, which groups all items from dataset at once.
// Each output sequence shares the same memory with the input sequence, so it state is invalidated
// on next iteration.
//
// Collect group to slice for later usage or use it immediately.
func Group[E comparable](seq iter.Seq[E]) iter.Seq[iter.Seq[E]] {
	return func(yield func(iter.Seq[E]) bool) {
		grouped := GroupByFn(seq, func(value E) E {
			return value
		})

		YieldFrom(yield, ValuesOf(grouped))
	}
}

func GroupByKey[K comparable, V any](seq iter.Seq2[K, V]) iter.Seq2[K, iter.Seq[V]] {
	return func(yield func(K, iter.Seq[V]) bool) {
		next, nextStop := iter.Pull2(seq)
		defer nextStop()

		key, start, ok := next()
		if !ok {
			return
		}

		for stop := false; !stop; {
			groupTail := Next(func() (V, bool) {
				var empty V

				nextKey, value, ok := next()

				if !ok {
					stop = true
					return empty, false
				}

				if nextKey != key {

					start = value
					key = nextKey
					return empty, false
				}

				return empty, true
			})

			group := Chain(
				One(start),
				groupTail,
			)

			ok = yield(key, group)

			Drain(group)

			if !ok {
				return
			}
		}
	}
}
