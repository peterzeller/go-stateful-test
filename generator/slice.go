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
func Slice[T, TR any](elemGen Generator[T, TR]) Generator[[]T, []TR] {
	return &sliceGen[T, TR]{
		elemGen: elemGen,
	}
}

type sliceGen[T, TR any] struct {
	elemGen Generator[T, TR]
}

func (s *sliceGen[T, TR]) Enumerate(depth int) geniterable.Iterable[[]TR] {
	return geniterable.NonExhaustive(EnumerateSlices(depth, depth, s.elemGen))
}

func EnumerateSlices[T, TR any](length, depth int, elemGen Generator[T, TR]) geniterable.Iterable[[]TR] {
	if length <= 0 {
		return geniterable.Singleton([]TR{})
	}
	smallerSlices := EnumerateSlices(length-1, depth, elemGen)
	return geniterable.Concat(
		smallerSlices,
		geniterable.FlatMap(
			smallerSlices,
			func(tail []TR) geniterable.Iterable[[]TR] {
				// append only to the longest lists
				if len(tail) < length-1 {
					return geniterable.Empty[[]TR]()
				}
				return geniterable.Map(
					elemGen.Enumerate(depth),
					func(head TR) []TR {
						return append([]TR{head}, tail...)
					})
			},
		),
	)
}

func (s *sliceGen[T, TR]) Name() string {
	return fmt.Sprintf("SliceGen[%s]", s.elemGen.Name())
}

func (s *sliceGen[T, TR]) RValue(elem []TR) ([]T, bool) {
	res := make([]T, len(elem))
	for i, rv := range elem {
		var ok bool
		res[i], ok = s.elemGen.RValue(rv)
		if !ok {
			return nil, false
		}
	}
	return res, true
}

func (s *sliceGen[T, TR]) Random(rnd Rand, size int) []TR {
	if size <= 0 {
		return []TR{}
	}
	l := rnd.R().Intn(size)
	res := make([]TR, l)
	for i := range res {
		res[i] = s.elemGen.Random(rnd, size-1)
	}
	return res
}

func (s *sliceGen[T, TR]) Shrink(elem []TR) iterable.Iterable[[]TR] {
	rvs := elem
	return iterable.Map(
		shrink.ShrinkList(
			linked.New(rvs...),
			func(rv TR) iterable.Iterable[TR] {
				return s.elemGen.Shrink(rv)
			}),
		func(l *linked.List[TR]) []TR {
			return l.ToSlice()
		})
}

func (s *sliceGen[T, TR]) Size(t []TR) *big.Int {
	var size big.Int
	rvs := t
	for _, rv := range rvs {
		size.Add(&size, s.elemGen.Size(rv))
	}
	return &size
}

// SliceDistinct generates slices with distinct elements.
func SliceDistinct[T, TR any](elemGen Generator[T, TR], eq equality.Equality[T]) Generator[[]T, []TR] {
	return &sliceDistinctGen[T, TR]{
		elemGen: elemGen,
		eq:      eq,
	}
}

type sliceDistinctGen[T, TR any] struct {
	elemGen Generator[T, TR]
	eq      equality.Equality[T]
}

func (s *sliceDistinctGen[T, TR]) Enumerate(depth int) geniterable.Iterable[[]TR] {
	// non-exhaustive, because there is no limit on length
	return geniterable.NonExhaustive(EnumerateSlicesDistinct(depth, depth, s.elemGen, eqRandomValue(s.elemGen.RValue, s.eq)))
}

func EnumerateSlicesDistinct[T, TR any](length, depth int, elemGen Generator[T, TR], eq equality.Equality[TR]) geniterable.Iterable[[]TR] {
	if length <= 0 {
		return geniterable.Singleton([]TR{})
	}
	smallerSlices := EnumerateSlicesDistinct(length-1, depth, elemGen, eq)
	return geniterable.Concat(
		smallerSlices,
		geniterable.FlatMap(
			smallerSlices,
			func(tail []TR) geniterable.Iterable[[]TR] {
				// append only to the longest lists
				if len(tail) < length-1 {
					return geniterable.Empty[[]TR]()
				}
				return geniterable.FlatMap(
					elemGen.Enumerate(depth),
					func(head TR) geniterable.Iterable[[]TR] {
						if slice.ContainsEq(tail, head, eq) {
							// skip if already in the list
							// (we could inverse the logic here to make this more efficient)
							return geniterable.Empty[[]TR]()
						}
						return geniterable.Singleton(append([]TR{head}, tail...))
					})
			},
		),
	)
}

func (s *sliceDistinctGen[T, TR]) Name() string {
	return fmt.Sprintf("SliceDistinct(%s)", s.elemGen.Name())
}

func (s *sliceDistinctGen[T, TR]) RValue(elem []TR) ([]T, bool) {
	res := make([]T, len(elem))
	for i, rv := range elem {
		var ok bool
		res[i], ok = s.elemGen.RValue(rv)
		if !ok {
			return nil, false
		}
	}
	return res, true
}

func (s *sliceDistinctGen[T, TR]) Random(rnd Rand, size int) []TR {
	if size <= 0 {
		return []TR{}
	}
	l := rnd.R().Intn(size)
	res := make([]TR, 0, l)
	resValues := make([]T, 0, l)
	for i := 0; i < l; i++ {
		vr := s.elemGen.Random(rnd, size-1)
		v, ok := s.elemGen.RValue(vr)
		if ok && !slice.ContainsEq(resValues, v, s.eq) {
			res = append(res, vr)
			resValues = append(resValues, v)
		}
	}
	return res
}

func (s *sliceDistinctGen[T, TR]) Shrink(elem []TR) iterable.Iterable[[]TR] {
	rvs := elem
	return iterable.Map(
		shrink.ShrinkList(
			linked.New(rvs...),
			func(rv TR) iterable.Iterable[TR] {
				return s.elemGen.Shrink(rv)
			}),
		func(l *linked.List[TR]) []TR {
			return l.ToSlice()
		})
}

func (s *sliceDistinctGen[T, TR]) Size(t []TR) *big.Int {
	var size big.Int
	rvs := t
	for _, rv := range rvs {
		size.Add(&size, s.elemGen.Size(rv))
	}
	return &size
}

func eqRandomValue[T, TR any](rValue func(r TR) (T, bool), orig equality.Equality[T]) equality.Equality[TR] {
	return equality.Fun[TR](func(ra, rb TR) bool {
		a, okA := rValue(ra)
		b, okB := rValue(rb)
		return okA && okB && orig.Equal(a, b)
	})
}

func SliceFixedLength[T, TR any](elemGen Generator[T, TR], length int) Generator[[]T, interface{}] {
	if length == 0 {
		return UntypedR(Constant([]T{}))
	}
	if length == 1 {
		return UntypedR(Map(elemGen, func(e T) []T {
			return []T{e}
		}))
	}
	return UntypedR(Zip(elemGen, SliceFixedLength(elemGen, length-1), func(e T, es []T) []T {
		return append([]T{e}, es...)
	}))
}
