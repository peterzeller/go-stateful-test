package paxos

import (
	"fmt"

	"github.com/peterzeller/go-fun/hash"
	"github.com/peterzeller/go-fun/iterable"
	"github.com/peterzeller/go-fun/set/hashset"
	"github.com/peterzeller/go-stateful-test/concurrency"
	"github.com/peterzeller/go-stateful-test/pick"
	"github.com/peterzeller/go-stateful-test/statefulTest"
)

type message struct {
	typ  messageType
	from process
	to   process
	bal  ballot
	val  voteValue
	vbal ballot
	vval voteValue
}

type process int

func (p process) Hash() int64 {
	return int64(p)
}

type messageType string

func (t messageType) Hash() int64 {
	return hash.String().Hash(string(t))
}

const messageType1a = "1a"
const messageType1b = "1b"
const messageType2a = "2a"
const messageType2b = "2b"

type ballot struct {
	number   int
	proposer process
}

func (b ballot) Hash() int64 {
	return hash.CombineHashes(
		int64(b.number),
		b.proposer.Hash(),
	)
}

func (b ballot) Less(other ballot) bool {
	return b.number < other.number ||
		b.number == other.number && b.proposer < other.proposer
}

func (b ballot) LessEq(other ballot) bool {
	return b.number < other.number ||
		b.number == other.number && b.proposer <= other.proposer
}

var NoBallot = ballot{-1, -1}

type voteValue string

func (v voteValue) Hash() int64 {
	return hash.String().Hash(string(v))
}

const NoValue voteValue = "no-value"

type system struct {
	scheduler concurrency.Scheduler
	t         statefulTest.T
	msgs      hashset.Set[message]
}

func (s *system) Proposer(self process) {
	// For each proposer, its ballot number.
	pBal := NoBallot
	// For each proposer, ballot of the highest registered vote.
	pVBal := NoBallot
	// For each proposer, value of the highest registered vote.
	pVVal := NoValue
	// Sets of acceptor ids for keeping record of "1b" and "2b"
	pQ1 := hashset.New(processHash())
	//   messages, repectively, received by each proposer.
	pQ2 := hashset.New(messageHash())
	// pWr[p] = Is proposer p's voted value pVVal written?
	pWr := false
	// Ballot of the last "2a" message sent by proposer.
	pLBal := NoBallot

p1:
	for {
		// Proposer step 1 [Set and send ballot].  Set the ballot number to
		// the current number plus one, store that number in pBal, and send a
		// "1a" message to all acceptors.
		//
		// A proposer p can be preempted.  Some acceptor may preempt the
		// execution of p by replying to p with a ballot number higher than
		// p's ballot.  In this case, p is enabled to execute action P1 again
		// and allowed to set a new ballot number.
		// Resetting the variables pVBal, pVVal, pQ1 and pQ2 is required in
		// case process prop was preempted.
		s.when("p1", func() bool {
			return !pWr
		})

		pBal = nextBallot(pBal, self)
		s.send(message{
			typ:  messageType1a,
			from: self,
			bal:  pBal,
		})
		pVBal = NoBallot
		pVVal = NoValue
		pQ1 = hashset.New(processHash())
		pQ2 = hashset.New(messageHash())

		// Proposer step 2a [Wait]. Receive and process one "1b" message at a
		// time, until a quorum of acceptors have replied. The messages must
		// satisfy the following conditions: p is the message's target, the
		// message has the same ballot as the proposer's. The sender ids (the
		// acceptors ids) are recorded in pQ1, until there is a majority of
		// acceptors in pQ1. If the message's ballot is higher than the
		// current ballot, the execution is aborted and restarted from P1.
		// The variables pVBal and pVVal store the ballot and vote of the
		// highest-seen ballot, discarding the votes that come with the lower
		// ballots.
		for !s.isQuorum(pQ1) {
			var m message
			s.when("p2", func() bool {
				relevantMessages := iterable.Filter[message](s.msgs, func(m message) bool {
					return m.typ == messageType1b &&
						m.to == self &&
						!pQ1.Contains(m.from)
				})
				// TODO inline cases
				cases := pick.Cases{}
				iterable.Foreach(relevantMessages, func(x message) {
					cases[x.String()] = func() {
						m = x
					}
				})
				if len(cases) > 0 {
					pick.Switch(s.t, cases)
					return true
				}
				return false
			})
			if m.bal == pBal {
				pQ1 = pQ1.Add(m.from)
				if pVBal.Less(m.vbal) {
					pVBal = m.vbal
					pVVal = m.vval
				}
			} else if pBal.Less(m.bal) {
				continue p1
			}
		}

	}

}

func (s *system) when(loc string, cond func() bool) {
	for !cond() {
		s.scheduler.Yield(loc)
	}
}

func processHash() hash.EqHash[process] {
	return hash.Map(hash.Num[int](), func(v process) int {
		return int(v)
	})
}

func messageHash() hash.EqHash[message] {
	return hash.Natural[message]()
}

func (m message) Equal(other message) bool {
	return m == other
}

func (m message) Hash() int64 {
	return hash.CombineHashes(
		m.typ.Hash(),
		m.from.Hash(),
		m.to.Hash(),
		m.bal.Hash(),
		m.val.Hash(),
		m.vbal.Hash(),
		m.vval.Hash(),
	)
}

func (m message) String() string {
	return fmt.Sprintf("%+v", m)
}
