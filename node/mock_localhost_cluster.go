package node

import (
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	xi "github.com/jddixon/xlattice_go/nodeID"
	xo "github.com/jddixon/xlattice_go/overlay"
	"github.com/jddixon/xlattice_go/rnglib"
	xt "github.com/jddixon/xlattice_go/transport"
)

var _ = fmt.Print

func MockLocalHostCluster(K int) (nodes []*Node, accs []*xt.TcpAcceptor) {

	rng := rnglib.MakeSimpleRNG()

	// Create K nodes, each with a NodeID, two RSA private keys (sig and
	// comms), and two RSA public keys.  Each node creates a TcpAcceptor
	// running on 127.0.0.1 and a random (= system-supplied) port.
	names := make([]string, K)
	nodeIDs := make([]*xi.NodeID, K)
	for i := 0; i < K; i++ {
		// TODO: MAKE NAMES UNIQUE
		names[i] = rng.NextFileName(4)
		val := make([]byte, xi.SHA1_LEN)
		rng.NextBytes(val)
		nodeIDs[i], _ = xi.NewNodeID(val)
	}
	nodes = make([]*Node, K)
	accs = make([]*xt.TcpAcceptor, K)
	accEndPoints := make([]*xt.TcpEndPoint, K)
	for i := 0; i < K; i++ {
		lfs := "tmp/" + hex.EncodeToString(nodeIDs[i].Value())
		nodes[i], _ = NewNew(names[i], nodeIDs[i], lfs)
	}
	// XXX We need this functionality in using code
	//	defer func() {
	//		for i := 0; i < K; i++ {
	//			if accs[i] != nil {
	//				accs[i].Close()
	//			}
	//		}
	//	}()

	// Collect the nodeID, public keys, and listening address from each
	// node.

	// all nodes on the same overlay
	ar, _ := xo.NewCIDRAddrRange("127.0.0.0/8")
	overlay, _ := xo.NewIPOverlay("XO", ar, "tcp", 1.0)

	// add an endpoint to each node
	for i := 0; i < K; i++ {
		ep, _ := xt.NewTcpEndPoint("127.0.0.1:0")
		nodes[i].AddEndPoint(ep)
		accs[i] = nodes[i].GetAcceptor(0).(*xt.TcpAcceptor)
		accEndPoints[i] = accs[i].GetEndPoint().(*xt.TcpEndPoint)
	}

	commsKeys := make([]*rsa.PublicKey, K)
	sigKeys := make([]*rsa.PublicKey, K)
	ctors := make([]*xt.TcpConnector, K)

	for i := 0; i < K; i++ {
		// we already have nodeIDs
		commsKeys[i] = nodes[i].GetCommsPublicKey()
		sigKeys[i] = nodes[i].GetSigPublicKey()
		ctors[i], _ = xt.NewTcpConnector(accEndPoints[i])
	}

	overlaySlice := []xo.OverlayI{overlay}
	peers := make([]*Peer, K)
	for i := 0; i < K; i++ {
		ctorSlice := []xt.ConnectorI{ctors[i]}
		_ = ctorSlice
		peers[i], _ = NewPeer(names[i], nodeIDs[i], commsKeys[i], sigKeys[i],
			overlaySlice, ctorSlice)
	}

	// Use the information collected to configure each node.
	for i := 0; i < K; i++ {
		for j := 0; j < K; j++ {
			if i != j {
				nodes[i].AddPeer(peers[j])
			}
		}
	}
	return
}
