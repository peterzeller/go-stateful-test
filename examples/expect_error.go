package examples

import (
	"fmt"
	"github.com/peterzeller/go-stateful-test/quickcheck"
	"strings"
	"testing"
)

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
