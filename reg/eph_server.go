package reg

// xlattice_go/reg/eph_server.go

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

type EphServer struct {
	acc    xt.AcceptorI
	Server *RegServer
}

// An ephemeral server is primarily intended for use in testing.
// It does not persist registry information to disk.  It listens
// on a random port, 127.0.0.1:0.

func NewEphServer() (ms *EphServer, err error) {

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
	rng.NextBytes(idBuf)
	lfs := "tmp/" + hex.EncodeToString(idBuf)
	id, err := xi.New(nil)
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
		}
		if err == nil {
			rn, err = NewRegNode(node, ckPriv, skPriv)
		}
		if err == nil {
			// a registry with no clusters and no logger
			opt := &RegOptions{
				EndPoint:  ep, // not used
				Ephemeral: true,
				Lfs:       lfs, // redundant (is in node's BaseNode)
				Logger:    nil,
				K:         DEFAULT_K,
				M:         DEFAULT_M,
			}
			reg, err = NewRegistry(nil, rn, opt)
		}
	}
	if err == nil {
		server, err = NewRegServer(reg, true, 1)
	}
	if err == nil {
		ms = &EphServer{
			acc:    rn.GetAcceptor(0),
			Server: server,
		}
	}
	return
}

// Start the ephemeral server running in a separate goroutine.

func (ms *EphServer) Run() (err error) {

	err = ms.Server.Run()
	return
}

func (ms *EphServer) Close() {
	ms.Server.Close()
}
