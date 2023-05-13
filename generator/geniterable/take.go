package geniterable

// Take the first n elements from the Iterable
func Take[T any](n int, i Iterable[T]) Iterable[T] {
	return IterableFun[T](func() Iterator[T] {
		count := 0
		it := i.Iterator()
		return Fun[T](func() NextResult[T] {
			if count >= n {
				return ResultNone[T](true)
			}
			count++
			return it.Next()
		})
	})
}

// TakeExhaustive takes the first n elements from the Iterable.
// The exhaustive flag is set to true if the original is exhaustive and all elements are taken.
func TakeExhaustive[T any](n int, i Iterable[T]) Iterable[T] {
	return IterableFun[T](func() Iterator[T] {
		count := 0
		exhaustive := false
		it := i.Iterator()
		return Fun[T](func() NextResult[T] {
			if count > n {
				return ResultNone[T](exhaustive)
			}
			count++
			r := it.Next()
			if !r.Present() {
				exhaustive = r.Exhaustive()
				count = n + 1
			}
			if count > n {
				return ResultNone[T](exhaustive)
			}
			return r
		})
	})
}

// TakeWhile takes elements from the Iterable, while the elements match the condition
func TakeWhile[T any](cond func(T) bool, i Iterable[T]) Iterable[T] {
	return IterableFun[T](func() Iterator[T] {
		it := i.Iterator()
		active := true
		return Fun[T](func() NextResult[T] {
			if !active {
				return ResultNone[T](true)
			}
			r := it.Next()
			if !r.Present() {
				active = false
				return ResultNone[T](r.Exhaustive())
			}
			if !cond(r.value) {
				active = false
				return ResultNone[T](true)
			}
			return r
		})
	})
}
