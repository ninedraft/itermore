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
		if !head && len(sep) > 0 {
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
		if !head && len(sep) > 0 {
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

	var buf []byte

	copyBuf := func(re io.Reader) error {
		n, err := io.CopyBuffer(wr, re, buf)
		written += n

		return err
	}

	var readFrom func(re io.Reader) error
	readFrom = func(re io.Reader) error {
		buf = make([]byte, defaultBufferSize)
		readFrom = copyBuf

		return readFrom(re)
	}

	// I expect that in 80% of cases the writer is either *bytes.Buffer or *bufio.Writer,
	// so we can use io.ReaderFrom without allocating a buffer.
	if reFrom, ok := wr.(io.ReaderFrom); ok {
		readFrom = func(re io.Reader) error {
			n, err := reFrom.ReadFrom(re)
			written += n

			return err
		}
	}

	for re := range seq {
		if !head && len(sep) > 0 {
			n, err := wr.Write(sep)
			written += int64(n)
			if err != nil {
				return written, err
			}
		}

		head = false

		err := readFrom(re)
		if err != nil {
			return written, err
		}
	}

	return written, nil
}
