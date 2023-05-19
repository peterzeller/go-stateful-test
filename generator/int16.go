package generator

import (
	"math"
)

func Int16() Generator[int16, int64] {
	return Int16Range(math.MinInt16, math.MaxInt16)
}

func Int16Range(min int16, max int16) Generator[int16, int64] {
	return Map(Int64Range(int64(min), int64(max)),
		func(i int64) int16 {
			return int16(i)
		})
}
