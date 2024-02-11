package generator

import (
	"math"
)

func UInt8() Generator[uint8, uint64] {
	return UInt8Range(0, math.MaxUint8)
}

func UInt8Range(min uint8, max uint8) Generator[uint8, uint64] {
	return Map(UInt64Range(uint64(min), uint64(max)),
		func(i uint64) uint8 {
			return uint8(i)
		})
}
