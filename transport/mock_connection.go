package transport

import (
	"bytes"
	"fmt"
	xc "github.com/jddixon/xlattice_go/crypto"
	"sync"
)

type MockConnection struct {
	State           int
	NearEnd, FarEnd *MockEndPoint
	a2bMsg, b2aMsg  [][]byte
	a2bMu, b2aMu    sync.Mutex
}

func NewNewMockConnection() (cnx *MockConnection, err error) {
	cnx = &MockConnection{
		State: CNX_UNBOUND,
	}
	return
}
func NewMockConnection(nearEnd, farEnd *MockEndPoint) (
	cnx *MockConnection, err error) {

	if nearEnd == nil || farEnd == nil {
		err = NilEndPoint
	} else {
		cnx = &MockConnection{
			NearEnd: nearEnd,
			FarEnd:  farEnd,
			State:   CNX_CONNECTED,
		}
	}
	return
}

// Returns a view from the other end of a MockConnection.  Given a client
// connection, this creates the same connection as seen by the server.

func NewReverseMockConnection(orig *MockConnection) (
	cnx *MockConnection, err error) {

	if orig == nil {
		err = NilConnection
	} else {
		cnx = &MockConnection{
			State:   orig.State,
			NearEnd: orig.FarEnd,
			FarEnd:  orig.NearEnd,
			a2bMsg:  orig.b2aMsg,
			b2aMsg:  orig.a2bMsg,
			a2bMu:   orig.b2aMu,
			b2aMu:   orig.a2bMu,
		}
	}
	return
}

// Return the current state index.
func (c *MockConnection) GetState() int {
	return c.State
}

// Set the near end point of a connection.  If either the
// near or far end point has already been set, this will
// return an error.  If successful, the connection's
// state becomes CNX_BOUND.
//
func (c *MockConnection) BindNearEnd(e EndPointI) (err error) {
	if c.State != CNX_UNBOUND {
		err = AlreadyBound
	} else {
		switch v := e.(type) {
		case *MockEndPoint:
			c.NearEnd = v
			c.State = CNX_BOUND
		default:
			err = NotAMockEndPoint
		}
	}
	return
}

// Set the far end point of a connection.  If the near end
// point has NOT been set or if the far end point has already
// been set -- in other words, if the connection is already
// beyond state BOUND -- this will cause an exception.
// If the operation is successful, the connection's state
// becomes either PENDING or CONNECTED.
//
// XXX The state should become CONNECTED if the far end is on
// XXX the same host and PENDING if it is on a remoted host.
//
func (c *MockConnection) BindFarEnd(e EndPointI) (err error) {
	if c.State == CNX_UNBOUND {
		err = NotBound
	} else if c.State > CNX_BOUND {
		err = AlreadyConnected
	} else {
		switch v := e.(type) {
		case *MockEndPoint:
			c.FarEnd = v
			c.State = CNX_CONNECTED
		default:
			err = NotAMockEndPoint
		}
	}
	return
}

// Bring the connection to the DISCONNECTED state.
//
// XXX This code allows you to close an UNBOUND or BOUND connection.
//
func (c *MockConnection) Close() (err error) {
	c.State = CNX_DISCONNECTED
	return
}

func (c *MockConnection) GetNearEnd() (ep EndPointI) {
	return c.NearEnd
}

// XXX 2013-07-20: this returns the near end instead !
func (c *MockConnection) GetFarEnd() (ep EndPointI) {
	return c.FarEnd
}

// Read from the connection.  In this implementation we have a queue of
// incoming messages, each a byte slice.  If it will fit, we read all of
// the first message into the output buffer b.  Otherwise, we read what
// will fit and leave the rest of the first message on the queue.
//
func (c *MockConnection) Read(b []byte) (count int, err error) {

	if len(c.b2aMsg) == 0 {
		// DEBUG
		fmt.Printf("Read: %d in b2a (but %d in a2b)\n",
			len(c.b2aMsg), len(c.a2bMsg))
		// END
		count = 0
	} else {
		lenB := len(b) // how many bytes we can return
		lenMsg := len(c.b2aMsg[0])
		if lenB >= lenMsg {
			buf := bytes.NewBuffer(b)
			c.b2aMu.Lock()
			count, err = buf.Write(c.b2aMsg[0])
			// DEBUG
			fmt.Printf("Read: buf.Write returns count = %d, len msg is %d\n",
				count, lenMsg)
			// END
			c.b2aMsg = c.b2aMsg[1:]
			c.b2aMu.Unlock()
		} else {
			// send what we can
			buf := bytes.NewBuffer(b)
			c.b2aMu.Lock()
			count, err = buf.Write(c.b2aMsg[0][0:lenB])
			c.b2aMsg[0] = c.b2aMsg[0][count:]
			c.b2aMu.Unlock()
		}
	}
	return
}

// Write msg b to the connection.  In this implementation we maintain
// a queue of output messages.  We will simply append this message to
// that queue, making no change to the message.
//
func (c *MockConnection) Write(b []byte) (count int, err error) {
	count = len(b)
	c.a2bMu.Lock()
	c.a2bMsg = append(c.a2bMsg, b)
	// DEBUG
	fmt.Printf("Write: after writing %d byte msg, %d msgs in buffer\n",
		count, len(c.a2bMsg))
	// END
	c.a2bMu.Unlock()
	return
}
func (c *MockConnection) IsBlocking() bool {
	// XXX STUB NotImplemented
	return false
}

// ///////////////////////////////////////////////////////////////////
// XXX CONFUSION BETWEEN PACKET vs STREAM AND BLOCKING vs NON-BLOCKING
// ///////////////////////////////////////////////////////////////////
// non-blocking

// blocking
//  GetInputStream(i *InputStream, e error)     // throws IOException
//  GetOutputStream(o *OutputStream, e error)   // throws IOException

// @return whether the connection is encrypted//
func (c *MockConnection) IsEncrypted() bool {
	// XXX STUB NotImplemented
	return false
}

//
// (Re)negotiate the Secret used to encrypt traffic over the
// connection.
//
// @param myKey  this Node's asymmetric key
// @param hisKey Peer's public key
//
func (c *MockConnection) Negotiate(myKey xc.KeyI, hisKey xc.PublicKeyI) (s xc.SecretI, e error) {
	// XXX STUB
	return nil, NotImplemented
}

func (c *MockConnection) Equal(any interface{}) bool {
	// XXX STUB NotImplemented
	return false
}

func (c *MockConnection) String() string {
	return fmt.Sprintf("Mock: %s --> %s",
		c.GetNearEnd().String(),
		c.GetFarEnd().String())
}
