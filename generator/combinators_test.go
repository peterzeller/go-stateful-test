package generator

import (
	"fmt"
	"github.com/peterzeller/go-stateful-test/generator/geniterable"
	"math/big"
)

// This example shows how to use the Map function to create a generator for big.Int from a generator for int.
func ExampleMap() {
	// Create a generator for ints
	intGen := IntRange(0, 10)

	// Create a generator for big.Ints
	bigIntGen := Map(intGen, func(i int) *big.Int {
		return big.NewInt(int64(i))
	})

	// enumerate some values:
	vals := geniterable.ToSlice(EnumerateValues(bigIntGen, 10))
	fmt.Printf("vals = %+v\n", vals)

	// Output: vals = [+0 +1 +2 +3 +4 +5 +6 +7 +8 +9]
}

type pair struct {
	n int
	s string
}

// This example shows how to use the FlatMap function to combine two generators into one.
// In the example, we create a generator for pairs consisting of an int 'n' and a string 's'.
// The integer 'n' is a number between 1 and 3, the second value 's' is a string consisting of the first 'n' letters of the alphabet.
func ExampleFlatMap() {
	// Create a generator for ints
	nGen := IntRange(1, 3)

	alphabet := []rune{'a', 'b', 'c'}

	// Create a generator for pairs
	pairGen := FlatMap(nGen, func(n int) Generator[pair, string] {
		sGen := String(alphabet[:n]...)
		return Map(sGen, func(s string) pair {
			return pair{n: n, s: s}
		})
	})

	// enumerate some values:
	vals := geniterable.ToSlice(EnumerateValues(pairGen, 3))
	for _, v := range vals {
		fmt.Printf("(%d, %#v)\n", v.n, v.s)
	}
	// Output: (1, "")
	// (1, "a")
	// (1, "aa")
	// (1, "aaa")
	// (2, "")
	// (2, "a")
	// (2, "b")
	// (2, "aa")
	// (2, "ab")
	// (2, "ba")
	// (2, "bb")
	// (2, "aaa")
	// (2, "aab")
	// (2, "aba")
	// (2, "abb")
	// (2, "baa")
	// (2, "bab")
	// (2, "bba")
	// (2, "bbb")
	// (3, "")
	// (3, "a")
	// (3, "b")
	// (3, "c")
	// (3, "aa")
	// (3, "ab")
	// (3, "ac")
	// (3, "ba")
	// (3, "bb")
	// (3, "bc")
	// (3, "ca")
	// (3, "cb")
	// (3, "cc")
	// (3, "aaa")
	// (3, "aab")
	// (3, "aac")
	// (3, "aba")
	// (3, "abb")
	// (3, "abc")
	// (3, "aca")
	// (3, "acb")
	// (3, "acc")
	// (3, "baa")
	// (3, "bab")
	// (3, "bac")
	// (3, "bba")
	// (3, "bbb")
	// (3, "bbc")
	// (3, "bca")
	// (3, "bcb")
	// (3, "bcc")
	// (3, "caa")
	// (3, "cab")
	// (3, "cac")
	// (3, "cba")
	// (3, "cbb")
	// (3, "cbc")
	// (3, "cca")
	// (3, "ccb")
	// (3, "ccc")
}

// This example shows how to use the Zip function to combine two generators into one.
// In the example, we create a generator for pairs consisting of an int 'n' and a string 's'.
// The integer 'n' is a number between 1 and 3, the second value 's' is a string consisting of the letters 'a' and 'b'.
func ExampleZip() {
	// Create a generator for ints
	nGen := IntRange(1, 3)

	// Create a generator for strings
	sGen := String('a', 'b')

	// Create a generator for pairs
	pairGen := Zip(nGen, sGen, func(n int, s string) pair {
		return pair{n: n, s: s}
	})

	// enumerate some values:
	vals := geniterable.ToSlice(EnumerateValues(pairGen, 3))
	for _, v := range vals {
		fmt.Printf("(%d, %#v)\n", v.n, v.s)
	}
	// Output: (1, "")
	// (1, "a")
	// (1, "b")
	// (1, "aa")
	// (1, "ab")
	// (1, "ba")
	// (1, "bb")
	// (1, "aaa")
	// (1, "aab")
	// (1, "aba")
	// (1, "abb")
	// (1, "baa")
	// (1, "bab")
	// (1, "bba")
	// (1, "bbb")
	// (2, "")
	// (2, "a")
	// (2, "b")
	// (2, "aa")
	// (2, "ab")
	// (2, "ba")
	// (2, "bb")
	// (2, "aaa")
	// (2, "aab")
	// (2, "aba")
	// (2, "abb")
	// (2, "baa")
	// (2, "bab")
	// (2, "bba")
	// (2, "bbb")
	// (3, "")
	// (3, "a")
	// (3, "b")
	// (3, "aa")
	// (3, "ab")
	// (3, "ba")
	// (3, "bb")
	// (3, "aaa")
	// (3, "aab")
	// (3, "aba")
	// (3, "abb")
	// (3, "baa")
	// (3, "bab")
	// (3, "bba")
	// (3, "bbb")
}

// This example shows how to use the Filter function to create a generator that only
// yields even numbers.
func ExampleFilter() {
	// Create a generator for ints
	intGen := IntRange(0, 10)

	// Create a generator for even numbers
	evenGen := Filter(intGen, func(i int) bool {
		return i%2 == 0
	})

	// enumerate some values:
	vals := geniterable.ToSlice(evenGen.Enumerate(10))
	fmt.Printf("vals = %v\n", vals)

	// output: vals = [0 2 4 6 8]
}
