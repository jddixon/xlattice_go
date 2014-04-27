package transport

// xlattice_go/transport/mock_connection_test.go

import (
	"bytes"
	"fmt"
	xr "github.com/jddixon/xlattice_go/rnglib"
	. "gopkg.in/check.v1"
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

	// testing whether the slice pointers are equal
	c.Assert(clientCnx.a2bMsg, Equals, serverCnx.b2aMsg)
	c.Assert(clientCnx.b2aMsg, Equals, serverCnx.a2bMsg)

	msg1Len := 32 + rng.Intn(32)
	msg2Len := 32 + rng.Intn(32)
	msg3Len := 32 + rng.Intn(32)

	msg1 := make([]byte, msg1Len)
	msg2 := make([]byte, msg2Len)
	msg3 := make([]byte, msg3Len)

	// the client writes three messages -----------------------------
	count, err := clientCnx.Write(msg1)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, msg1Len)

	count, err = clientCnx.Write(msg2)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, msg2Len)

	count, err = clientCnx.Write(msg3)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, msg3Len)

	// the server reads the three messages --------------------------
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

	// the server echoes the three messages back -----------------------
	count, err = serverCnx.Write(sBuf1)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, msg1Len)

	count, err = serverCnx.Write(sBuf2)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, msg2Len)

	count, err = serverCnx.Write(sBuf3)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, msg3Len)

	// the client reads back the messages and compares with original
	cBuf1 := make([]byte, msg1Len)
	cBuf2 := make([]byte, msg2Len)
	cBuf3 := make([]byte, msg3Len)

	count, err = clientCnx.Read(cBuf1)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, msg1Len)
	bytes.Equal(msg1, cBuf1)

	count, err = clientCnx.Read(cBuf2)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, msg2Len)
	bytes.Equal(msg2, cBuf2)

	count, err = clientCnx.Read(cBuf3)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, msg3Len)
	bytes.Equal(msg3, cBuf2)

}
