# Go Stateful Test

This is an easy-to-use library for property based testing in Go.

Status:

 - Work in progress. Don't use this :)


Features:

 - Type safe generators (using Go 1.18 generics)
 - Integration with standard Go tests and assertions
 - Designed for testing stateful components
 - Quickcheck support (random testing)
 - Smallcheck support (exhaustive bounded testing)

## Smallcheck Example

The following function is supposed to compute the maximum out of 3 integers.
Unfortunately it contains a bug.
Can you find a counter example where it would fail?

    func max3(x, y, z int) int {
        if x > y && x > z {
            return x
        } else if y > x && y > z {
            return y
        } else {
            return z
        }
    }

With SmallCheck it is easy to write a test that finds a minimal counter example:

    func TestMax3(t *testing.T) {
        smallcheck.Run(t, smallcheck.Config{}, func(t statefulTest.T) {
            x := pick.Val(t, generator.Int())
            y := pick.Val(t, generator.Int())
            z := pick.Val(t, generator.Int())
            res := max3(x, y, z)
            t.Logf("min3(%d, %d, %d) = %d", x, y, z, res)
            assert.True(t, res >= x, "res >= x")
            assert.True(t, res >= y, "res >= y")
            assert.True(t, res >= z, "res >= z")
        })
    }

As the example shows, the test is a simple Go unit test and standard assertions can be used.
When we run the test with `go test`, we get the following error with a minimized example:

     === RUN   TestMax3
        run.go:48: Test failed:
            min3(1, 1, 0) = 0
                Error Trace:	smallcheck_example_test.go:33
                                            run.go:36
                                            state.go:40
                                            state.go:41
                                            run.go:46
                                            smallcheck_example_test.go:27
                Error:      	Should be true
                Messages:   	res >= x
            
                Error Trace:	smallcheck_example_test.go:34
                                            run.go:36
                                            state.go:40
                                            state.go:41
                                            run.go:46
                                            smallcheck_example_test.go:27
                Error:      	Should be true
                Messages:   	res >= y
    --- FAIL: TestMax3 (0.00s)
    
    
## Quickcheck Example

The API for Quickcheck follows the same conventions as the SmallCheck API.
The 

    func TestMax3Quick(t *testing.T) {
        quickcheck.Run(t, quickcheck.Config{}, func(t statefulTest.T) {
            x := pick.Val(t, generator.Int())
            y := pick.Val(t, generator.Int())
            z := pick.Val(t, generator.Int())
            res := max3(x, y, z)
            t.Logf("min3(%d, %d, %d) = %d", x, y, z, res)
            assert.True(t, res >= x, "res >= x")
            assert.True(t, res >= y, "res >= y")
            assert.True(t, res >= z, "res >= z")
        })
    }

This gives the following error when run with `go test`:

    === RUN   TestMax3Quick
        run.go:51: Found error, shrinking testcase ...
        run.go:59: Shrunk Test Run:
            min3(3, 3, 0) = 0
                Error Trace:	example_test.go:58
                                            run.go:37
                                            run.go:45
                                            run.go:79
                                            run.go:43
                                            example_test.go:52
                Error:      	Should be true
                Messages:   	res >= x
            
                Error Trace:	example_test.go:59
                                            run.go:37
                                            run.go:45
                                            run.go:79
                                            run.go:43
                                            example_test.go:52
                Error:      	Should be true
                Messages:   	res >= y
    --- FAIL: TestMax3Quick (0.00s)