package generator

import (
	"github.com/peterzeller/go-fun/equality"
	"github.com/peterzeller/go-fun/iterable"
	"github.com/peterzeller/go-fun/list/linked"
	"github.com/peterzeller/go-fun/slice"
	"github.com/peterzeller/go-stateful-test/generator/geniterable"
	"github.com/peterzeller/go-stateful-test/generator/shrink"
	"math/big"
	"strings"
)

func String(chars ...rune) Generator[string] {
	if len(chars) == 0 {
		chars = []rune{'a', 'b'}
	}
	return genString{
		chars: chars,
	}
}

type genString struct {
	chars []rune
}

func (g genString) Name() string {
	return "genString"
}

func (g genString) Random(rnd Rand, size int) RandomValue[string] {
	r := rnd.R()
	length := r.Intn(size + 1)
	var s strings.Builder
	for i := 0; i < length; i++ {
		s.WriteRune(g.chars[r.Intn(len(g.chars))])
	}
	return R(s.String())
}

func (g genString) Enumerate(depth int) geniterable.Iterable[string] {
	return geniterable.FlatMap(
		geniterable.NonExhaustive(geniterable.Range(0, depth+1)),
		func(length int) geniterable.Iterable[string] {
			return enumerateStrings(length, g.chars)
		})
}

func enumerateStrings(length int, chars []rune) geniterable.Iterable[string] {
	if length <= 0 {
		return geniterable.Singleton("")
	} else {
		smaller := enumerateStrings(length-1, chars)
		return geniterable.FlatMap(smaller, func(a string) geniterable.Iterable[string] {
			return geniterable.Map(
				geniterable.FromSlice(chars),
				func(r rune) string {
					var s strings.Builder
					s.WriteString(a)
					s.WriteRune(r)
					return s.String()
				})
		})
	}
}

func (g genString) Shrink(elem RandomValue[string]) iterable.Iterable[RandomValue[string]] {
	runes := linked.FromIterable(iterable.FromString(elem.Get()))
	return iterable.Map(shrink.ShrinkList(runes, g.shrinkRune),
		func(runes *linked.List[rune]) RandomValue[string] {
			var s strings.Builder
			for it := iterable.Start[rune](runes); it.HasNext(); it.Next() {
				s.WriteRune(it.Current())
			}
			return R(s.String())
		})
}

func (g genString) shrinkRune(r rune) iterable.Iterable[rune] {
	index := slice.IndexOf(r, g.chars, equality.Default[rune]())
	switch {
	case index < 0:
		return iterable.Singleton(g.chars[0])
	case index == 0:
		return iterable.Empty[rune]()
	case index <= 5:
		return iterable.Singleton(g.chars[index-1])
	default:
		return iterable.New(g.chars[index/2], g.chars[index-1])
	}

}

func (g genString) Size(t RandomValue[string]) *big.Int {
	var sum big.Int
	for _, r := range t.Get() {
		index := slice.IndexOf(r, g.chars, equality.Default[rune]())
		if index < 0 {
			index = len(g.chars)
		}
		sum.Add(&sum, big.NewInt(int64(index)))
	}
	return &sum
}

func (g genString) RValue(elem RandomValue[string]) (string, bool) {
	return elem.Get(), true
}
