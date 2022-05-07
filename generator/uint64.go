package generator

import (
	"math"
	"math/big"

	"github.com/peterzeller/go-fun/iterable"
)

func UInt64() Generator[uint64] {
	return genUInt64{
		min: 0,
		max: math.MaxUint64,
	}
}

func UInt64Range(min uint64, max uint64) Generator[uint64] {
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

func (g genUInt64) Random(rnd Rand, size int) RandomValue[uint64] {
	r := rnd.R()
	p := r.Float64()
	n := 1 + g.max - g.min
	switch {
	// higher probability for 0
	case p < 0.05 && g.min < 0 && g.max < 0:
		return R(uint64(0))
	// higher probability for boundary cases
	case p < 0.1:
		return R(g.min)
	case p < 0.15:
		return R(g.max)
	case p < 0.5:
		// normal distribution around 0
		res := uint64(math.Abs(r.NormFloat64()) * 3)
		if res > g.max {
			res = g.max
		}
		return R(res)
	default:
		if n < math.MaxInt64 {
			return R(g.min + uint64(r.Int63n(int64(n))))
		}
		v := g.min + (r.Uint64() % n)
		return R(v)
	}
}

func (g genUInt64) Enumerate(depth int) iterable.Iterable[uint64] {
	return iterable.Take(depth, iterable.RangeI(g.min, g.max))
}

func (g genUInt64) Shrink(r RandomValue[uint64]) iterable.Iterable[RandomValue[uint64]] {
	elem := r.Get()
	if elem == 0 {
		return iterable.Empty[RandomValue[uint64]]()
	}
	return iterable.Map(iterable.New(elem/2, elem-1), R[uint64])
}

func (g genUInt64) RValue(r RandomValue[uint64]) (uint64, bool) {
	if g.min > g.max {
		return 0, false
	}
	elem := r.Get()
	if elem < g.min {
		return g.min, true
	}
	if elem > g.max {
		return g.max, true
	}
	return elem, true
}

func (g genUInt64) Size(r RandomValue[uint64]) *big.Int {
	elem := r.Get()
	if elem < 0 {
		r := big.NewInt(int64(elem))
		r.Abs(r)
		r.Add(r, big.NewInt(1))
		return r
	}
	return big.NewInt(int64(elem))
}
