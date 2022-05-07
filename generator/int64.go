package generator

import (
	"math"
	"math/big"

	"github.com/peterzeller/go-fun/iterable"
)

func Int64() Generator[int64] {
	return genInt64{
		min: math.MinInt64,
		max: math.MaxInt64,
	}
}

func Int64Range(min int64, max int64) Generator[int64] {
	return genInt64{
		min: min,
		max: max,
	}
}

type genInt64 struct {
	min int64
	max int64
}

func (g genInt64) Name() string {
	return "genInt64"
}

func (g genInt64) Random(rnd Rand, size int) RandomValue[int64] {
	r := rnd.R()
	p := r.Float64()
	n := 1 + g.max - g.min
	switch {
	// higher probability for 0
	case p < 0.05 && g.min < 0 && g.max < 0:
		return R(int64(0))
	// higher probability for boundary cases
	case p < 0.1:
		return R(g.min)
	case p < 0.15:
		return R(g.max)
	case p < 0.8 && g.min <= 0 && 0 < g.max:
		// normal distribution around 0
		res := int64(math.Abs(r.NormFloat64()) * 3)
		if res > g.max {
			res = g.max
		}
		return R(res)
	default:
		if n > 0 {
			// uniform distribution
			return R(g.min + r.Int63n(n))
		}
		if r.Float64() < 0.1 {
			// negative number
			return R(-1 + r.Int63n(-1+g.min))
		}
		return R(1 + r.Int63n(g.max))
	}
}

func (g genInt64) Enumerate(depth int) iterable.Iterable[int64] {
	if g.min < 0 && g.max > 0 {
		return iterable.Take(depth,
			iterable.FlatMap(
				iterable.Generate(0, func(i int64) int64 { return i + 1 }),
				func(a int64) iterable.Iterable[int64] {
					if a == 0 {
						return iterable.Singleton(int64(0))
					}
					res := make([]int64, 0, 2)
					if a <= g.max {
						res = append(res, a)
					}
					if -a >= g.min {
						res = append(res, -a)
					}
					return iterable.New(res...)
				}))
	}
	return iterable.Take(depth, iterable.RangeI(g.min, g.max))
}

func (g genInt64) Shrink(r RandomValue[int64]) iterable.Iterable[RandomValue[int64]] {
	elem := r.Get()
	if elem == 0 {
		return iterable.Empty[RandomValue[int64]]()
	}
	if elem < 0 {
		return iterable.Map(iterable.New(elem/2, -elem, elem+1), R[int64])
	} else {
		return iterable.Map(iterable.New(elem/2, elem-1), R[int64])
	}
}

func (g genInt64) RValue(r RandomValue[int64]) (int64, bool) {
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

func (g genInt64) Size(r RandomValue[int64]) *big.Int {
	elem := r.Get()
	if elem < 0 {
		r := big.NewInt(elem)
		r.Abs(r)
		r.Add(r, big.NewInt(1))
		return r
	}
	return big.NewInt(elem)
}
