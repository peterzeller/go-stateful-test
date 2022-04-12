package generator

import (
	"math/big"

	"github.com/peterzeller/go-fun/iterable"
)

// Constant generator that always returns the same value
func Constant[T any](c T) Generator[T] {
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

func (c constGenerator[T]) Random(rnd Rand, size int) RandomValue[T] {
	return R(c.c)
}

func (c constGenerator[T]) Enumerate(depth int) iterable.Iterable[T] {
	return iterable.Singleton(c.c)
}

func (c constGenerator[T]) Shrink(elem RandomValue[T]) iterable.Iterable[RandomValue[T]] {
	return iterable.Empty[RandomValue[T]]()
}

func (c constGenerator[T]) Size(t RandomValue[T]) *big.Int {
	return big.NewInt(1)
}

func (c constGenerator[T]) RValue(elem RandomValue[T]) (T, bool) {
	return c.c, true
}
