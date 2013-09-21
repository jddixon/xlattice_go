package reg

// xlattice_go/reg/mock_client.go

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
	Err      error   // run information
	Client   *Client // the real client
	*xn.Node         // dummy node providing keys, etc
}

// A Mock Client for use in testing.  Given contact information for a
// registry and the name of a cluster, it joins the cluster, collects
// information on the other members, and terminates when it has info
// on the entire membership.

func NewMockClient(
	serverName string, serverID *xi.NodeID, serverEnd xt.EndPointI,
	clusterName string, clusterID *xi.NodeID, size int, endPointCount int) (
	mc *MockClient, err error) {

	if endPointCount < 1 {
		err = ClientMustHaveEndPoint
		return
	}

	var ckPriv, skPriv *rsa.PrivateKey
	var cn *xn.Node
	var ep []xt.EndPointI
	var client *Client

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
		client, err = NewClient(serverName, serverID, serverEnd,
			&ckPriv.PublicKey,
			clusterName, clusterID, size, ep, cn)
	}
	if err == nil {
		// THIS IS WRONG: We create a Client first
		mc = &MockClient{
			Client: client,
			Node:   cn,
		}
	}
	return
}

// Start the client running in separate goroutine, so that this function
// is non-blocking.

func (mc *MockClient) Run() (err error) {

	fmt.Println("mock starting real client")
	go func() {
		mc.Client.Run()
	}()
	return
}
