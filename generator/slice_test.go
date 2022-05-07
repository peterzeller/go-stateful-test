package generator

import (
	"fmt"
	"testing"

	"github.com/peterzeller/go-fun/equality"
	"github.com/peterzeller/go-fun/iterable"
	"github.com/stretchr/testify/require"
)

func TestSliceEnumerate(t *testing.T) {
	g := Slice(IntRange(1, 3))
	enumerated := iterable.ToSlice(g.Enumerate(3))
	require.Equal(t, [][]int{
		{},
		{1},
		{2},
		{3},
		{1, 1},
		{2, 1},
		{3, 1},
		{1, 2},
		{2, 2},
		{3, 2},
		{1, 3},
		{2, 3},
		{3, 3},
		{1, 1, 1},
		{2, 1, 1},
		{3, 1, 1},
		{1, 2, 1},
		{2, 2, 1},
		{3, 2, 1},
		{1, 3, 1},
		{2, 3, 1},
		{3, 3, 1},
		{1, 1, 2},
		{2, 1, 2},
		{3, 1, 2},
		{1, 2, 2},
		{2, 2, 2},
		{3, 2, 2},
		{1, 3, 2},
		{2, 3, 2},
		{3, 3, 2},
		{1, 1, 3},
		{2, 1, 3},
		{3, 1, 3},
		{1, 2, 3},
		{2, 2, 3},
		{3, 2, 3},
		{1, 3, 3},
		{2, 3, 3},
		{3, 3, 3}},
		enumerated)
}

func TestSliceRandom(t *testing.T) {
	g := Slice(IntRange(1, 5))
	rnd := newTestRand(2)
	rv := g.Random(rnd, 10)
	t.Logf("rv = %+v", rv)
	v, ok := g.RValue(rv)
	require.True(t, ok)
	require.Equal(t, []int{3, 5, 1, 3, 2, 4}, v)
}

func TestSliceShrink(t *testing.T) {
	g := Slice(IntRange(0, 10))
	rv := sliceRandomValue([]int64{4, 5, 6, 7, 8})
	shrinks := iterable.ToSlice(iterable.Map(g.Shrink(rv), func(rv RandomValue[[]int]) []int {
		v, ok := g.RValue(rv)
		require.True(t, ok)
		return v
	}))
	require.Equal(t, [][]int{
		// remove all 5
		{},
		// remove 2
		{4, 5, 8},
		{6, 7, 8},
		// remove 1
		{4, 5, 6, 7},
		{4, 5, 6, 8},
		{4, 5, 7, 8},
		{4, 6, 7, 8},
		{5, 6, 7, 8},
		// shrink single elements
		{2, 5, 6, 7, 8},
		{3, 5, 6, 7, 8},
		{4, 2, 6, 7, 8},
		{4, 4, 6, 7, 8},
		{4, 5, 3, 7, 8},
		{4, 5, 5, 7, 8},
		{4, 5, 6, 3, 8},
		{4, 5, 6, 6, 8},
		{4, 5, 6, 7, 4},
		{4, 5, 6, 7, 7}},
		shrinks)
}

func sliceRandomValue(elem []int64) RandomValue[[]int] {
	rvs := make([]RandomValue[int], len(elem))
	for i, e := range elem {
		rvs[i] = RandomValue[int]{
			Value: e,
		}
	}
	return RandomValue[[]int]{
		Value: rvs,
	}
}

func ExampleSliceFixedLength() {
	g := SliceFixedLength(IntRange(1, 3), 3)
	for it := iterable.Start(g.Enumerate(100)); it.HasNext(); it.Next() {
		fmt.Printf("%+v\n", it.Current())
	}
	// Output: [1 1 1]
	// [1 1 2]
	// [1 1 3]
	// [1 2 1]
	// [1 2 2]
	// [1 2 3]
	// [1 3 1]
	// [1 3 2]
	// [1 3 3]
	// [2 1 1]
	// [2 1 2]
	// [2 1 3]
	// [2 2 1]
	// [2 2 2]
	// [2 2 3]
	// [2 3 1]
	// [2 3 2]
	// [2 3 3]
	// [3 1 1]
	// [3 1 2]
	// [3 1 3]
	// [3 2 1]
	// [3 2 2]
	// [3 2 3]
	// [3 3 1]
	// [3 3 2]
	// [3 3 3]
}

func ExampleSliceDistinct() {
	g := SliceDistinct(IntRange(1, 3), equality.Default[int]())
	for it := iterable.Start(g.Enumerate(3)); it.HasNext(); it.Next() {
		fmt.Printf("%+v\n", it.Current())
	}
	// Output: []
	// [1]
	// [2]
	// [3]
	// [2 1]
	// [3 1]
	// [1 2]
	// [3 2]
	// [1 3]
	// [2 3]
	// [3 2 1]
	// [2 3 1]
	// [3 1 2]
	// [1 3 2]
	// [2 1 3]
	// [1 2 3]
}
