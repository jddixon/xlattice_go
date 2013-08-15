package node

// xlattice_go/node/localHost_test.go

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	xo "github.com/jddixon/xlattice_go/overlay"
	"github.com/jddixon/xlattice_go/rnglib"
	xt "github.com/jddixon/xlattice_go/transport"
	"io/ioutil"
	. "launchpad.net/gocheck"
	"os"
	"path"
	"time"
)

var _ = fmt.Print
var _ = xo.NewIPOverlay

const (
	MIN_LEN = 1024
	MAX_LEN = 2048
	// SHA1_LEN	= 20		// declared elsewhere
	// SHA3_LEN	= 32
	Q = 64 // "too many open files" if 64
)

var (
	ANY_END_POINT, _ = xt.NewTcpEndPoint("127.0.0.1:0")
)

// See cluster_test.go for a general description of these tests.
//
// This test involves nodes executing on a single machine, with accessor
// IP addresses 127.0.0.1:P, where P represents a system-assigned unique
// port number.

// Accept connections from peers until a message is received on stopCh.
// For each message received from a peer, calculate its SHA1 hash,
// send that as a reply, and close the connection.  Send on stoppedCh
// when all replies have been sent.
func (s *XLSuite) nodeAsServer(c *C, node *Node, stopCh, stoppedCh chan bool) {
	acceptor := node.acceptors[0]
	go func() {
		for {
			conn, err := acceptor.Accept()
			if err != nil {
				break
			}
			if conn != nil {
				cnx := conn.(*xt.TcpConnection)
				go func() {
					defer cnx.Close()
					buf := make([]byte, MAX_LEN)
					count, err := cnx.Read(buf)
					if err == nil {
						buf = buf[:count]
						// calculate hash
						dig := sha1.New()
						dig.Write(buf)
						hash := dig.Sum(nil)
						count, err = cnx.Write(hash)
						_, _ = count, err

						time.Sleep(time.Millisecond)
					}
				}()
			}
		}
	}()

	<-stopCh
	acceptor.Close()
	stoppedCh <- true
}

// Send q messages to each peer, expecting to receive an SHA1 hash
// back.  When all are received and verified, send on doneCh.

func (s *XLSuite) nodeAsClient(c *C, node *Node, q int, doneCh chan bool) {
	P := node.SizePeers()

	for i := 0; i < q; i++ {
		for j := 0; j < P; j++ {
			go func(j int) {
				var err error
				var count int
				rng := rnglib.MakeSimpleRNG()

				// open cnx to peer j
				peer := node.GetPeer(j)
				ctor := peer.GetConnector(0)
				cnx, err := ctor.Connect(ANY_END_POINT)
				c.Assert(err, Equals, nil)
				c.Assert(cnx, Not(IsNil))
				defer cnx.Close()
				tcpCnx := cnx.(*xt.TcpConnection)

				// make msg, random length, random content
				msgLen := 1024 + rng.Intn(1024)
				buf := make([]byte, msgLen)
				rng.NextBytes(&buf)

				// send msg
				count, err = tcpCnx.Write(buf)
				c.Assert(err, IsNil)
				c.Assert(count, Equals, msgLen)

				// calculate hash
				dig := sha1.New()
				dig.Write(buf)
				hash := dig.Sum(nil)
				hashBuf := make([]byte, SHA1_LEN)

				// wait for reply
				count, err = tcpCnx.Read(hashBuf)
				c.Assert(err, IsNil)
				c.Assert(count, Equals, SHA1_LEN)

				// complain if hash differs from reply
				hashOut := hex.EncodeToString(hash)
				hashIn := hex.EncodeToString(hashBuf)
				c.Assert(hashOut, Equals, hashIn)
			}(j)
		}
	}
	doneCh <- true

}

func (s *XLSuite) TestLocalHostTcpCluster(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_LOCAL_HOST_TCP_CLUSTER")
	}
	var err error
	const K = 5
	var doneCh, stopCh, stoppedCh []chan (bool)
	for i := 0; i < K; i++ {
		doneCh = append(doneCh, make(chan bool))
		stopCh = append(stopCh, make(chan bool))
		stoppedCh = append(stoppedCh, make(chan bool))
	}
	nodes, accs := MockLocalHostCluster(K)
	defer func() {
		for i := 0; i < K; i++ {
			if accs[i] != nil {
				accs[i].Close()
			}
		}
	}()

	// AT THIS POINT we have K nodes, each with K-1 peers.
	// Save the configurations
	pathsToCfg := make([]string, K)
	for i := 0; i < K; i++ {
		hexNodeID := hex.EncodeToString(nodes[i].GetNodeID().Value())
		pathsToCfg[i] = path.Join("tmp", hexNodeID, ".xlattice")

		err = os.MkdirAll(pathsToCfg[i], 0755)
		c.Assert(err, IsNil)
		cfgFileName := path.Join(pathsToCfg[i], "config")

		fmt.Printf("WRITING CONFIG FILE %s\n", cfgFileName)

		cfg := nodes[i].String()
		err = ioutil.WriteFile(cfgFileName, []byte(cfg), 0644)
		c.Assert(err, IsNil)
	}

	// Start each node running in a separate goroutine.

	// Each node will in a somewhat randomized fashion send N messages
	// to every other node, expecting to receive back from the peer a
	// digital signature for the message.  As each response = digital
	// signature comes back it is validated.  When all messages have
	// been validated, the node sends a 'done' message on a boolean
	// channel to the supervisor.
	for i := 0; i < K; i++ {
		go s.nodeAsClient(c, nodes[i], Q, doneCh[i])
		go s.nodeAsServer(c, nodes[i], stopCh[i], stoppedCh[i])
	}
	// When all nodes have signaled that they are done, the supervisor
	// sends on stopCh, the stop command channel.
	for i := 0; i < K; i++ {
		<-doneCh[i]
	}
	for i := 0; i < K; i++ {
		stopCh[i] <- true
	}

	// Each node will send a reply to the supervisor on stoppedCh.
	// and then terminate.
	for i := 0; i < K; i++ {
		<-stoppedCh[i]
	}

	// When the supervisor has received stopped signals from all nodes,
	// it summarize results and terminates.
	// XXX STUB XXX

}
