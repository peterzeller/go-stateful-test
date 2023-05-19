package generator

import (
	"github.com/peterzeller/go-fun/iterable"
	"github.com/peterzeller/go-fun/zero"
	"github.com/peterzeller/go-stateful-test/generator/geniterable"
	"math/big"
)

// Empty generator that generates no values
func Empty[T, R any]() Generator[T, R] {
	return emptyGenerator[T, R]{}
}

type emptyGenerator[T, R any] struct {
}

func (c emptyGenerator[T, R]) Name() string {
	return "empty"
}

func (c emptyGenerator[T, R]) Random(rnd Rand, size int) (res R) {
	return
}

func (c emptyGenerator[T, R]) Enumerate(depth int) geniterable.Iterable[R] {
	return geniterable.Empty[R]()
}

func (c emptyGenerator[T, R]) Shrink(elem R) iterable.Iterable[R] {
	return iterable.Empty[R]()
}

func (c emptyGenerator[T, R]) Size(t R) *big.Int {
	return big.NewInt(1)
}

func (c emptyGenerator[T, R]) RValue(elem R) (T, bool) {
	return zero.Value[T](), false
}
