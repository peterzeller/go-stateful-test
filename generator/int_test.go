package generator

import (
	"testing"

	"github.com/peterzeller/go-fun/iterable"
	"github.com/stretchr/testify/require"
)

func TestGenInt_Enumerate(t *testing.T) {
	g := IntRange(1, 3)
	require.Equal(t, []int{1, 2, 3}, iterable.ToSlice(g.Enumerate(10)))
}
