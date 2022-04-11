package generator

import (
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

func (r *testRand) Fork(name string) Rand {
	return r
}
func (r *testRand) HasMore() bool {
	return false
}

// R is the underlying random number generator
func (r *testRand) R() *rand.Rand {
	return r.rnd
}
