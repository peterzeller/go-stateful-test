package generator_test

import (
	"github.com/peterzeller/go-fun/hash"
	"github.com/peterzeller/go-stateful-test/generator"
	"github.com/peterzeller/go-stateful-test/generator/geniterable"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSetInt(t *testing.T) {
	s := generator.Set(generator.IntRange(0, 2), hash.Num[int]())
	require.Equal(t, "[[], [1], [0], [0, 1], ...]", geniterable.String(s.Enumerate(2)))
	require.Equal(t, "[[], [2], [1], [1, 2], [0], [0, 2], [0, 1], [0, 1, 2]]", geniterable.String(s.Enumerate(3)))
}
