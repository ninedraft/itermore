package itermore_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"iter"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/ninedraft/itermore"
)

func ExampleCollectJoin() {
	iter := itermore.Items("a", "b", "c")

	str := &strings.Builder{}

	_, _ = itermore.CollectJoin(str, iter, ", ")

	fmt.Println(str.String())

	// Output: a, b, c
}

func TestCollectJoin(t *testing.T) {
	t.Parallel()

	t.Run("happy", func(t *testing.T) {
		t.Parallel()

		iter := itermore.Items("a", "b", "c")
		str := &strings.Builder{}

		n, err := itermore.CollectJoin(str, iter, ", ")

		if err != nil {
			t.Fatal(err)
		}

		if n != 7 {
			t.Fatalf("got %d, want %d", n, 7)
		}

		got := str.String()
		want := "a, b, c"

		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		errTest := errors.New("test error")

		iter := itermore.Items("a", "b", "c")
		wr := &errWriter{err: errTest}

		_, err := itermore.CollectJoin(wr, iter, ", ")

		if !errors.Is(err, err) {
			t.Fatalf("got %v, want %v", err, err)
		}
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		iter := itermore.None[string]
		str := &strings.Builder{}

		n, err := itermore.CollectJoin(str, iter, ", ")

		if err != nil {
			t.Fatal(err)
		}

		if n != 0 {
			t.Errorf("got %d, want %d", n, 0)
		}

		if str.String() != "" {
			t.Errorf("got %q, want empty string", str.String())
		}
	})
}

func ExampleCollectJoinBytes() {
	iter := itermore.Items([]byte("a"), []byte("b"), []byte("c"))

	buf := &strings.Builder{}

	_, _ = itermore.CollectJoinBytes(buf, iter, []byte(", "))

	fmt.Println(buf.String())

	// Output: a, b, c
}

func TestCollectJoinBytes(t *testing.T) {
	t.Parallel()

	t.Run("happy", func(t *testing.T) {
		t.Parallel()

		iter := itermore.Items([]byte("a"), []byte("b"), []byte("c"))
		buf := &strings.Builder{}

		n, err := itermore.CollectJoinBytes(buf, iter, []byte(", "))

		if err != nil {
			t.Fatal(err)
		}

		if n != 7 {
			t.Fatalf("got %d, want %d", n, 7)
		}

		got := buf.String()
		want := "a, b, c"

		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		errTest := errors.New("test error")

		iter := itermore.Items([]byte("a"), []byte("b"), []byte("c"))
		wr := &errWriter{err: errTest}

		_, err := itermore.CollectJoinBytes(wr, iter, []byte(", "))

		if !errors.Is(err, err) {
			t.Fatalf("got %v, want %v", err, err)
		}
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		iter := itermore.None[[]byte]
		buf := &strings.Builder{}

		n, err := itermore.CollectJoinBytes(buf, iter, []byte(", "))

		if err != nil {
			t.Fatal(err)
		}

		if n != 0 {
			t.Errorf("got %d, want %d", n, 0)
		}

		if buf.String() != "" {
			t.Errorf("got %q, want empty string", buf.String())
		}
	})
}

func ExampleCollectJoinReaders() {
	iter := itermore.Items(strings.NewReader("a"), strings.NewReader("b"), strings.NewReader("c"))

	buf := &strings.Builder{}
	sep := []byte(", ")

	_, _ = itermore.CollectJoinReaders(buf, iter, sep)

	fmt.Println(buf.String())

	// Output: a, b, c
}

func TestCollectJoinReaders(t *testing.T) {
	t.Parallel()

	makeReaders := func(n int) iter.Seq[*strings.Reader] {
		return func(yield func(*strings.Reader) bool) {
			for i := range n {
				re := strings.NewReader(fmt.Sprintf("reader%d", i))
				if !yield(re) {
					return
				}
			}
		}
	}

	t.Run("happy", func(t *testing.T) {
		t.Parallel()

		iter := makeReaders(3)
		buf := &strings.Builder{}
		sep := []byte(", ")
		expected := "reader0, reader1, reader2"

		n, err := itermore.CollectJoinReaders(buf, iter, sep)
		if err != nil {
			t.Fatal(err, "while collecting join readers")
		}

		if n != int64(len(expected)) {
			t.Errorf("got %d, want %d", n, len(expected))
		}

		if expected != buf.String() {
			t.Errorf("got %q, want %q", buf.String(), expected)
		}
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		errTest := errors.New("test error")

		iter := itermore.Items(strings.NewReader("a"), strings.NewReader("b"), strings.NewReader("c"))
		wr := &errWriter{err: errTest}

		_, err := itermore.CollectJoinReaders(wr, iter, []byte(", "))

		if !errors.Is(err, errTest) {
			t.Fatalf("got %v, want %v", err, errTest)
		}
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		iter := itermore.None[*strings.Reader]
		buf := &strings.Builder{}
		sep := []byte(", ")

		n, err := itermore.CollectJoinReaders(buf, iter, sep)
		if err != nil {
			t.Fatal(err, "while collecting join readers")
		}

		if n != 0 {
			t.Errorf("got %d, want %d", n, 0)
		}

		if buf.String() != "" {
			t.Errorf("got %q, want empty string", buf.String())
		}
	})
}

type errWriter struct{ err error }

func (ew *errWriter) Write(p []byte) (int, error) {
	return 0, ew.err
}

func (ew *errWriter) WriteString(s string) (int, error) {
	return 0, ew.err
}

