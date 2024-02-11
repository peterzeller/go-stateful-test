package generator

import (
	"math"
)

func UInt() Generator[uint, uint64] {
	return UIntRange(0, math.MaxUint)
}

func UIntRange(min uint, max uint) Generator[uint, uint64] {
	return Map(UInt64Range(uint64(min), uint64(max)),
		func(i uint64) uint {
			return uint(i)
		})
}
