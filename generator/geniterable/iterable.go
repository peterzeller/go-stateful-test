package geniterable

import (
	"fmt"
	"github.com/peterzeller/go-fun/iterable"
)

type Iterable[T any] interface {
	Iterator() Iterator[T]
}

type Iterator[T any] interface {
	Next() NextResult[T]
}

type NextResult[T any] struct {
	// status:
	// 0 - does not exist and not exhaustive
	// 1 - does not exist and exhaustive
	// 2 - does exist
	status byte
	value  T
}

func ResultSome[T any](v T) NextResult[T] {
	return NextResult[T]{
		status: 2,
		value:  v,
	}
}

func ResultNone[T any](exhaustive bool) NextResult[T] {
	status := byte(0)
	if exhaustive {
		status = 1
	}
	return NextResult[T]{
		status: status,
	}
}

func (n NextResult[T]) String() string {
	if n.status == 2 {
		return fmt.Sprintf("Some(%v)", n.value)
	} else if n.status == 1 {
		return "done (exhaustive)"
	} else {
		return "done (not exhaustive)"
	}
}

func (n NextResult[T]) Value() T {
	return n.value
}

func (n NextResult[T]) Present() bool {
	return n.status == 2
}

func (n NextResult[T]) Exhaustive() bool {
	return n.status == 1
}

type Fun[T any] func() NextResult[T]

func (f Fun[T]) Next() NextResult[T] {
	return f()
}

type IterableFun[T any] func() Iterator[T]

func (f IterableFun[T]) Iterator() Iterator[T] {
	return f()
}

type emptyIterable[T any] struct {
}

type emptyIterator[T any] struct {
}

func (e emptyIterator[T]) Next() NextResult[T] {
	return ResultNone[T](true)
}

func (e emptyIterable[T]) Iterator() Iterator[T] {
	return emptyIterator[T]{}
}

func Empty[T any]() Iterable[T] {
	return emptyIterable[T]{}
}

func Singleton[T any](x T) Iterable[T] {
	return IterableFun[T](func() Iterator[T] {
		first := true
		return Fun[T](func() NextResult[T] {
			if first {
				first = false
				return ResultSome(x)
			}
			return ResultNone[T](true)
		})
	})
}

// NonExhaustive makes the iterable non-exhaustive
func NonExhaustive[T any](orig Iterable[T]) Iterable[T] {
	return IterableFun[T](func() Iterator[T] {
		it := orig.Iterator()
		return Fun[T](func() NextResult[T] {
			next := it.Next()
			if next.Present() {
				return next
			}
			return ResultNone[T](false)
		})
	})
}

// IsExhaustive checks whether the iterable is exhaustive.
// Warning: this enumerates all elements.
func IsExhaustive[T any](i Iterable[T]) bool {
	it := i.Iterator()
	for {
		r := it.Next()
		if !r.Present() {
			return r.Exhaustive()
		}
	}
}

func ToIterable[T any](orig Iterable[T]) iterable.Iterable[T] {
	return iterable.IterableFun[T](func() iterable.Iterator[T] {
		it := orig.Iterator()
		return iterable.Fun[T](func() (T, bool) {
			r := it.Next()
			return r.Value(), r.Present()
		})
	})
}
