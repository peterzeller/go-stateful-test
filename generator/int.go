package generator

import (
	"math"
	"math/big"

	"github.com/peterzeller/go-fun/iterable"
)

func Int() Generator[int] {
	return genInt{
		min: math.MinInt,
		max: math.MaxInt,
	}
}

func IntRange(min int, max int) Generator[int] {
	return genInt{
		min: min,
		max: max,
	}
}

type genInt struct {
	min int
	max int
}

func (g genInt) Name() string {
	return "genInt"
}

func (g genInt) Random(rnd Rand, size int) RandomValue[int] {
	r := rnd.R()
	p := r.Float64()
	n := 1 + g.max - g.min
	switch {
	// higher probability for 0
	case p < 0.05 && g.min < 0 && g.max < 0:
		return R(0)
	// higher probability for error cases
	case p < 0.1:
		return R(g.min)
	case p < 0.15:
		return R(g.max)
	case p < 0.8 && g.min <= 0 && 0 < g.max:
		// normal distribution around 0
		res := int(math.Abs(r.NormFloat64()) * 3)
		if res > g.max {
			res = g.max
		}
		return R(res)
	default:
		if n > 0 {
			// uniform distribution
			return R(g.min + r.Intn(n))
		}
		if r.Float64() < 0.1 {
			// negative number
			return R(-1 + r.Intn(-1+g.min))
		}
		return R(1 + r.Intn(g.max))
	}
}

func (g genInt) Enumerate(depth int) iterable.Iterable[int] {
	if g.min == math.MinInt && g.max == math.MaxInt {
		return iterable.Take(depth,
			iterable.FlatMap(
				iterable.Generate(0, func(i int) int { return i + 1 }),
				func(a int) iterable.Iterable[int] {
					if a == 0 {
						return iterable.Singleton(0)
					}
					return iterable.New(a, -a)
				}))
	}
	return iterable.Take(depth, iterable.RangeI(g.min, g.max))
}

func (g genInt) Shrink(r RandomValue[int]) iterable.Iterable[RandomValue[int]] {
	elem := r.Get()
	if elem == 0 {
		return iterable.Empty[RandomValue[int]]()
	}
	if elem < 0 {
		return iterable.Map(iterable.New(elem/2, -elem, elem+1), R[int])
	} else {
		return iterable.Map(iterable.New(elem/2, elem-1), R[int])
	}
}

func (g genInt) RValue(r RandomValue[int]) (int, bool) {
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

func (g genInt) Size(r RandomValue[int]) *big.Int {
	elem := r.Get()
	if elem < 0 {
		r := big.NewInt(int64(elem))
		r.Abs(r)
		r.Add(r, big.NewInt(1))
		return r
	}
	return big.NewInt(int64(elem))
}
