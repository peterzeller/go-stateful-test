package geniterable_test

import (
	"fmt"
	"github.com/peterzeller/go-stateful-test/generator/geniterable"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	parts := geniterable.TakeWhile(
		func(x int) bool {
			return x > 0
		}, geniterable.Generate[int](
			10,
			func(x int) int {
				return x / 2
			}))

	require.Equal(t, []int{10, 5, 2, 1}, geniterable.ToSlice(parts))
}

func ExampleGenerate() {
	it := geniterable.Generate[int](
		10,
		func(x int) int {
			return x / 2
		})

	fmt.Printf("it = %s", geniterable.String(geniterable.Take(5, it)))
	// output: it = [10, 5, 2, 1, 0]
}

func ExampleGenerateState() {
	it := geniterable.GenerateState[int](
		0,
		func(state int) (newState int, r geniterable.NextResult[int]) {
			if state <= 10 {
				r = geniterable.ResultSome(10 * state)
				newState = state + 2
			} else {
				r = geniterable.ResultNone[int](true)
			}
			return
		})

	fmt.Printf("it = %s", geniterable.String(it))
	// output: it = [0, 20, 40, 60, 80, 100]
}
