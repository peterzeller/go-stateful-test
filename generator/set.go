package generator

import (
	"fmt"
	"github.com/peterzeller/go-fun/iterable"
	"github.com/peterzeller/go-stateful-test/generator/geniterable"
	"math/big"

	"github.com/peterzeller/go-fun/hash"
	"github.com/peterzeller/go-fun/list/linked"
	"github.com/peterzeller/go-fun/set/hashset"
	"github.com/peterzeller/go-stateful-test/generator/shrink"
)

// Set is a generator for immutable hashsets.
func Set[T any](gen Generator[T], h hash.EqHash[T]) Generator[hashset.Set[T]] {
	return &setGenerator[T]{
		gen: gen,
		h:   h,
	}
}

type setGenerator[T any] struct {
	gen Generator[T]
	h   hash.EqHash[T]
}

// Enumerate implements Generator
func (s *setGenerator[T]) Enumerate(depth int) geniterable.Iterable[hashset.Set[T]] {
	it := s.gen.Enumerate(depth).Iterator()
	elems := linked.New[T]()
	exhaustive := false
	for {
		r := it.Next()
		if !r.Present() {
			exhaustive = r.Exhaustive()
			break
		}
		elems = linked.Cons(r.Value(), elems)
	}
	sets := enumerateSets(elems.Reversed(), s.h)
	if !exhaustive {
		sets = geniterable.NonExhaustive(sets)
	}
	return sets
}

func enumerateSets[T any](elems *linked.List[T], h hash.EqHash[T]) geniterable.Iterable[hashset.Set[T]] {
	if elems == nil {
		// return empty set
		return geniterable.Singleton(hashset.New(h))
	}
	tailSets := enumerateSets(elems.Tail(), h)
	return geniterable.Concat(
		tailSets,
		geniterable.Map(
			tailSets,
			func(tail hashset.Set[T]) hashset.Set[T] {
				return tail.Add(elems.Head())
			},
		))
}

// Random implements Generator
func (s *setGenerator[T]) Random(rnd Rand, size int) RandomValue[hashset.Set[T]] {
	n := rnd.R().Intn(size)
	set := hashset.New(s.rvHash())
	for i := 0; i < n; i++ {
		set = set.Add(s.gen.Random(rnd, size))
	}
	return RandomValue[hashset.Set[T]]{
		Value: set,
	}
}

func (s *setGenerator[T]) rvHash() hash.EqHash[RandomValue[T]] {
	return &hash.Fun[RandomValue[T]]{
		H: func(elem RandomValue[T]) int64 {
			v, _ := s.gen.RValue(elem)
			return s.h.Hash(v)
		},
		Eq: func(rv1, rv2 RandomValue[T]) bool {
			v1, _ := s.gen.RValue(rv1)
			v2, _ := s.gen.RValue(rv2)
			return s.h.Equal(v1, v2)
		},
	}
}

// Shrink implements Generator
func (s *setGenerator[T]) Shrink(elem RandomValue[hashset.Set[T]]) iterable.Iterable[RandomValue[hashset.Set[T]]] {
	elemSet := elem.Value.(hashset.Set[RandomValue[T]])
	asList := linked.FromIterable[RandomValue[T]](elemSet)
	return iterable.Map(
		shrink.ShrinkList(asList, s.gen.Shrink),
		func(l *linked.List[RandomValue[T]]) RandomValue[hashset.Set[T]] {
			return RandomValue[hashset.Set[T]]{
				Value: hashset.New(s.rvHash(), l.ToSlice()...),
			}
		})
}

// RValue implements Generator
func (s *setGenerator[T]) RValue(elem RandomValue[hashset.Set[T]]) (hashset.Set[T], bool) {
	elemSet := elem.Value.(hashset.Set[RandomValue[T]])
	res := hashset.New(s.h)
	for it := iterable.Start[RandomValue[T]](elemSet); it.HasNext(); it.Next() {
		v, ok := s.gen.RValue(it.Current())
		if ok {
			res = res.Add(v)
		}
	}
	return res, true
}

// Size implements Generator
func (s *setGenerator[T]) Size(t RandomValue[hashset.Set[T]]) *big.Int {
	elemSet := t.Value.(hashset.Set[RandomValue[T]])
	var size big.Int
	for it := iterable.Start[RandomValue[T]](elemSet); it.HasNext(); it.Next() {
		size.Add(&size, s.gen.Size(it.Current()))
	}
	return &size
}

func (s *setGenerator[T]) Name() string {
	return fmt.Sprintf("Set(%s)", s.gen.Name())
}

// SetMut is a generator for mutable sets encoded as a map[T]bool.
func SetMut[T comparable](gen Generator[T], h hash.EqHash[T]) Generator[map[T]bool] {
	return Map(Set(gen, h),
		func(a hashset.Set[T]) map[T]bool {
			res := make(map[T]bool)
			for it := iterable.Start[T](a); it.HasNext(); it.Next() {
				res[it.Current()] = true
			}
			return res
		})
}
