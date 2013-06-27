package consensus

// xlattice_go/consensus/cmdBuffer.go

import (
	"container/heap"
	// "fmt" // DEBUG
	"sync"
)

// Mediates between a producer and a consumer, where the producer
// sends a series of <number, command> pairs which will generally
// be unordered and may contain pairs with the same sequence
// number.  The consumer expects a stream of  commands in ascending
// numerical order, with no duplicates and no gaps.  This code
// effects that requirement by discarding duplicates, buffering up
// pairs until gaps are filled, and releasing pairs from the
// internal sorted buffer in ascending order as soon as possible.

// 2013-06-26 Added a BypassCh which allows a second producer to
// inject <number, command> pairs.  This might be used in testing;
// in fact it is used to bypass the consensus mechanism, an ill-
// advised practice.

type NumberedCmd struct {
	Seqn int64
	Cmd  string
}

type cmdPlus struct {
	pair  *NumberedCmd
	index int // used by heap logic
}

type pairQ []*cmdPlus

var qSy sync.RWMutex

func (q pairQ) Len() int { // not in heap interface
	// Q length is volatile, so we lock, copy, unlock, return the copy.
	// Expected use is in testing.
	qSy.RLock()
	length := len(q)
	qSy.RUnlock()
	return length
}

// implementation of the heap interface /////////////////////////////
// These functions are invoke like heap.Push(&q, &whatever)

func (q pairQ) Less(i, j int) bool {
	qSy.RLock()
	whether := q[i].pair.Seqn < q[j].pair.Seqn
	qSy.RUnlock()
	return whether
}

func (q pairQ) Swap(i, j int) {
	qSy.Lock()
	q[i], q[j] = q[j], q[i]
	q[i].index = i // remember this is post-swap
	q[j].index = j
	qSy.Unlock()
}

func (q *pairQ) Push(x interface{}) {
	qSy.Lock()
	n := len(*q)
	thePair := x.(*cmdPlus) // a cast
	thePair.index = n
	*q = append(*q, thePair)
	qSy.Unlock()
}

/////////////////////////////////////////////////////////////////////
// XXX if the length is zero, we get a panic in the heap code. //////
/////////////////////////////////////////////////////////////////////
func (q *pairQ) Pop() interface{} {
	qSy.Lock()
	nowQ := *q
	n := len(nowQ)
	if n == 0 {
		return nil
	}
	lastPair := nowQ[n-1] // last element
	lastPair.index = -1   // doesn't matter
	*q = nowQ[0 : n-1]
	qSy.Unlock()
	return lastPair
}

type CmdBufferI interface {
	InCh() chan NumberedCmd
	BypassCh() chan NumberedCmd
	Init(out chan NumberedCmd, StopCh chan bool, lastSeqn int64, bufSize int)
	Len() int
	Running() bool
	Run()
}
type CmdBuffer struct {
	inCh     chan NumberedCmd
	bypassCh chan NumberedCmd
	OutCh    chan NumberedCmd
	StopCh   chan bool
	q        pairQ
	sy       sync.Mutex
	lastSeqn int64
	running  bool
}

func (c *CmdBuffer) BypassCh() chan NumberedCmd {
	return c.bypassCh
}
func (c *CmdBuffer) InCh() chan NumberedCmd {
	return c.inCh
}

func (c *CmdBuffer) Init(out chan NumberedCmd, StopCh chan bool, lastSeqn int64, bufSize int) {
	c.q = pairQ{}
	c.inCh = make(chan NumberedCmd, bufSize)     // buffered
	c.bypassCh = make(chan NumberedCmd, bufSize) // buffered
	c.OutCh = out                                // should also be buffered
	c.StopCh = StopCh
	c.lastSeqn = lastSeqn
}

// This is synched at the q level.
func (c *CmdBuffer) Len() int {
	return c.q.Len()
}
func (c *CmdBuffer) Running() bool {
	// c.running is volatile, so we lock, copy, unlock, return the copy
	c.sy.Lock()
	whether := c.running
	c.sy.Unlock()
	return whether
}

func (c *CmdBuffer) handleCmd(inPair NumberedCmd, ok bool) bool {
	if !ok {
		// channel is closed and empty
		c.sy.Lock()
		c.running = false
		c.sy.Unlock()
		return false // we are done
	}
	seqN := inPair.Seqn
	// fmt.Printf("RECEIVED PAIR %v\n", seqN)
	if seqN <= c.lastSeqn { // already sent, so discard
		// fmt.Printf("    ALREADY SEEN, DISCARDING\n")
		return true // get next command
	} else if seqN == c.lastSeqn+1 {
		c.OutCh <- inPair
		c.lastSeqn += 1
		// fmt.Printf("    SEQN %v MATCHED LAST + 1, SENDING\n", seqN)
		for c.q.Len() > 0 {
			first := c.q[0]
			if first.pair.Seqn <= c.lastSeqn {
				//	// fmt.Printf("        Q: DISCARDING %v, DUPE\n",
				//							first.pair.Seqn)
				// a duplicate, so discard
				_ = heap.Pop(&c.q).(*cmdPlus)
			} else if first.pair.Seqn == c.lastSeqn+1 {
				pp := heap.Pop(&c.q).(*cmdPlus)
				c.OutCh <- *pp.pair
				c.lastSeqn += 1
				//	// fmt.Printf("        Q: SENT %v\n", c.lastSeqn)
			} else {
				//	// fmt.Printf("        Q: LEAVING %v IN Q\n",
				//							first.pair.Seqn)
				break
			}
		}
	} else {
		// seqN > c.lastSeqn + 1, so buffer
		// fmt.Printf("    HIGH SEQN %v, SO BUFFERING\n", seqN)
		pp := &cmdPlus{pair: &inPair}
		heap.Push(&c.q, pp)
	}
	return true // get next command
}

func (c *CmdBuffer) Run() {
	c.running = true
	for {
		c.sy.Lock()
		whether := c.running
		c.sy.Unlock()
		if !whether {
			break
		}
		select {
		case inPair, ok := <-c.inCh: // get the next command
			if c.handleCmd(inPair, ok) {
				continue
			} else {
				break
			}
		case inPair, ok := <-c.bypassCh: // get the next command
			if c.handleCmd(inPair, ok) {
				continue
			} else {
				break
			}
		case <-c.StopCh:
			c.sy.Lock()
			c.running = false
			c.sy.Unlock()
			// fmt.Println("c.running has been set to false")
		}
	}
}
