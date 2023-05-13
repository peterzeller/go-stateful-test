package geniterable

// LoopIterator is a different form of iterator better suited for use in for-loops:
//
//	for it := iterable.Start(x); it.HasNext(); it.Next() {
//	   doSomethingWith(it.Current())
//	}
type LoopIterator[T any] struct {
	it      Iterator[T]
	current NextResult[T]
}

func Start[T any](i Iterable[T]) LoopIterator[T] {
	it := i.Iterator()
	r := it.Next()
	return LoopIterator[T]{
		it:      it,
		current: r,
	}
}

func (l *LoopIterator[T]) HasNext() bool {
	return l.current.Present()
}

func (l *LoopIterator[T]) Next() {
	l.current = l.it.Next()
}

func (l *LoopIterator[T]) Current() T {
	return l.current.Value()
}

// Foreach runs the function f on each element of the iterable.
func Foreach[T any](i Iterable[T], f func(elem T)) {
	it := i.Iterator()
	for {
		r := it.Next()
		if !r.Present() {
			return
		}
		f(r.Value())
	}
}
