package transport

// xlattice_go/transport/tcp_connector_test.go

import (
	. "launchpad.net/gocheck"
)

func (s *XLSuite) TestCtorInterface(c *C) {
	ep, err := NewTcpEndPoint("127.0.0.1:80")
	c.Assert(err, Equals, nil)

	ctor, err := NewTcpConnector(ep)
	c.Assert(err, Equals, nil)

	ep2 := ctor.GetFarEnd()
	c.Assert(ep.Equal(ep2), Equals, true)

	foo := ConnectorI(ctor) // fails
	_ = foo
}
