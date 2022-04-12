package examples

import (
	"testing"

	"github.com/peterzeller/go-stateful-test/generator"
	"github.com/peterzeller/go-stateful-test/pick"
	"github.com/peterzeller/go-stateful-test/quickcheck"
	"github.com/peterzeller/go-stateful-test/smallcheck"
	"github.com/peterzeller/go-stateful-test/statefulTest"
	"github.com/stretchr/testify/require"
)

// Queue implements integer queue with a fixed maximum size.
type Queue struct {
	buf []int
	in  int
	out int
}

func NewQueue(n int) *Queue {
	return &Queue{
		buf: make([]int, n+1),
	}
}

// Precondition: Size() > 0.
func (q *Queue) Get() int {
	i := q.buf[q.out]
	q.out = (q.out + 1) % len(q.buf)
	return i
}

// Precondition: Size() < n.
func (q *Queue) Put(i int) {
	q.buf[q.in] = i
	q.in = (q.in + 1) % len(q.buf)
}

func (q *Queue) Size() int {
	return (q.in - q.out) % len(q.buf)
}

func QueueProperty(t statefulTest.T) {
	// init a queue with given size
	n := pick.Val(t, generator.IntRange(1, 10))
	q := NewQueue(n)
	t.Logf("Initialize NewQueue(%d)", n)

	// model the state of the queue with a slice
	var model []int

	// repeat commands
	for t.HasMore() {
		pick.Switch(t, pick.Cases{
			"get": func() {
				// Test q.Get:
				if q.Size() == 0 {
					// skip if queue is empty
					return
				}
				i := q.Get()
				t.Logf("Calling q.Get() -> %d", i)
				require.Equal(t, model[0], i, "result of q.Get()")
				model = model[1:]
			},
			"put": func() {
				// Test q.Put
				if q.Size() >= n {
					// skip if queue is full
					return
				}
				i := pick.Val(t, generator.Int())
				t.Logf("Calling q.Put(%d)", i)
				q.Put(i)
				model = append(model, i)
			},
		})
		// check invariant
		require.Equal(t, len(model), q.Size(), "invariant: queue size")
	}
}

func TestQueueWithQuickCheck(t *testing.T) {
	expectError(t, func(t quickcheck.TestingT) {
		quickcheck.Run(t, quickcheck.Config{}, QueueProperty)
	})
}

func TestQueueWithSmallCheck(t *testing.T) {
	expectError(t, func(t quickcheck.TestingT) {
		smallcheck.Run(t, smallcheck.Config{}, QueueProperty)
	})
}
