package quickcheck

import (
	"testing"

	"github.com/peterzeller/go-fun/iterable"
	"github.com/peterzeller/go-fun/linkedlist"
)

func TestRemoves(t *testing.T) {
	list := linkedlist.New(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	shrinks := removes(5, 10, list)

	t.Logf("original = %+v", list)
	for it := iterable.Start(shrinks); it.HasNext(); it.Next() {
		t.Logf("list = %v", it.Current())
	}
}

func TestShrinkList(t *testing.T) {
	list := linkedlist.New(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	shrinks := shrinkList(list, func(t int) iterable.Iterable[int] {
		return iterable.Singleton(t / 2)
	})

	t.Logf("original = %+v", list)
	for it := iterable.Start(shrinks); it.HasNext(); it.Next() {
		t.Logf("list = %v", it.Current())
	}

}
