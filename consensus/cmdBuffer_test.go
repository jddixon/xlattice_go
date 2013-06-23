package consensus

import (
	"container/heap"
	"fmt"
	"github.com/jddixon/xlattice_go/rnglib"
	. "launchpad.net/gocheck"
	"testing"
	"time"
)

// gocheck tie-in /////////////////////
func Test(t *testing.T) { TestingT(t) }

type XLSuite struct{}

var _ = Suite(&XLSuite{})

// end gocheck setup //////////////////

func (s *XLSuite) makeSimpleRNG() *rnglib.PRNG {
	t := time.Now().Unix()
	rng := rnglib.NewSimpleRNG(t)
	return rng
}

func (s *XLSuite) TestCmdQ(c *C) {
	q := pairQ{}
	heap.Init(&q)
	c.Assert(q.Len(), Equals, 0)

	pair0 := CmdPair{Seqn: 42, Cmd: "foo"}
	pair1 := CmdPair{Seqn: 1, Cmd: "bar"}
	pair2 := CmdPair{Seqn: 99, Cmd: "baz"}

	pp0 := pairPlus{pair: &pair0}
	pp1 := pairPlus{pair: &pair1}
	pp2 := pairPlus{pair: &pair2}

	heap.Push(&q, &pp0)
	heap.Push(&q, &pp1)
	heap.Push(&q, &pp2)
	c.Assert(q.Len(), Equals, 3)

	out := heap.Pop(&q).(*pairPlus)
	c.Assert(out.pair.Seqn, Equals, int64(1))
	c.Assert(out.pair.Cmd, Equals, "bar")

	out = heap.Pop(&q).(*pairPlus)
	c.Assert(out.pair.Seqn, Equals, int64(42))
	c.Assert(out.pair.Cmd, Equals, "foo")

	out = heap.Pop(&q).(*pairPlus)
	c.Assert(out.pair.Seqn, Equals, int64(99))
	c.Assert(out.pair.Cmd, Equals, "baz")

	c.Assert(q.Len(), Equals, 0)
	// XXX THIS PANICS - so if popping from a heap, always check
	// the length first.
	//zzz		:= heap.Pop(&q)
	// c.Assert(zzz, Equals, nil)
}
func (s *XLSuite) TestCmdBuffer(c *C) {
	var pairMap = map[int64]string{
		1: "foo",
		2: "bar",
		3: "baz",
		4: "it's me!",
		5: "my chance will come soon",
		6: "it's my turn now",
		7: "wait for me",
	}
	order := [...]int{1, 2, 3, 6, 6, 5, 4, 1, 7}
	var buf CmdBuffer
	p := &buf
	var out = make(chan CmdPair, len(order)+1) // must exceed len(order)
	var stopCh = make(chan bool, 1)
	p.Init(out, stopCh, 0)               // was Init()
	c.Assert(p.Running(), Equals, false) // HACK - was false

	go p.Run()

	for !p.Running() {
		fmt.Print(".")
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println("")
	c.Assert(p.Running(), Equals, true)

	fmt.Println("sending pairs")
	for n := 0; n < len(order); n++ {
		cmd := pairMap[int64(n)]
		pair := CmdPair{Seqn: int64(order[n]), Cmd: cmd}
		p.InCh <- pair
	}

	time.Sleep(time.Millisecond)
	fmt.Println("collecting results")
	var results [7]CmdPair
	for n := 0; n < 7; n++ {
		results[n] = <-out
		c.Assert(results[n].Seqn, Equals, int64(n+1))
		fmt.Printf("GOT RESULT %v\n", n+1)
	}

	time.Sleep(time.Millisecond)
	stopCh <- true
	fmt.Println("stop has been sent")
}
