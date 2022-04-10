package generator

import (
	"github.com/peterzeller/go-fun/iterable"
	"github.com/peterzeller/go-fun/zero"
	"math/big"
)

func OneOf[T any](gs ...Generator[T]) Generator[T] {
	if len(gs) == 0 {
		return Constant(zero.Value[T]())
	}
	return &AnonGenerator[T]{
		GenName: "OneOf",
		GenRandom: func(rnd Rand, size int) RandomValue[T] {
			n := rnd.R().Intn(len(gs))
			g := gs[n]
			return RandomValue[T]{
				Value: OneOfRandom[T]{
					generator: g,
					value:     g.Random(rnd, size),
				},
			}
		},
		GenEnumerate: func(depth int) iterable.Iterable[T] {
			return iterable.FlatMapBreadthFirst(iterable.FromSlice(gs),
				func(g Generator[T]) iterable.Iterable[T] {
					return g.Enumerate(depth)
				})
		},
		GenShrink: func(elem RandomValue[T]) iterable.Iterable[RandomValue[T]] {
			r := elem.Value.(OneOfRandom[T])
			return iterable.Map(
				r.generator.Shrink(r.value),
				func(v RandomValue[T]) RandomValue[T] {
					return RandomValue[T]{
						Value: OneOfRandom[T]{
							generator: r.generator,
							value:     v,
						},
					}
				})
		},
		GenSize: func(rv RandomValue[T]) *big.Int {
			r := rv.Value.(OneOfRandom[T])
			return r.generator.Size(r.value)
		},
	}
}

type OneOfRandom[T any] struct {
	generator Generator[T]
	value     RandomValue[T]
}
