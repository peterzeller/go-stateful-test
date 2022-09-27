package shrink

import (
	"github.com/peterzeller/go-fun/iterable"
	"github.com/peterzeller/go-fun/list/linked"
)

func ShrinkListTail[T any](list *linked.List[T], shrinkFun func(t T) iterable.Iterable[T]) iterable.Iterable[*linked.List[T]] {
	if list == nil {
		return iterable.Empty[*linked.List[T]]()
	}
	tailShrinks := ShrinkList(list.Tail(), shrinkFun)
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

func ShrinkList[T any](list *linked.List[T], shrinkFun func(t T) iterable.Iterable[T]) iterable.Iterable[*linked.List[T]] {
	listLen := list.Length()
	toRemoveLengths := iterable.TakeWhile(
		func(x int) bool {
			return x > 0
		}, iterable.Generate(
			listLen,
			func(x int) int {
				return x / 2
			}))

	var partsRemoved iterable.Iterable[*linked.List[T]] = iterable.FlatMap(toRemoveLengths,
		func(k int) iterable.Iterable[*linked.List[T]] {
			return Removes(k, listLen, list)
		})

	shrinkOnes := ListShrinkOne(list, shrinkFun)

	return iterable.Concat(
		partsRemoved,
		shrinkOnes)
}

func Removes[T any](k int, n int, list *linked.List[T]) iterable.Iterable[*linked.List[T]] {
	if k > n {
		return iterable.Empty[*linked.List[T]]()
	}
	xs1 := list.Limit(k)
	xs2 := list.Skip(k)
	if xs2 == nil {
		return iterable.Singleton(linked.New[T]())
	}
	return iterable.Concat(
		iterable.Map(Removes(k, n-k, xs2),
			func(xs2r *linked.List[T]) *linked.List[T] {
				return xs1.Append(xs2r)
			}),
		iterable.Singleton(xs2),
	)
}

func ListShrinkOne[T any](list *linked.List[T], shrinkFun func(t T) iterable.Iterable[T]) iterable.Iterable[*linked.List[T]] {
	if list == nil {
		return iterable.Empty[*linked.List[T]]()
	}

	headShrinks := iterable.Map(shrinkFun(list.Head()),
		func(h T) *linked.List[T] {
			return linked.New(h).Append(list.Tail())
		})

	tailShrinks := iterable.Map(ListShrinkOne(list.Tail(), shrinkFun),
		func(t *linked.List[T]) *linked.List[T] {
			return linked.New(list.Head()).Append(t)
		})

	return iterable.Concat(
		headShrinks,
		tailShrinks)
}
