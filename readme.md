# Go Stateful Test

This is an easy-to-use library for property based testing in Go.

Status:

 - Work in progress. Don't use this :)


Features:

 - Type safe generators (using Go 1.18 generics)
 - Integration with standard Go tests and assertions
 - Designed for testing stateful components
 - Quickcheck support (random testing)
 - Planned: Smallcheck support (exhaustive bounded testing)

## Basic Example

The following excerpt from [example_test.go](./examples/example_test.go) shows a simple test case that generates two integers and tests a whether they satisfy a condition:

    func TestInts(t *testing.T) {
        quickcheck.Run(t, quickcheck.Config{}, func(t statefulTest.T) {
            x := pick.Val(t, generator.Int())
            y := pick.Val(t, generator.Int())
            t.Logf("x = %d, y = %d", x, y)
            require.True(t, x+y < 10)
        })
    }

As the example shows, the test is a simple Go unit test and standard assertions can be used.
When we run the test with `go test`, we get the following error with a minimized example:

     example_test.go:47: error as expected: test failed
        Found error, shrinking testcase ...
        Shrunk Test Run:
        x = 0, y = 10
        	Error Trace:	example_test.go:66
        	            				run.go:37
        	            				shrink.go:36
        	            				shrink.go:12
        	            				run.go:56
        	            				example_test.go:62
        	            				example_test.go:57
        	            				example_test.go:61
        	Error:      	Should be true
