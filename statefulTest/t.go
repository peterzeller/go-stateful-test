package statefulTest

// T is the testing interface used for writing stateful tests
type T interface {
	Errorf(format string, args ...interface{})
	FailNow()
	Logf(format string, args ...any)
}
