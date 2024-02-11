package generator

import (
	"github.com/peterzeller/go-fun/iterable"
	"github.com/peterzeller/go-stateful-test/generator/geniterable"
	"math"
	"math/big"
)

func UInt64() Generator[uint64, uint64] {
	return genUInt64{
		min: 0,
		max: math.MaxUint64,
	}
}

func UInt64Range(min uint64, max uint64) Generator[uint64, uint64] {
	return genUInt64{
		min: min,
		max: max,
	}
}

type genUInt64 struct {
	min uint64
	max uint64
}

func (g genUInt64) Name() string {
	return "genUInt64"
}

func (g genUInt64) Random(rnd Rand, size int) uint64 {
	r := rnd.R()
	p := r.Float64()
	n := 1 + g.max - g.min
	switch {
	// higher probability for boundary cases
	case p < 0.1:
		return (g.min)
	case p < 0.15:
		return (g.max)
	case p < 0.5:
		// normal distribution around 0
		res := uint64(math.Abs(r.NormFloat64()) * 3)
		if res > g.max {
			res = g.max
		}
		return (res)
	default:
		if n < math.MaxInt64 {
			return (g.min + uint64(r.Int63n(int64(n))))
		}
		v := g.min + (r.Uint64() % n)
		return (v)
	}
}

func (g genUInt64) Enumerate(depth int) geniterable.Iterable[uint64] {
	return geniterable.TakeExhaustive(depth, geniterable.RangeI(g.min, g.max))
}

func (g genUInt64) Shrink(r uint64) iterable.Iterable[uint64] {
	elem := r
	if elem == 0 {
		return iterable.Empty[uint64]()
	}
	return iterable.New(elem/2, elem-1)
}

func (g genUInt64) RValue(r uint64) (uint64, bool) {
	if g.min > g.max {
		return 0, false
	}
	elem := r
	if elem < g.min {
		return g.min, true
	}
	if elem > g.max {
		return g.max, true
	}
	return elem, true
}

func (g genUInt64) Size(r uint64) *big.Int {
	elem := r
	return big.NewInt(int64(elem))
}
