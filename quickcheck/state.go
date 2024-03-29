package quickcheck

import (
	"fmt"
	"github.com/peterzeller/go-fun/iterable"
	"github.com/peterzeller/go-stateful-test/generator/geniterable"
	"math/big"
	"math/rand"
	"strings"

	"github.com/peterzeller/go-fun/list/linked"
	"github.com/peterzeller/go-stateful-test/generator"
	"github.com/peterzeller/go-stateful-test/quickcheck/tree"
	"github.com/peterzeller/go-stateful-test/statefulTest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type state struct {
	mainFork *fork
	// initialized to false and set to true when the test has failed
	failed bool
	// buffer for log messages.
	// As we only want to print the log for the last failed test run, we cannot write directly to standard out.
	log     strings.Builder
	cleanup []func()
	cfg     Config
}

func (s *state) Cleanup(f func()) {
	s.cleanup = append(s.cleanup, f)
}

// fork of a state
type fork struct {
	parent *state
	// record the generated values in a tree-like structure, which we use for shrinking
	genTree *tree.MutableGenNode
	// optional tree with preset values that is used for shrinking.
	// If a value exists in the presetTree we take it from there.
	// Otherwise, we generate a new one.
	presetTree *tree.GenNode
	// maxSize is the maximum size to generate when picking random values
	maxSize int
}

func (f *fork) String() string {
	return fmt.Sprintf("fork{genTree: %v, presetTree: %v}", f.genTree, f.presetTree)
}

type forkGenerator struct {
	origin *fork
	name   string
}

var _ generator.Generator[*fork, *fork] = &forkGenerator{}

func (f forkGenerator) Name() string {
	return f.name
}

func (f forkGenerator) Random(rnd generator.Rand, size int) *fork {
	seed := rnd.R().Int63()
	return &fork{
		parent:     f.origin.parent,
		genTree:    tree.NewGenNode(seed),
		presetTree: nil,
		maxSize:    size,
	}
}

func (f forkGenerator) Size(elem *fork) *big.Int {
	return elem.Size()
}

func (f forkGenerator) Enumerate(depth int) geniterable.Iterable[*fork] {
	panic("enumerate is not implemented for quickcheck")
}

func (f forkGenerator) Shrink(elem *fork) iterable.Iterable[*fork] {
	shrinks := shrinkTree(elem.presetTree)
	return iterable.Map(shrinks,
		func(t *tree.GenNode) *fork {
			return &fork{
				parent:     f.origin.parent,
				genTree:    tree.NewGenNode(0),
				presetTree: t,
				maxSize:    0,
			}
		})
}

func (f forkGenerator) RValue(rv *fork) (*fork, bool) {
	return rv, true
}

func (f *fork) Fork(name string) generator.Rand {
	gen := generator.ToUntyped[*fork, *fork](forkGenerator{name: name, origin: f})
	child := f.PickValue(gen).Value.(*fork)
	return child
}

func (f *fork) R() *rand.Rand {
	return f.genTree.Rand
}

func (f *fork) PickValue(gen generator.UntypedGenerator) generator.UV {
	// check if we have a preset value in the presetTree
	genName := gen.Name()
	var picked tree.GeneratedValue
	foundPreset := false
	if f.presetTree != nil {
		generatedValues := f.presetTree.GeneratedValues()
		v, newGVhead, ok := generatedValues.Head().FindAndRemove(func(gv tree.GeneratedValue) bool {
			return gv.Generator.Name() == genName
		})
		if ok {
			picked = v
			f.presetTree = f.presetTree.With(linked.Cons(newGVhead, generatedValues.Tail()))
			foundPreset = true
		}
	}
	if !foundPreset {
		// generate new random value
		v := gen.Random(f, f.maxSize)
		picked = tree.GeneratedValue{
			Generator: gen,
			Value:     v,
		}
	}
	// store generated value
	lastIndex := len(f.genTree.GeneratedValues) - 1
	if lastIndex < 0 {
		f.genTree.GeneratedValues = [][]tree.GeneratedValue{{}}
		lastIndex = 0
	}
	f.genTree.GeneratedValues[lastIndex] = append(f.genTree.GeneratedValues[lastIndex], picked)
	repaired, ok := gen.RValue(picked.Value)
	if !ok {
		panic(fmt.Errorf("quickcheck: could not convert generated value %v to generator %v", picked.Value, gen.Name()))
	}
	return repaired
}

func (f *fork) HasMore() bool {
	result := false
	if f.presetTree != nil {
		// replay from presetTree
		length := f.presetTree.GeneratedValues().Length()
		if length > 0 {
			result = length > 1
			// move to next section
			old := f.presetTree
			f.presetTree = tree.New(old.GeneratedValues().Tail())
		}
	} else {
		if f.genTree.Rand.Float64()*float64(f.maxSize) > 1 {
			result = true
		}
	}
	// record HasMore in genTree by appending a new section
	if result {
		f.genTree.GeneratedValues = append(f.genTree.GeneratedValues, []tree.GeneratedValue{})
	}
	return result
}

func (f *fork) Size() *big.Int {
	return f.genTree.Size()
}

func (s *state) PickValue(gen generator.UntypedGenerator) generator.UV {
	return s.mainFork.PickValue(gen)
}

func (s *state) HasMore() bool {
	return s.mainFork.HasMore()
}

func (s *state) Logf(format string, args ...any) {
	if s.cfg.PrintAllLogs {
		fmt.Printf(format, args...)
		fmt.Printf("\n")
		return
	}
	// TODO implement like in real Log and add source code line to message?
	_, _ = fmt.Fprintf(&s.log, format, args...)
	s.log.WriteRune('\n')
}

var _ statefulTest.T = &state{}
var _ assert.TestingT = &state{}
var _ require.TestingT = &state{}

func (s *state) Errorf(format string, args ...interface{}) {
	s.failed = true
	_, _ = fmt.Fprintf(&s.log, format, args...)
}

var errTestFailed = fmt.Errorf("test failed")

func (s *state) FailNow() {
	s.failed = true
	panic(errTestFailed)
}

func (s *state) Failed() bool {
	return s.failed
}

func (s *state) GetLog() string {
	return s.log.String()
}

func (s *state) runCleanups() {
	for _, f := range s.cleanup {
		f()
	}
	s.cleanup = nil
}

func initState(cfg Config, seed int64) *state {
	s := &state{
		mainFork: &fork{
			parent:     nil,
			genTree:    tree.NewGenNode(seed),
			presetTree: nil,
			maxSize:    100, // TODO init differently
		},
		failed: false,
		log:    strings.Builder{},
		cfg:    cfg,
	}
	s.mainFork.parent = s
	return s
}
