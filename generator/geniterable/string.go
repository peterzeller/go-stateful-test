package geniterable

func FromString(s string) Iterable[rune] {
	var runes []rune
	for _, r := range s {
		runes = append(runes, r)
	}
	return FromSlice(runes)
}

func FromStringBytes(s string) Iterable[byte] {
	return IterableFun[byte](func() Iterator[byte] {
		pos := 0
		return Fun[byte](func() NextResult[byte] {
			if pos >= len(s) {
				return ResultNone[byte](true)
			}
			b := s[pos]
			pos++
			return ResultSome(b)
		})
	})
}
