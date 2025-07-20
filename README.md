# itermore 

[![Go Reference](https://pkg.go.dev/badge/github.com/ninedraft/itermore.svg)](https://pkg.go.dev/github.com/ninedraft/itermore)

range iterables with more features

## Description

This package provides a set of functions to work with iterables in Go. It is partially inspired by Python's itertools module.

Available utils include for now:
- iterator combinators: `Chain`, `Zip`, `Skip`, `Take`, etc.
- iterator destructors: `Collect`, `CollectJoin`, `CollectMap`, etc.
- iterator constructors: `Items`, `For`, `ChanCtx`, etc.

To use this package you need go 1.24.22 or later.

## Examples

```go
func ExampleCollectJoin() {
	iter := itermore.Items("a", "b", "c")

	str := &strings.Builder{}

	itermore.CollectJoin(str, iter, ", ")

	fmt.Println(str.String())

	// Output: a, b, c
}
```

```go
func ExampleChain() {
	one := itermore.One(1)
	xx := itermore.Items(10, 20, 30)

	for x := range itermore.Chain(one, xx) {
		fmt.Println(x)
	}

	// Output: 1
	// 10
	// 20
	// 30
}
```

```go
func ExampleZip() {
	xx := itermore.Items(10, 20, 30)
	yy := itermore.Items("a", "b", "c")

	for x, y := range itermore.Zip(xx, yy) {
		fmt.Println(x, y)
	}
	// Output: 10 a
	// 20 b
	// 30 c
}
```

## Roadmap

- [ ] group by
- [ ] iterable versions of standard library functions
