package examples

import (
	"github.com/peterzeller/go-stateful-test/generator"
	"github.com/peterzeller/go-stateful-test/pick"
	"github.com/peterzeller/go-stateful-test/quickcheck"
	"github.com/peterzeller/go-stateful-test/smallcheck"
	"github.com/peterzeller/go-stateful-test/statefulTest"
	"github.com/stretchr/testify/assert"
	"testing"
)

//The following function is supposed to compute the maximum out of 3 integers.
//Unfortunately it contains a bug.
//Can you find a counter example where it would fail?
func max3(x, y, z int) int {
	if x > y && x > z {
		return x
	} else if y > x && y > z {
		return y
	} else {
		return z
	}
}

func TestMax3(t *testing.T) {
	expectError(t, func(t quickcheck.TestingT) {
		smallcheck.Run(t, smallcheck.Config{}, func(t statefulTest.T) {
			x := pick.Val(t, generator.Int())
			y := pick.Val(t, generator.Int())
			z := pick.Val(t, generator.Int())
			res := max3(x, y, z)
			t.Logf("min3(%d, %d, %d) = %d", x, y, z, res)
			assert.True(t, res >= x, "res >= x")
			assert.True(t, res >= y, "res >= y")
			assert.True(t, res >= z, "res >= z")
		})
	})
}
