package geniterable

func Generate[T any](start T, next func(prev T) T) Iterable[T] {
	return IterableFun[T](func() Iterator[T] {
		current := start
		first := true
		return Fun[T](func() NextResult[T] {
			if first {
				first = false
				return ResultSome(current)
			}
			current = next(current)
			return ResultSome(current)
		})
	})
}

func GenerateState[S, T any](initialState S, next func(state S) (S, NextResult[T])) Iterable[T] {
	return IterableFun[T](func() Iterator[T] {
		state := initialState
		return Fun[T](func() NextResult[T] {
			var r NextResult[T]
			state, r = next(state)
			return r
		})
	})
}
