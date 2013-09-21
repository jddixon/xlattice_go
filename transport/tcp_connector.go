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
var NotTcpEndPoint = errors.New("not a Tcp endpoint")

func NewTcpConnector(farEnd EndPointI) (*TcpConnector, error) {
	switch v := farEnd.(type) {
	case *TcpEndPoint:
		_ = v
	default:
		return nil, NotTcpEndPoint
	}
	tcpFarEnd := farEnd.(*TcpEndPoint)

	// copy the far end
	ep2, err := tcpFarEnd.Clone()
	if err == nil {
		ctor := TcpConnector{ep2.(*TcpEndPoint)}
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

func (c *TcpConnector) Connect(nearEnd EndPointI) (ConnectionI, error) {
	var tcpNearEnd *TcpEndPoint
	if nearEnd == nil {
		tcpNearEnd = ANY_TCP_END_POINT
	} else {
		// XXX CHECK TYPE
		tcpNearEnd = nearEnd.(*TcpEndPoint)
	}
	tcpConn, err := net.DialTCP("tcp", tcpNearEnd.GetTcpAddr(),
		c.farEnd.GetTcpAddr())
	if err == nil {
		cnx := TcpConnection{tcpConn, CONNECTED}
		return &cnx, nil
	} else {
		return nil, err
	}
}

// return the Acceptor EndPoint that this Connector is used to
//          establish connections to

func (c *TcpConnector) GetFarEnd() EndPointI {
	// XXX Should return copy
	return c.farEnd
}

func (c *TcpConnector) String() string {
	// farEnd serialization begins with "TcpEndPoint: "
	return "TcpConnector: " + c.farEnd.String()[13:]
}
