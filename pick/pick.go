package pick

import (
	"github.com/peterzeller/go-fun/dict/hashdict"
	"github.com/peterzeller/go-fun/hash"
	"github.com/peterzeller/go-fun/reducer"
	"github.com/peterzeller/go-stateful-test/generator"
	"github.com/peterzeller/go-stateful-test/statefulTest"
)

func Val[V any](t statefulTest.T, g generator.Generator[V]) V {
	v := t.PickValue(generator.ToUntyped(g))
	return v.(V)
}

type Cases map[string]func()

// Switch executes one function from the map
func Switch(t statefulTest.T, of Cases) {
	// collect keys from the map
	keys := hashdict.FromMap(hash.String(), of).Keys()
	compareStrings := func(a, b string) bool {
		return a < b
	}
	keysSorted := reducer.Sorted(compareStrings, reducer.ToSlice[string]()).Apply(keys)
	// pick key
	key := Val(t, generator.OneConstantOf(keysSorted...))
	// execute function
	of[key]()
}
