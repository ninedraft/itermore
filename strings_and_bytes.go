package itermore

import (
	"errors"
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

// MultiReader returns an io.ReadCloser that's the logical concatenation of the
// readers from the provided sequence. They're read sequentially. Once all
// readers have returned EOF, Read will return EOF. If any of the readers
// return a non-nil, non-EOF error, Read will return that error.
// It is an analog of io.MultiReader, but for iter.Seq[io.Reader].
// The caller should call Close on the returned reader to release resources
// associated with the iterator.
func MultiReader(rere iter.Seq[io.Reader]) io.ReadCloser {
	next, stop := iter.Pull(rere)

	return &multiReader{
		next: next,
		stop: stop,
	}
}

type multiReader struct {
	re   io.Reader
	next func() (io.Reader, bool)
	stop func()
}

func (mr *multiReader) Close() error {
	mr.stop()
	return nil
}

func (mr *multiReader) Read(p []byte) (n int, err error) {
next:
	if mr.re == nil {
		re, ok := mr.next()
		if !ok {
			// stop the iterator earlier to help GC
			mr.stop()
			return n, io.EOF
		}

		mr.re = re
	}

	n, err = mr.re.Read(p)
	isEOF := errors.Is(err, io.EOF)

	if isEOF {
		mr.re = nil
	}

	if isEOF && n == 0 {
		goto next
	}

	if isEOF && n != 0 {
		return n, nil
	}

	return n, err
}
