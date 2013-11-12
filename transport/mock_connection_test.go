package transport

// xlattice_go/transport/mock_connection_test.go

import (
	"bytes"
	"fmt"
	xr "github.com/jddixon/xlattice_go/rnglib"
	. "launchpad.net/gocheck"
	"reflect"
	"unsafe"
)

var _ = fmt.Print

func (s *XLSuite) TestMockConnection(c *C) {

	var (
		err error
	)

	rng := xr.MakeSimpleRNG()

	aEnd := NewMockEndPoint("T", "A").(*MockEndPoint)
	c.Assert(aEnd, NotNil)
	bEnd := NewMockEndPoint("T", "B").(*MockEndPoint)
	c.Assert(bEnd, NotNil)

	clientCnx, err := NewMockConnection(aEnd, bEnd)
	c.Assert(err, IsNil)

	serverCnx, err := NewReverseMockConnection(clientCnx)
	c.Assert(err, IsNil)

	c.Assert(clientCnx.State, Equals, CNX_CONNECTED)
	c.Assert(serverCnx.State, Equals, clientCnx.State)
	c.Assert(clientCnx.NearEnd, Equals, serverCnx.FarEnd)
	c.Assert(clientCnx.FarEnd, Equals, serverCnx.NearEnd)

	cA2BHdr := (*reflect.SliceHeader)(unsafe.Pointer(&clientCnx.a2bMsg))
	cB2AHdr := (*reflect.SliceHeader)(unsafe.Pointer(&clientCnx.b2aMsg))
	sA2BHdr := (*reflect.SliceHeader)(unsafe.Pointer(&serverCnx.a2bMsg))
	sB2AHdr := (*reflect.SliceHeader)(unsafe.Pointer(&serverCnx.b2aMsg))

	// XXX This doesn't work because the arrays are empty.  That is, the
	// tests succeed but are misleading.
	c.Assert(cA2BHdr.Data, Equals, sB2AHdr.Data)
	c.Assert(cB2AHdr.Data, Equals, sA2BHdr.Data)

	msg1Len := 32 + rng.Intn(32)
	msg2Len := 32 + rng.Intn(32)
	msg3Len := 32 + rng.Intn(32)

	msg1 := make([]byte, msg1Len)
	msg2 := make([]byte, msg2Len)
	msg3 := make([]byte, msg3Len)

	count, err := clientCnx.Write(msg1)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, msg1Len)

	count, err = clientCnx.Write(msg2)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, msg2Len)

	count, err = clientCnx.Write(msg3)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, msg3Len)

	sBuf1 := make([]byte, msg1Len)
	sBuf2 := make([]byte, msg2Len)
	sBuf3 := make([]byte, msg3Len)

	count, err = serverCnx.Read(sBuf1)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, msg1Len)
	bytes.Equal(msg1, sBuf1)

	count, err = serverCnx.Read(sBuf2)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, msg2Len)
	bytes.Equal(msg2, sBuf2)

	count, err = serverCnx.Read(sBuf3)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, msg3Len)
	bytes.Equal(msg3, sBuf3)

	_, _ = err, rng // DEBUG
	_, _ = clientCnx, serverCnx
}
