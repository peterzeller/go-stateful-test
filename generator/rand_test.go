package generator

import (
	"github.com/peterzeller/go-stateful-test/quickcheck/randomsource"
	"math/rand"
)

func newTestRand(seed int64) *testRand {
	return &testRand{
		rnd: rand.New(rand.NewSource(seed)),
	}
}

// testRand is an implementation of Rand for testing purposes
type testRand struct {
	rnd *rand.Rand
}

func (r *testRand) UseHeuristics() bool {
	return true
}

func (r *testRand) Fork(name string) Rand {
	return r
}
func (r *testRand) HasMore() bool {
	return false
}

// R is the underlying random number generator
func (r *testRand) R() randomsource.RandomStream {
	return randomsource.FromSeed(r.rnd.Int63()).Iterator()
}
