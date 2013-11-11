package transport

import (
	"fmt"
	xc "github.com/jddixon/xlattice_go/crypto"
	"sync"
)

type MockConnection struct {
	State           int
	NearEnd, FarEnd *MockEndPoint
	a2bMsg, b2aMsg  [][]byte
	a2bMu, b2aMu    *sync.Mutex
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

func (c *MockConnection) Read(b []byte) (count int, err error) {

	if len(c.b2aMsg) == 0 {
		count = 0
	} else {
		maxCount := len(b) // how many bytes we can return
		lenMsg := len(c.b2aMsg[0])
		if maxCount >= lenMsg {
			// EXPERIMENT
			// XXX STUB
		}

	}
	return
}

// This is seen from the client's view.
func (c *MockConnection) Write(b []byte) (count int, err error) {
	count = len(b)
	c.a2bMu.Lock()
	c.a2bMsg = append(c.a2bMsg, b)
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
