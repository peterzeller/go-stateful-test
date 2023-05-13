package geniterable

import "github.com/peterzeller/go-fun/zero"

// Find an element in an Iterable
func Find[T any](i Iterable[T], cond func(T) bool) (T, bool) {
	it := i.Iterator()
	for {
		r := it.Next()
		if !r.Present() {
			return zero.Value[T](), false
		}
		if cond(r.value) {
			return r.value, true
		}
	}
}
