package generator

import (
	"fmt"
	"math/big"

	"github.com/peterzeller/go-fun/iterable"
	"github.com/peterzeller/go-fun/linked"
	"github.com/peterzeller/go-stateful-test/generator/shrink"
)

func Slice[T any](elemGen Generator[T]) Generator[[]T] {
	return &sliceGen[T]{
		elemGen: elemGen,
	}
}

type sliceGen[T any] struct {
	elemGen Generator[T]
}

func (s *sliceGen[T]) Enumerate(depth int) iterable.Iterable[[]T] {
	return EnumerateSlices(depth, depth, s.elemGen)
}

func EnumerateSlices[T any](length, depth int, elemGen Generator[T]) iterable.Iterable[[]T] {
	if length <= 0 {
		return iterable.Singleton([]T{})
	}
	smallerSlices := EnumerateSlices(length-1, depth, elemGen)
	return iterable.Concat(
		smallerSlices,
		iterable.FlatMap(
			smallerSlices,
			func(tail []T) iterable.Iterable[[]T] {
				// append only to the longest lists
				if len(tail) < length-1 {
					return iterable.Empty[[]T]()
				}
				return iterable.Map(
					elemGen.Enumerate(depth),
					func(head T) []T {
						return append([]T{head}, tail...)
					})
			},
		),
	)
}

func (s *sliceGen[T]) Name() string {
	return fmt.Sprintf("SliceGen[%s]", s.elemGen.Name())
}

func (s *sliceGen[T]) RValue(elem RandomValue[[]T]) ([]T, bool) {
	rvs, ok := elem.Value.([]RandomValue[T])
	if !ok {
		return nil, false
	}
	res := make([]T, len(rvs))
	for i, rv := range rvs {
		res[i], ok = s.elemGen.RValue(rv)
		if !ok {
			return nil, false
		}
	}
	return res, true
}

func (s *sliceGen[T]) Random(rnd Rand, size int) RandomValue[[]T] {
	if size <= 0 {
		return RandomValue[[]T]{Value: []RandomValue[T]{}}
	}
	l := rnd.R().Intn(size)
	res := make([]RandomValue[T], l)
	for i := range res {
		res[i] = s.elemGen.Random(rnd, size-1)
	}
	return RandomValue[[]T]{Value: res}
}

func (s *sliceGen[T]) Shrink(elem RandomValue[[]T]) iterable.Iterable[RandomValue[[]T]] {
	rvs := elem.Value.([]RandomValue[T])
	return iterable.Map(
		shrink.ShrinkList(
			linked.New(rvs...),
			func(rv RandomValue[T]) iterable.Iterable[RandomValue[T]] {
				return s.elemGen.Shrink(rv)
			}),
		func(l *linked.List[RandomValue[T]]) RandomValue[[]T] {
			return RandomValue[[]T]{Value: l.ToSlice()}
		})
}

func (s *sliceGen[T]) Size(t RandomValue[[]T]) *big.Int {
	var size big.Int
	rvs := t.Value.([]RandomValue[T])
	for _, rv := range rvs {
		size.Add(&size, s.elemGen.Size(rv))
	}
	return &size
}
