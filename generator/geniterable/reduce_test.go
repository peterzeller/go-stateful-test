package geniterable_test

import (
	"fmt"
	"github.com/peterzeller/go-stateful-test/generator/geniterable"
)

func ExampleLength() {
	it := geniterable.New(4, 5, 6, 7)
	l := geniterable.Length(it)
	fmt.Printf("l = %d\n", l)
	// output: l = 4
}
