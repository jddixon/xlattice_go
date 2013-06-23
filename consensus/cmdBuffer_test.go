package consensus

import (
	"container/heap"
	// "fmt"
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
