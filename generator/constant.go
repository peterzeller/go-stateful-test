package generator

import (
	"github.com/peterzeller/go-fun/iterable"
	"github.com/peterzeller/go-stateful-test/generator/geniterable"
	"math/big"
)

// Constant generator that always returns the same value
func Constant[T any](c T) Generator[T, T] {
	return constGenerator[T]{
		c: c,
	}
}

type constGenerator[T any] struct {
	c T
}

func (c constGenerator[T]) Name() string {
	return "const"
}

func (c constGenerator[T]) Random(rnd Rand, size int) T {
	return c.c
}

func (c constGenerator[T]) Enumerate(depth int) geniterable.Iterable[T] {
	return geniterable.Singleton(c.c)
}

func (c constGenerator[T]) Shrink(elem T) iterable.Iterable[T] {
	return iterable.Empty[T]()
}

func (c constGenerator[T]) Size(t T) *big.Int {
	return big.NewInt(1)
}

func (c constGenerator[T]) RValue(elem T) (T, bool) {
	return c.c, true
}
