package fuzzcheck

import (
	"github.com/peterzeller/go-stateful-test/quickcheck"
	"github.com/peterzeller/go-stateful-test/quickcheck/randomsource"
	"github.com/peterzeller/go-stateful-test/statefulTest"
	"testing"
	"time"
)

type Config struct {
	MaxShrinkDuration time.Duration
	PrintAllLogs      bool
	DisableHeuristics bool
}

func Run(t TestingT, cfg Config, f func(t statefulTest.T)) {
	t.Fuzz(func(t *testing.T, generatorString []byte) {
		quickcheck.Run(t, quickcheck.Config{
			// only one run, because the fuzzing framework controls the runs
			NumberOfRuns:      1,
			MaxShrinkDuration: cfg.MaxShrinkDuration,
			PrintAllLogs:      cfg.PrintAllLogs,
			FixedRandomSource: randomsource.RandomSourceFromBytes(generatorString),
			DisableHeuristics: cfg.DisableHeuristics,
		}, f)
	})
}

type TestingT interface {
	quickcheck.TestingT
	Fuzz(f any)
}
