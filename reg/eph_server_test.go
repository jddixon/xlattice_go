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
	c.Assert(es, Not(IsNil))

	server := es.Server

	c.Assert(&server.RegNode.ckPriv.PublicKey,
		DeepEquals, server.GetCommsPublicKey())
	serverName := server.GetName()
	serverID := server.GetNodeID()
	serverEnd := server.GetEndPoint(0)
	serverCK := server.GetCommsPublicKey()
	serverSK := server.GetSigPublicKey()
	c.Assert(serverEnd, Not(IsNil))

	// start the mock server ------------------------------
	err = es.Run()
	c.Assert(err, IsNil)

	// 2. create a random cluster name and size ---------------------
	clusterName := rng.NextFileName(8)
	K := 2 + rng.Intn(6) // so the size is 2 .. 7

	// 3. create an AdminClient, use it to get the clusterID
	an, err := NewAdminClient(serverName, serverID, serverEnd,
		serverCK, serverSK, clusterName, K, 1)
	c.Assert(err, IsNil)

	// DEBUG
	fmt.Printf("Server name %s, id %s\n", serverName,
		hex.EncodeToString(serverID.Value()))
	fmt.Printf("cluster name %s; K is %d\n", clusterName, K)
	// END

	an.Run()
	cn := &an.ClientNode // a bit ugly, this ...
	<-cn.doneCh

	// DEBUG
	fmt.Printf("we're back! cluster ID is %s\n",
		hex.EncodeToString(cn.clusterID.Value()))
	// END

	// 4. create K clients ------------------------------------------
	//
	//	// MUST REPLACE MockClient BY UserClient
	//	uc := make([]*MockClient, K)
	//	for i := 0; i < K; i++ {
	//		uc[i], err = NewMockClient(rng,
	//			serverName, serverID, serverEnd, serverCK,
	//			clusterName, clusterID, K, 1) // 1 is endPoint count
	//		c.Assert(err, IsNil)
	//		c.Assert(uc[i], Not(IsNil))
	//	}

	// 5. start the K clients, each in a separate goroutine ---------
	//	for i := 0; i < K; i++ {
	//		err = uc[i].Run()
	//		c.Assert(err, IsNil)
	//	}
	//
	//	// wait until all clients are done ------------------------------
	//	// XXX MUST REPLACE OldClient BY ClientNode
	//	for i := 0; i < K; i++ {
	//		<-uc[i].OldClient.doneCh
	//	}
	//
	//	// stop the server by closing its acceptor ----------------------
	//	es.Close()
	//
	//	// verify that results are as expected --------------------------
	//
	//	// XXX STUB XXX

	// DEBUG
	_, _, _, _, _ = serverName, serverID, serverEnd, serverCK, serverSK
	_, _ = clusterName, K
}
