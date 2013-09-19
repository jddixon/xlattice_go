package reg

// xlattice_go/reg/mock_registry.go

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

type MockServer struct {
	stopCh      chan bool
	stoppedCh   chan bool
	clusterName string
	clusterID   *xi.NodeID

	RegNode
}

func NewMockServer(clusterName string, clusterID *xi.NodeID) (
	ms *MockServer, err error) {

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
		ms = &MockServer{
			stopCh:      make(chan bool, 1),
			stoppedCh:   make(chan bool, 1),
			clusterName: clusterName,
			clusterID:   clusterID,
			RegNode:     *rn,
		}
	}
	return
}

func (ms *MockServer) Run() {
	// XXX STUB XXX
}
