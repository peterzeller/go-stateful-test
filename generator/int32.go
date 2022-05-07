package generator

import (
	"math"
)

func Int32() Generator[int32] {
	return Int32Range(math.MinInt32, math.MaxInt32)
}

func Int32Range(min int32, max int32) Generator[int32] {
	return Map(Int64Range(int64(min), int64(max)),
		func(i int64) int32 {
			return int32(i)
		})
}
