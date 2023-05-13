package generator

import (
	"fmt"
	"github.com/peterzeller/go-fun/iterable"
	"github.com/peterzeller/go-stateful-test/generator/geniterable"
	"math/big"

	"github.com/peterzeller/go-fun/zero"
)

// Map transfers a generator for type A to a generator of type B
func Map[A, B any](aGen Generator[A], toB func(a A) B) Generator[B] {
	return &AnonGenerator[B]{
		GenName: fmt.Sprintf("Map(%s)", aGen.Name()),
		GenRandom: func(rnd Rand, size int) RandomValue[B] {
			rv := aGen.Random(rnd, size)
			return RandomValue[B](rv)
		},
		GenShrink: func(rv RandomValue[B]) iterable.Iterable[RandomValue[B]] {
			return iterable.Map(
				aGen.Shrink(RandomValue[A](rv)),
				func(rv RandomValue[A]) RandomValue[B] {
					return RandomValue[B](rv)
				})
		},
		GenSize: func(rv RandomValue[B]) *big.Int {
			return aGen.Size(RandomValue[A](rv))
		},
		GenRValue: func(rv RandomValue[B]) (B, bool) {
			a, ok := aGen.RValue(RandomValue[A](rv))
			if !ok {
				return zero.Value[B](), false
			}
			return toB(a), true
		},
		GenEnumerate: func(depth int) geniterable.Iterable[B] {
			return geniterable.Map(
				aGen.Enumerate(depth),
				toB)
		},
	}
}

// FlatMap works like Map but uses a function that returns another generator.
// This allows combining two generators into one, similar to a dependent product.
// If the two generators are independent, use the Zip function instead.
func FlatMap[A, B any](aGen Generator[A], toB func(a A) Generator[B]) Generator[B] {
	return &AnonGenerator[B]{
		GenName: fmt.Sprintf("FlatMap(%s)", aGen.Name()),
		GenRandom: func(rnd Rand, size int) RandomValue[B] {
			aRv := aGen.Random(rnd, size)
			v, ok := aGen.RValue(aRv)
			if !ok {
				panic(fmt.Errorf("invalid random value generated: %v", aRv))
			}
			bRv := toB(v).Random(rnd, size)
			return RandomValue[B]{
				Value: flatmapRv[A, B]{
					aRv: aRv,
					bRv: bRv,
				},
			}
		},
		GenShrink: func(rv RandomValue[B]) iterable.Iterable[RandomValue[B]] {
			fRv := rv.Value.(flatmapRv[A, B])
			aShrinks := iterable.FlatMap(
				aGen.Shrink(fRv.aRv),
				func(aRv RandomValue[A]) iterable.Iterable[RandomValue[B]] {
					av, ok := aGen.RValue(fRv.aRv)
					if !ok {
						return iterable.Empty[RandomValue[B]]()
					}
					bGen := toB(av)
					return iterable.Map(
						bGen.Shrink(fRv.bRv),
						func(bRv RandomValue[B]) RandomValue[B] {
							return RandomValue[B]{
								Value: flatmapRv[A, B]{
									aRv: aRv,
									bRv: bRv,
								},
							}
						})
				})
			av, ok := aGen.RValue(fRv.aRv)
			if ok {
				bShrinks := iterable.Map(
					toB(av).Shrink(fRv.bRv),
					func(bRv RandomValue[B]) RandomValue[B] {
						return RandomValue[B]{
							Value: flatmapRv[A, B]{
								aRv: fRv.aRv,
								bRv: bRv,
							},
						}
					})
				return iterable.Concat(aShrinks, bShrinks)
			}
			return aShrinks
		},
		GenSize: func(rv RandomValue[B]) *big.Int {
			fRv := rv.Value.(flatmapRv[A, B])
			av, ok := aGen.RValue(fRv.aRv)
			if !ok {
				return big.NewInt(0)
			}
			bGen := toB(av)
			res := aGen.Size(fRv.aRv)
			res.Add(res, bGen.Size(fRv.bRv))
			return res
		},
		GenRValue: func(rv RandomValue[B]) (B, bool) {
			fRv := rv.Value.(flatmapRv[A, B])
			av, ok := aGen.RValue(fRv.aRv)
			if !ok {
				return zero.Value[B](), false
			}
			bGen := toB(av)
			return bGen.RValue(fRv.bRv)
		},
		GenEnumerate: func(depth int) geniterable.Iterable[B] {
			return geniterable.FlatMap(
				aGen.Enumerate(depth),
				func(a A) geniterable.Iterable[B] {
					bGen := toB(a)
					return bGen.Enumerate(depth)
				})
		},
	}
}

type flatmapRv[A, B any] struct {
	aRv RandomValue[A]
	bRv RandomValue[B]
}

type zipRv[A, B any] struct {
	aRv RandomValue[A]
	bRv RandomValue[B]
}

