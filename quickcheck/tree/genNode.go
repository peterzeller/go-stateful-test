package tree

import (
	"fmt"
	"github.com/peterzeller/go-fun/list/linked"
	"strings"
)

type GenNode struct {
	generatedValues *linked.List[*linked.List[GeneratedValue]]
}

func New(gvs *linked.List[*linked.List[GeneratedValue]]) *GenNode {
	return &GenNode{generatedValues: gvs}
}

func (n GenNode) GeneratedValues() *linked.List[*linked.List[GeneratedValue]] {
	return n.generatedValues
}

func (n GenNode) With(gvs *linked.List[*linked.List[GeneratedValue]]) *GenNode {
	return &GenNode{generatedValues: gvs}
}

func (n GenNode) String() string {
	var s strings.Builder
	s.WriteString("G[")
	for i, value := range n.generatedValues.ToSlice() {
		if i > 0 {
			s.WriteString(", ")
		}
		s.WriteString("[")
		for j, v := range value.ToSlice() {
			if j > 0 {
				s.WriteString(", ")
			}
			s.WriteString(fmt.Sprintf("%s: %v", v.Generator.Name(), v.Value))
		}
		s.WriteString("]")
	}
	s.WriteString("]")
	return s.String()
}
