package reg

// xlattice_go/reg/mock_client.go

//////////////////////////
// THIS IS BEING REPLACED.
//////////////////////////

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	//xc "github.com/jddixon/xlattice_go/crypto"
	xn "github.com/jddixon/xlattice_go/node"
	xi "github.com/jddixon/xlattice_go/nodeID"
	xr "github.com/jddixon/xlattice_go/rnglib"
	xt "github.com/jddixon/xlattice_go/transport"
)

var _ = fmt.Print

type MockClient struct {
	// Err      error   // run information
	OldClient *OldClient // the real client
	*xn.Node             // node providing keys, etc
}

// Given contact information for a registry and the name of a cluster,
// a MockClient joins the cluster, collects information on the other
// members, and terminates when it has info on the entire membership.

func NewMockClient(
	rng *xr.PRNG,
	serverName string, serverID *xi.NodeID, serverEnd xt.EndPointI,
	serverCK *rsa.PublicKey,
	clusterName string, clusterID *xi.NodeID, size int, endPointCount int) (
	mc *MockClient, err error) {

	if endPointCount < 1 {
		err = ClientMustHaveEndPoint
		return
	}

	var ckPriv, skPriv *rsa.PrivateKey
	var cn *xn.Node
	var ep []xt.EndPointI
	var client *OldClient

	name := rng.NextFileName(16)
	idBuf := make([]byte, SHA1_LEN)
	rng.NextBytes(&idBuf)
	lfs := "tmp/" + hex.EncodeToString(idBuf)
	id, err := xi.New(idBuf)
	if err == nil {
		// XXX key must be large enough to hold data to be encrypted
		ckPriv, err = rsa.GenerateKey(rand.Reader, 1024)
		if err == nil {
			skPriv, err = rsa.GenerateKey(rand.Reader, 512)
		}
	}
	if err == nil {
		for i := 0; i < endPointCount; i++ {
			var endPoint xt.EndPointI
			endPoint, err = xt.NewTcpEndPoint("127.0.0.1:0")
			if err != nil {
				break
			}
			ep = append(ep, endPoint)
		}
	}
	if err == nil {
		cn, err = xn.New(name, id, lfs, ckPriv, skPriv, nil, ep, nil)
	}
	if err == nil {
		client, err = NewOldClient(serverName, serverID, serverEnd,
			serverCK,
			clusterName, clusterID, size, ep, cn)
	}
	if err == nil {
		// THIS IS WRONG: We create a OldClient first
		mc = &MockClient{
			OldClient: client,
			Node:      cn,
		}
	}
	return
}

// The client's Run() runs in separate goroutine, so that this function
// is non-blocking.

func (mc *MockClient) Run() (err error) {
	mc.OldClient.Run()
	return
}
