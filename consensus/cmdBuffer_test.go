package consensus

import (
	"container/heap"
	"fmt"
	"github.com/jddixon/xlattice_go/rnglib"
	"io/ioutil"
	. "launchpad.net/gocheck"
	"os"
	"testing"
	"time"
)

// gocheck tie-in /////////////////////
func Test(t *testing.T) { TestingT(t) }

type XLSuite struct{}

var _ = Suite(&XLSuite{})

// end gocheck setup //////////////////

const (
	TEST_PAIR_COUNT = 7
)

func (s *XLSuite) makeSimpleRNG() *rnglib.PRNG {
	t := time.Now().Unix()
	rng := rnglib.NewSimpleRNG(t)
	return rng
}

func (s *XLSuite) TestCmdQ(c *C) {
	q := pairQ{}
	heap.Init(&q)
	c.Assert(q.Len(), Equals, 0)

	pair0 := NumberedCmd{Seqn: 42, Cmd: "foo"}
	pair1 := NumberedCmd{Seqn: 1, Cmd: "bar"}
	pair2 := NumberedCmd{Seqn: 99, Cmd: "baz"}

	pp0 := cmdPlus{pair: &pair0}
	pp1 := cmdPlus{pair: &pair1}
	pp2 := cmdPlus{pair: &pair2}

	heap.Push(&q, &pp0)
	heap.Push(&q, &pp1)
	heap.Push(&q, &pp2)
	c.Assert(q.Len(), Equals, 3)

	out := heap.Pop(&q).(*cmdPlus)
	c.Assert(out.pair.Seqn, Equals, int64(1))
	c.Assert(out.pair.Cmd, Equals, "bar")

	out = heap.Pop(&q).(*cmdPlus)
	c.Assert(out.pair.Seqn, Equals, int64(42))
	c.Assert(out.pair.Cmd, Equals, "foo")

	out = heap.Pop(&q).(*cmdPlus)
	c.Assert(out.pair.Seqn, Equals, int64(99))
	c.Assert(out.pair.Cmd, Equals, "baz")

	c.Assert(q.Len(), Equals, 0)
	// XXX THIS PANICS - so if popping from a heap, always check
	// the length first.
	//zzz		:= heap.Pop(&q)
	// c.Assert(zzz, Equals, nil)
}
func (s *XLSuite) doTestCmdBufferI(c *C, p CmdBufferI, logging bool) {
	var cmdMap = map[int64]string{
		1: "foo",
		2: "bar",
		3: "baz",
		4: "it's me!",
		5: "my chance will come soon",
		6: "it's my turn now",
		7: "wait for me",
	}
	c.Assert(len(cmdMap), Equals, TEST_PAIR_COUNT)

	// we send the messages somewhat out of order, with some duplicates
	order := [...]int{1, 2, 3, 6, 3, 2, 6, 5, 4, 1, 7}
	var out = make(chan NumberedCmd, len(order)+1) // must exceed len(order)
	var logFile string
	if logging {
		logFile = "tmp/logFile"
	}
	stopCh := p.Init(out, 0, 4, logFile, 0, false) // 4 is chan bufSize
	c.Assert(p.Running(), Equals, false)

	fmt.Println("  starting p loop ...")
	// XXX Run() can return an error, which must be nil
	go p.Run()
	for !p.Running() {
		time.Sleep(time.Millisecond)
	}
	c.Assert(p.Running(), Equals, true)
	if logging {
		_, err := os.Stat(logFile) // created by Run()
		c.Assert(err, Equals, nil)
	}

	for n := 0; n < len(order); n++ {
		which := order[n]
		cmd := cmdMap[int64(which)]
		pair := NumberedCmd{Seqn: int64(which), Cmd: cmd}
		// DEBUG
		// fmt.Printf("sending %d : %s\n", order[n], cmd)
		// END
		p.InCh() <- pair
	}

	var results [7]NumberedCmd
	for n := 0; n < 7; n++ {
		results[n] = <-out
		c.Assert(results[n].Seqn, Equals, int64(n+1))
	}

	c.Assert(p.Running(), Equals, true)
	stopCh <- true
	time.Sleep(time.Microsecond)
	c.Assert(p.Running(), Equals, false)

	if logging {
		var expected string
		for i := 1; i <= TEST_PAIR_COUNT; i++ {
			n := int64(i)
			cmd := cmdMap[n]
			line := fmt.Sprintf("%d %s\n", n, cmd)
			expected += line
		}
		raw, err := ioutil.ReadFile(logFile)
		c.Assert(nil, Equals, err)
		actual := string(raw)
		c.Assert(actual, Equals, expected)
	}
} // GEEP
func (s *XLSuite) TestCmdBuffer(c *C) {
	var buf CmdBuffer
	fmt.Println("running test without logging")
	s.doTestCmdBufferI(c, &buf, false)
	fmt.Println("running test -with- logging")
	s.doTestCmdBufferI(c, &buf, true)
}

func (s *XLSuite) TestLogBufferOverflow(c *C) {
	var buf CmdBuffer
	p := &buf
	var bufSize = LOG_BUFFER_SIZE
	N := int(1.25 * float64(bufSize) / float64(5+96))
	fmt.Printf("TEST_LOG_BUFFER_OVERFLOW WITH %d RECORDS\n", N)
	rng := s.makeSimpleRNG()

	var out = make(chan NumberedCmd, 4) //
	logFile := "tmp/overflows.log"
	stopCh := p.Init(out, 0, 4, logFile, 0, false) // 4 is bufSize
	c.Assert(p.Running(), Equals, false)

	fmt.Println("  starting p loop for overflow test ...")
	// XXX Run() can return an error, which must be nil
	go p.Run()
	for !p.Running() {
		time.Sleep(time.Millisecond)
	}
	c.Assert(p.Running(), Equals, true)

	// Run() will have created the log file
	_, err := os.Stat(logFile)
	c.Assert(err, Equals, nil)

	cmds := make([]*string, N)
	results := make([]*NumberedCmd, N)
	var expectedInFile string
	for n := 0; n < N; n++ {
		seqN := int64(n + 1)
		cmdLen := 64 + rng.Intn(64)
		raw := make([]byte, cmdLen)
		rng.NextBytes(raw)
		cmd := string(raw)
		pair := NumberedCmd{Seqn: seqN, Cmd: cmd}
		p.InCh() <- pair
		txt := fmt.Sprintf("%d %s\n", n+1, cmd)
		cmds[n] = &txt
		expectedInFile += txt

		nextResult := <-out
		results[n] = &nextResult
		//       expected    ==      actual
		c.Assert(int64(n+1), Equals, nextResult.Seqn)
		c.Assert(cmd, Equals, nextResult.Cmd)

		// fmt.Printf("%d sent and received\n", n)
	}

	c.Assert(p.Running(), Equals, true)
	stopCh <- true
	time.Sleep(time.Microsecond)
	c.Assert(p.Running(), Equals, false)

	// compare actual data in log file with expected
	raw, err := ioutil.ReadFile(logFile)
	c.Assert(nil, Equals, err)
	actual := string(raw)
	c.Assert(expectedInFile, Equals, actual)

}
