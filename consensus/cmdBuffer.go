package consensus

// xlattice_go/consensus/cmdBuffer.go

import (
	"container/heap"
	//"sync"
)

// Mediates between a producer and a consumer, where the producer
// sends a series of <number, command> pairs which will generally
// be unordered and may contain pairs with the same sequence
// number.  The consumer expects a stream of  commands in ascending
// numerical order, with no duplicates and no gaps.  This code
// effects that requirement by discarding duplicates, buffering up
// pairs until gaps are filled, and releasing pairs from the
// internal sorted buffer in ascending order as soon as possible.

type CmdPair struct {
	Seqn int64
	Cmd  string
}

type pairPlus struct {
	pair  *CmdPair
	index int // used by heap logic
}

type pairQ []*pairPlus

func (q pairQ) Len() int { // not in heap interface
	return len(q)
}

// implementation of the heap interface /////////////////////////////
// These functions are invoke like heap.Push(&q, &whatever)

////////////////////////////////////////////////////
// XXX SO FAR, THIS IMPLEMENATION ACCEPTS DUPLICATES
////////////////////////////////////////////////////

func (q pairQ) Less(i, j int) bool {
	return q[i].pair.Seqn < q[j].pair.Seqn
}

func (q pairQ) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].index = i // remember this is post-swap
	q[j].index = j
}

func (q *pairQ) Push(x interface{}) {
	n := len(*q)
	thePair := x.(*pairPlus) // a cast
	thePair.index = n
	*q = append(*q, thePair)
}

/////////////////////////////////////////////////////////////////////
// XXX if the length is zero, we get a panic in the heap code. //////
/////////////////////////////////////////////////////////////////////
func (q *pairQ) Pop() interface{} {
	nowQ := *q
	n := len(nowQ)
	if n == 0 {
		return nil
	}
	lastPair := nowQ[n-1] // last element
	lastPair.index = -1   // doesn't matter
	*q = nowQ[0 : n-1]
	return lastPair
}

//
type CmdBuffer struct {
	InCh     chan CmdPair
	outCh    chan CmdPair
	stop     chan bool
	q        pairQ
	lastSeqn int64
}

func (c CmdBuffer) Init(out chan CmdPair, stop chan bool, lastSeqn int64) {
	c.q = pairQ{}
	c.InCh = make(chan CmdPair, 4) // buffered

	c.outCh = out // should also be buffered
	c.stop = stop
	c.lastSeqn = lastSeqn
}
func (c CmdBuffer) Run() {
	for running := true; running; {
		select {
		case <-c.stop:
			running = false
		case inPair := <-c.InCh: // get the next command
			seqN := inPair.Seqn
			if seqN <= c.lastSeqn { // already sent, so discard
				continue
			} else if seqN == c.lastSeqn+1 {
				c.outCh <- inPair
				c.lastSeqn += 1
				for c.q.Len() > 0 {
					first := c.q[0]
					if first.pair.Seqn <= c.lastSeqn {
						// a duplicate, so discard
						_ = heap.Pop(&c.q).(*pairPlus)
					} else if first.pair.Seqn == c.lastSeqn+1 {
						pp := heap.Pop(&c.q).(*pairPlus)
						c.outCh <- *pp.pair
						c.lastSeqn += 1
					}
				}
			} else {
				// seqN > c.lastSeqn + 1, so buffer
				pp := pairPlus{pair: &inPair}
				heap.Push(&c.q, pp)
			}
		}
	}
}
