package shrink

import (
	"github.com/peterzeller/go-fun/iterable"
	"github.com/peterzeller/go-fun/list/linked"
	"testing"

	"github.com/stretchr/testify/require"
)

// toSlices you say?
func toSlices[T any, I iterable.Iterable[T]](l iterable.Iterable[I]) [][]T {
	return iterable.ToSlice(
		iterable.Map(l, func(i I) []T {
			return iterable.ToSlice[T](i)
		}))
}

func TestRemoves(t *testing.T) {
	list := linked.New(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	require.Equal(t, [][]int{
		{1, 2, 3, 4, 5},
		{6, 7, 8, 9, 10}}, toSlices[int](Removes(5, 10, list)))

	require.Equal(t, [][]int{
		{1, 2, 3, 4, 5, 6, 7, 8},
		{1, 2, 3, 4, 5, 6, 9, 10},
		{1, 2, 3, 4, 7, 8, 9, 10},
		{1, 2, 5, 6, 7, 8, 9, 10},
		{3, 4, 5, 6, 7, 8, 9, 10}}, toSlices[int](Removes(2, 10, list)))
}

func TestShrinkList(t *testing.T) {
	list := linked.New(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	shrinks := ShrinkList(list, func(t int) iterable.Iterable[int] {
		return iterable.Singleton(t / 2)
	})

	require.Equal(t, [][]int{
		// remove all 10
		{},
		// remove 5
		{1, 2, 3, 4, 5},
		{6, 7, 8, 9, 10},
		// remove 2
		{1, 2, 3, 4, 5, 6, 7, 8},
		{1, 2, 3, 4, 5, 6, 9, 10},
		{1, 2, 3, 4, 7, 8, 9, 10},
		{1, 2, 5, 6, 7, 8, 9, 10},
		{3, 4, 5, 6, 7, 8, 9, 10},
		// remove 1
		{1, 2, 3, 4, 5, 6, 7, 8, 9},
		{1, 2, 3, 4, 5, 6, 7, 8, 10},
		{1, 2, 3, 4, 5, 6, 7, 9, 10},
		{1, 2, 3, 4, 5, 6, 8, 9, 10},
		{1, 2, 3, 4, 5, 7, 8, 9, 10},
		{1, 2, 3, 4, 6, 7, 8, 9, 10},
		{1, 2, 3, 5, 6, 7, 8, 9, 10},
		{1, 2, 4, 5, 6, 7, 8, 9, 10},
		{1, 3, 4, 5, 6, 7, 8, 9, 10},
		{2, 3, 4, 5, 6, 7, 8, 9, 10},
		// shrink one element
		{0, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		{1, 1, 3, 4, 5, 6, 7, 8, 9, 10},
		{1, 2, 1, 4, 5, 6, 7, 8, 9, 10},
		{1, 2, 3, 2, 5, 6, 7, 8, 9, 10},
		{1, 2, 3, 4, 2, 6, 7, 8, 9, 10},
		{1, 2, 3, 4, 5, 3, 7, 8, 9, 10},
		{1, 2, 3, 4, 5, 6, 3, 8, 9, 10},
		{1, 2, 3, 4, 5, 6, 7, 4, 9, 10},
		{1, 2, 3, 4, 5, 6, 7, 8, 4, 10},
		{1, 2, 3, 4, 5, 6, 7, 8, 9, 5}}, toSlices[int](shrinks))
}
