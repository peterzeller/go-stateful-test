package examples

import (
	"fmt"
	"strings"
	"testing"

	"github.com/peterzeller/go-stateful-test/generator"
	"github.com/peterzeller/go-stateful-test/pick"
	"github.com/peterzeller/go-stateful-test/quickcheck"
	"github.com/peterzeller/go-stateful-test/statefulTest"
	"github.com/stretchr/testify/require"
)

type logT struct {
	log    strings.Builder
	failed bool
}

func (l *logT) Errorf(format string, args ...interface{}) {
	l.failed = true
	l.log.WriteString(fmt.Sprintf(format, args...))
	l.log.WriteString("\n")
}

func (l *logT) FailNow() {
	l.failed = true
	panic(fmt.Errorf("test failed"))
}

func (l *logT) Failed() bool {
	return l.failed
}

func (l *logT) Logf(format string, args ...interface{}) {
	l.log.WriteString(fmt.Sprintf(format, args...))
	l.log.WriteString("\n")
}

func expectError(t *testing.T, f func(testingT quickcheck.TestingT)) {
	testingT := &logT{}

	defer func() {
		r := recover()
		if r != nil {
			if err, ok := r.(error); ok {
				t.Logf("error as expected: %v\n%s", err, testingT.log.String())
				return
			}
			panic(r)
		}
		if !testingT.Failed() {
			t.Logf("Expected quickcheck to find an error, but no error found")
			t.Fail()
		}
	}()
	f(testingT)
}

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
