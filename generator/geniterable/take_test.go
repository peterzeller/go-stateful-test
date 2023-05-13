package geniterable_test

import (
	"fmt"
	"github.com/peterzeller/go-stateful-test/generator/geniterable"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTake(t *testing.T) {
	it := geniterable.Take(2, geniterable.FromSlice([]int{1, 2, 3, 4, 5}))
	require.Equal(t, []int{1, 2}, geniterable.ToSlice(it))
}

func ExampleTakeWhile_empty() {
	it := geniterable.TakeWhile(func(x int) bool {
		return x <= 3
	}, geniterable.FromSlice([]int{5}))
	fmt.Printf("it = %s\n", geniterable.String(it))
	// output: it = []
}

func ExampleTakeWhile() {
	it := geniterable.TakeWhile(func(x int) bool {
		return x <= 3
	}, geniterable.FromSlice([]int{1, 2, 3, 4, 3, 2, 1}))
	fmt.Printf("it = %s\n", geniterable.String(it))
	// output: it = [1, 2, 3]
}

func TestTakeWhile(t *testing.T) {
	tw := geniterable.TakeWhile(func(x int) bool {
		return x <= 3
	}, geniterable.New(1, 2, 3, 4, 3, 2, 1))
	it := tw.Iterator()
	v := it.Next()
	require.Equal(t, 1, v.Value())
	require.True(t, v.Present())
	v = it.Next()
	require.Equal(t, 2, v.Value())
	require.True(t, v.Present())
	v = it.Next()
	require.Equal(t, 3, v.Value())
	require.True(t, v.Present())
	v = it.Next()
	require.Equal(t, 0, v.Value())
	require.False(t, v.Present())
	v = it.Next()
	require.Equal(t, 0, v.Value())
	require.False(t, v.Present())
}
