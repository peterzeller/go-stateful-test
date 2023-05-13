package geniterable

import (
	"fmt"
	"strings"
)

func String[T any](i Iterable[T]) string {
	var res strings.Builder
	res.WriteString("[")
	first := true
	it := i.Iterator()
	for {
		r := it.Next()
		if !r.Present() {
			if !r.Exhaustive() {
				if !first {
					res.WriteString(", ")
				}
				res.WriteString("...")
			}
			break
		}
		if !first {
			res.WriteString(", ")
		}
		res.WriteString(fmt.Sprintf("%+v", r.Value()))
		first = false
	}
	res.WriteString("]")
	return res.String()
}

func IteratorToSlice[T any](it Iterator[T]) []T {
	var res []T
	for {
		r := it.Next()
		if !r.Present() {
			return res
		}
		res = append(res, r.value)
	}
}
