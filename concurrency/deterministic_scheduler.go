package concurrency

import (
	"fmt"

	"github.com/peterzeller/go-fun/equality"
	"github.com/peterzeller/go-fun/list"
)

// DeterministicScheduler executes at most one thread at a time, which guarantees deterministic execution
type DeterministicScheduler struct {
	currentThread *Thread
	threads       list.List[*Thread]
	onYield       func(location string)
}

type Thread struct {
	name string
	// waitingChan is a channel on which threads are waiting for their turn to run
	waitingChan chan bool
	status      threadStatus
}

type threadStatus int

const running = 0
const waiting = 1
const terminated = 2

var _ Scheduler = &DeterministicScheduler{}

func NewDeterministicScheduler(onYield func(location string)) *DeterministicScheduler {
	return &DeterministicScheduler{
		currentThread: nil,
		threads:       list.New[*Thread](),
		onYield:       onYield,
	}
}

// Go starts a new go routine on the scheduler.
// This is equivalent
func (s *DeterministicScheduler) Go(threadName string, f func()) {
	t := &Thread{
		name:        threadName,
		waitingChan: make(chan bool),
	}
	s.threads = s.threads.AppendElems(t)
	go func() {
		defer func() {
			t.status = terminated
			// when the thread is completed, remove it from the scheduler
			s.threads = s.threads.RemoveFirst(t, equality.Default[*Thread]())
			// and treat it as a yield
			s.onYield("completed_thread_" + threadName)
		}()

		// wait for thread to be activated
		<-t.waitingChan
		// run the thread implementation
		f()
	}()

	s.Yield("started_thread_" + threadName)
}

func (s *DeterministicScheduler) Yield(location string) {
	if s.currentThread == nil {
		return
	}
	t := s.currentThread
	s.currentThread = nil
	t.status = waiting

	s.onYield(location)
	// wait for the thread to be activated again
	<-t.waitingChan
}

func (s *DeterministicScheduler) Threads() list.List[*Thread] {
	return s.threads
}

func (s *DeterministicScheduler) ContinueThread(thread *Thread) {
	if s.currentThread != nil {
		panic(fmt.Errorf("ContinueThread called while thread '%s' is still active", s.currentThread.name))
	}
	if !s.threads.Contains(thread, equality.Default[*Thread]()) {
		panic(fmt.Errorf("ContinueThread called with invalid thread '%s'", thread.name))
	}
	s.currentThread = thread
	thread.waitingChan <- true
}
