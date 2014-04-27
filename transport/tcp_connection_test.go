package transport

// xlattice_go/transport/tcp_connector_test.go

import (
	"fmt"
	. "gopkg.in/check.v1"
)

var _ = fmt.Print

func (s *XLSuite) TestCnxInterface(c *C) {
	acc, err := NewTcpAcceptor("127.0.0.1:0")
	c.Assert(err, Equals, nil)
	defer acc.Close() // just in case
	accEndPoint := acc.GetEndPoint()
	// fmt.Printf("tcp_connector_test acceptor listening on %s\n",
	//	accEndPoint.String())
	go func() {
		for {
			_, err := acc.Accept()
			if err != nil {
				break
			}
			// otherwise just ignore the connection
		}
	}()
	ctor, err := NewTcpConnector(accEndPoint)
	c.Assert(err, Equals, nil)
	nearEnd, err := NewTcpEndPoint("127.0.0.1:0")
	c.Assert(err, Equals, nil)
	cnx, err := ctor.Connect(nearEnd)
	c.Assert(err, Equals, nil)
	defer cnx.Close()

	// We have a good TcpConnection
	c.Assert(cnx.GetState(), Equals, CNX_CONNECTED)

	// But the port number on nearEnd is 0, the dynamically assigned
	// port number is returned.
	//c.Assert( nearEnd.String(), Equals, cnx.GetNearEnd().String())

	acc.Close()
	foo := ConnectionI(cnx)
	_ = foo
}
