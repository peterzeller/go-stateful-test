package generator

import (
	"fmt"
	"math/big"
	"math/rand"

	"github.com/peterzeller/go-fun/iterable"
)

type Generator[T any] interface {
	// Name of the generator (used for shrinking)
	Name() string
	// Random element with a given maximum size
	Random(rnd Rand, size int) T
	// Enumerate all elements of this type up to the given size
	Enumerate(depth int) iterable.Iterable[T]
	// Shrink an element for reducing test cases - produces zero or more elements that are smaller than the original
	Shrink(elem T) iterable.Iterable[T]
	// Size gives a size estimate for the given value.
	// Elements returned by shrink must have smaller values
	Size(t T) *big.Int
}

// UntypedGenerator is a workaround for Go not having existential types.
// It wraps a typed generator and removes the type parameter, so that we can use it in heterogeneous contexts.
// The type is only used internally and not exposed in the API.
type UntypedGenerator interface {
	// Name of the generator (used for shrinking)
	Name() string
	// Random element with a given maximum size
	Random(rnd Rand, size int) interface{}
	// Enumerate all elements of this type up to the given size
	Enumerate(depth int) iterable.Iterable[interface{}]
	// Shrink an element for reducing test cases - produces zero or more elements that are smaller than the original
	Shrink(elem interface{}) iterable.Iterable[interface{}]
	// Size gives a size estimate for the given value.
	// Elements returned by shrink must have smaller values
	Size(value interface{}) *big.Int
}

// untypedGen is the canonical implementation for UntypedGenerator.
type untypedGen struct {
	name      func() string
	random    func(rnd Rand, size int) interface{}
	enumerate func(depth int) iterable.Iterable[interface{}]
	shrink    func(elem interface{}) iterable.Iterable[interface{}]
	size      func(elem interface{}) *big.Int
}

func (u untypedGen) Size(value interface{}) *big.Int {
	return u.size(value)
}

func (u untypedGen) Name() string {
	return u.name()
}

func (u untypedGen) Random(rnd Rand, size int) interface{} {
	return u.random(rnd, size)
}

func (u untypedGen) Enumerate(depth int) iterable.Iterable[interface{}] {
	return u.enumerate(depth)
}

func (u untypedGen) Shrink(elem interface{}) iterable.Iterable[interface{}] {
	return u.shrink(elem)
}

func ToUntyped[T any](gen Generator[T]) UntypedGenerator {
	return untypedGen{
		name: gen.Name,
		random: func(rnd Rand, size int) interface{} {
			return gen.Random(rnd, size)
		},
		enumerate: func(depth int) iterable.Iterable[interface{}] {
			return iterable.Map(gen.Enumerate(depth),
				func(e T) interface{} {
					return e
				})
		},
		shrink: func(elem interface{}) iterable.Iterable[interface{}] {
			elemT, ok := elem.(T)
			if !ok {
				panic(fmt.Errorf("could not convert element %#v", elemT))
			}
			return iterable.Map(
				gen.Shrink(elemT),
				func(e T) interface{} {
					return e
				})
		},
		size: func(elem interface{}) *big.Int {
			elemT, ok := elem.(T)
			if !ok {
				panic(fmt.Errorf("could not convert element %#v", elemT))
			}
			return gen.Size(elemT)
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
