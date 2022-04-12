package generator

import (
	"math/big"

	"github.com/peterzeller/go-fun/equality"
	"github.com/peterzeller/go-fun/iterable"
	"github.com/peterzeller/go-fun/slice"
	"github.com/peterzeller/go-fun/zero"
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
		GenRValue: func(rv RandomValue[T]) (T, bool) {
			r := rv.Value.(OneOfRandom[T])
			// TODO check if the generator is still valid:

			return r.generator.RValue(r.value)
		},
	}
}

type OneOfRandom[T any] struct {
	generator Generator[T]
	value     RandomValue[T]
}

func OneConstantOf[T comparable](values ...T) Generator[T] {
	if len(values) == 0 {
		return Constant(zero.Value[T]())
	}
	return &AnonGenerator[T]{
		GenName: "OneConstantOf",
		GenRandom: func(rnd Rand, size int) RandomValue[T] {
			n := rnd.R().Intn(len(values))
			g := values[n]
			return RandomValue[T]{
				Value: g,
			}
		},
		GenEnumerate: func(depth int) iterable.Iterable[T] {
			return iterable.Take(depth, iterable.FromSlice(values))
		},
		GenShrink: func(elem RandomValue[T]) iterable.Iterable[RandomValue[T]] {
			v := elem.Get()
			i := slice.IndexOf(v, values, equality.Default[T]())
			if i >= len(values) {
				i = len(values) - 1
			}
			if i <= 0 {
				return iterable.Empty[RandomValue[T]]()
			}
			return iterable.Singleton(R(values[i]))
		},
		GenSize: func(rv RandomValue[T]) *big.Int {
			v := rv.Get()
			i := slice.IndexOf(v, values, equality.Default[T]())
			return big.NewInt(int64(i))
		},
		GenRValue: func(rv RandomValue[T]) (T, bool) {
			v := rv.Get()
			if !slice.Contains(values, v) {
				if len(values) > 0 {
					// fall back to first option
					return values[0], true
				}
				return zero.Value[T](), false
			}
			return v, true
		},
	}
}
