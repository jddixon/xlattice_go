package transport

/**
 * Used to establish a Connection with another entity (Node).
 *
 * The notion is that a node has a collection of Connectors used
 * for establishing Connections with Peers, neighboring nodes.
 *
 * @author Jim Dixon
 */

import (
	"errors"
	"net"
)

type TcpConnector struct {
	farEnd *TcpEndPoint
}

var NilEndPoint = errors.New("nil endpoint")

func NewTcpConnector(farEnd TcpEndPoint) (*TcpConnector, error) {
	ep2, err := farEnd.Clone()
	if err == nil {
		ctor := TcpConnector{ep2}
		if err == nil {
			return &ctor, nil
		}
	}
	return nil, err
}

/**
 * Establish a Connection with another entity using the transport
 * and address in the EndPoint.
 *
 * @param nearEnd  local end point to use for connection
 * @param blocking whether the new Connection is to be blocking
 */

func (c *TcpConnector) Connect(nearEnd *TcpEndPoint) (*TcpConnection, error) {

	tcpConn, err := net.DialTCP("tcp", nearEnd.GetTcpEndPoint(),
		c.farEnd.GetTcpEndPoint())
	if err == nil {
		cnx := TcpConnection{tcpConn, CONNECTED}
		return &cnx, nil
	} else {
		return nil, err
	}
}

// return the Acceptor EndPoint that this Connector is used to
//          establish connections to

func (c *TcpConnector) GetFarEnd() *TcpEndPoint {
	// XXX Should return copy
	return c.farEnd
}

func (c *TcpConnector) String() string {
	return "TCPConnector: " + c.farEnd.String()
}
