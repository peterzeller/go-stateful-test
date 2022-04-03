package quickcheck

import (
	"context"
	"github.com/peterzeller/go-fun/iterable"
	"github.com/peterzeller/go-fun/linked"
	"github.com/peterzeller/go-stateful-test/generator/shrink"
	"github.com/peterzeller/go-stateful-test/quickcheck/tree"
)

func shrinkState(ctx context.Context, s *state, runState func(*state) (result *state)) *state {
	for ctx.Err() == nil {
		s2 := shrinkOne(ctx, s, runState)
		if s2 == nil || s2 == s {
			// no further shrink possible -> return last state
			return s
		}
		// continue loop with smaller state and try again:
		//log.Printf("found smaller shrink:\n%s", s2.mainFork)
		s = s2
	}
	return s
}

func shrinkOne(ctx context.Context, s *state, runState func(*state) (result *state)) *state {
	gTree := s.mainFork.genTree.ToImmutable()
	iterator := shrinkTree(gTree).Iterator()
	for ctx.Err() == nil {
		currentShrink, ok := iterator.Next()
		if !ok {
			// could not find better shrink -> return original state
			return s
		}
		// Try to run with current shrink:
		s2 := initState(s.mainFork.genTree.Seed)
		s2.mainFork.presetTree = currentShrink
		res := runState(s2)
		if res != nil && res.failed {
			newSize := res.mainFork.genTree.Size()
			oldSize := s.mainFork.genTree.Size()
			if newSize.Cmp(oldSize) < 0 {
				// found a smaller execution that also fails
				return res
			}
		}
	}
	return s
}

func shrinkTree(t *tree.GenNode) iterable.Iterable[*tree.GenNode] {
	listShrinks := shrink.ShrinkListTail(t.GeneratedValues(), shrinkGeneratedValues)
	return iterable.Map(listShrinks,
		func(l *linked.List[*linked.List[tree.GeneratedValue]]) *tree.GenNode {
			return tree.New(l)
		})
}

func shrinkGeneratedValues(values *linked.List[tree.GeneratedValue]) iterable.Iterable[*linked.List[tree.GeneratedValue]] {
	return shrink.ListShrinkOne(values, func(gv tree.GeneratedValue) iterable.Iterable[tree.GeneratedValue] {
		shrinks := gv.Generator.Shrink(gv.Value)
		return iterable.Map(shrinks,
			func(v interface{}) tree.GeneratedValue {
				return tree.GeneratedValue{
					Generator: gv.Generator,
					Value:     v,
				}
			})
	})
}
