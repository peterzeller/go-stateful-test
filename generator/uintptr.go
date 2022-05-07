package generator

import (
	"math"
)

func Uintptr() Generator[uintptr] {
	return UintptrRange(0, math.MaxUint64)
}

func UintptrRange(min uintptr, max uintptr) Generator[uintptr] {
	return Map(UInt64Range(uint64(min), uint64(max)),
		func(i uint64) uintptr {
			return uintptr(i)
		})
}
