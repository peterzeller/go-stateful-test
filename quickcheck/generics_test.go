package quickcheck

import (
	"reflect"
	"testing"
)

type Gen[T any] interface {
	Give() T
}

type intGen struct {
}

func (i intGen) Give() int {
	return 42
}

type stringGen struct {
}

func (i stringGen) Give() string {
	return "foo"
}

func TestBlub(t *testing.T) {
	gens := []interface{}{
		intGen{},
		stringGen{},
	}
	for i, gen := range gens {
		t.Logf("i = %d", i)
		switch gen.(type) {
		case Gen[interface{}]:
			t.Logf("adfsa")
		case Gen[string]:
			t.Logf("string")
		//case Gen[int]:
		//	t.Logf("int")
		default:
			t.Logf("other: %v", reflect.TypeOf(gen))
		}
	}

	t.Logf("gens = %#v", gens)
}
