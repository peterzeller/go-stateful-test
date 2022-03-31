package quickcheck

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/peterzeller/go-fun/iterable"
	"github.com/peterzeller/go-fun/zero"
	"github.com/peterzeller/go-stateful-test/statefulTest"
)

// Run the given testing function `f` multiple times.
// Each run will be executed with different values in generators.
// When an error is found, we try to shrink the test run before showing the final error.
// Only the error message in the shrunk execution and the logs from this run will be shown.
func Run(t TestingT, cfg Config, f func(t statefulTest.T)) {
	runState := func(s *state) (result *state) {
		defer func() {
			// handle panics
			r := recover()
			if r != nil {
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
	if !s.Failed() {
		// all ok
		return
	}
	t.Logf("Found error, shrinking testcase ...")

	// start shrinking:
	ctx, cancel := context.WithTimeout(context.Background(), cfg.MaxShrinkDuration)
	defer cancel()
	shrunkS := shrink(ctx, s, runState)

	if shrunkS.failed {
		t.FailNow()
	} else {
		// print original error
		t.Logf("Could not reproduce error while shrinking (flaky test?)")
		fmt.Printf("%s", s.GetLog())
	}
}

func shrink(ctx context.Context, s *state, runState func(*state) (result *state)) *state {
	for ctx.Err() == nil {
		s2 := shrinkOne(ctx, s, runState)
		if s2 == nil || s2 == s {
			// no further shrink possible -> return last state
			return s
		}
		// continue loop with smaller state and try again:
		s = s2
	}
	return s
}

func shrinkOne(ctx context.Context, s *state, runState func(*state) (result *state)) *state {
	tree := s.genTree.toImmutable()
	iterator := shrinkIterator(tree)
	for ctx.Err() == nil {
		currentShrink, ok := iterator.Next()
		if !ok {
			// could not find better shrink -> return original state
			return s
		}
		// Try to run with current shrink:
		s2 := initState(s.genTree.seed)
		s2.presetTree = currentShrink
		res := runState(s2)
		if res.failed && res.size < s.size {
			// found a smaller execution that also fails
			return res
		}
	}
	return s
}

func shrinkIterator(tree *genNode) iterable.Iterator[*genNode] {
	var iterators []iterable.Iterator[*genNode]
	// 1. shrink size of the list
	if tree.hasMoreCount > 0 {
		// shrink list by halfing
		listShrinks := iterable.Fun[*genNode](func() (*genNode, bool) {
			return nil, false
		})
	}
	// 2. shrink individual list elements
	// 3. shrink individual children

	return appendIterators(iterators...)
}

func shrinkList() {

}

//shrinkList :: (a -> [a]) -> [a] -> [[a]]
//shrinkList shr xs = concat [ removes k n xs | k <- takeWhile (>0) (iterate (`div`2) n) ]
//++ shrinkOne xs
//where
//n = length xs
//
//shrinkOne []     = []
//shrinkOne (x:xs) = [ x':xs | x'  <- shr x ]
//++ [ x:xs' | xs' <- shrinkOne xs ]
//
//removes k n xs
//| k > n     = []
//| null xs2  = [[]]
//| otherwise = xs2 : map (xs1 ++) (removes k (n-k) xs2)
//where
//xs1 = take k xs
//xs2 = drop k xs

func appendIterators[T any](iterators ...iterable.Iterator[T]) iterable.Iterator[T] {
	pos := 0
	return iterable.Fun[T](func() (T, bool) {
		for pos < len(iterators) {
			n, ok := iterators[pos].Next()
			if ok {
				return n, true
			}
			pos++
		}
		return zero.Value[T](), false
	})
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
