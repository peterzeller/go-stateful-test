package generator

import (
	"fmt"
	"github.com/peterzeller/go-fun/iterable"
	"github.com/peterzeller/go-stateful-test/generator/geniterable"
	"math/big"

	"github.com/peterzeller/go-fun/zero"
)

// UntypedR returns a generator where the R value is untyped (using interface{}).
func UntypedR[T, R any](g Generator[T, R]) Generator[T, interface{}] {
	toInterface := func(r R) interface{} {
		return r
	}
	return &AnonGenerator[T, interface{}]{
		GenName: fmt.Sprintf("UntypedR(%s)", g.Name()),
		GenRandom: func(rnd Rand, size int) interface{} {
			return g.Random(rnd, size)
		},
		GenShrink: func(elem interface{}) iterable.Iterable[interface{}] {
			return iterable.Map(g.Shrink(elem.(R)), toInterface)
		},
		GenSize: func(t interface{}) *big.Int {
			return g.Size(t.(R))
		},
		GenRValue: func(r interface{}) (T, bool) {
			return g.RValue(r.(R))
		},
		GenEnumerate: func(depth int) geniterable.Iterable[interface{}] {
			return geniterable.Map(g.Enumerate(depth), toInterface)
		},
	}
}

// Map transfers a generator for type A to a generator of type B
func Map[A, B, RA any](aGen Generator[A, RA], toB func(a A) B) Generator[B, RA] {
	return &AnonGenerator[B, RA]{
		GenName: fmt.Sprintf("Map(%s)", aGen.Name()),
		GenRandom: func(rnd Rand, size int) RA {
			rv := aGen.Random(rnd, size)
			return rv
		},
		GenShrink: func(rv RA) iterable.Iterable[RA] {
			return aGen.Shrink(rv)
		},
		GenSize: func(rv RA) *big.Int {
			return aGen.Size(rv)
		},
		GenRValue: func(rv RA) (B, bool) {
			a, ok := aGen.RValue(rv)
			if !ok {
				return zero.Value[B](), false
			}
			return toB(a), true
		},
		GenEnumerate: func(depth int) geniterable.Iterable[RA] {
			return aGen.Enumerate(depth)
		},
	}
}

// FlatMap works like Map but uses a function that returns another generator.
// This allows combining two generators into one, similar to a dependent product.
// If the two generators are independent, use the Zip function instead.
func FlatMap[A, RA, B, RB any](aGen Generator[A, RA], toB func(a A) Generator[B, RB]) Generator[B, flatmapRv[RA, RB]] {
	return &AnonGenerator[B, flatmapRv[RA, RB]]{
		GenName: fmt.Sprintf("FlatMap(%s)", aGen.Name()),
		GenRandom: func(rnd Rand, size int) flatmapRv[RA, RB] {
			aRv := aGen.Random(rnd, size)
			v, ok := aGen.RValue(aRv)
			if !ok {
				panic(fmt.Errorf("invalid random value generated: %v", aRv))
			}
			bRv := toB(v).Random(rnd, size)
			return flatmapRv[RA, RB]{
				aRv: aRv,
				bRv: bRv,
			}
		},
		GenShrink: func(rv flatmapRv[RA, RB]) iterable.Iterable[flatmapRv[RA, RB]] {
			fRv := rv
			aShrinks := iterable.FlatMap(
				aGen.Shrink(fRv.aRv),
				func(aRv RA) iterable.Iterable[flatmapRv[RA, RB]] {
					av, ok := aGen.RValue(fRv.aRv)
					if !ok {
						return iterable.Empty[flatmapRv[RA, RB]]()
					}
					bGen := toB(av)
					return iterable.Map(
						bGen.Shrink(fRv.bRv),
						func(bRv RB) flatmapRv[RA, RB] {
							return flatmapRv[RA, RB]{
								aRv: aRv,
								bRv: bRv,
							}
						})
				})
			av, ok := aGen.RValue(fRv.aRv)
			if ok {
				bShrinks := iterable.Map(
					toB(av).Shrink(fRv.bRv),
					func(bRv RB) flatmapRv[RA, RB] {
						return flatmapRv[RA, RB]{
							aRv: fRv.aRv,
							bRv: bRv,
						}
					})
				return iterable.Concat(aShrinks, bShrinks)
			}
			return aShrinks
		},
		GenSize: func(rv flatmapRv[RA, RB]) *big.Int {
			av, ok := aGen.RValue(rv.aRv)
			if !ok {
				return big.NewInt(0)
			}
			bGen := toB(av)
			res := aGen.Size(rv.aRv)
			res.Add(res, bGen.Size(rv.bRv))
			return res
		},
		GenRValue: func(rv flatmapRv[RA, RB]) (B, bool) {
			av, ok := aGen.RValue(rv.aRv)
			if !ok {
				return zero.Value[B](), false
			}
			bGen := toB(av)
			return bGen.RValue(rv.bRv)
		},
		GenEnumerate: func(depth int) geniterable.Iterable[flatmapRv[RA, RB]] {
			return geniterable.FlatMap(
				aGen.Enumerate(depth),
				func(aRv RA) geniterable.Iterable[flatmapRv[RA, RB]] {
					fRv := aRv
					av, ok := aGen.RValue(fRv)
					if !ok {
						return geniterable.Empty[flatmapRv[RA, RB]]()
					}
					bGen := toB(av)
					return geniterable.Map(bGen.Enumerate(depth), func(bRv RB) flatmapRv[RA, RB] {
						return flatmapRv[RA, RB]{
							aRv: aRv,
							bRv: bRv,
						}
					})
				})
		},
	}
}

