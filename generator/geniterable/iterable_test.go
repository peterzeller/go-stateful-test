package geniterable_test

import (
	"fmt"
	"github.com/peterzeller/go-stateful-test/generator/geniterable"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFun(t *testing.T) {
	r := geniterable.IterableFun[int](func() geniterable.Iterator[int] {
		current := 0
		return geniterable.Fun[int](func() geniterable.NextResult[int] {
			if current < 5 {
				current++
				return geniterable.ResultSome(current)
			}
			return geniterable.ResultNone[int](true)
		})
	})

	require.Equal(t, []int{1, 2, 3, 4, 5}, geniterable.ToSlice[int](r))
}

func TestLoop(t *testing.T) {
	s := []int{1, 42, 7}
	is := geniterable.FromSlice(s)

	var c []int
	for it := geniterable.Start(is); it.HasNext(); it.Next() {
		c = append(c, it.Current())
	}
	require.Equal(t, s, c)
}

func TestMap(t *testing.T) {
	s := geniterable.FromSlice([]int{1, 2, 3})
	s2 := geniterable.Map(s, func(x int) string { return fmt.Sprintf("x%d", x) })
	require.Equal(t, []string{"x1", "x2", "x3"}, geniterable.ToSlice(s2))
}

func TestMapIterator(t *testing.T) {
	s := geniterable.FromSlice([]int{1, 2, 3})
	it := geniterable.MapIterator(s.Iterator(), func(x int) string { return fmt.Sprintf("x%d", x) })
	a := it.Next()
	require.True(t, a.Present())
	require.Equal(t, "x1", a.Value())
	b := it.Next()
	require.True(t, b.Present())
	require.Equal(t, "x2", b.Value())
	c := it.Next()
	require.True(t, c.Present())
	require.Equal(t, "x3", c.Value())
	d := it.Next()
	require.False(t, d.Present())
}

func TestToString(t *testing.T) {
	s := geniterable.New(1, 2, 3)
	require.Equal(t, "[1, 2, 3]", geniterable.String(s))
}

func TestWhere(t *testing.T) {
	a := geniterable.FromSlice([]int{1, 2, 3, 4, 5, 6})
	isEven := func(x int) bool { return x%2 == 0 }
	require.Equal(t, []int{2, 4, 6}, geniterable.ToSlice(geniterable.Filter(a, isEven)))

}

func TestRange(t *testing.T) {
	require.Equal(t, []int{1, 2, 3}, geniterable.ToSlice(geniterable.Range(1, 4)))
}

func TestRangeI(t *testing.T) {
	require.Equal(t, []int{1, 2, 3, 4}, geniterable.ToSlice(geniterable.RangeI(1, 4)))
}

func TestRangeStep(t *testing.T) {
	require.Equal(t, []int{1, 4, 7, 10}, geniterable.ToSlice(geniterable.RangeStep(1, 13, 3)))
}

func TestRangeIStep(t *testing.T) {
	require.Equal(t, []int{1, 4, 7, 10, 13}, geniterable.ToSlice(geniterable.RangeIStep(1, 13, 3)))
}

func TestRangeStepRev(t *testing.T) {
	require.Equal(t, []int{13, 10, 7, 4}, geniterable.ToSlice(geniterable.RangeStep(13, 1, -3)))
}

func TestRangeIStepRev(t *testing.T) {
	require.Equal(t, []int{13, 10, 7, 4, 1}, geniterable.ToSlice(geniterable.RangeIStep(13, 1, -3)))
}

func ExampleEmpty() {
	it := geniterable.Empty[int]()
	fmt.Printf("it = %s", geniterable.String(it))
	// output: it = []
}

func ExampleSingleton() {
	it := geniterable.Singleton(42)
	fmt.Printf("it = %s", geniterable.String(it))
	// output: it = [42]
}

func ExampleForeach() {
	it := geniterable.New(1, 2, 3, 4, 5)
	geniterable.Foreach(it, func(i int) {
		fmt.Printf("%d, ", i)
	})
	// output: 1, 2, 3, 4, 5,
}
