package smallcheck

import (
	"fmt"
	"github.com/peterzeller/go-stateful-test/generator"
	"github.com/peterzeller/go-stateful-test/pick"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestExplore(t *testing.T) {
	rs := rState{
		stack:           nil,
		continueAtDepth: 0,
		maxDepth:        3,
		done:            false,
	}

	explored := make(map[string]bool)

	rs.exploreStates(func(s *state) {
		x := pick.Val(s, generator.Int())
		y := pick.Val(s, generator.Int())
		z := pick.Val(s, generator.Int())
		key := fmt.Sprintf("[%d, %d, %d]", x, y, z)
		explored[key] = true
	})

	// test that we explored some samples:
	require.Contains(t, explored, "[0, 0, 0]")
	require.Contains(t, explored, "[-1, 0, 1]")
	require.Contains(t, explored, "[1, 1, 1]")
	require.Contains(t, explored, "[-1, -1, -1]")
}
