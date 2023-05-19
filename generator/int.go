package generator

import (
	"math"
)

func Int() Generator[int, int64] {
	return IntRange(math.MinInt, math.MaxInt)
}

func IntRange(min int, max int) Generator[int, int64] {
	return Map(Int64Range(int64(min), int64(max)),
		func(i int64) int {
			return int(i)
		})
}
