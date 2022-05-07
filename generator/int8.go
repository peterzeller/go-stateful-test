package generator

import (
	"math"
)

func Int8() Generator[int8] {
	return Int8Range(math.MinInt8, math.MaxInt8)
}

func Int8Range(min int8, max int8) Generator[int8] {
	return Map(Int64Range(int64(min), int64(max)),
		func(i int64) int8 {
			return int8(i)
		})
}
