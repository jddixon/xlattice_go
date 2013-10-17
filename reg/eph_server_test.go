package reg

// xlattice_go/reg/eph_server_test.go
//   replaces mock_server_test.go AKA old_server_test.go

// BEING MODIFIED to follow the new approach, whereby
// 1.  we create an ephemeral registry using NewMockServer()
// 2.  we generate a random cluster name and size
// 3.  run an AdminClient to register the cluster with the registry (which
//       gives us a cluster ID), and then
// 4.  create the appropriate number K of UserClients
// 5.  do test run in which the K clients exchange details through the registry

import (
	// "crypto/rsa"
	"encoding/hex"
	"fmt"
	//xc "github.com/jddixon/xlattice_go/crypto"
	//xn "github.com/jddixon/xlattice_go/node"
	// xi "github.com/jddixon/xlattice_go/nodeID"
	xr "github.com/jddixon/xlattice_go/rnglib"
	xt "github.com/jddixon/xlattice_go/transport"
	. "launchpad.net/gocheck"
)

func (s *XLSuite) TestServer(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_SERVER")
	}

	rng := xr.MakeSimpleRNG()

	// 1.  create a new ephemeral server ----------------------------
	es, err := NewEphServer()
	c.Assert(err, IsNil)
	c.Assert(es, NotNil)

	server := es.Server

	c.Assert(&server.RegNode.ckPriv.PublicKey,
		DeepEquals, server.GetCommsPublicKey())
	serverName := server.GetName()
	serverID := server.GetNodeID()
	serverEnd := server.GetEndPoint(0)
	serverCK := server.GetCommsPublicKey()
	serverSK := server.GetSigPublicKey()
	c.Assert(serverEnd, NotNil)

	// start the mock server ------------------------------
	err = es.Run()
	c.Assert(err, IsNil)

	// 2. create a random cluster name and size ---------------------
	clusterName := rng.NextFileName(8)
	clusterAttrs := uint64(rng.Int63())
	K := 2 + rng.Intn(6) // so the size is 2 .. 7

	// 3. create an AdminClient, use it to get the clusterID
	an, err := NewAdminClient(serverName, serverID, serverEnd,
		serverCK, serverSK, clusterName, clusterAttrs, K, 1, nil)
	c.Assert(err, IsNil)

	// DEBUG
	fmt.Printf("Server name %s, id %s\n", serverName,
		hex.EncodeToString(serverID.Value()))
	fmt.Printf("cluster name %s; K is %d\n", clusterName, K)
	// END

	an.Run()
	cn := &an.ClientNode // a bit ugly, this ...
	<-cn.doneCh

	c.Assert(cn.clusterID, NotNil)          // the purpose of the exercise
	c.Assert(cn.epCount, Equals, uint32(1)) // FAILS

	// DEBUG
	fmt.Printf("AdminClient has registered a cluster of size %d\n    cluster ID is %s\n",
		K, hex.EncodeToString(cn.clusterID.Value()))
	// END

	// 4. create K clients ------------------------------------------

	uc := make([]*UserClient, K)
	ucNames := make([]string, K)
	for i := 0; i < K; i++ {
		var ep *xt.TcpEndPoint
		ep, err = xt.NewTcpEndPoint("127.0.0.1:0")
		c.Assert(err, IsNil)
		e := []xt.EndPointI{ep}
		ucNames[i] = rng.NextFileName(8) // not guaranteed to be unique
		uc[i], err = NewUserClient(ucNames[i], "",
			serverName, serverID, serverEnd, serverCK, serverSK,
			clusterName, cn.clusterAttrs, cn.clusterID,
			K, 1, e) //1 is endPoint count
		c.Assert(err, IsNil)
		c.Assert(uc[i], NotNil)
		c.Assert(uc[i].clusterID, NotNil) // FAILS :-)
	}

	// 5. start the K clients, each in a separate goroutine ---------
	for i := 0; i < K; i++ {
		err = uc[i].Run()
		c.Assert(err, IsNil)
	}

	// wait until all clients are done ------------------------------
	for i := 0; i < K; i++ {
		<-uc[i].ClientNode.doneCh
	}
	//
	// stop the server by closing its acceptor ----------------------
	es.Close()
	//
	// verify that results are as expected --------------------------
	//
	// XXX STUB XXX

}
