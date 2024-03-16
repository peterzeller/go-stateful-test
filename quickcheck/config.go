package quickcheck

import (
	"github.com/peterzeller/go-stateful-test/quickcheck/randomsource"
	"testing"
	"time"
)

type Config struct {
	NumberOfRuns      int
	MaxShrinkDuration time.Duration
	PrintAllLogs      bool
	FixedRandomSource randomsource.RandomSource
	DisableHeuristics bool
}

type TestingT interface {
	Errorf(format string, args ...interface{})
	FailNow()
	Failed() bool
	Logf(format string, args ...interface{})
}

var _ TestingT = &testing.T{}
