package pick

import (
	"github.com/peterzeller/go-stateful-test/generator"
	"github.com/peterzeller/go-stateful-test/statefulTest"
)

func Val[V any](t statefulTest.T, g generator.Generator[V]) V {
	v := t.PickValue(generator.ToUntyped(g))
	return v.(V)
}

type Of map[string]func()

//
func One(of Of) {

}
