package generator

import (
	"github.com/peterzeller/go-stateful-test/generator/geniterable"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenInt_Enumerate(t *testing.T) {
	g := IntRange(1, 3)
	require.Equal(t, []int{1, 2, 3}, geniterable.ToSlice(g.Enumerate(10)))
}

func TestGenInt_EnumerateExhaustive(t *testing.T) {
	g := IntRange(1, 3)
	require.Equal(t, "[1, 2, 3]", geniterable.String(g.Enumerate(3)))
}

func TestGenInt_EnumerateExhaustive2(t *testing.T) {
	g := IntRange(1, 10)
	require.Equal(t, "[1, 2, 3, ...]", geniterable.String(g.Enumerate(3)))
}
