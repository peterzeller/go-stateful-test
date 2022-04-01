package quickcheck

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"

	"github.com/peterzeller/go-fun/dict/hashdict"
	"github.com/peterzeller/go-fun/hash"
	"github.com/peterzeller/go-fun/linkedlist"
	"github.com/peterzeller/go-stateful-test/statefulTest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type state struct {
	// record the generated values in a tree-like structure, which we use for shrinking
	genTree *mutableGenNode
	// optional tree with preset values that is used for shrinking.
	// If a value exists in the presetTree we take it from there.
	// Otherwise, we generate a new one.
	presetTree *genNode
	// initialized to false and set to true when the test has failed
	failed bool
	// size is the total size of all generated values in the run
	size int64
	// buffer for log messages.
	// As we only want to print the log for the last failed test run, we cannot write directly to standard out.
	log strings.Builder
	mut sync.Mutex
}

func (s *state) Logf(format string, args ...any) {
	s.mut.Lock()
	defer s.mut.Unlock()
	// TODO implement like in real Log and add source code line to message?
	_, _ = fmt.Fprintf(&s.log, format, args...)
}

var _ statefulTest.T = &state{}
var _ assert.TestingT = &state{}
var _ require.TestingT = &state{}

func (s *state) Errorf(format string, args ...interface{}) {
	s.mut.Lock()
	defer s.mut.Unlock()
	s.failed = true
	_, _ = fmt.Fprintf(&s.log, format, args...)
}

var testFailedErr = fmt.Errorf("test failed")

func (s *state) FailNow() {
	s.mut.Lock()
	defer s.mut.Unlock()
	s.failed = true
	panic(testFailedErr)
}

func (s *state) Failed() bool {
	s.mut.Lock()
	defer s.mut.Unlock()
	return s.failed
}

func (s *state) GetLog() string {
	s.mut.Lock()
	defer s.mut.Unlock()
	return s.log.String()
}

func initState(seed int64) *state {
	return &state{
		genTree: newGenNode(seed),
		failed:  false,
		log:     strings.Builder{},
		mut:     sync.Mutex{},
	}
}

type mutableGenNode struct {
	children  map[string]*mutableGenNode
	generated map[string][]interface{}
	rand      *rand.Rand
	seed      int64
	// hasMoreCount counts how often HasMore was called.
	hasMoreCount int
}

func newGenNode(seed int64) *mutableGenNode {
	return &mutableGenNode{
		children:     map[string]*mutableGenNode{},
		generated:    map[string][]interface{}{},
		rand:         rand.New(rand.NewSource(seed)),
		seed:         seed,
		hasMoreCount: 0,
	}
}

func (m *mutableGenNode) toImmutable() *genNode {
	cs := hashdict.New[string, *genNode](hash.String())
	for k, v := range m.children {
		cs = cs.Set(k, v.toImmutable())
	}
	generated := hashdict.New[string, *linkedlist.LinkedList[interface{}]](hash.String())
	for k, v := range m.generated {
		generated = generated.Set(k, linkedlist.New(v...))
	}
	return &genNode{
		children:     cs,
		generated:    generated,
		hasMoreCount: m.hasMoreCount,
	}
}

type genNode struct {
	children     hashdict.Dict[string, *genNode]
	generated    hashdict.Dict[string, *linkedlist.LinkedList[interface{}]]
	hasMoreCount int
}