// Zip combines two generators into one using a combine-function.
func Zip[A, B, C any](aGen Generator[A], bGen Generator[B], combine func(a A, b B) C) Generator[C] {
	return &AnonGenerator[C]{
		GenName: fmt.Sprintf("Zip(%s, %s)", aGen.Name(), bGen.Name()),
		GenRandom: func(rnd Rand, size int) RandomValue[C] {
			aRv := aGen.Random(rnd, size)
			bRv := bGen.Random(rnd, size)
			return RandomValue[C]{
				Value: zipRv[A, B]{
					aRv: aRv,
					bRv: bRv,
				},
			}
		},
		GenShrink: func(rv RandomValue[C]) iterable.Iterable[RandomValue[C]] {
			fRv := rv.Value.(zipRv[A, B])
			aShrinks := iterable.Map(
				aGen.Shrink(fRv.aRv),
				func(aRv RandomValue[A]) RandomValue[C] {
					return RandomValue[C]{
						Value: zipRv[A, B]{
							aRv: aRv,
							bRv: fRv.bRv,
						},
					}
				})
			bShrinks := iterable.Map(
				bGen.Shrink(fRv.bRv),
				func(bRv RandomValue[B]) RandomValue[C] {
					return RandomValue[C]{
						Value: zipRv[A, B]{
							aRv: fRv.aRv,
							bRv: bRv,
						},
					}
				})
			return iterable.Concat(aShrinks, bShrinks)
		},
		GenSize: func(rv RandomValue[C]) *big.Int {
			fRv := rv.Value.(zipRv[A, B])
			res := aGen.Size(fRv.aRv)
			res.Add(res, bGen.Size(fRv.bRv))
			return res
		},
		GenRValue: func(rv RandomValue[C]) (C, bool) {
			fRv := rv.Value.(zipRv[A, B])
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
		GenEnumerate: func(depth int) geniterable.Iterable[C] {
			return geniterable.FlatMap(
				aGen.Enumerate(depth),
				func(a A) geniterable.Iterable[C] {
					return geniterable.Map(
						bGen.Enumerate(depth),
						func(b B) C {
							return combine(a, b)
						})
				})
		},
	}
}

// Filter a generator by a predicate.
func Filter[A any](gen Generator[A], predicate func(a A) bool) Generator[A] {
	return &AnonGenerator[A]{
		GenName: fmt.Sprintf("Filter(%s)", gen.Name()),
		GenRandom: func(rnd Rand, size int) RandomValue[A] {
			for i := 0; i < 1000; i++ {
				rv := gen.Random(rnd, size)
				v, ok := gen.RValue(rv)
				if ok && predicate(v) {
					return rv
				}
			}
			// give up and return unrestriced value (will be caught in RValue later)
			return gen.Random(rnd, size)
		},
		GenShrink: func(rv RandomValue[A]) iterable.Iterable[RandomValue[A]] {
			return iterable.Filter(
				gen.Shrink(rv),
				func(rv2 RandomValue[A]) bool {
					v, ok := gen.RValue(rv2)
					return ok && predicate(v)
				})
		},
		GenSize: func(rv RandomValue[A]) *big.Int {
			return gen.Size(rv)
		},
		GenRValue: func(rv RandomValue[A]) (A, bool) {
			v, ok := gen.RValue(rv)
			return v, ok && predicate(v)
		},
		GenEnumerate: func(depth int) geniterable.Iterable[A] {
			return geniterable.Filter(
				gen.Enumerate(depth),
				predicate)
		},
	}
}

// Map transfers a generator for type A to a generator of type B
func FilterMap[A, B any](aGen Generator[A], toB func(a A) (B, bool)) Generator[B] {
	return &AnonGenerator[B]{
		GenName: fmt.Sprintf("Map(%s)", aGen.Name()),
		GenRandom: func(rnd Rand, size int) RandomValue[B] {
			rv := aGen.Random(rnd, size)
			return RandomValue[B](rv)
		},
		GenShrink: func(rv RandomValue[B]) iterable.Iterable[RandomValue[B]] {
			return iterable.Map(
				aGen.Shrink(RandomValue[A](rv)),
				func(rv RandomValue[A]) RandomValue[B] {
					return RandomValue[B](rv)
				})
		},
		GenSize: func(rv RandomValue[B]) *big.Int {
			return aGen.Size(RandomValue[A](rv))
		},
		GenRValue: func(rv RandomValue[B]) (B, bool) {
			a, ok := aGen.RValue(RandomValue[A](rv))
			if !ok {
				return zero.Value[B](), false
			}
			return toB(a)
		},
		GenEnumerate: func(depth int) geniterable.Iterable[B] {
			return geniterable.FlatMap(
				aGen.Enumerate(depth),
				func(a A) geniterable.Iterable[B] {
					b, ok := toB(a)
					if !ok {
						return geniterable.Empty[B]()
					}
					return geniterable.Singleton(b)
				})
		},
	}
}
