package transport

// xlattice_go/transport/server_test.go

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/jddixon/xlattice_go/rnglib"
	. "launchpad.net/gocheck"
	// "regexp"
	"time"
)

// Start an acceptor running on a random port.  Create K*N blocks of
// random data.  These will be sent by K clients to the server.  The
// server will reply to each with the SHA1 hash of the block.  Run
// the clients in separate goroutines.  After all clients have sent
// all of their messages, verify that the hashes received back are
// correct.

const (
	// XXX 2013-07-20 test hangs if K=16,N=32 and K increasd to 32 OR
	// N increased from 32 to 64
	K        = 32   // number of clients
	N        = 16   // number of messages for each client
	MIN_LEN  = 1024 // minimum length of message
	MAX_LEN  = 2048 // maximum
	SHA1_LEN = 20
)

var rng = rnglib.MakeSimpleRNG()

func (s *XLSuite) handleMsg(cnx *TcpConnection) error {
	defer cnx.Close()
	buf := make([]byte, MAX_LEN)

	// read the message
	count, err := cnx.Read(buf)
	buf = buf[:count] // ESSENTIAL
	if err == nil {
		// calculate its hash
		d := sha1.New()
		d.Write(buf)
		digest := d.Sum(nil) // a binary value

		// send the digest as a reply
		count, err = cnx.Write(digest)

		_ = count // XXX verify length of 20
	}
	// XXX allow the other end to read the reply; it would be
	// better to loop until a 'closed connection' error is returned
	time.Sleep(100 * time.Millisecond)
	return err
}

func (s *XLSuite) TestHashingServer(c *C) {
	ANY_END_POINT, _ := NewTcpEndPoint("127.0.0.1:0")
	SERVER_ADDR := "127.0.0.1:0"

	// -- setup  -----------------------------------------------------
	fmt.Println("building messages")
	var messages [][][]byte = make([][][]byte, K)
	var hashes [][][]byte = make([][][]byte, K)
	for i := 0; i < K; i++ {
		messages[i] = make([][]byte, N)
		for j := 0; j < N; j++ {
			msgLen := MIN_LEN + rng.Intn(MAX_LEN-MIN_LEN)
			messages[i][j] = make([]byte, msgLen)
			rng.NextBytes(&messages[i][j])
		}
		hashes[i] = make([][]byte, N)
		for j := 0; j < N; j++ {
			hashes[i][j] = make([]byte, SHA1_LEN)
		}
	}

	// -- create and start server -----------------------------------
	acc, err := NewTcpAcceptor(SERVER_ADDR)
	c.Assert(err, Equals, nil)
	defer acc.Close()
	accEndPoint := acc.GetEndPoint()
	fmt.Printf("server_test acceptor listening on %s\n", accEndPoint.String())
	go func() {
		for {
			cnx, err := acc.Accept()
			if err != nil { // ESSENTIAL
				break
			}
			if cnx != nil {
				go func(cnx *TcpConnection) {
					_ = s.handleMsg(cnx)
					// c.Assert(err, Equals, nil)
				}(cnx)
			}
		}
	}()

	// -- create K client connectors --------------------------------
	ktors := make([]*TcpConnector, K)
	for i := 0; i < K; i++ {
		ktors[i], err = NewTcpConnector(accEndPoint)
		c.Assert(err, Equals, nil)
	}

	// -- start the clients -----------------------------------------
	var clientDone [K]chan bool
	for i := 0; i < K; i++ {
		clientDone[i] = make(chan bool)
	}
	for i := 0; i < K; i++ {
		go func(i int) {
			for j := 0; j < N; j++ {
				// the client sends N messages, expecting an SHA1 back
				var count int
				cnx, err := ktors[i].Connect(ANY_END_POINT)
				c.Assert(err, Equals, nil)
				tcpCnx := cnx.(*TcpConnection)
				count, err = tcpCnx.Write(messages[i][j])
				if err != nil {
					fmt.Printf("error writing [%d][%d]: %v\n", i, j, err)
				}
				count, err = tcpCnx.Read(hashes[i][j])
				if err != nil {
					fmt.Printf("error reading [%d][%d]: %v\n", i, j, err)
				}
				cnx.Close()

				_, _ = count, err // XXX
			}
			clientDone[i] <- true
		}(i)
	}
	// -- when all clients have completed, shut down server ---------
	for i := 0; i < K; i++ {
		<-clientDone[i]
	}
	if acc != nil {
		if err = acc.Close(); err != nil {
			fmt.Printf("unexpected error closing acceptor: %v\n", err)
		} else {
			fmt.Printf("acceptor closed successfully\n")
		}
	}
	// -- calculate and verify K*N hashes ---------------------------
	for i := 0; i < K; i++ {
		for j := 0; j < N; j++ {
			d := sha1.New()
			d.Write(messages[i][j])
			digest := d.Sum(nil)                // a binary value
			hashX := hex.EncodeToString(digest) // DEBUG
			hashY := hex.EncodeToString(hashes[i][j])
			c.Assert(hashX, Equals, hashY)
		}
	}
}
