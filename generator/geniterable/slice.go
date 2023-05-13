package geniterable

type slicegeniterable[T any] struct {
	slice []T
}

func (s slicegeniterable[T]) Length() int {
	return len(s.slice)
}

type sliceIterator[T any] struct {
	slice []T
}

func FromSlice[T any](slice []T) Iterable[T] {
	return slicegeniterable[T]{slice}
}

func New[T any](slice ...T) Iterable[T] {
	return slicegeniterable[T]{slice}
}

func (s slicegeniterable[T]) Iterator() Iterator[T] {
	return &sliceIterator[T]{s.slice}
}

func (s *sliceIterator[T]) Next() NextResult[T] {
	if len(s.slice) == 0 {
		return ResultNone[T](true)
	}
	next := s.slice[0]
	s.slice = s.slice[1:]
	return ResultSome(next)
}

func ToSlice[T any](i Iterable[T]) []T {
	res := make([]T, 0, Length(i))
	it := i.Iterator()
	for {
		if r := it.Next(); r.Present() {
			res = append(res, r.Value())
			continue
		}
		return res
	}
}
