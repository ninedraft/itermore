package itermore

import (
	"iter"
	"slices"
)

// Slice creates a sequence that yields values from the given slice.
// If slice is nil or empty, the sequence will be empty.
func Slice[E any](items []E) iter.Seq[E] {
	return func(yield func(E) bool) {
		for _, value := range items {
			if !yield(value) {
				return
			}
		}
	}
}

// Items creates a sequence that yields values from the given variadic arguments.
// If no arguments are provided, the sequence will be empty.
// It's a shortcut for Slice([]E{items...}).
func Items[E any](items ...E) iter.Seq[E] {
	return Slice(items)
}

// Loop forever yields values from the given slice in the order they appear in the slice.
// After the last value is yielded, the sequence will start from the beginning.
func Loop[E any](items []E) iter.Seq[E] {
	if len(items) == 0 {
		return None[E]
	}

	return func(yield func(E) bool) {
		for i := 0; ; i = (i + 1) % len(items) {
			if !yield(items[i]) {
				return
			}
		}
	}
}

// Collect writes values from provided sequence to the given slice.
// If dst slice is nil, Collect will create a new slice.
// It can return a new slice or the same slice that was passed as dst following the same rules as append function.
func Collect[S ~[]E, E any](dst S, seq iter.Seq[E]) S {
	return slices.AppendSeq(dst, seq)
}

// FlattenSlices walks through provided sequence and emits each element from sequenced slice.
//
// Example:
//
//	[1, 2, 3] [4, 5, 6] -> 1, 2, 3, 4, 5, 6
func FlattenSlices[E any](seq iter.Seq[[]E]) iter.Seq[E] {
	return func(yield func(E) bool) {
		for slice := range seq {
			if !YieldFrom(yield, Slice(slice)) {
				return
			}
		}
	}
}
