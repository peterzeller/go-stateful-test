package generator

import (
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

type Rand interface {
	// Fork this random number generator for controlling a subgroup of the test
	Fork(name string) Rand
}
