package quickcheck

import (
	"context"
	"github.com/peterzeller/go-fun/iterable"
	"github.com/peterzeller/go-fun/linked"
	"github.com/peterzeller/go-stateful-test/quickcheck/tree"
)

func shrink(ctx context.Context, s *state, runState func(*state) (result *state)) *state {
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
	listShrinks := shrinkListTail(t.GeneratedValues(), shrinkGeneratedValues)
	return iterable.Map(listShrinks,
		func(l *linked.List[*linked.List[tree.GeneratedValue]]) *tree.GenNode {
			return tree.New(l)
		})
}

func shrinkListTail[T any](list *linked.List[T], shrinkFun func(t T) iterable.Iterable[T]) iterable.Iterable[*linked.List[T]] {
	if list == nil {
		return iterable.Empty[*linked.List[T]]()
	}
	tailShrinks := shrinkList(list.Tail(), shrinkFun)
	headShrinks := shrinkFun(list.Head())

	return iterable.Concat(
		iterable.Map(tailShrinks,
			func(t *linked.List[T]) *linked.List[T] {
				return linked.Cons(list.Head(), t)
			}),
		iterable.Map(headShrinks,
			func(h T) *linked.List[T] {
				return linked.Cons(h, list.Tail())
			}),
	)
}

func shrinkGeneratedValues(values *linked.List[tree.GeneratedValue]) iterable.Iterable[*linked.List[tree.GeneratedValue]] {
	return listShrinkOne(values, func(gv tree.GeneratedValue) iterable.Iterable[tree.GeneratedValue] {
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

func shrinkList[T any](list *linked.List[T], shrinkFun func(t T) iterable.Iterable[T]) iterable.Iterable[*linked.List[T]] {
	listLen := list.Length()
	toRemoveLengths := iterable.TakeWhile(
		func(x int) bool {
			return x > 0
		}, iterable.Generate[int](
			listLen,
			func(x int) int {
				return x / 2
			}))

	var partsRemoved iterable.Iterable[*linked.List[T]] = iterable.FlatMap(toRemoveLengths,
		func(k int) iterable.Iterable[*linked.List[T]] {
			return removes(k, listLen, list)
		})

	shrinkOnes := listShrinkOne(list, shrinkFun)

	return iterable.Concat[*linked.List[T]](
		partsRemoved,
		shrinkOnes)
}

func removes[T any](k int, n int, list *linked.List[T]) iterable.Iterable[*linked.List[T]] {
	if k > n {
		return iterable.Empty[*linked.List[T]]()
	}
	xs1 := list.Limit(k)
	xs2 := list.Skip(k)
	if xs2 == nil {
		return iterable.Singleton(linked.New[T]())
	}
	return iterable.Concat(
		iterable.Map(removes(k, n-k, xs2),
			func(xs2r *linked.List[T]) *linked.List[T] {
				return xs1.Append(xs2r)
			}),
		iterable.Singleton(xs2),
	)
}

func listShrinkOne[T any](list *linked.List[T], shrinkFun func(t T) iterable.Iterable[T]) iterable.Iterable[*linked.List[T]] {
	if list == nil {
		return iterable.Empty[*linked.List[T]]()
	}

	headShrinks := iterable.Map(shrinkFun(list.Head()),
		func(h T) *linked.List[T] {
			return linked.New(h).Append(list.Tail())
		})

	tailShrinks := iterable.Map(listShrinkOne(list.Tail(), shrinkFun),
		func(t *linked.List[T]) *linked.List[T] {
			return linked.New(list.Head()).Append(t)
		})

	return iterable.Concat(
		headShrinks,
		tailShrinks)
}
