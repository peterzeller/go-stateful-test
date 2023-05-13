package geniterable_test

import (
	"github.com/peterzeller/go-stateful-test/generator/geniterable"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConcat(t *testing.T) {
	it := geniterable.Concat(
		geniterable.FromSlice([]int{1, 2, 3}),
		geniterable.FromSlice([]int{}),
		geniterable.FromSlice([]int{4, 5}),
	)

	require.Equal(t, []int{1, 2, 3, 4, 5}, geniterable.ToSlice(it))
}

func TestConcatIterators(t *testing.T) {
	it := geniterable.ConcatIterators(
		geniterable.FromSlice([]int{1, 2, 3}).Iterator(),
		geniterable.FromSlice([]int{}).Iterator(),
		geniterable.FromSlice([]int{4, 5}).Iterator(),
	)

	require.Equal(t, []int{1, 2, 3, 4, 5}, geniterable.IteratorToSlice(it))
}