type flatmapRv[RA, RB any] struct {
	aRv RA
	bRv RB
}

type zipRv[RA, RB any] struct {
	aRv RA
	bRv RB
}

// Zip combines two generators into one using a combine-function.
func Zip[A, B, C, RA, RB any](aGen Generator[A, RA], bGen Generator[B, RB], combine func(a A, b B) C) Generator[C, zipRv[RA, RB]] {
	return &AnonGenerator[C, zipRv[RA, RB]]{
		GenName: fmt.Sprintf("Zip(%s, %s)", aGen.Name(), bGen.Name()),
		GenRandom: func(rnd Rand, size int) zipRv[RA, RB] {
			aRv := aGen.Random(rnd, size)
			bRv := bGen.Random(rnd, size)
			return zipRv[RA, RB]{
				aRv: aRv,
				bRv: bRv,
			}
		},
		GenShrink: func(rv zipRv[RA, RB]) iterable.Iterable[zipRv[RA, RB]] {
			fRv := rv
			aShrinks := iterable.Map(
				aGen.Shrink(fRv.aRv),
				func(aRv RA) zipRv[RA, RB] {
					return zipRv[RA, RB]{
						aRv: aRv,
						bRv: fRv.bRv,
					}
				})
			bShrinks := iterable.Map(
				bGen.Shrink(fRv.bRv),
				func(bRv RB) zipRv[RA, RB] {
					return zipRv[RA, RB]{
						aRv: fRv.aRv,
						bRv: bRv,
					}
				})
			return iterable.Concat(aShrinks, bShrinks)
		},
		GenSize: func(rv zipRv[RA, RB]) *big.Int {
			fRv := rv
			res := aGen.Size(fRv.aRv)
			res.Add(res, bGen.Size(fRv.bRv))
			return res
		},
		GenRValue: func(rv zipRv[RA, RB]) (C, bool) {
			fRv := rv
			av, ok := aGen.RValue(fRv.aRv)
			if !ok {
				return zero.Value[C](), false
			}
			bv, ok := bGen.RValue(fRv.bRv)
			if !ok {
				return zero.Value[C](), false
			}
			return combine(av, bv), true
		},
		GenEnumerate: func(depth int) geniterable.Iterable[zipRv[RA, RB]] {
			return geniterable.FlatMap(
				aGen.Enumerate(depth),
				func(a RA) geniterable.Iterable[zipRv[RA, RB]] {
					return geniterable.Map(
						bGen.Enumerate(depth),
						func(b RB) zipRv[RA, RB] {
							return zipRv[RA, RB]{
								aRv: a,
								bRv: b,
							}
						})
				})
		},
	}
}

// Filter a generator by a predicate.
func Filter[A, RA any](gen Generator[A, RA], predicate func(a A) bool) Generator[A, RA] {
	return &AnonGenerator[A, RA]{
		GenName: fmt.Sprintf("Filter(%s)", gen.Name()),
		GenRandom: func(rnd Rand, size int) RA {
			for i := 0; i < 1000; i++ {
				rv := gen.Random(rnd, size)
				v, ok := gen.RValue(rv)
				if ok && predicate(v) {
					return rv
				}
			}
			// give up and return unrestricted value (will be caught in RValue later)
			return gen.Random(rnd, size)
		},
		GenShrink: func(rv RA) iterable.Iterable[RA] {
			return iterable.Filter(
				gen.Shrink(rv),
				func(rv2 RA) bool {
					v, ok := gen.RValue(rv2)
					return ok && predicate(v)
				})
		},
		GenSize: func(rv RA) *big.Int {
			return gen.Size(rv)
		},
		GenRValue: func(rv RA) (A, bool) {
			v, ok := gen.RValue(rv)
			return v, ok && predicate(v)
		},
		GenEnumerate: func(depth int) geniterable.Iterable[RA] {
			return geniterable.Filter(
				gen.Enumerate(depth),
				func(ra RA) bool {
					a, ok := gen.RValue(ra)
					return ok && predicate(a)
				})
		},
	}
}

// FilterMap transfers a generator for type A to a generator of type B
func FilterMap[A, RA, B any](aGen Generator[A, RA], toB func(a A) (B, bool)) Generator[B, RA] {
	return Map(
		Filter(aGen, func(a A) bool {
			_, ok := toB(a)
			return ok
		}), func(a A) B {
			b, _ := toB(a)
			return b
		})
}
