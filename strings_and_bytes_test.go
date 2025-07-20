package itermore_test

import (
	"errors"
	"fmt"
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
}

type errWriter struct{ err error }

func (ew *errWriter) Write(p []byte) (int, error) {
	return 0, ew.err
}

func (ew *errWriter) WriteString(s string) (int, error) {
	return 0, ew.err
}
