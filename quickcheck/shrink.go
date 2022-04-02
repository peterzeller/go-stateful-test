package quickcheck

import (
	"context"
	"log"

	"github.com/peterzeller/go-fun/iterable"
	"github.com/peterzeller/go-fun/linkedlist"
)

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
		iterators = append(iterators, listShrinks)
	}
	// 2. shrink individual list elements
	// 3. shrink individual children

	return iterable.ConcatIterators(iterators...)
}

func shrinkList[T any](list *linkedlist.LinkedList[T], shrinkFun func(t T) iterable.Iterable[T]) iterable.Iterable[*linkedlist.LinkedList[T]] {
	listLen := list.Length()
	toRemoveLengths := iterable.TakeWhile(
		func(x int) bool {
			return x > 0
		}, iterable.Generate[int](
			listLen,
			func(x int) int {
				return x / 2
			}))

	log.Printf("toRemoveLengths = %v", iterable.ToSlice(toRemoveLengths))

	var partsRemoved iterable.Iterable[*linkedlist.LinkedList[T]] = iterable.FlatMap(
		func(k int) iterable.Iterable[*linkedlist.LinkedList[T]] {
			return removes(k, listLen, list)
		})(toRemoveLengths)

	shrinkOnes := listShrinkOne(list, shrinkFun)

	return iterable.Concat[*linkedlist.LinkedList[T]](
		partsRemoved,
		shrinkOnes)
}

func removes[T any](k int, n int, list *linkedlist.LinkedList[T]) iterable.Iterable[*linkedlist.LinkedList[T]] {
	if k > n {
		return iterable.Empty[*linkedlist.LinkedList[T]]()
	}
	xs1 := list.Limit(k)
	xs2 := list.Skip(k)
	if xs2 == nil {
		return iterable.Singleton(linkedlist.New[T]())
	}
	return iterable.Concat(
		iterable.Map(func(xs2r *linkedlist.LinkedList[T]) *linkedlist.LinkedList[T] {
			return xs1.Append(xs2r)
		})(removes(k, n-k, xs2)),
		iterable.Singleton(xs2),
	)
}

func listShrinkOne[T any](list *linkedlist.LinkedList[T], shrinkFun func(t T) iterable.Iterable[T]) iterable.Iterable[*linkedlist.LinkedList[T]] {
	if list == nil {
		return iterable.Empty[*linkedlist.LinkedList[T]]()
	}

	headShrinks := iterable.Map(func(h T) *linkedlist.LinkedList[T] {
		return linkedlist.New(h).Append(list.Tail())
	})(shrinkFun(list.Head()))

	tailShrinks := iterable.Map(func(t *linkedlist.LinkedList[T]) *linkedlist.LinkedList[T] {
		return linkedlist.New(list.Head()).Append(t)
	})(listShrinkOne(list.Tail(), shrinkFun))

	return iterable.Concat(
		headShrinks,
		tailShrinks)
}
