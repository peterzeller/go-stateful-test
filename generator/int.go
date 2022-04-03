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

func (g genInt) Random(rnd Rand, size int) int {
	r := rnd.R()
	p := r.Float64()
	n := 1 + g.max - g.min
	switch {
	// higher probability for 0
	case p < 0.05 && g.min < 0 && g.max < 0:
		return 0
	// higher probability for error cases
	case p < 0.1:
		return g.min
	case p < 0.15:
		return g.max
	default:
		if n > 0 {
			// uniform distribution
			return g.min + r.Intn(n)
		}
		if r.Float64() < 0.1 {
			// negative number
			return -1 + r.Intn(-1+g.min)
		}
		return 1 + r.Intn(g.max)
	}
}

func (g genInt) Enumerate(depth int) iterable.Iterable[int] {
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

func (g genInt) Shrink(elem int) iterable.Iterable[int] {
	if elem == 0 {
		return iterable.Empty[int]()
	}
	if elem < 0 {
		return iterable.New(elem/2, -elem, elem+1)
	} else {
		return iterable.New(elem/2, elem-1)
	}
}

func (g genInt) Size(elem int) *big.Int {
	if elem < 0 {
		return big.NewInt(int64(1 - elem))
	}
	return big.NewInt(int64(elem))
}
