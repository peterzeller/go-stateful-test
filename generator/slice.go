package generator

import (
	"fmt"
	"github.com/peterzeller/go-fun/iterable"
	"github.com/peterzeller/go-stateful-test/generator/geniterable"
	"math/big"

	"github.com/peterzeller/go-fun/equality"
	"github.com/peterzeller/go-fun/list/linked"
	"github.com/peterzeller/go-fun/slice"
	"github.com/peterzeller/go-stateful-test/generator/shrink"
)

// Slice is a generator for slices.
func Slice[T any](elemGen Generator[T]) Generator[[]T] {
	return &sliceGen[T]{
		elemGen: elemGen,
	}
}

type sliceGen[T any] struct {
	elemGen Generator[T]
}

func (s *sliceGen[T]) Enumerate(depth int) geniterable.Iterable[[]T] {
	return geniterable.NonExhaustive(EnumerateSlices(depth, depth, s.elemGen))
}

func EnumerateSlices[T any](length, depth int, elemGen Generator[T]) geniterable.Iterable[[]T] {
	if length <= 0 {
		return geniterable.Singleton([]T{})
	}
	smallerSlices := EnumerateSlices(length-1, depth, elemGen)
	return geniterable.Concat(
		smallerSlices,
		geniterable.FlatMap(
			smallerSlices,
			func(tail []T) geniterable.Iterable[[]T] {
				// append only to the longest lists
				if len(tail) < length-1 {
					return geniterable.Empty[[]T]()
				}
				return geniterable.Map(
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

// SliceDistinct generates slices with distinct elements.
func SliceDistinct[T any](elemGen Generator[T], eq equality.Equality[T]) Generator[[]T] {
	return &sliceDistinctGen[T]{
		elemGen: elemGen,
		eq:      eq,
	}
}

type sliceDistinctGen[T any] struct {
	elemGen Generator[T]
	eq      equality.Equality[T]
}

func (s *sliceDistinctGen[T]) Enumerate(depth int) geniterable.Iterable[[]T] {
	// non-exhaustive, because there is no limit on length
	return geniterable.NonExhaustive(EnumerateSlicesDistinct(depth, depth, s.elemGen, s.eq))
}

func EnumerateSlicesDistinct[T any](length, depth int, elemGen Generator[T], eq equality.Equality[T]) geniterable.Iterable[[]T] {
	if length <= 0 {
		return geniterable.Singleton([]T{})
	}
	smallerSlices := EnumerateSlicesDistinct(length-1, depth, elemGen, eq)
	return geniterable.Concat(
		smallerSlices,
		geniterable.FlatMap(
			smallerSlices,
			func(tail []T) geniterable.Iterable[[]T] {
				// append only to the longest lists
				if len(tail) < length-1 {
					return geniterable.Empty[[]T]()
				}
				return geniterable.FlatMap(
					elemGen.Enumerate(depth),
					func(head T) geniterable.Iterable[[]T] {
						if slice.ContainsEq(tail, head, eq) {
							// skip if already in the list
							// (we could inverse the logic here to make this more efficient)
							return geniterable.Empty[[]T]()
						}
						return geniterable.Singleton(append([]T{head}, tail...))
					})
			},
		),
	)
}

func (s *sliceDistinctGen[T]) Name() string {
	return fmt.Sprintf("SliceDistinct(%s)", s.elemGen.Name())
}

func (s *sliceDistinctGen[T]) RValue(elem RandomValue[[]T]) ([]T, bool) {
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

func (s *sliceDistinctGen[T]) Random(rnd Rand, size int) RandomValue[[]T] {
	if size <= 0 {
		return RandomValue[[]T]{Value: []RandomValue[T]{}}
	}
	l := rnd.R().Intn(size)
	res := make([]RandomValue[T], 0, l)
	resValues := make([]T, 0, l)
	for i := 0; i < l; i++ {
		vr := s.elemGen.Random(rnd, size-1)
		v, ok := s.elemGen.RValue(vr)
		if ok && !slice.ContainsEq(resValues, v, s.eq) {
			res = append(res, vr)
			resValues = append(resValues, v)
		}
	}
	return RandomValue[[]T]{Value: res}
}

func (s *sliceDistinctGen[T]) Shrink(elem RandomValue[[]T]) iterable.Iterable[RandomValue[[]T]] {
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

func (s *sliceDistinctGen[T]) Size(t RandomValue[[]T]) *big.Int {
	var size big.Int
	rvs := t.Value.([]RandomValue[T])
	for _, rv := range rvs {
		size.Add(&size, s.elemGen.Size(rv))
	}
	return &size
}

func SliceFixedLength[T any](elemGen Generator[T], length int) Generator[[]T] {
	if length == 0 {
		return Constant([]T{})
	}
	if length == 1 {
		return Map(elemGen, func(e T) []T {
			return []T{e}
		})
	}
	return Zip(elemGen, SliceFixedLength(elemGen, length-1), func(e T, es []T) []T {
		return append([]T{e}, es...)
	})
}
