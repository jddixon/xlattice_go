package consensus

// xlattice_go/consensus/cmdBuffer.go

import (
	"container/heap"
	"fmt"
	"os"
	"strings"
	"sync"
)

var _ = fmt.Println

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

const LOG_BUFFER_SIZE = 16 * 256 * 256

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
	Init(chan NumberedCmd, chan bool, int64, int, string, int, bool)
	Len() int
	Running() bool
	Run() error
}
type CmdBuffer struct {
	inCh      chan NumberedCmd
	bypassCh  chan NumberedCmd
	OutCh     chan NumberedCmd // output only
	StopCh    chan bool        // receive only
	q         pairQ
	sy        sync.Mutex
	lastSeqn  int64
	running   bool
	pathToLog string // if none, we aren't logging
	fd        *os.File
	b         *logBuffer

	verbose  int  // :-(
	disabled bool // :-( :-(
}

func (c *CmdBuffer) BypassCh() chan NumberedCmd {
	return c.bypassCh
}
func (c *CmdBuffer) InCh() chan NumberedCmd {
	return c.inCh
}

func (c *CmdBuffer) Init(out chan NumberedCmd, StopCh chan bool,
	lastSeqn int64, bufSize int, pathToLog string,
	verbose int, disabled bool) {
	c.q = pairQ{}
	c.inCh = make(chan NumberedCmd, bufSize)     // buffered
	c.bypassCh = make(chan NumberedCmd, bufSize) // buffered
	c.OutCh = out                                // should also be buffered
	c.StopCh = StopCh
	c.lastSeqn = lastSeqn
	c.pathToLog = pathToLog
	c.verbose = verbose
	c.disabled = disabled
	if c.verbose > 0 {
		fmt.Println("\n* CmdBuffer initialized *")
	}
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
	if c.verbose > 0 {
		fmt.Printf("RECEIVED PAIR %v\n", seqN)
	}
	// DEBUG
	if c.disabled {
		c.OutCh <- inPair
		return true
	}
	// END
	if seqN <= c.lastSeqn { // already sent, so discard
		if c.verbose == 1 {
			fmt.Printf("    ALREADY SEEN, DISCARDING\n")
		} else if c.verbose > 1 {
			fmt.Printf("    %s ALREADY SEEN, DISCARDING\n", inPair.Cmd)
		}
		return true // get next command
	} else if seqN == c.lastSeqn+1 {
		c.OutCh <- inPair
		if c.b != nil {
			c.b.copyAndLog(seqN, inPair.Cmd)
		}
		c.lastSeqn += 1
		if c.verbose > 0 {
			fmt.Printf("    SEQN %v MATCHED LAST + 1, SENDING\n", seqN)
		}
		for c.q.Len() > 0 {
			first := c.q[0]
			if first.pair.Seqn <= c.lastSeqn {
				if c.verbose > 0 {
					fmt.Printf("        Q: DISCARDING %v, DUPE\n",
						first.pair.Seqn)
				}
				// a duplicate, so discard
				_ = heap.Pop(&c.q).(*cmdPlus)
			} else if first.pair.Seqn == c.lastSeqn+1 {
				pp := heap.Pop(&c.q).(*cmdPlus)
				c.OutCh <- *pp.pair
				if c.b != nil {
					c.b.copyAndLog(pp.pair.Seqn, pp.pair.Cmd)
				}
				c.lastSeqn += 1
				if c.verbose > 0 {
					fmt.Printf("        Q: SENT %v\n", c.lastSeqn)
				}
			} else {
				if c.verbose > 0 {
					fmt.Printf("        Q: LEAVING %v IN Q\n",
						first.pair.Seqn)
				}
				break
			}
		}
	} else {
		// seqN > c.lastSeqn + 1, so buffer
		if c.verbose > 0 {
			fmt.Printf("    HIGH SEQN %v, SO BUFFERING\n", seqN)
		}
		pp := &cmdPlus{pair: &inPair}
		heap.Push(&c.q, pp)
	}
	return true // get next command
}

func (c *CmdBuffer) Run() (err error) {
	if c.pathToLog != "" {
		// fmt.Println("    we have a log file")
		parts := strings.Split(c.pathToLog, "/")
		n := len(parts)
		if n > 1 {
			pathToDir := strings.Join(parts[:n-1], "/")
			if err := os.MkdirAll(pathToDir, 0775); err != nil {
				panic(err)
			}
		}
		// fmt.Printf("    opening %s\n", c.pathToLog)
		c.fd, err = os.Create(c.pathToLog)
		if err != nil {
			fmt.Printf("    * ERROR * opening log file %v\n", err)
			return
		}
		defer c.fd.Close()
		// fmt.Println("    initializing log buffer")
		var buf logBuffer
		c.b = &buf
		c.b.init(c.fd)
	}
	c.sy.Lock()
	c.running = true
	c.sy.Unlock()
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
			if c.verbose > 0 {
				fmt.Println("c.running has been set to false")
			}
		}
	}
	return
}

type logBuffer struct {
	fd         *os.File
	buffer     []byte
	c          *sync.Cond
	begin, end int
}

func (b *logBuffer) init(fd *os.File) {
	var rwM sync.RWMutex
	b.fd = fd
	b.buffer = make([]byte, LOG_BUFFER_SIZE)
	b.c = sync.NewCond(&rwM)
	// fmt.Println("* buffer initialized *")
	go b.writeToDisk() // sets up writes
}

func (b *logBuffer) copyAndLog(seqn int64, cmd string) {
	txt := fmt.Sprintf("%v %s\n", seqn, cmd)
	size := len(txt)
	b.c.L.Lock() // LOCK ////////////////////////////////////////////
	var from = b.end
	var to = b.end + size
	if to <= LOG_BUFFER_SIZE {
		b.end = to
		// the normal case ////////////////////////////////
		//            dest     <--     src
		count := copy(b.buffer[from:], txt)
		if count != size {
			ohMy := fmt.Sprintf("tried to copy %d bytes but only %d copied",
				size, count)
			panic(ohMy)
		}
		b.c.L.Unlock() // UNLOCK ////////////////////////////////////
		b.c.Signal()
	} else {
		// block while writing
		// (a) flush the buffer while holding lock
		var count int
		count, _ = b.fd.Write(b.buffer[b.begin:b.end])
		// XXX need to check for errors, check count for consistency

		// (b) write the string that would have overflowed
		count, _ = b.fd.Write([]byte(txt))
		// XXX need to check for errors, check count for consistency
		_ = count

		b.begin = 0
		b.end = 0

		b.c.L.Unlock() // UNLOCK ////////////////////////////////////
	}
}

// This runs in a separate goroutine and starts when b.c.Signal is
// invoked.
func (b *logBuffer) writeToDisk() {
	for {
		b.c.L.Lock()
		b.c.Wait()
		for b.begin >= b.end { // XXX WON"T HANDLE CIRCULAR BUFFER
			b.c.Wait()
		}
		// we have the lock again
		from := b.begin
		to := b.end
		b.begin = to
		b.c.L.Unlock()

		// this will block
		if to > from {
			count, _ := b.fd.Write(b.buffer[from:to])
			// XXX need to check for errors, check count for consistency
			_ = count
		}

	}
}
func (b *logBuffer) FlushLog() {

}
