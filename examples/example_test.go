package examples

import (
	"testing"

	"github.com/peterzeller/go-stateful-test/generator"
	"github.com/peterzeller/go-stateful-test/pick"
	"github.com/peterzeller/go-stateful-test/quickcheck"
	"github.com/peterzeller/go-stateful-test/statefulTest"
	"github.com/stretchr/testify/require"
)

func TestInts(t *testing.T) {
	expectError(t, func(t quickcheck.TestingT) {
		quickcheck.Run(t, quickcheck.Config{}, func(t statefulTest.T) {
			x := pick.Val(t, generator.Int())
			y := pick.Val(t, generator.Int())
			t.Logf("x = %d, y = %d", x, y)
			require.True(t, x+y < 10)
		})
	})
}

func TestSequence(t *testing.T) {
	expectError(t, func(t quickcheck.TestingT) {
		quickcheck.Run(t, quickcheck.Config{}, func(t statefulTest.T) {
			var s []int
			for t.HasMore() {
				x := pick.Val(t, generator.IntRange(0, 9))
				s = append(s, x)
			}
			//log.Printf("s = %#v", s)
			t.Logf("s = %#v", s)
			require.True(t, len(s) < 5)
		})
	})
}
