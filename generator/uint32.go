package generator

import (
	"math"
)

func UInt32() Generator[uint32, uint64] {
	return UInt32Range(0, math.MaxUint32)
}

func UInt32Range(min uint32, max uint32) Generator[uint32, uint64] {
	return Map(UInt64Range(uint64(min), uint64(max)),
		func(i uint64) uint32 {
			return uint32(i)
		})
}
