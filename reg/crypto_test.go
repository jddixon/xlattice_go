package reg

// xlattice_go/msg/reg_test.go

import (
	"crypto/aes"
	"fmt"
	//xc "github.com/jddixon/xlattice_go/crypto"
	xn "github.com/jddixon/xlattice_go/node"
	xi "github.com/jddixon/xlattice_go/nodeID"
	xr "github.com/jddixon/xlattice_go/rnglib"
	. "launchpad.net/gocheck"
)

func (s *XLSuite) TestCrytpo(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_CRYPTO")
	}

	rng := xr.MakeSimpleRNG()
	id := make([]byte, SHA3_LEN)
	rng.NextBytes(&id)
	nodeID, err := xi.New(id)
	c.Assert(err, IsNil)
	c.Assert(nodeID, Not(IsNil))
	node, err := xn.NewNew("foo", nodeID)
	c.Assert(err, IsNil)
	c.Assert(node, Not(IsNil))

	// Generate 16-byte AES IV, 32-byte AES key, and 8-byte salt.
	// For testing purposes these need not be crypto grade.
	iv1 := make([]byte, aes.BlockSize)
	rng.NextBytes(&iv1)
	key1 := make([]byte, 2*aes.BlockSize) // so 32 bytes
	rng.NextBytes(&key1)

	iv2 := make([]byte, aes.BlockSize)
	rng.NextBytes(&iv2)
	key2 := make([]byte, 2*aes.BlockSize)
	rng.NextBytes(&key2)

	salt1 := make([]byte, 8)
	rng.NextBytes(&salt1)

	salt2 := make([]byte, 8)
	rng.NextBytes(&salt2)

	engine1a, err := aes.NewCipher(key1)
	c.Assert(err, IsNil)
	c.Assert(engine1a, Not(IsNil))

	engine1b, err := aes.NewCipher(key1)
	c.Assert(err, IsNil)
	c.Assert(engine1b, Not(IsNil))

	engine2a, err := aes.NewCipher(key2)
	c.Assert(err, IsNil)
	c.Assert(engine2a, Not(IsNil))

	engine2b, err := aes.NewCipher(key2)
	c.Assert(err, IsNil)
	c.Assert(engine2b, Not(IsNil))

	// -- HELLO -----------------------------------------------------
	// create and marshal a hello message containing AES iv+key1
	// XXX STUB XXX

	// encrypt the hello using the node's public comms key
	// XXX STUB XXX

	// decrypt the hello using the node's private comms key
	// XXX STUB XXX

	// verify that iv1, key1, and salt1 are the same
	// XXX STUB XXX

	// -- HELLO REPLY -----------------------------------------------
	// create and marshal a hello reply containing salt, salt2, iv+key2
	// XXX STUB XXX

	// encrypt the reply using engine1a = iv1, key1
	// XXX STUB XXX

	// decrypt the reply using engine1b = iv1, key1
	// XXX STUB XXX

	// verify that salt, salt2, iv2, and key2 survive the trip unchanged
	// XXX STUB XXX

	// -- JOIN ------------------------------------------------------
	// create and marshal a join containing id, ck, sk, myEnd*
	// XXX STUB XXX

	// encrypt the join using engine2a = iv2, key2
	// XXX STUB XXX

	// decrypt the join using engine2b = iv2, key2
	// XXX STUB XXX

	// verify that id, ck, sk, myEnd* survive the trip unchanged
	// XXX STUB XXX

	// -- MEMBERS ---------------------------------------------------
	// create and marshal a set of 3-5 tokens each containing attrs,
	// nodeID, clusterID
	// XXX STUB XXX

	// encrypt the join using engine2a = iv2, key2
	// XXX STUB XXX

	// decrypt the join using engine2b = iv2, key2
	// XXX STUB XXX

	// verify that id, ck, sk, myEnd* survive the trip unchanged
	// XXX STUB XXX

	// -- MEMBER LIST  ----------------------------------------------

	// LOOP ANOTHER N=16 TIMES WITH BLOCKS OF RANDOM DATA -----------
	for i := 0; i < 16; i++ {
		// create block of data
		// XXX STUB XXX

		// encrypt it with engine2a (iv2, key2)
		// XXX STUB XXX

		// decrypt with engine2b (iv2, key2)
		// XXX STUB XXX

		// verify same data
		// XXX STUB XXX

	}
}
