package generator

import (
	"github.com/peterzeller/go-fun/dict/hashdict"
	"github.com/peterzeller/go-fun/equality"
	"github.com/peterzeller/go-fun/hash"
)

// Dict is a generator for immutable dictionaries.
func Dict[K, RK, V, RV any](keyGen Generator[K, RK], valueGen Generator[V, RV], h hash.EqHash[K]) Generator[hashdict.Dict[K, V], flatmapRv[[]RK, interface{}]] {
	keys := SliceDistinct[K](keyGen, h)
	return FlatMap(keys, func(keys []K) Generator[hashdict.Dict[K, V], interface{}] {
		values := SliceFixedLength(valueGen, len(keys))
		return Map(values, func(values []V) hashdict.Dict[K, V] {
			m := hashdict.New[K, V](h)
			for i, key := range keys {
				m = m.Set(key, values[i])
			}
			return m
		})
	})
}

// DictMut is a generator for mutable dictionaries (maps).
func DictMut[K comparable, RK, V, RV any](keyGen Generator[K, RK], valueGen Generator[V, RV]) Generator[map[K]V, flatmapRv[[]RK, interface{}]] {
	keys := SliceDistinct(keyGen, equality.Default[K]())
	return FlatMap(keys, func(keys []K) Generator[map[K]V, interface{}] {
		values := SliceFixedLength(valueGen, len(keys))
		return Map(values, func(values []V) map[K]V {
			m := make(map[K]V)
			for i, key := range keys {
				m[key] = values[i]
			}
			return m
		})
	})
}
