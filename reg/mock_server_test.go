package reg

// xlattice_go/reg/mock_server_test.go

import (
	// "crypto/rsa"
	"fmt"
	//xc "github.com/jddixon/xlattice_go/crypto"
	//xn "github.com/jddixon/xlattice_go/node"
	xi "github.com/jddixon/xlattice_go/nodeID"
	xr "github.com/jddixon/xlattice_go/rnglib"
	. "launchpad.net/gocheck"
)

func (s *XLSuite) TestMockServer(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_MOCK_SERVER")
	}

	rng := xr.MakeSimpleRNG()

	// create the mock server and collect its attributes ------------
	clusterName := rng.NextFileName(8)
	clusterID, err := xi.New(nil) // creates a random ID
	c.Assert(err, IsNil)
	K := 2 + rng.Intn(6) // so the size is 2 .. 7
	ms, err := NewMockServer(clusterName, clusterID, K)
	c.Assert(err, IsNil)
	c.Assert(ms, Not(IsNil))

	c.Assert(&ms.RegNode.ckPriv.PublicKey, DeepEquals, ms.GetCommsPublicKey())

	serverName := ms.GetName()
	serverID := ms.GetNodeID()
	serverEnd := ms.GetEndPoint(0)
	c.Assert(serverEnd, Not(IsNil))

	// creake K clients ---------------------------------------------
	mc := make([]*MockClient, K)
	for i := 0; i < K; i++ {
		mc[i], err = NewMockClient(serverName, serverID, serverEnd,
			clusterName, clusterID, K, 1) // 1 is endPoint count
		c.Assert(err, IsNil)
		c.Assert(mc[i], Not(IsNil))
	}
	// start the mock server ----------------------------------------
	err = ms.Run()
	c.Assert(err, IsNil)

	// start the K clients, each in a separate goroutine ------------
	for i := 0; i < K; i++ {
		err = mc[i].Run()
		c.Assert(err, IsNil)
	}

	// wait until all clients are done ------------------------------
	for i := 0; i < K; i++ {
		<-mc[i].Client.doneCh
		fmt.Printf("mock client %d says that it's done\n", i)
	}

	// stop the server by closing its acceptor ----------------------
	ms.Close()

	// verify that results are as expected --------------------------

	// XXX STUB XXX
}
