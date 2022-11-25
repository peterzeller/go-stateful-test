package smallcheck

import (
	"fmt"
	"github.com/peterzeller/go-fun/iterable"
	"github.com/peterzeller/go-stateful-test/generator"
	"strings"
)

// rState is the state over several runs
type rState struct {
	stack           []*stackEntry
	continueAtDepth int
	maxDepth        int
	done            bool
	cfg             Config
}

func (rs *rState) exploreStates(runState func(s *state)) *state {
	for !rs.done {
		s := &state{
			parent: rs,
			log:    strings.Builder{},
			failed: false,
		}

		func() {
			defer func() {
				r := recover()
				if r != nil {
					if _, ok := r.(emptyIterator); ok {
						// ignore empty iterator error
						s.failed = false
						return
					}
					// propagate other errors
					panic(r)
				}
			}()
			defer s.runCleanups()

			runState(s)
		}()
		if rs.cfg.PrintAllLogs {
			fmt.Printf("\n%s---\n", s.GetLog())
		}
		if s.failed {
			// found a failed testcase
			return s
		}
		rs.advanceStack(s.depth - 1)
	}
	return nil
}

func (rs *rState) advanceStack(depth int) {
	if depth < 0 {
		rs.done = true
		return
	}
	entry := rs.stack[depth]
	newCurrent, ok := entry.iterator.Next()
	if ok {
		// if we have a next element, continue at this level
		entry.current = newCurrent
		rs.continueAtDepth = depth
		return
	}
	// otherwise, we need to advance the stack one position below
	rs.stack[depth] = nil
	rs.advanceStack(depth - 1)
}

type stackEntry struct {
	current  interface{}
	iterator iterable.Iterator[interface{}]
}

// state for a single iteration
type state struct {
	parent       *rState
	log          strings.Builder
	failed       bool
	depth        int
	hasMoreCalls int
	cleanup      []func()
}

func (s *state) Cleanup(f func()) {
	s.cleanup = append(s.cleanup, f)
}

func (s *state) Errorf(format string, args ...interface{}) {
	s.failed = true
	_, _ = fmt.Fprintf(&s.log, format, args...)
	s.log.WriteRune('\n')
}

func (s *state) FailNow() {
	s.failed = true
	panic(testFailedErr)
}

func (s *state) Failed() bool {
	return s.failed
}

func (s *state) Logf(format string, args ...any) {
	if s.parent.cfg.PrintAllLogs {
		fmt.Printf(format, args...)
		fmt.Printf("\n")
		return
	}
	// TODO implement like in real Log and add source code line to message?
	_, _ = fmt.Fprintf(&s.log, format, args...)
	s.log.WriteRune('\n')
}

func (s *state) GetLog() string {
	return s.log.String()
}

func (s *state) PickValue(gen generator.UntypedGenerator) interface{} {
	rs := s.parent
	if s.depth < len(rs.stack) && s.depth <= rs.continueAtDepth {
		// We already have an iterator.
		// Return the current value and move to the next.
		entry := rs.stack[s.depth]
		s.depth++
		return entry.current
	}
	// We don't have an iterator for this stack level yet
	it := gen.Enumerate(rs.maxDepth).Iterator()
	current, ok := it.Next()
	if !ok {
		panic(emptyIterator{depth: s.depth})
	}
	newEntry := &stackEntry{
		iterator: it,
		current:  current,
	}
	if s.depth < len(rs.stack) {
		rs.stack[s.depth] = newEntry
	} else {
		rs.stack = append(rs.stack, newEntry)
	}

	s.depth++
	return current
}

type emptyIterator struct {
	depth int
}

func (s *state) HasMore() bool {
	s.hasMoreCalls++
	return s.hasMoreCalls < s.parent.maxDepth
}

func (s *state) runCleanups() {
	for _, f := range s.cleanup {
		f()
	}
	s.cleanup = nil
}
