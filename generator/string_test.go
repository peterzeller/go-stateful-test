package generator

import (
	"testing"

	"github.com/peterzeller/go-fun/iterable"
	"github.com/stretchr/testify/require"
)

func TestGenString_Size(t *testing.T) {
	s := String('a', 'b', 'c')
	require.Equal(t, int64(3), s.Size(R("abc")).Int64())
}

func TestGenString_Shrink(t *testing.T) {
	s := String('a', 'b', 'c')
	require.Equal(t, []string{"", "cb", "ca", "ba", "bba", "caa"}, iterable.ToSlice(ShrinkValues(s, "cba")))
}

func TestGenString_Enumerate(t *testing.T) {
	s := String('a', 'b')
	require.Equal(t, []string{"", "a", "b", "aa", "ab", "ba", "bb", "aaa", "aab", "aba", "abb", "baa", "bab", "bba", "bbb"}, iterable.ToSlice(s.Enumerate(3)))
}
