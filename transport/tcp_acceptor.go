package transport

/**
 * An Acceptor is used by a Node or Peer to accept connection requests.
 * It is an advertisement for a service within a Overlay, that is,
 * within a given address space and using a particular transport
 * protocol.
 *
 * An Acceptor is an abstraction of a TCP/IP ServerSocket.  It is a
 * single EndPoint whose Address may be well known.  Other entities on
 * the network send messages to the Acceptor in order to establish
 * Connections.  The Acceptor may in some cases NOT be one of the
 * EndPoints involved in the new Connection; the Connection might
 * be between the requesting remote EndPoint and a new, ephemeral
 * local EndPoint.
 *
 * The transport protocol understood by the Acceptor need not be
 * the same as the transport protocol of Connections created.  That is,
 * the new Connection need not be in the same Overlay as the Acceptor.
 *
 * @author Jim Dixon
 */

import (
	"net"
)

type TcpAcceptor struct {
	closed   bool
	endPoint *TcpEndPoint
	listener *net.TCPListener
}

// XXX IMPEDANCE MISMATCH: Go's Listener returns simplified stream
// connections; ListenIP appears to return packet-by-packet
// connections.  Investigating.

func NewTcpAcceptor(strAddr string) (*TcpAcceptor, error) {
	var err error
	var listener *net.TCPListener
	var tcpAddr *net.TCPAddr
	if tcpAddr, err = net.ResolveTCPAddr("tcp", strAddr); err == nil {
		listener, err = net.ListenTCP("tcp", tcpAddr)
	}
	if err == nil {
		a := TcpAcceptor{}
		a.endPoint, _ = NewTcpEndPoint(strAddr)
		a.listener = listener
		return &a, nil
	} else {
		return nil, err
	}
}
func (a *TcpAcceptor) Accept() (cnx *TcpConnection, err error) {
	conn, err := a.listener.AcceptTCP()
	if err == nil {
		cnx, err = NewTcpConnection(conn)
	}
	return
}
func (a *TcpAcceptor) Close() error {
	a.closed = true
	return a.listener.Close()
}
func (a *TcpAcceptor) IsClosed() bool {
	return a.closed
}
func (a *TcpAcceptor) GetEndPoint() *TcpEndPoint {
	return a.endPoint

}
func (a *TcpAcceptor) String() string {
	return "TcpAcceptor: " + a.endPoint.String()
}
