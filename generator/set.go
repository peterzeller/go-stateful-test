package generator

import (
	"fmt"
	"github.com/peterzeller/go-fun/iterable"
	"github.com/peterzeller/go-stateful-test/generator/geniterable"
	"github.com/peterzeller/go-stateful-test/quickcheck/randomsource"
	"math/big"

	"github.com/peterzeller/go-fun/hash"
	"github.com/peterzeller/go-fun/list/linked"
	"github.com/peterzeller/go-fun/set/hashset"
	"github.com/peterzeller/go-stateful-test/generator/shrink"
)

// Set is a generator for immutable hashsets.
func Set[T, RT any](gen Generator[T, RT], h hash.EqHash[T]) Generator[hashset.Set[T], hashset.Set[RT]] {
	return &setGenerator[T, RT]{
		gen: gen,
		h:   h,
	}
}

type setGenerator[T, RT any] struct {
	gen Generator[T, RT]
	h   hash.EqHash[T]
}

// Enumerate implements Generator
func (s *setGenerator[T, RT]) Enumerate(depth int) geniterable.Iterable[hashset.Set[RT]] {
	it := s.gen.Enumerate(depth).Iterator()
	elems := linked.New[RT]()
	exhaustive := false
	for {
		r := it.Next()
		if !r.Present() {
			exhaustive = r.Exhaustive()
			break
		}
		elems = linked.Cons(r.Value(), elems)
	}
	sets := enumerateSets(elems.Reversed(), s.rvHash())
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
func (s *setGenerator[T, RT]) Random(rnd Rand, size int) hashset.Set[RT] {
	n := randomsource.IntN(rnd.R(), size)
	set := hashset.New(s.rvHash())
	for i := 0; i < n; i++ {
		set = set.Add(s.gen.Random(rnd, size))
	}
	return set
}

func (s *setGenerator[T, RT]) rvHash() hash.EqHash[RT] {
	return &hash.Fun[RT]{
		H: func(elem RT) int64 {
			v, _ := s.gen.RValue(elem)
			return s.h.Hash(v)
		},
		Eq: func(rv1, rv2 RT) bool {
			v1, _ := s.gen.RValue(rv1)
			v2, _ := s.gen.RValue(rv2)
			return s.h.Equal(v1, v2)
		},
	}
}

// Shrink implements Generator
func (s *setGenerator[T, RT]) Shrink(elem hashset.Set[RT]) iterable.Iterable[hashset.Set[RT]] {
	elemSet := elem
	asList := linked.FromIterable[RT](elemSet)
	return iterable.Map(
		shrink.ShrinkList(asList, s.gen.Shrink),
		func(l *linked.List[RT]) hashset.Set[RT] {
			return hashset.New(s.rvHash(), l.ToSlice()...)
		})
}

// RValue implements Generator
func (s *setGenerator[T, RT]) RValue(elem hashset.Set[RT]) (hashset.Set[T], bool) {
	elemSet := elem
	res := hashset.New(s.h)
	for it := iterable.Start[RT](elemSet); it.HasNext(); it.Next() {
		v, ok := s.gen.RValue(it.Current())
		if ok {
			res = res.Add(v)
		}
	}
	return res, true
}

// Size implements Generator
func (s *setGenerator[T, RT]) Size(t hashset.Set[RT]) *big.Int {
	elemSet := t
	var size big.Int
	for it := iterable.Start[RT](elemSet); it.HasNext(); it.Next() {
		size.Add(&size, s.gen.Size(it.Current()))
	}
	return &size
}

func (s *setGenerator[T, RT]) Name() string {
	return fmt.Sprintf("Set(%s)", s.gen.Name())
}

// SetMut is a generator for mutable sets encoded as a map[T]bool.
func SetMut[T, RT comparable](gen Generator[T, RT], h hash.EqHash[T]) Generator[map[T]bool, hashset.Set[RT]] {
	return Map(Set(gen, h),
		func(a hashset.Set[T]) map[T]bool {
			res := make(map[T]bool)
			for it := iterable.Start[T](a); it.HasNext(); it.Next() {
				res[it.Current()] = true
			}
			return res
		})
}
