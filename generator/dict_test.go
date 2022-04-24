package generator

import (
	"fmt"

	"github.com/peterzeller/go-fun/hash"
	"github.com/peterzeller/go-fun/iterable"
)

func ExampleDict() {
	g := Dict(IntRange(1, 2), IntRange(3, 4), hash.Num[int]())
	for it := iterable.Start(g.Enumerate(3)); it.HasNext(); it.Next() {
		fmt.Println(it.Current())
	}
	// Output: []
	// [1 -> 3]
	// [1 -> 4]
	// [2 -> 3]
	// [2 -> 4]
	// [1 -> 3, 2 -> 3]
	// [1 -> 4, 2 -> 3]
	// [1 -> 3, 2 -> 4]
	// [1 -> 4, 2 -> 4]
	// [1 -> 3, 2 -> 3]
	// [1 -> 3, 2 -> 4]
	// [1 -> 4, 2 -> 3]
	// [1 -> 4, 2 -> 4]
}

func ExampleDictMut() {
	g := DictMut(IntRange(1, 2), IntRange(3, 4))
	for it := iterable.Start(g.Enumerate(3)); it.HasNext(); it.Next() {
		fmt.Println(it.Current())
	}
	// Output: map[]
	// map[1:3]
	// map[1:4]
	// map[2:3]
	// map[2:4]
	// map[1:3 2:3]
	// map[1:4 2:3]
	// map[1:3 2:4]
	// map[1:4 2:4]
	// map[1:3 2:3]
	// map[1:3 2:4]
	// map[1:4 2:3]
	// map[1:4 2:4]
}
