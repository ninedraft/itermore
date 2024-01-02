package itermore

import "iter"

// Map creates a sequence that yields key-value pairs from the given map.
// If map nil or empty, the sequence will be empty.
func Map[K comparable, V any](m map[K]V) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range m {
			if !yield(k, v) {
				return
			}
		}
	}
}

// CollectMap writes key-value pairs from provided sequence to the given map.
// It will panic if dst is nil.
func CollectMap[M ~map[K]V, K comparable, V any](dst M, seq iter.Seq2[K, V]) {
	for k, v := range seq {
		dst[k] = v
	}
}

// CollectKeys writes keys from provided sequence to the given map and assigns them the given value.
// It can be used to create a set from a sequence.
// It will panic if dst is nil.
func CollectKeys[S ~map[K]V, K comparable, V any](dst S, value V, seq iter.Seq[K]) {
	for k := range seq {
		dst[k] = value
	}
}
