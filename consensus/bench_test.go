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
		case nc := <-from:
			_ = nc
		}
	}
}

func doBenchmarkCmdBuffer(b *testing.B, pathToLog string) {
	NCmds := makeSomeCommands(b.N)
	out := make(chan NumberedCmd)
	go sink(out)
	stopCh := make(chan bool)
	var buf CmdBuffer
	p := &buf
	p.Init(out, stopCh, 0, 4, pathToLog) // 4 is bufSize, "" means no log
	go p.Run()
	for !p.Running() {
		time.Sleep(10 * time.Millisecond)
		fmt.Print(".")
	}
	fmt.Println()
	b.ResetTimer()
	// typically 28.4 ns/op
	for i := 0; i < b.N; i++ {
		p.InCh() <- *NCmds[i]
	}
	stopCh <- true
	close(out)
}

// without log to disk, 2 million ops, 942 ns/op
func BenchmarkCmdBufferWithoutLog(b *testing.B) {
	fmt.Println("BENCHMARK WITHOUT LOGGING")
	doBenchmarkCmdBuffer(b, "")
}

// then do it with a log
func BenchmarkCmdBufferWithLog(b *testing.B) {
	fmt.Println("BENCHMARK *WITH* LOGGING")
	doBenchmarkCmdBuffer(b, "tmp/benchMark.log")
}
