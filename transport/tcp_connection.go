package transport

import (
	"fmt"
	xc "github.com/jddixon/xlattice_go/crypto"
	"net"
)

type TcpConnection struct {
	conn  *net.TCPConn
	state int
}

func NewTcpConnection(conn *net.TCPConn) (cnx *TcpConnection, err error) {
	if conn == nil {
		err = NilConnection
	} else {
		cnx = &TcpConnection{conn: conn, state: CNX_CONNECTED}
	}
	return
}

// Return the current state index.
func (c *TcpConnection) GetState() int {
	return c.state
}

// Set the near end point of a connection.  If either the
// near or far end point has already been set, this will
// cause an exception.  If successful, the connection's
// state becomes BOUND.
//
func (c *TcpConnection) BindNearEnd(e EndPointI) (err error) {
	return NotImplemented
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
func (c *TcpConnection) BindFarEnd(e EndPointI) (err error) {
	return NotImplemented
}

// Bring the connection to the DISCONNECTED state.
//
func (c *TcpConnection) Close() (err error) {
	c.state = CNX_DISCONNECTED
	return c.conn.Close()
}

// XXX 2013-07-20: this returns the far end instead !
func (c *TcpConnection) GetNearEnd() (ep EndPointI) {
	ep, _ = NewTcpEndPoint(c.conn.LocalAddr().String())
	return ep
}

// XXX 2013-07-20: this returns the near end instead !
func (c *TcpConnection) GetFarEnd() (ep EndPointI) {
	ep, _ = NewTcpEndPoint(c.conn.RemoteAddr().String())
	return ep
}

func (c *TcpConnection) Read(b []byte) (int, error) {
	return c.conn.Read(b)
}
func (c *TcpConnection) Write(b []byte) (int, error) {
	return c.conn.Write(b)
}
func (c *TcpConnection) IsBlocking() bool {
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
func (c *TcpConnection) IsEncrypted() bool {
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
func (c *TcpConnection) Negotiate(myKey xc.KeyI, hisKey xc.PublicKeyI) (s xc.SecretI, e error) {
	// XXX STUB
	return nil, NotImplemented
}

func (c *TcpConnection) Equal(any interface{}) bool {
	// XXX STUB NotImplemented
	return false
}

func (c *TcpConnection) String() string {
	return fmt.Sprintf("Tcp: %s --> %s",
		c.GetNearEnd().String(),
		c.GetFarEnd().String())
}
