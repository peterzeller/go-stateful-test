package generator

import (
	"math"
)

func Int() Generator[int] {
	return IntRange(math.MinInt, math.MaxInt)
}

func IntRange(min int, max int) Generator[int] {
	return Map(Int64Range(int64(min), int64(max)),
		func(i int64) int {
			return int(i)
		})
}
