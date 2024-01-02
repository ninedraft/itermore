package itermore

import "iter"

// Zip creates a sequence that yields pairs of values from the given sequences.
// If any of the sequences is stopped, the sequence will stop.
func Zip[A, B any](a iter.Seq[A], b iter.Seq[B]) iter.Seq2[A, B] {
	return func(yield func(A, B) bool) {
		nextA, cancelA := iter.Pull(a)
		defer cancelA()

		nextB, cancelB := iter.Pull(b)
		defer cancelB()

		for {
			av, okA := nextA()
			bv, okB := nextB()

			if !okA || !okB {
				break
			}

			if !yield(av, bv) {
				break
			}
		}
	}
}

// ZipLongest creates a sequence that yields pairs of values from the given sequences.
// If any of the sequences is stopped, the sequence will continue to emit zero values for that sequence.
// If both sequences are stopped, the sequence will stop.
func ZipLongest[A, B any](a iter.Seq[A], b iter.Seq[B]) iter.Seq2[A, B] {
	return func(yield func(A, B) bool) {
		nextA, cancelA := iter.Pull(a)
		defer cancelA()

		nextB, cancelB := iter.Pull(b)
		defer cancelB()

		for {
			av, okA := nextA()
			bv, okB := nextB()

			if !okA && !okB {
				break
			}

			if !yield(av, bv) {
				break
			}
		}
	}
}
