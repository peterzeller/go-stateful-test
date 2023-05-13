package geniterable

func Concat[T any](geniterables ...Iterable[T]) Iterable[T] {
	return IterableFun[T](func() Iterator[T] {
		pos := 0
		var current Iterator[T]
		exhaustive := true
		return Fun[T](func() NextResult[T] {
			for {
				if current == nil {
					if pos >= len(geniterables) {
						return ResultNone[T](exhaustive)
					}
					current = geniterables[pos].Iterator()
					pos++
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

func ConcatIterators[T any](iterators ...Iterator[T]) Iterator[T] {
	pos := 0
	exhaustive := true
	return Fun[T](func() NextResult[T] {
		for pos < len(iterators) {
			r := iterators[pos].Next()
			if r.Present() {
				return r
			}
			exhaustive = exhaustive && r.Exhaustive()
			pos++
		}
		return ResultNone[T](exhaustive)
	})
}
