// Package generator defines the interface for generators, as well as generators for common Go types and
// combinators that simplify writing new generators.
//
// # Defining a generator
//
// To define a new generator, you need to implement the Generator interface.
// This can be done by implementing a new type and with the required methods or
// directly inline by creating a AnonGenerator instance.
//
// It is also possible to create a gnerator from existing generators using combinators such as Map FlatMap or Filter.
package generator

import (
	"github.com/peterzeller/go-fun/iterable"
	"github.com/peterzeller/go-fun/zero"
	"github.com/peterzeller/go-stateful-test/generator/geniterable"
	"github.com/peterzeller/go-stateful-test/quickcheck/randomsource"
	"math/big"
)

// Generator is an interface for generating values of type T with internal value representation R.
// The type R can be equal to type T, but it can also store additional information.
// The method RValue can be used to convert from R to T.
type Generator[T any, R any] interface {
	// Name of the generator (used for shrinking)
	Name() string
	// Random element with a given maximum size
	Random(rnd Rand, size int) R
	// Shrink an element for reducing test cases - produces zero or more elements that are smaller than the original
	Shrink(elem R) iterable.Iterable[R]
	// RValue takes a generated random element and tries to repair it to satisfy the bounds of the generator.
	RValue(elem R) (T, bool)
	// Size gives a size estimate for the given value.
	// Elements returned by shrink must have smaller values
	Size(t R) *big.Int
	// Enumerate all elements of this type up to the given size
	Enumerate(depth int) geniterable.Iterable[R]
}

// UV is the value type for untyped generators.
type UV struct {
	Value interface{}
}

// UR is the representation type for untyped generators.
type UR struct {
	value interface{}
}

// UntypedGenerator is a workaround for Go not having existential types.
// It wraps a typed generator and removes the type parameter, so that we can use it in heterogeneous contexts.
// The type is only used internally and not exposed in the API.
type UntypedGenerator Generator[UV, UR]

// untypedGen is the canonical implementation for UntypedGenerator.
type untypedGen struct {
	name      func() string
	random    func(rnd Rand, size int) UR
	shrink    func(elem UR) iterable.Iterable[UR]
	rvalue    func(elem UR) (UV, bool)
	size      func(elem UR) *big.Int
	enumerate func(depth int) geniterable.Iterable[UR]
}

func (u untypedGen) Size(value UR) *big.Int {
	return u.size(value)
}

func (u untypedGen) Name() string {
	return u.name()
}

func (u untypedGen) Random(rnd Rand, size int) UR {
	return u.random(rnd, size)
}

func (u untypedGen) Enumerate(depth int) geniterable.Iterable[UR] {
	return u.enumerate(depth)
}

func (u untypedGen) RValue(elem UR) (UV, bool) {
	return u.rvalue(elem)
}

func (u untypedGen) Shrink(elem UR) iterable.Iterable[UR] {
	return u.shrink(elem)
}

func ToUntyped[T, R any](gen Generator[T, R]) UntypedGenerator {
	wrapR := func(e R) UR {
		return UR{e}
	}
	unwrapR := func(i UR) R {
		return i.value.(R)
	}
	return untypedGen{
		name: gen.Name,
		random: func(rnd Rand, size int) UR {
			return UR{gen.Random(rnd, size)}
		},
		enumerate: func(depth int) geniterable.Iterable[UR] {
			return geniterable.Map(gen.Enumerate(depth), wrapR)
		},
		shrink: func(elem UR) iterable.Iterable[UR] {
			return iterable.Map(
				gen.Shrink(unwrapR(elem)),
				wrapR)
		},
		rvalue: func(elem UR) (UV, bool) {
			value, ok := gen.RValue(unwrapR(elem))
			return UV{value}, ok
		},
		size: func(elem UR) *big.Int {
			return gen.Size(unwrapR(elem))
		},
	}
}

type Rand interface {
	// Fork this random number generator for controlling a subgroup of the test
	Fork(name string) Rand
	// HasMore to generate sequences of elements and if there are more elements
	HasMore() bool
	// R is the underlying random number generator
	R() randomsource.RandomStream
	// UseHeuristics is true when the search should use heuristics to find elements. For example try small numbers first.
	UseHeuristics() bool
}

func ShrinkValues[T any](gen Generator[T, T], v T) iterable.Iterable[T] {
	return iterable.Map(
		gen.Shrink(v),
		func(rv T) T {
			return rv
		},
	)
}

// ToTypedGenerator converts an UntypedGenerator to a typed Generator.
func ToTypedGenerator[T, R any](g UntypedGenerator) Generator[T, R] {
	toR := func(i UR) R {
		return i.value.(R)
	}
	return &AnonGenerator[T, R]{
		GenName: g.Name(),
		GenRandom: func(rnd Rand, size int) R {
			return toR(g.Random(rnd, size))
		},
		GenShrink: func(elem R) iterable.Iterable[R] {
			return iterable.Map(
				g.Shrink(UR{elem}),
				toR,
			)
		},
		GenSize: func(t R) *big.Int {
			return g.Size(UR{t})
		},
		GenRValue: func(r R) (T, bool) {
			v, ok := g.RValue(UR{r})
			if !ok {
				return zero.Value[T](), false
			}
			return v.Value.(T), ok
		},
		GenEnumerate: func(depth int) geniterable.Iterable[R] {
			return geniterable.Map(
				g.Enumerate(depth),
				toR,
			)
		},
	}
}

// EnumerateValues enumerates the values of a generator up to the given depth.
func EnumerateValues[T, R any](gen Generator[T, R], depth int) geniterable.Iterable[T] {
	return geniterable.FlatMap(gen.Enumerate(depth), func(r R) geniterable.Iterable[T] {
		value, ok := gen.RValue(r)
		if ok {
			return geniterable.Singleton(value)
		}
		return geniterable.Empty[T]()
	})
}
