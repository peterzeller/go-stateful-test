package generator_test

import (
	"fmt"
	"github.com/peterzeller/go-stateful-test/generator/geniterable"
	"testing"

	"github.com/peterzeller/go-stateful-test/generator"
	"github.com/stretchr/testify/require"
)

type Pair struct {
	A string
	B int
}

func ExampleReflectionGen() {
	// This example creates a generator for the following "Pair" struct type:
	// type Pair struct {
	// 	A string
	// 	B int
	// }
	genOpts := generator.ReflectionGenDefaultOpts()
	g := generator.ReflectionGen[Pair](genOpts)
	for it := geniterable.Start(generator.EnumerateValues(g, 3)); it.HasNext(); it.Next() {
		fmt.Printf("%+v\n", it.Current())
	}
	// Output: {A: B:0}
	// {A: B:1}
	// {A: B:-1}
	// {A:a B:0}
	// {A:a B:1}
	// {A:a B:-1}
	// {A:b B:0}
	// {A:b B:1}
	// {A:b B:-1}
	// {A:aa B:0}
	// {A:aa B:1}
	// {A:aa B:-1}
	// {A:ab B:0}
	// {A:ab B:1}
	// {A:ab B:-1}
	// {A:ba B:0}
	// {A:ba B:1}
	// {A:ba B:-1}
	// {A:bb B:0}
	// {A:bb B:1}
	// {A:bb B:-1}
	// {A:aaa B:0}
	// {A:aaa B:1}
	// {A:aaa B:-1}
	// {A:aab B:0}
	// {A:aab B:1}
	// {A:aab B:-1}
	// {A:aba B:0}
	// {A:aba B:1}
	// {A:aba B:-1}
	// {A:abb B:0}
	// {A:abb B:1}
	// {A:abb B:-1}
	// {A:baa B:0}
	// {A:baa B:1}
	// {A:baa B:-1}
	// {A:bab B:0}
	// {A:bab B:1}
	// {A:bab B:-1}
	// {A:bba B:0}
	// {A:bba B:1}
	// {A:bba B:-1}
	// {A:bbb B:0}
	// {A:bbb B:1}
	// {A:bbb B:-1}
}

type GenericPair[A, B any] struct {
	A A
	B B
}

func TestReflectionGenGenerics(t *testing.T) {
	genOpts := generator.ReflectionGenDefaultOpts()
	g := generator.ReflectionGen[GenericPair[string, int]](genOpts)
	values := geniterable.ToSlice(generator.EnumerateValues(g, 2))
	require.Equal(t, []GenericPair[string, int]{{A: "", B: 0}, {A: "", B: 1}, {A: "a", B: 0}, {A: "a", B: 1}, {A: "b", B: 0}, {A: "b", B: 1}, {A: "aa", B: 0}, {A: "aa", B: 1}, {A: "ab", B: 0}, {A: "ab", B: 1}, {A: "ba", B: 0}, {A: "ba", B: 1}, {A: "bb", B: 0}, {A: "bb", B: 1}}, values)
}

func TestReflectionGenSlice(t *testing.T) {
	genOpts := generator.ReflectionGenDefaultOpts()
	g := generator.ReflectionGen[[]int](genOpts)
	values := geniterable.ToSlice(generator.EnumerateValues(g, 2))
	require.Equal(t, [][]int([][]int{{}, {0}, {1}, {0, 0}, {1, 0}, {0, 1}, {1, 1}}), values)
}

func ExampleReflectionGeneratorOptions_RegisterConstructor() {
	genOpts := generator.ReflectionGenDefaultOpts()
	// generate a custom constructor for strings
	genOpts.RegisterConstructor(func(x, y int) string {
		return fmt.Sprintf("(%d, %d)", x, y)
	})
	g := generator.ReflectionGen[string](genOpts)
	for it := geniterable.Start(generator.EnumerateValues(g, 3)); it.HasNext(); it.Next() {
		fmt.Printf("%s\n", it.Current())
	}
	// Output: (0, 0)
	// (0, 1)
	// (0, -1)
	// (1, 0)
	// (1, 1)
	// (1, -1)
	// (-1, 0)
	// (-1, 1)
	// (-1, -1)
}
