package reg

// xlattice_go/reg/client_node.go

// this used to be called mock_client.go

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

type ClientNode struct {
	// Err      error   // run information
	Client   *Client // the real client
	*xn.Node         // node providing keys, etc
}

// Given contact information for a registry and the name of a cluster, 
// a ClientNode joins the cluster, collects information on the other 
// members, and terminates when it has info on the entire membership.

func NewClientNode(
	rng *xr.PRNG,
	serverName string, serverID *xi.NodeID, serverEnd xt.EndPointI,
	serverCK *rsa.PublicKey,
	clusterName string, clusterID *xi.NodeID, size int, endPointCount int) (
	mc *ClientNode, err error) {

	if endPointCount < 1 {
		err = ClientMustHaveEndPoint
		return
	}

	var ckPriv, skPriv *rsa.PrivateKey
	var cn *xn.Node
	var ep []xt.EndPointI
	var client *Client

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
		client, err = NewClient(serverName, serverID, serverEnd,
			serverCK,
			clusterName, clusterID, size, ep, cn)
	}
	if err == nil {
		// THIS IS WRONG: We create a Client first
		mc = &ClientNode{
			Client: client,
			Node:   cn,
		}
	}
	return
}

// The client's Run() runs in separate goroutine, so that this function
// is non-blocking.

func (mc *ClientNode) Run() (err error) {
	mc.Client.Run()
	return
}
