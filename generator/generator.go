package generator

import (
	"fmt"

	"github.com/peterzeller/go-fun/iterable"
)

type Generator[T any] interface {
	// Name of the generator (used for shrinking)
	Name() string
	// Random element with a given maximum size
	Random(rnd Rand, size int) T
	// Enumerate all elements of this type up to the given size
	Enumerate(depth int) iterable.Iterator[T]
	// Shrink an element for reducing test cases - produces zero or more elements that are smaller than the original
	Shrink(elem T) iterable.Iterator[T]
}

type UntypedGenerator interface {
	// Name of the generator (used for shrinking)
	Name() string
	// Random element with a given maximum size
	Random(rnd Rand, size int) interface{}
	// Enumerate all elements of this type up to the given size
	Enumerate(depth int) iterable.Iterator[interface{}]
	// Shrink an element for reducing test cases - produces zero or more elements that are smaller than the original
	Shrink(elem interface{}) iterable.Iterator[interface{}]
}

type gen[T any] struct {
	name      func() string
	random    func(rnd Rand, size int) T
	enumerate func(depth int) iterable.Iterator[T]
	shrink    func(elem interface{}) iterable.Iterator[T]
}

type untypedGen struct {
	name      func() string
	random    func(rnd Rand, size int) interface{}
	enumerate func(depth int) iterable.Iterator[interface{}]
	shrink    func(elem interface{}) iterable.Iterator[interface{}]
}

func (u untypedGen) Name() string {
	return u.name()
}

func (u untypedGen) Random(rnd Rand, size int) interface{} {
	return u.random(rnd, size)
}

func (u untypedGen) Enumerate(depth int) iterable.Iterator[interface{}] {
	return u.enumerate(depth)
}

func (u untypedGen) Shrink(elem interface{}) iterable.Iterator[interface{}] {
	return u.shrink(elem)
}

func ToUntyped[T any](gen Generator[T]) UntypedGenerator {
	return untypedGen{
		name: gen.Name,
		random: func(rnd Rand, size int) interface{} {
			return gen.Random(rnd, size)
		},
		enumerate: func(depth int) iterable.Iterator[interface{}] {
			return iterable.MapIterator(func(e T) interface{} {
				return e
			})(gen.Enumerate(depth))
		},
		shrink: func(elem interface{}) iterable.Iterator[interface{}] {
			elemT, ok := elem.(T)
			if !ok {
				panic(fmt.Errorf("could not convert element %#v", elemT))
			}
			return iterable.MapIterator(func(e T) interface{} {
				return e
			})(gen.Shrink(elemT))
		},
	}
}

type Rand interface {
	// Fork this random number generator for controlling a subgroup of the test
	Fork(name string) Rand
}
