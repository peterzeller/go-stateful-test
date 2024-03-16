package examples

import (
	"github.com/peterzeller/go-stateful-test/fuzzcheck"
	"github.com/peterzeller/go-stateful-test/generator"
	"github.com/peterzeller/go-stateful-test/pick"
	"github.com/peterzeller/go-stateful-test/statefulTest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"sync"
	"testing"
)

// FuzzMax3Quick tests the buggy max3 function with fuzzing.
func FuzzMax3Quick(t *testing.F) {
	// We skip the test, because it would find the error and I could not figure out how to write a fuzz test that expects an error
	t.SkipNow()
	fuzzcheck.Run(t, fuzzcheck.Config{}, func(t statefulTest.T) {
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

var mut sync.Mutex

// AppendStringToFile appends the given text to the file specified by filename.
// If the file does not exist, it will be created.
func appendStringToFile(filename, text string) (err error) {
	mut.Lock()
	defer mut.Unlock()
	// Open the file in append mode. If it doesn't exist, create it with permissions 0644.
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func() {
		err = file.Close()
	}()

	// Write the text to the file
	_, err = file.WriteString(text)
	if err != nil {
		return err
	}

	return nil
}

// FuzzComplexIfs tests code with some complex branching logic.
// You need exactly the right value 5 times to get to the bug.
// Both quickcheck and smallcheck would fail to find this bug, but coverage guided fuzzing finds it.
func FuzzComplexIfs(t *testing.F) {
	// We skip the test, because it would find the error and I could not figure out how to write a fuzz test that expects an error
	t.SkipNow()
	fuzzcheck.Run(t, fuzzcheck.Config{DisableHeuristics: true}, func(t statefulTest.T) {
		a := pick.Val(t, generator.IntRange(0, 20))
		b := pick.Val(t, generator.IntRange(0, 20))
		c := pick.Val(t, generator.IntRange(0, 20))
		d := pick.Val(t, generator.IntRange(0, 20))
		e := pick.Val(t, generator.IntRange(0, 20))
		f := pick.Val(t, generator.IntRange(0, 20))
		t.Logf("%d, %d, %d, %d, %d, %d\n", a, b, c, d, e, f)
		//require.NoError(t, appendStringToFile("log.txt", fmt.Sprintf("%d, %d, %d, %d, %d, %d\n", a, b, c, d, e, f)))
		if a != 11 {
			return
		}
		if b != 12 {
			return
		}
		if c != 13 {
			return
		}
		if d != 14 {
			return
		}
		if e != 15 {
			return
		}
		require.NotEqual(t, 16, f)
	})
}
