package geniterable

type hasLength interface {
	Length() int
}

type hasSize interface {
	Size() int
}

// Length calculates the number of elements in an Iterable.
// This operation takes linear time, unless the Iterable implements a Length or Size method
func Length[T any](i Iterable[T]) (size int) {
	if h, ok := i.(hasLength); ok {
		return h.Length()
	}
	if h, ok := i.(hasSize); ok {
		return h.Size()
	}
	it := i.Iterator()
	for {
		n := it.Next()
		if !n.Present() {
			return
		}
		size++
	}
}
