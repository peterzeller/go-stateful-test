package generator

import (
	"math"
)

func UInt16() Generator[uint16] {
	return UInt16Range(0, math.MaxUint16)
}

func UInt16Range(min uint16, max uint16) Generator[uint16] {
	return Map(UInt64Range(uint64(min), uint64(max)),
		func(i uint64) uint16 {
			return uint16(i)
		})
}
