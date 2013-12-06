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
	"fmt"
	xr "github.com/jddixon/xlattice_go/rnglib"
	xt "github.com/jddixon/xlattice_go/transport"
	. "launchpad.net/gocheck"
)

func (s *XLSuite) TestEphServer(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_EPH_SERVER")
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
	defer es.Close() // stop the server by closing its acceptor

	// verify Bloom filter is running
	reg := es.Server.Registry
	c.Assert(reg, NotNil)
	regID := reg.GetNodeID()
	c.Assert(reg.IDCount(), Equals, uint(1)) // the registry's own ID
	found, err := reg.ContainsID(regID)
	c.Assert(found, Equals, true)

	// 2. create a random cluster name and size ---------------------
	clusterName := rng.NextFileName(8)
	clusterAttrs := uint64(rng.Int63())
	K := 2 + rng.Intn(6) // so the size is 2 .. 7

	// 3. create an AdminClient, use it to get the clusterID
	an, err := NewAdminClient(serverName, serverID, serverEnd,
		serverCK, serverSK, clusterName, clusterAttrs, K, 1, nil)
	c.Assert(err, IsNil)

	an.Run()
	cn := &an.ClientNode // a bit ugly, this ...
	<-cn.DoneCh

	c.Assert(cn.ClusterID, NotNil) // the purpose of the exercise
	c.Assert(cn.EpCount, Equals, uint32(1))

	anID := reg.GetNodeID()
	c.Assert(reg.IDCount(), Equals, uint(3)) // regID + anID + clusterID

	found, err = reg.ContainsID(anID)
	c.Assert(err, IsNil)
	c.Assert(found, Equals, true)
	// may be redundant...
	found, err = reg.ContainsID(an.ClusterID)
	c.Assert(err, IsNil)
	c.Assert(found, Equals, true)
	found, err = reg.ContainsID(cn.ClusterID)
	c.Assert(err, IsNil)
	c.Assert(found, Equals, true)

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
			nil, nil, // private RSA keys are generated if nil
			serverName, serverID, serverEnd, serverCK, serverSK,
			clusterName, cn.ClusterAttrs, cn.ClusterID,
			K, 1, e) //1 is endPoint count
		c.Assert(err, IsNil)
		c.Assert(uc[i], NotNil)
		c.Assert(uc[i].ClusterID, NotNil)
	}

	// 5. start the K clients, each in a separate goroutine ---------
	for i := 0; i < K; i++ {
		uc[i].Run()
	}

	// wait until all clients are done ------------------------------
	for i := 0; i < K; i++ {
		success := <-uc[i].ClientNode.DoneCh
		c.Assert(success, Equals, true)
		// if false, should check cn.Err for error

		// XXX NEXT LINE APPARENTLY DOES NOT WORK
		// nodeID := uc[i].ClientNode.GetNodeID()
		nodeID := uc[i].clientID
		c.Assert(nodeID, NotNil)
		found, err := reg.ContainsID(nodeID)
		c.Assert(err, IsNil)
		c.Assert(found, Equals, true)
	}
	c.Assert(reg.IDCount(), Equals, uint(3+K)) // regID + anID + clusterID + K

	//
	// verify that results are as expected --------------------------
	//
	// XXX STUB XXX

}
