package datakeyed

// xlattice_go/overlay/datakeyed/memCache_test.go

import (
	"crypto/sha1"
	"fmt"
	xi "github.com/jddixon/xlattice_go/nodeID"
	xr "github.com/jddixon/xlattice_go/rnglib"
	. "gopkg.in/check.v1"
)

var _ = fmt.Print

func (s *XLSuite) TestSomeInserts(c *C) {

	mc, err := NewMemCache(uint64(1024*1024), uint(1024))
	c.Assert(err, IsNil)
	c.Assert(mc, NotNil)
	c.Assert(mc.ItemCount(), Equals, uint(0))
	c.Assert(mc.ByteCount(), Equals, uint64(0))

	rng := xr.MakeSimpleRNG()
	count := uint(16 + rng.Intn(17))
	values := make([][]byte, count)
	hashes := make([][]byte, count)
	ids := make([]*xi.NodeID, count)

	for i := uint(0); i < count; i++ {
		size := 256 + rng.Intn(257)
		values[i] = make([]byte, size)
		rng.NextBytes(values[i])
		d := sha1.New()
		d.Write(values[i])
		hashes[i] = d.Sum(nil)
		id, err := xi.NewNodeID(hashes[i])
		c.Assert(err, IsNil)
		ids[i] = id
	}

	var totalBytes uint64
	var totalItems uint
	for i := uint(0); i < count; i++ {
		err = mc.Add(ids[i], values[i])
		c.Assert(err, IsNil)
		totalItems++
		totalBytes += uint64(len(values[i]))

		c.Assert(mc.ItemCount(), Equals, totalItems)
		c.Assert(mc.ByteCount(), Equals, totalBytes)
	}
	// test idempotence
	for i := uint(0); i < count; i++ {
		err = mc.Add(ids[i], values[i])
		c.Assert(err, IsNil)

		c.Assert(mc.ItemCount(), Equals, totalItems)
		c.Assert(mc.ByteCount(), Equals, totalBytes)
	}

}
