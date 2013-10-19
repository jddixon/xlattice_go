package reg

// xlattice_go/reg/reg_cred_test.go

import (
	// "crypto/rsa"
	"fmt"
	xn "github.com/jddixon/xlattice_go/node"
	xi "github.com/jddixon/xlattice_go/nodeID"
	xr "github.com/jddixon/xlattice_go/rnglib"
	xt "github.com/jddixon/xlattice_go/transport"
	xu "github.com/jddixon/xlattice_go/util"
	. "launchpad.net/gocheck"
)

func (s *XLSuite) makeRegCred(c *C, rng *xr.PRNG) (rc *RegCred) {

	name := rng.NextFileName(8)
	nodeID, _ := xi.New(nil)
	node, err := xn.NewNew(name, nodeID, "") // "" is LFS
	c.Assert(err, IsNil)
	epCount := rng.Intn(4)
	var e []xt.EndPointI
	for i := 0; i < epCount; i++ {
		port := 1024 + rng.Intn(256*256-1024)
		strAddr := fmt.Sprintf("127.0.0.1:%d", port)
		ep, err := xt.NewTcpEndPoint(strAddr)
		c.Assert(err, IsNil)
		e = append(e, ep)
	}
	version := xu.DecimalVersion(uint32(rng.Int31()))
	rc = &RegCred{
		Name:        name,
		ID:          nodeID,
		CommsPubKey: node.GetCommsPublicKey(),
		SigPubKey:   node.GetSigPublicKey(),
		EndPoints:   e,
		Version:     version,
	}
	return
}
func (s *XLSuite) TestRegCred(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_REG_CRED_TEST")
	}
	rng := xr.MakeSimpleRNG()
	for i := 0; i < 4; i++ {
		rc := s.makeRegCred(c, rng)
		serialized := rc.String()
		backAgain, err := ParseRegCred(serialized)
		c.Assert(err, IsNil)
		serialized2 := backAgain.String()
		c.Assert(serialized2, Equals, serialized)
	}
}
