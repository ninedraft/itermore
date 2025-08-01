package itermore_test

import (
	"fmt"
	"regexp"
	"slices"

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
