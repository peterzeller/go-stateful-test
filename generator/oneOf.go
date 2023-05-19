package generator

import (
	"github.com/peterzeller/go-fun/iterable"
	"github.com/peterzeller/go-stateful-test/generator/geniterable"
	"math/big"

	"github.com/peterzeller/go-fun/equality"
	"github.com/peterzeller/go-fun/slice"
	"github.com/peterzeller/go-fun/zero"
)

func OneOf[T, R any](gs ...Generator[T, R]) Generator[T, OneOfRandom[R]] {
	if len(gs) == 0 {
		return Empty[T, OneOfRandom[R]]()
	}
	return &AnonGenerator[T, OneOfRandom[R]]{
		GenName: "OneOf",
		GenRandom: func(rnd Rand, size int) OneOfRandom[R] {
			n := rnd.R().Intn(len(gs))
			g := gs[n]
			return OneOfRandom[R]{
				generator: n,
				value:     g.Random(rnd, size),
			}
		},
		GenEnumerate: func(depth int) geniterable.Iterable[OneOfRandom[R]] {
			return geniterable.FlatMapBreadthFirst(geniterable.Range(0, len(gs)),
				func(i int) geniterable.Iterable[OneOfRandom[R]] {
					return geniterable.Map(gs[i].Enumerate(depth), func(a R) OneOfRandom[R] {
						return OneOfRandom[R]{
							generator: i,
							value:     a,
						}
					})
				})
		},
		GenShrink: func(elem OneOfRandom[R]) iterable.Iterable[OneOfRandom[R]] {
			if elem.generator < 0 || elem.generator >= len(gs) {
				return iterable.Empty[OneOfRandom[R]]()
			}
			r := elem
			g := gs[elem.generator]
			return iterable.Map(
				g.Shrink(r.value),
				func(v R) OneOfRandom[R] {
					return OneOfRandom[R]{
						generator: r.generator,
						value:     v,
					}
				})
		},
		GenSize: func(rv OneOfRandom[R]) *big.Int {
			if rv.generator < 0 || rv.generator >= len(gs) {
				return big.NewInt(0)
			}
			r := rv
			return gs[r.generator].Size(r.value)
		},
		GenRValue: func(rv OneOfRandom[R]) (T, bool) {
			if rv.generator < 0 || rv.generator >= len(gs) {
				return zero.Value[T](), false
			}
			return gs[rv.generator].RValue(rv.value)
		},
	}
}

type OneOfRandom[R any] struct {
	generator int
	value     R
}

func OneConstantOf[T comparable](values ...T) Generator[T, T] {
	if len(values) == 0 {
		return Empty[T, T]()
	}
	return &AnonGenerator[T, T]{
		GenName: "OneConstantOf",
		GenRandom: func(rnd Rand, size int) T {
			n := rnd.R().Intn(len(values))
			g := values[n]
			return g
		},
		GenEnumerate: func(depth int) geniterable.Iterable[T] {
			return geniterable.TakeExhaustive(depth, geniterable.FromSlice(values))
		},
		GenShrink: func(elem T) iterable.Iterable[T] {
			v := elem
			i := slice.IndexOf(v, values, equality.Default[T]())
			if i >= len(values) {
				i = len(values) - 1
			}
			if i <= 0 {
				return iterable.Empty[T]()
			}
			return iterable.Singleton(values[i])
		},
		GenSize: func(rv T) *big.Int {
			v := rv
			i := slice.IndexOf(v, values, equality.Default[T]())
			return big.NewInt(int64(i))
		},
		GenRValue: func(rv T) (T, bool) {
			v := rv
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

// Bool returns a generator for Boolean values.
func Bool() Generator[bool, bool] {
	return OneConstantOf(false, true)
}
