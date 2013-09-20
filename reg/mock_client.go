package reg

// xlattice_go/reg/mock_client.go

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	//xc "github.com/jddixon/xlattice_go/crypto"
	//xn "github.com/jddixon/xlattice_go/node"
	xi "github.com/jddixon/xlattice_go/nodeID"
	xr "github.com/jddixon/xlattice_go/rnglib"
	xt "github.com/jddixon/xlattice_go/transport"
)

var _ = fmt.Print

type MockClient struct {
	serverName string
	serverID   *xi.NodeID
	serverAcc  xt.AcceptorI

	clusterName string
	clusterID   *xi.NodeID
	size        int

	// run information
	doneCh chan bool
	err    error

	// information on other cluster members
	others []*ClusterMember

	// information on this cluster member
	acc *xt.AcceptorI
	RegNode
}

// A Mock Client for use in testing.  Given contact information for a
// registry and the name of a cluster, it joins the cluster, collects
// information on the other members, and terminates when it has info
// on the entire membership.

func NewMockClient(
	serverName string, serverID *xi.NodeID, serverAcc xt.AcceptorI,
	clusterName string, clusterID *xi.NodeID, size int) (
	mc *MockClient, err error) {

	// sanity checks on parameter list
	if serverName == "" || serverID == nil || serverAcc == nil {
		err = MissingServerInfo
	} else if clusterName == "" || clusterID == nil {
		err = MissingClusterNameOrID
	} else if size < 2 {
		err = ClusterMustHaveTwo
	}

	if err != nil {
		return
	}

	var ckPriv, skPriv *rsa.PrivateKey
	var rn *RegNode
	var ep *xt.TcpEndPoint

	rng := xr.MakeSimpleRNG()
	name := rng.NextFileName(16)
	idBuf := make([]byte, SHA1_LEN)
	rng.NextBytes(&idBuf)
	lfs := "tmp/" + hex.EncodeToString(idBuf)
	id, err := xi.New(idBuf)
	if err == nil {
		// XXX cheap keys, not meant for any serious use
		ckPriv, err = rsa.GenerateKey(rand.Reader, 512)
		if err == nil {
			skPriv, err = rsa.GenerateKey(rand.Reader, 512)
		}
	}
	if err == nil {
		ep, err = xt.NewTcpEndPoint("127.0.0.1:0")
	}
	if err == nil {
		rn, err = NewRegNode(name, id, lfs, ckPriv, skPriv, nil, ep)
	}
	if err == nil {
		mc = &MockClient{
			doneCh:      make(chan bool, 1),
			clusterName: clusterName,
			clusterID:   clusterID,
			size:        size,
			RegNode:     *rn,
		}
	}
	return
}

// Start the client running in separate goroutine, so that this function
// is non-blocking.

func (mc *MockClient) Run() (err error) {

	go func() {

		// XXX STUB XXX

		mc.doneCh <- true
	}()
	return
}
