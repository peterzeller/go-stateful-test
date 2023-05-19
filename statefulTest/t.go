package statefulTest

import (
	"github.com/peterzeller/go-stateful-test/generator"
)

// T is the testing interface used for writing stateful tests
type T interface {
	Errorf(format string, args ...interface{})
	FailNow()
	Logf(format string, args ...any)
	// PickValue returns the next value
	PickValue(untyped generator.UntypedGenerator) generator.UV
	// HasMore is used for generating a sequence of values
	HasMore() bool
	// Cleanup runs a function when the test is done
	Cleanup(f func())
}
