package geniterable

import (
	"github.com/peterzeller/go-fun/slice"
)

func Map[A, B any](base Iterable[A], f func(A) B) Iterable[B] {
	return &mapgeniterable[A, B]{base, f}
}

func MapIterator[A, B any](base Iterator[A], f func(A) B) Iterator[B] {
	return &mapIterator[A, B]{base, f}
}

type mapgeniterable[A, B any] struct {
	base Iterable[A]
	f    func(A) B
}

// Length of a mapgeniterable is the same as the length of the base
func (i mapgeniterable[A, B]) Length() int {
	return Length(i.base)
}

type mapIterator[A, B any] struct {
	base Iterator[A]
	f    func(A) B
}

func (i mapgeniterable[A, B]) Iterator() Iterator[B] {
	return &mapIterator[A, B]{i.base.Iterator(), i.f}
}

func (i *mapIterator[A, B]) Next() NextResult[B] {
	var r NextResult[A]
	if r = i.base.Next(); r.Present() {
		return ResultSome(i.f(r.Value()))
	}
	return ResultNone[B](r.Exhaustive())
}

func FlatMap[A, B any](base Iterable[A], f func(A) Iterable[B]) Iterable[B] {
	return IterableFun[B](func() Iterator[B] {
		it := base.Iterator()
		var current Iterator[B]
		exhaustive := true
		return Fun[B](func() NextResult[B] {
			for {
				if current == nil {
					r := it.Next()
					if !r.Present() {
						return ResultNone[B](exhaustive && r.Exhaustive())
					}
					current = f(r.Value()).Iterator()
				}
				r := current.Next()
				if r.Present() {
					return r
				}
				exhaustive = exhaustive && r.Exhaustive()
				current = nil
			}
		})
	})
}

func FlatMapBreadthFirst[A, B any](base Iterable[A], f func(A) Iterable[B]) Iterable[B] {
	return IterableFun[B](func() Iterator[B] {
		it := base.Iterator()
		firstPass := true
		var iterators []Iterator[B]
		pos := 0
		exhaustive := true
		return Fun[B](func() NextResult[B] {
			for {
				if !firstPass && len(iterators) == 0 {
					return ResultNone[B](exhaustive)
				}
				if pos >= len(iterators) {
					if firstPass {
						// get next element from base iterator
						r := it.Next()
						if r.Present() {
							iterators = append(iterators, f(r.Value()).Iterator())
						} else {
							// no more element in base iterator
							firstPass = false
							pos = 0
							exhaustive = exhaustive && r.Exhaustive()
							continue
						}
					} else {
						pos = 0
						continue
					}
				}
				r := iterators[pos].Next()
				if r.Present() {
					pos++
					return r
				}
				exhaustive = exhaustive && r.Exhaustive()
				// remove iterator from iterators list and try with next position
				iterators = slice.RemoveAt(iterators, pos)
			}
		})
	})
}
