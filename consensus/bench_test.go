package consensus

import (
	"fmt"
	xc "github.com/jddixon/xlattice_go"
	"testing"
	"time"
)

func makeSomeCommands(n int) []*NumberedCmd {
	rng := xc.MakeSimpleRNG()
	s := make([]*NumberedCmd, n)
	for i := 0; i < n; i++ {
		size := 64 + rng.Intn(64)
		cmd := make([]byte, size)
		rng.NextBytes(&cmd)
		s[i] = &NumberedCmd{int64(i), string(cmd)}
	}
	return s
}
func sink(from chan NumberedCmd) {
	// should use range and close the channel to shut down
	for {
		select {
		case _, ok := <-from:
			if !ok {
				return
			}
		}
	}
}

func doBenchmarkCmdBuffer(b *testing.B, pathToLog string, verbosity int) {
	NCmds := makeSomeCommands(b.N)
	out := make(chan NumberedCmd)
	go sink(out)
	var buf CmdBuffer
	p := &buf
	// 4 is bufSize, "" means no log
	stopCh := p.Init(out, 0, 4, pathToLog, verbosity, false)
	if verbosity > 1 {
		fmt.Println("initialization complete")
	}
	go p.Run()
	if verbosity > 1 {
		fmt.Println("goroutine started ...")
	}
	for !p.Running() {
		time.Sleep(10 * time.Millisecond)
		if verbosity > 0 {
			fmt.Print(".")
		}
	}
	if verbosity > 1 {
		fmt.Println("\ncalling ResetTimer") //
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.InCh() <- *NCmds[i]
	}
	if verbosity > 1 {
		fmt.Println("sending to stopCh") // DEBUG
	}
	stopCh <- true
	close(out) // should kill sink()
}

// without log to disk, 2 million ops, 942 ns/op
func BenchmarkCmdBufferWithoutLog(b *testing.B) {
	// About 1025 - 1100 ns/op over 1 million ops.
	doBenchmarkCmdBuffer(b, "", 0) // 0 => not verbose
}

// then do it with a log - benchmark appears to hang
func BenchmarkCmdBufferWithLog(b *testing.B) {
	// about 4100 ns/op over 500,000 ops
	doBenchmarkCmdBuffer(b, "tmp/benchMark.log", 0)
}
