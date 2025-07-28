package itermore

import (
	"io"
	"iter"
)

// CollectJoin writes values from provided sequence to the given writer.
func CollectJoin[S ~string](wr io.StringWriter, seq iter.Seq[S], sep string) (int64, error) {
	head := true
	written := int64(0)
	for value := range seq {
		if !head {
			n, err := wr.WriteString(sep)
			written += int64(n)
			if err != nil {
				return written, err
			}
		}

		head = false

		n, err := wr.WriteString(string(value))
		written += int64(n)
		if err != nil {
			return written, err
		}
	}

	return written, nil
}

// CollectJoinBytes writes values from provided sequence to the given writer.
func CollectJoinBytes[P ~[]byte](wr io.Writer, seq iter.Seq[P], sep []byte) (int64, error) {
	head := true
	written := int64(0)
	for value := range seq {
		if !head {
			n, err := wr.Write(sep)
			written += int64(n)
			if err != nil {
				return written, err
			}
		}

		head = false

		n, err := wr.Write(value)
		written += int64(n)
		if err != nil {
			return written, err
		}
	}

	return written, nil
}

const defaultBufferSize = 32 * 1024

func CollectJoinReaders[R io.Reader](wr io.Writer, seq iter.Seq[R], sep []byte) (int64, error) {
	head := true
	written := int64(0)

	buf := make([]byte, defaultBufferSize)

	for re := range seq {
		if !head {
			n, err := wr.Write(sep)
			written += int64(n)
			if err != nil {
				return written, err
			}
		}

		head = false

		n, err := io.CopyBuffer(wr, re, buf)
		written += n
		if err != nil {
			return written, err
		}
	}

	return written, nil
}
