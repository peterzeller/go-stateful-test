package quickcheck

import (
	"testing"
	"time"
)

type Config struct {
	NumberOfRuns      int
	MaxShrinkDuration time.Duration
	PrintAllLogs      bool
}

type TestingT interface {
	Errorf(format string, args ...interface{})
	FailNow()
	Failed() bool
	Logf(format string, args ...interface{})
}

var _ TestingT = &testing.T{}
