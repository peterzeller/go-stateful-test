package generator

import (
	"github.com/peterzeller/go-fun/iterable"
	"math/big"
)

type AnonGenerator[T any] struct {
	GenName      string
	GenRandom    func(rnd Rand, size int) RandomValue[T]
	GenShrink    func(elem RandomValue[T]) iterable.Iterable[RandomValue[T]]
	GenSize      func(t RandomValue[T]) *big.Int
	GenRValue    func(r RandomValue[T]) (T, bool)
	GenEnumerate func(depth int) iterable.Iterable[T]
}

func (a *AnonGenerator[T]) Name() string {
	return a.GenName
}

func (a *AnonGenerator[T]) Random(rnd Rand, size int) RandomValue[T] {
	return a.GenRandom(rnd, size)
}

func (a *AnonGenerator[T]) Enumerate(depth int) iterable.Iterable[T] {
	return a.GenEnumerate(depth)
}

func (a *AnonGenerator[T]) Shrink(elem RandomValue[T]) iterable.Iterable[RandomValue[T]] {
	return a.GenShrink(elem)
}

func (a *AnonGenerator[T]) Size(t RandomValue[T]) *big.Int {
	return a.GenSize(t)
}

func (a *AnonGenerator[T]) RValue(t RandomValue[T]) (T, bool) {
	return a.GenRValue(t)
}
