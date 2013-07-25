package transport

// xlattice_go/transport/tcp_endpoint_test.go

import (
	. "launchpad.net/gocheck"
)

func (s *XLSuite) TestEndPointInterface(c *C) {
	ep, err := NewTcpEndPoint("127.0.0.1:80")
	c.Assert(err, Equals, nil)

	addr := ep.Address()
	c.Assert(addr.String(), Equals, "127.0.0.1:80")

	x, err := ep.Clone()
	c.Assert(err, Equals, nil)
	c.Assert(ep.String(), Equals, x.String())
	c.Assert(ep.Equal(x), Equals, true)

	c.Assert(ep.Transport(), Equals, "tcp")

	foo := EndPointI(ep) // compiler accepts
	// bar := EndPointI(*ep)		// compiler rejects

	_ = foo
	// _,_ = foo, bar
}
