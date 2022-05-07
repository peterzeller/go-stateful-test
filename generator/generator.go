// The generator package defines the interface for generators, as well as generators for common Go types and
// combinators that simplify writing new generators.
//
// Defining a generator
//
// To define a new generator, you need to implement the Generator interface.
// This can be done by implementing a new type and with the required methods or
// directly inline by creating a AnonGenerator instance.
//
// It is also possible to create
package generator

import (
	"fmt"
	"math/big"
	"math/rand"
	"reflect"

	"github.com/peterzeller/go-fun/iterable"
	"github.com/peterzeller/go-fun/zero"
)

type Generator[T any] interface {
	// Name of the generator (used for shrinking)
	Name() string
	// Random element with a given maximum size
	Random(rnd Rand, size int) RandomValue[T]
	// Shrink an element for reducing test cases - produces zero or more elements that are smaller than the original
	Shrink(elem RandomValue[T]) iterable.Iterable[RandomValue[T]]
	// RValue takes a generated random element and tries to repair it it to satisfy the bounds of the generator.
	RValue(elem RandomValue[T]) (T, bool)
	// Size gives a size estimate for the given value.
	// Elements returned by shrink must have smaller values
	Size(t RandomValue[T]) *big.Int
	// Enumerate all elements of this type up to the given size
	Enumerate(depth int) iterable.Iterable[T]
}

// RandomValue generated by a generator.
// It may contain additional metadata that is used for shrinking.
type RandomValue[T any] struct {
	// Value is the data stored in the RandomValue.
	// It will often be of type T, but some generators will want to store additional metadata in this field.
	// Therefore, the Generator.RValue method should be used to extract the value.
	Value interface{}
}

// R is a shorthand for creating a RandomValue from a value.
func R[T any](elem T) RandomValue[T] {
	return RandomValue[T]{Value: elem}
}

// Get the value stored in the RandomValue if it is of type T.
// This function should be used only in specific generators.
// For the general case, use Generator.RValue.
func (r RandomValue[T]) Get() T {
	v, ok := r.Value.(T)
	if !ok {
		panic(fmt.Errorf("RandomValue.Get: type mismatch, was %v, expected %v", reflect.TypeOf(r.Value), reflect.TypeOf(zero.Value[T]())))
	}
	return v
}

func (r RandomValue[T]) Untyped() RandomValue[interface{}] {
	return RandomValue[interface{}]{Value: r.Value}
}

func RTyped[T any](elem RandomValue[interface{}]) RandomValue[T] {
	return RandomValue[T]{Value: elem.Value}
}

// UntypedGenerator is a workaround for Go not having existential types.
// It wraps a typed generator and removes the type parameter, so that we can use it in heterogeneous contexts.
// The type is only used internally and not exposed in the API.
type UntypedGenerator interface {
	// Name of the generator (used for shrinking)
	Name() string
	// Random element with a given maximum size
	Random(rnd Rand, size int) RandomValue[interface{}]
	// Shrink an element for reducing test cases - produces zero or more elements that are smaller than the original
	Shrink(elem RandomValue[interface{}]) iterable.Iterable[RandomValue[interface{}]]
	// Elements returned by shrink must have smaller values
	// RValue takes a generated random element and tries to repair it it to satisfy the bounds of the generator.
	RValue(elem RandomValue[interface{}]) (interface{}, bool)
	// Size gives a size estimate for the given value.
	Size(value RandomValue[interface{}]) *big.Int
	// Enumerate all elements of this type up to the given size
	Enumerate(depth int) iterable.Iterable[interface{}]
}

// untypedGen is the canonical implementation for UntypedGenerator.
type untypedGen struct {
	name      func() string
	random    func(rnd Rand, size int) RandomValue[interface{}]
	shrink    func(elem RandomValue[interface{}]) iterable.Iterable[RandomValue[interface{}]]
	rvalue    func(elem RandomValue[interface{}]) (interface{}, bool)
	size      func(elem RandomValue[interface{}]) *big.Int
	enumerate func(depth int) iterable.Iterable[interface{}]
}

func (u untypedGen) Size(value RandomValue[interface{}]) *big.Int {
	return u.size(value)
}

func (u untypedGen) Name() string {
	return u.name()
}

func (u untypedGen) Random(rnd Rand, size int) RandomValue[interface{}] {
	return u.random(rnd, size).Untyped()
}

func (u untypedGen) Enumerate(depth int) iterable.Iterable[interface{}] {
	return u.enumerate(depth)
}

func (u untypedGen) RValue(elem RandomValue[interface{}]) (interface{}, bool) {
	return u.rvalue(elem)
}

func (u untypedGen) Shrink(elem RandomValue[interface{}]) iterable.Iterable[RandomValue[interface{}]] {
	return u.shrink(elem)
}

func ToUntyped[T any](gen Generator[T]) UntypedGenerator {
	return untypedGen{
		name: gen.Name,
		random: func(rnd Rand, size int) RandomValue[interface{}] {
			return gen.Random(rnd, size).Untyped()
		},
		enumerate: func(depth int) iterable.Iterable[interface{}] {
			return iterable.Map(gen.Enumerate(depth),
				func(e T) interface{} {
					return e
				})
		},
		shrink: func(elem RandomValue[interface{}]) iterable.Iterable[RandomValue[interface{}]] {
			return iterable.Map(
				gen.Shrink(RTyped[T](elem)),
				func(e RandomValue[T]) RandomValue[interface{}] {
					return e.Untyped()
				})
		},
		rvalue: func(elem RandomValue[interface{}]) (interface{}, bool) {
			return gen.RValue(RTyped[T](elem))
		},
		size: func(elem RandomValue[interface{}]) *big.Int {
			return gen.Size(RTyped[T](elem))
		},
	}
}

type Rand interface {
	// Fork this random number generator for controlling a subgroup of the test
	Fork(name string) Rand
	// HasMore to generate sequences of elements and if there are more elements
	HasMore() bool
	// R is the underlying random number generator
	R() *rand.Rand
}

func ShrinkValues[T any](gen Generator[T], v T) iterable.Iterable[T] {
	return iterable.Map(
		gen.Shrink(R(v)),
		func(rv RandomValue[T]) T {
			return rv.Get()
		},
	)
}

// ToTypedGenerator converts an UntypedGenerator to a typed Generator.
func ToTypedGenerator[T any](g UntypedGenerator) Generator[T] {
	return &AnonGenerator[T]{
		GenName: g.Name(),
		GenRandom: func(rnd Rand, size int) RandomValue[T] {
			return RandomValue[T]{Value: g.Random(rnd, size).Value}
		},
		GenShrink: func(elem RandomValue[T]) iterable.Iterable[RandomValue[T]] {
			return iterable.Map(
				g.Shrink(elem.Untyped()),
				func(rv RandomValue[interface{}]) RandomValue[T] {
					return RandomValue[T]{Value: rv.Value}
				},
			)
		},
		GenSize: func(t RandomValue[T]) *big.Int {
			return g.Size(t.Untyped())
		},
		GenRValue: func(r RandomValue[T]) (T, bool) {
			v, ok := g.RValue(r.Untyped())
			if !ok {
				return zero.Value[T](), false
			}
			return v.(T), ok
		},
		GenEnumerate: func(depth int) iterable.Iterable[T] {
			return iterable.Map(
				g.Enumerate(depth),
				func(e interface{}) T {
					return e.(T)
				},
			)
		},
	}
}
