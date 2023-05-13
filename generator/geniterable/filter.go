package geniterable

func Filter[A any](base Iterable[A], cond func(A) bool) Iterable[A] {
	return &wheregeniterable[A]{base, cond}
}

type wheregeniterable[A any] struct {
	base Iterable[A]
	cond func(A) bool
}

type whereIterator[A any] struct {
	base Iterator[A]
	cond func(A) bool
}

func (i *wheregeniterable[A]) Iterator() Iterator[A] {
	return &whereIterator[A]{i.base.Iterator(), i.cond}
}

func (i *whereIterator[A]) Next() NextResult[A] {
	for {
		var r NextResult[A]
		if r = i.base.Next(); r.Present() {
			if i.cond(r.Value()) {
				return r
			}
			continue
		}
		return r
	}
}
