package generator

import (
	"github.com/peterzeller/go-fun/iterable"
	"github.com/peterzeller/go-stateful-test/generator/geniterable"
	"github.com/peterzeller/go-stateful-test/quickcheck/randomsource"
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

// findMSBPosition uses a binary search approach to find the position of the most significant bit.
func findMSBPosition(n uint64) int {
	if n == 0 {
		return 0
	}

	position := 0
	if n >= 1<<32 {
		n >>= 32
		position += 32
	}
	if n >= 1<<16 {
		n >>= 16
		position += 16
	}
	if n >= 1<<8 {
		n >>= 8
		position += 8
	}
	if n >= 1<<4 {
		n >>= 4
		position += 4
	}
	if n >= 1<<2 {
		n >>= 2
		position += 2
	}
	if n >= 1<<1 {
		position += 1
	}
	return position + 1 // Add 1 because positions start at 1.
}

func (g genUInt64) Random(rnd Rand, size int) uint64 {
	if !rnd.UseHeuristics() {
		interval := g.max - g.min
		bitSize := findMSBPosition(interval)
		i := randomsource.Uint64B(rnd.R(), 1+(bitSize-1)/8)
		i = i % interval
		i = i + g.min
		return i
	}
	p := randomsource.Float64(rnd.R())
	n := 1 + g.max - g.min
	switch {
	// higher probability for boundary cases
	case p < 0.1:
		return (g.min)
	case p < 0.15:
		return (g.max)
	case p < 0.5:
		// normal distribution around 0
		res := uint64(math.Abs(randomsource.NormFloat64(rnd.R())) * 3)
		if res > g.max {
			res = g.max
		}
		return (res)
	default:
		if n < math.MaxInt64 {
			return g.min + uint64(randomsource.Int64N(rnd.R(), int64(n)))
		}
		v := g.min + (randomsource.Uint64(rnd.R()) % n)
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
