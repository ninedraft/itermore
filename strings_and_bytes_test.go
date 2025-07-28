package itermore_test

import (
	"errors"
	"fmt"
	"iter"
	"strings"
	"testing"

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
