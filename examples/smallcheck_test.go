package examples

import (
	"github.com/peterzeller/go-stateful-test/generator"
	"github.com/peterzeller/go-stateful-test/pick"
	"github.com/peterzeller/go-stateful-test/smallcheck"
	"github.com/peterzeller/go-stateful-test/statefulTest"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestSmallCheckNonDet tests that we can have nondeterministic tests without returning wrong values.
func TestSmallCheckNonDet(t *testing.T) {
	// use a variable outside the test to simulate nondeterministic execution
	run := 0

	smallcheck.Run(t, smallcheck.Config{
		Depth:         100,
		PrintAllLogs:  true,
		PrintLiveLogs: true,
	}, func(t statefulTest.T) {
		run++
		var x int
		if run < 5 {
			x = pick.Val(t, generator.IntRange(0, 10))
		} else {
			x = pick.Val(t, generator.IntRange(11, 15))
			// require.Less(t, 10, x)
		}
		t.Logf("run %d: x = %d", run, x)
	})
}

// TestSmallExhaustiveStop tests that we stop early once all values are covered.
func TestSmallExhaustiveStop(t *testing.T) {
	generatedValues := make([]int, 0)

	cfg := smallcheck.Config{
		Depth: 2000,
	}
	smallcheck.Run(t, cfg, func(t statefulTest.T) {
		x := pick.Val(t, generator.IntRange(0, 3))
		y := pick.Val(t, generator.IntRange(10, 13))
		generatedValues = append(generatedValues, x, y)
	})

	//t.Logf("%#v", generatedValues)

	lastValueFound := 0
	for i := 0; i < len(generatedValues)-1; i++ {
		if generatedValues[i] == 3 && generatedValues[i+1] == 13 {
			lastValueFound++
			//t.Logf("found last value at position %d", i)
		}
	}
	require.Equal(t, 1, lastValueFound, "the last value should only be found once, due to the exhaustiveness check")
}
