package smallcheck

import (
	"errors"
	"fmt"
	"github.com/peterzeller/go-stateful-test/statefulTest"
	"runtime/debug"
)

func Run(t TestingT, cfg Config, f func(t statefulTest.T)) {
	cfg = setDefaults(cfg)

	runState := func(s *state) {
		defer func() {
			// handle panics
			r := recover()
			if r != nil {
				if err, ok := r.(error); ok {
					if errors.Is(err, testFailedErr) {
						s.failed = true
						return
					}
				}
				if _, ok := r.(emptyIterator); ok {
					// ignore error and just continue with next iteration
					s.failed = false
					return
				}

				stackTrace := debug.Stack()
				s.Errorf("Panic in test:\n%v\n%s", r, stackTrace)
				s.failed = true
			}
		}()

		f(s)
	}

	for depth := 1; depth < cfg.Depth; depth++ {
		rs := &rState{
			stack:           nil,
			continueAtDepth: 0,
			maxDepth:        depth,
			done:            false,
			cfg:             cfg,
		}

		s := rs.exploreStates(runState)
		if s != nil && s.failed {
			t.Errorf("Test failed at depth %d:\n%s", depth, s.GetLog())
			return
		}
	}
}

var testFailedErr = fmt.Errorf("test failed")
