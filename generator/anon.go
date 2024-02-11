package generator

import (
	"github.com/peterzeller/go-fun/iterable"
	"github.com/peterzeller/go-stateful-test/generator/geniterable"
	"math/big"
)

type AnonGenerator[T, R any] struct {
	GenName      string
	GenRandom    func(rnd Rand, size int) R
	GenShrink    func(elem R) iterable.Iterable[R]
	GenSize      func(t R) *big.Int
	GenRValue    func(r R) (T, bool)
	GenEnumerate func(depth int) geniterable.Iterable[R]
}

func (a *AnonGenerator[T, R]) Name() string {
	return a.GenName
}

func (a *AnonGenerator[T, R]) Random(rnd Rand, size int) R {
	return a.GenRandom(rnd, size)
}

func (a *AnonGenerator[T, R]) Enumerate(depth int) geniterable.Iterable[R] {
	return a.GenEnumerate(depth)
}

func (a *AnonGenerator[T, R]) Shrink(elem R) iterable.Iterable[R] {
	return a.GenShrink(elem)
}

func (a *AnonGenerator[T, R]) Size(t R) *big.Int {
	return a.GenSize(t)
}

func (a *AnonGenerator[T, R]) RValue(t R) (T, bool) {
	return a.GenRValue(t)
}
