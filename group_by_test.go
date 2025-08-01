package itermore_test

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/ninedraft/itermore"
)

func ExampleGroupByFn() {
	groups := itermore.GroupByFn(
		itermore.Items("apple", "apricot", "banana", "blueberry"),
		regexp.MustCompile(`^([a-z])`).FindString,
	)

	for key, group := range groups {
		collected := slices.Collect(group)
		fmt.Println(key, collected)
	}

	// Output:
	// a [apple apricot]
	// b [banana blueberry]
}

type testGroup struct {
	key   string
	items []string
}

func newTg(key string, items ...string) testGroup {
	return testGroup{
		key:   key,
		items: items,
	}
}

func (tg testGroup) String() string {
	return fmt.Sprintf("%s=>[%v]", tg.key, strings.Join(tg.items, ", "))
}

func TestGroupByFn(t *testing.T) {
	toKey := func(str string) rune {
		ru, _ := utf8.DecodeRuneInString(str)
		return ru
	}

	tc := func(name string, input []string, expected []testGroup) {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var got []testGroup

			groups := itermore.GroupByFn(itermore.Slice(input), toKey)
			for key, gotGroup := range groups {
				collected := slices.Collect(gotGroup)
				got = append(got, testGroup{
					key:   string(key),
					items: collected,
				})
			}

			if len(got) != len(expected) {
				t.Errorf("want: %d groups", len(expected))
				t.Errorf("got: %d groups", len(got))
				t.Errorf("want: %+v", expected)
				t.Errorf("got:  %+v", got)
			}
		})
	}

	tc("empty", []string{}, []testGroup{})
	tc("nil", nil, []testGroup{})

	tc("single", []string{"apple"}, []testGroup{
		newTg("a", "apple"),
	})

	tc("alternating",
		[]string{"apple", "banana", "apricot", "blueberry"},
		[]testGroup{
			newTg("a", "apple"),
			newTg("b", "banana"),
			newTg("a", "apricot"),
			newTg("b", "blueberry"),
		},
	)

	tc("except last one",
		[]string{"apple", "apricot", "banana"},
		[]testGroup{
			newTg("a", "apple", "apricot"),
			newTg("b", "banana"),
		},
	)

	assertBreak2(t, itermore.GroupByFn(
		itermore.Items("apple", "apricot", "banana", "blueberry"),
		toKey))

	t.Run("stop early", func(t *testing.T) {
		t.Parallel()

		input := itermore.Items("apple", "apricot", "banana", "blueberry", "cherry")

		grouped := itermore.GroupByFn(input, toKey)

		i := 0
		for range grouped {
			i++
			if i == 3 {
				return
			}
		}

		i = 0
		for key, group := range grouped {
			i++

			collected := slices.Collect(group)
			if key != toKey("c") {
				t.Errorf("expected key 'c', got %q", key)
			}

			if !slices.Equal(collected, []string{"cherry"}) {
				t.Errorf("expected group to contain 'cherry', got %v", collected)
			}

			if i == 2 {
				t.Errorf("expected to stop after 1 iteration, but continued")
			}
		}
	})

}

func BenchmarkGropByFn(b *testing.B) {
	words := []string{
		"apple", "apricot", "banana", "blueberry",
	}

	const emitBytes = 10 * 1024 * 1024
	emitted := 0
	source := func(yield func(string) bool) {
		for {
			for _, word := range words {
				if !yield(word) {
					return
				}
				emitted += len(word)
				if emitted >= emitBytes {
					return
				}
			}
		}
	}

	toKey := func(str string) byte {
		return str[0]
	}

	b.ResetTimer()

	for b.Loop() {
		emitted = 0
		for _, group := range itermore.GroupByFn(source, toKey) {
			for range group {
				// consume group
			}
		}
	}
}
