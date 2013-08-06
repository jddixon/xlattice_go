package transport

// xlattice_go/transport/tcp_connector_test.go

import (
	"fmt"
	"github.com/jddixon/xlattice_go/rnglib"
	. "launchpad.net/gocheck"
)

func (s *XLSuite) TestCtorInterface(c *C) {
	ep, err := NewTcpEndPoint("127.0.0.1:80")
	c.Assert(err, Equals, nil)

	ctor, err := NewTcpConnector(ep)
	c.Assert(err, Equals, nil)

	ep2 := ctor.GetFarEnd()
	c.Assert(ep.Equal(ep2), Equals, true)

	foo := ConnectorI(ctor)
	_ = foo
}

func (s *XLSuite) TestSerialization(c *C) {
	rng := rnglib.MakeSimpleRNG()

	a := rng.Intn(256)
	b := rng.Intn(256)
	_c := rng.Intn(256)
	d := rng.Intn(256)
	port := rng.Intn(256 * 256)

	addr := fmt.Sprintf("%d.%d.%d.%d:%d", a, b, _c, d, port)

	ep, err := NewTcpEndPoint(addr)
	c.Assert(err, Equals, nil)
	ctor, err := NewTcpConnector(ep)
	c.Assert(err, Equals, nil)

	serialized := ctor.String()
	backAgain, err := ParseConnector(serialized)
	c.Assert(err, Equals, nil)
	reserialized := backAgain.String()
	c.Assert(serialized, Equals, reserialized)
}