// --------------------------------
// Test suite for MultiReader, adapted from Golang multi_test.go
// Licensed under BSD-style license, which is located in contrib/go.LICENSE
// --------------------------------

// helper to build itermore.MultiReader over a variadic list of readers (returns io.ReadCloser)
func mkMultiReaderRC(rs ...io.Reader) io.ReadCloser {
	return itermore.MultiReader(itermore.Slice(rs))
}

// byteAndEOFReader reads one byte and returns EOF in a single Read call.
type byteAndEOFReader byte

func (b byteAndEOFReader) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		panic("unexpected call: zero-length Read")
	}
	p[0] = byte(b)
	return 1, io.EOF
}

func TestMultiReader(t *testing.T) {
	var mr io.ReadCloser
	var buf []byte
	nread := 0

	withFooBar := func(tests func()) {
		r1 := strings.NewReader("foo ")
		r2 := strings.NewReader("")
		r3 := strings.NewReader("bar")
		mr = mkMultiReaderRC(r1, r2, r3)
		t.Cleanup(func() { _ = mr.Close() })

		buf = make([]byte, 20)
		tests()
	}

	expectRead := func(size int, expected string, eerr error) {
		nread++
		n, gerr := mr.Read(buf[0:size])
		if n != len(expected) {
			t.Errorf("#%d, expected %d bytes; got %d", nread, len(expected), n)
		}
		got := string(buf[0:n])
		if got != expected {
			t.Errorf("#%d, expected %q; got %q", nread, expected, got)
		}
		if gerr != eerr {
			t.Errorf("#%d, expected error %v; got %v", nread, eerr, gerr)
		}
		buf = buf[n:]
	}

	withFooBar(func() {
		expectRead(2, "fo", nil)
		expectRead(5, "o ", nil)
		expectRead(5, "bar", nil)
		expectRead(5, "", io.EOF)
	})
	withFooBar(func() {
		expectRead(4, "foo ", nil)
		expectRead(1, "b", nil)
		expectRead(3, "ar", nil)
		expectRead(1, "", io.EOF)
	})
	withFooBar(func() {
		expectRead(5, "foo ", nil)
	})
}

// This guards an old bug: reading forever when a reader returns (1, EOF).
func TestMultiReaderSingleByteWithEOF(t *testing.T) {
	mr := mkMultiReaderRC(byteAndEOFReader('a'), byteAndEOFReader('b'))
	t.Cleanup(func() { _ = mr.Close() })

	got, err := io.ReadAll(io.LimitReader(mr, 10))
	if err != nil {
		t.Fatal(err)
	}
	const want = "ab"
	if string(got) != want {
		t.Errorf("got %q; want %q", got, want)
	}
}

// A reader returning (n, EOF) at the end continues to return EOF on its final read.
func TestMultiReaderFinalEOF(t *testing.T) {
	mr := mkMultiReaderRC(bytes.NewReader(nil), byteAndEOFReader('a'))
	t.Cleanup(func() { _ = mr.Close() })

	buf := make([]byte, 2)
	n, err := mr.Read(buf)
	if n != 1 || err != io.EOF {
		t.Errorf("got %v, %v; want 1, EOF", n, err)
	}
}

func TestMultiReaderFreesExhaustedReaders(t *testing.T) {
	var mr io.ReadCloser
	closed := make(chan struct{})

	// Avoid a live reference to buf1 on our stack after MultiReader is inlined.
	func() {
		buf1 := bytes.NewReader([]byte("foo"))
		buf2 := bytes.NewReader([]byte("bar"))
		mr = mkMultiReaderRC(buf1, buf2)
	}()
	t.Cleanup(func() { _ = mr.Close() })

	// Arrange cleanup callback tied to buf1's GC.
	func() {
		buf1 := bytes.NewReader([]byte("foo"))
		runtime.AddCleanup(buf1, func(ch chan struct{}) { close(ch) }, closed)
		_ = buf1 // ensure compiled in
	}()

	buf := make([]byte, 4)
	if n, err := io.ReadFull(mr, buf); err != nil || string(buf) != "foob" {
		t.Fatalf(`ReadFull = %d (%q), %v; want 4, "foob", nil`, n, buf[:n], err)
	}

	runtime.GC()
	select {
	case <-closed:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for collection of buf1")
	}

	if n, err := io.ReadFull(mr, buf[:2]); err != nil || string(buf[:2]) != "ar" {
		t.Fatalf(`ReadFull = %d (%q), %v; want 2, "ar", nil`, n, buf[:n], err)
	}
}

func TestInterleavedMultiReader(t *testing.T) {
	r1 := strings.NewReader("123")
	r2 := strings.NewReader("45678")

	mr1 := mkMultiReaderRC(r1, r2)
	t.Cleanup(func() { _ = mr1.Close() })
	mr2 := mkMultiReaderRC(mr1)
	t.Cleanup(func() { _ = mr2.Close() })

	buf := make([]byte, 4)

	// Consume via mr2 (which internally uses mr1's readers):
	n, err := io.ReadFull(mr2, buf)
	if got := string(buf[:n]); got != "1234" || err != nil {
		t.Errorf(`ReadFull(mr2) = (%q, %v), want ("1234", nil)`, got, err)
	}

	// Now consume the rest via mr1; should not panic even though mr2 advanced it:
	n, err = io.ReadFull(mr1, buf)
	if got := string(buf[:n]); got != "5678" || err != nil {
		t.Errorf(`ReadFull(mr1) = (%q, %v), want ("5678", nil)`, got, err)
	}
}
