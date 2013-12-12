package chunks

// xlattice_go/protocol/chunks/chunklist_test.go

import (
	"bytes"
	"code.google.com/p/go.crypto/sha3"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	//xi "github.com/jddixon/xlattice_go/nodeID"
	xr "github.com/jddixon/xlattice_go/rnglib"
	. "launchpad.net/gocheck"
)

var _ = fmt.Print

func (s *XLSuite) TestChunkList(c *C) {
	rng := xr.MakeSimpleRNG()

	dataLen := 1 + rng.Intn(3*MAX_CHUNK_BYTES)
	data := make([]byte, dataLen)
	reader := bytes.NewReader(data)
	d := sha3.NewKeccak256()
	d.Write(data)
	hash := d.Sum(nil)
	skPriv, err := rsa.GenerateKey(rand.Reader, 1024) // cheap key

	sk := &skPriv.PublicKey
	title := rng.NextFileName(8)
	timestamp := rng.Int63()

	cl, err := NewChunkList(sk, title, timestamp, reader, int64(dataLen), hash)
	c.Assert(err, IsNil)
	c.Assert(cl, NotNil)

	_, _ = data, reader
}
