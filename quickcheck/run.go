package quickcheck

import (
	"context"
	"errors"
	"runtime/debug"
	"time"

	"github.com/peterzeller/go-stateful-test/statefulTest"
)

// Run the given testing function `f` multiple times.
// Each run will be executed with different values in generators.
// When an error is found, we try to shrink the test run before showing the final error.
// Only the error message in the shrunk execution and the logs from this run will be shown.
func Run(t TestingT, cfg Config, f func(t statefulTest.T)) {
	cfg = setDefaults(cfg)

	runState := func(s *state) (result *state) {
		defer func() {
			// handle panics
			r := recover()
			if r != nil {
				if err, ok := r.(error); ok {
					if errors.Is(err, testFailedErr) {
						result = s
						return
					}
				}

				stackTrace := debug.Stack()
				s.Errorf("Panic in test:\n%v\n%s", r, stackTrace)
				result = s
			}
		}()

		f(s)
		if s.Failed() {
			return s
		}
		return nil
	}
	var s *state = firstNotNil[state](cfg, func(iteration int) *state {
		s := initState(int64(iteration))
		return runState(s)
	})
	if s == nil || !s.Failed() {
		// all ok
		return
	}
	t.Logf("Found error, shrinking testcase ...")

	// start shrinking:
	ctx, cancel := context.WithTimeout(context.Background(), cfg.MaxShrinkDuration)
	defer cancel()
	shrunkS := shrinkState(ctx, s, runState)

	if shrunkS.failed {
		t.Logf("Shrunk Test Run:\n%s", shrunkS.GetLog())
		t.FailNow()
	} else {
		// print original error
		t.Logf("Could not reproduce error while shrinking (flaky test?)\n%s", shrunkS.GetLog())
	}
}

func setDefaults(cfg Config) Config {
	if cfg.NumberOfRuns == 0 {
		cfg.NumberOfRuns = 100
	}
	if cfg.MaxShrinkDuration == 0 {
		cfg.MaxShrinkDuration = 30 * time.Second
	}
	return cfg
}

func firstNotNil[X any](cfg Config, f func(iteration int) *X) *X {
	for i := 0; i < cfg.NumberOfRuns; i++ {
		res := f(i)
		if res != nil {
			return res
		}
	}
	return nil
}
