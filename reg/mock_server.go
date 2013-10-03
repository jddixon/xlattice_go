package reg

// xlattice_go/reg/mock_server.go

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	xn "github.com/jddixon/xlattice_go/node"
	xi "github.com/jddixon/xlattice_go/nodeID"
	xr "github.com/jddixon/xlattice_go/rnglib"
	xt "github.com/jddixon/xlattice_go/transport"
)

var _ = fmt.Print

type MockServer struct {
	acc         xt.AcceptorI
	clusterName string
	clusterID   *xi.NodeID
	size        int
	Server      *RegServer
}

// A Mock Server is primarily intended for use in testing.  It contains
// a registry which handles one and only one cluster of a fixed size.

func NewMockServer(clusterName string, clusterID *xi.NodeID, size int) (
	ms *MockServer, err error) {

	if clusterName == "" || clusterID == nil {
		err = MissingClusterNameOrID
	} else if size < 2 {
		err = ClusterMustHaveTwo
	}
	if err != nil {
		return
	}

	// Create an XLattice node with quasi-random parameters including
	// low-quality keys and an endPoint in 127.0.0.1, localhost.

	var (
		ckPriv, skPriv *rsa.PrivateKey
		rn             *RegNode
		ep             *xt.TcpEndPoint
		node           *xn.Node
		reg            *Registry
		server         *RegServer
	)

	rng := xr.MakeSimpleRNG()
	name := rng.NextFileName(16)
	idBuf := make([]byte, SHA1_LEN)
	rng.NextBytes(&idBuf)
	lfs := "tmp/" + hex.EncodeToString(idBuf)
	id, err := xi.New(idBuf)
	if err == nil {
		// XXX cheap keys, too weak for any serious use
		ckPriv, err = rsa.GenerateKey(rand.Reader, 1024)
		if err == nil {
			skPriv, err = rsa.GenerateKey(rand.Reader, 1024)
		}
	}
	if err == nil {
		ep, err = xt.NewTcpEndPoint("127.0.0.1:0")
		eps := []xt.EndPointI{ep}
		if err == nil {
			node, err = xn.New(name, id, lfs, ckPriv, skPriv, nil, eps, nil)
			// a registry with no clusters and no logger
			reg, err = NewRegistry(nil, node, ckPriv, skPriv, nil)
		}
	}
	// DEBUG
	if reg.ClustersByID == nil {
		fmt.Println("NewMockServer: CLUSTERS_BY_ID IS NIL")
	}
	// END

	if err == nil {
		server, err = NewRegServer(reg, true, 1)
	}

	if err == nil {
		rn = &reg.RegNode
		ms = &MockServer{
			acc:         rn.GetAcceptor(0),
			clusterName: clusterName,
			clusterID:   clusterID,
			size:        size,
			Server:      server,
		}
	}
	return
}

// Start the mock server running in a separate goroutine.

func (ms *MockServer) Run() (err error) {

	err = ms.Server.Run()
	return
}

func (ms *MockServer) Close() {
	ms.Server.Close()
}
