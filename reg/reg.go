package reg

// xlattice_go/reg/reg.go

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	xm "github.com/jddixon/xlattice_go/msg"
	xn "github.com/jddixon/xlattice_go/node"
	xi "github.com/jddixon/xlattice_go/nodeID"
	xo "github.com/jddixon/xlattice_go/overlay"
	xt "github.com/jddixon/xlattice_go/transport"
)

var _ = fmt.Print

// bit flags
const (
	EPHEMERAL = 1 << iota
	FOO
	BAR
)

type RegNode struct {
	Acc          *xt.TcpAcceptor
	StopCh       chan bool
	StoppedCh    chan bool
	privCommsKey *rsa.PrivateKey // duplicated here to provide
	privSigKey   *rsa.PrivateKey // visibility in this package
	xn.Node
}

// options normally set from the command line or derived from those
type RegOptions struct {
	Name     string
	Lfs      string
	Address  string
	Port     int
	EndPoint xt.EndPointI // derived from Address, Port
	Testing  bool
	Verbose  bool
}

// Create a registry node.

func New(name, lfs string, id *xi.NodeID,
	cKey, sKey *rsa.PrivateKey,
	overlay *xo.OverlayI,
	endPoint xt.EndPointI) (rn *RegNode, err error) {

	var acc xt.AcceptorI
	var n *xn.Node

	if name == "" {
		name = "xlReg"
	}
	if lfs == "" {
		lfs = "/var/app/xlReg"
	}
	if err == nil && id == nil {
		id, err = xi.New(nil)
	}
	if err == nil && cKey == nil {
		cKey, err = rsa.GenerateKey(rand.Reader, 2048)
	}
	if err == nil && sKey == nil {
		sKey, err = rsa.GenerateKey(rand.Reader, 2048)
	}
	if err == nil && endPoint == nil {
		endPoint, err = xt.NewTcpEndPoint("127.0.0.1:44444")
	}
	if err == nil {
		endPoints := []xt.EndPointI{endPoint}
		n, err = xn.New(name, id, lfs, cKey, sKey, nil, endPoints, nil)
	}
	if err == nil {
		acc = n.GetAcceptor(0) // XXX should be open
		if acc == nil {
			err = xm.AcceptorNotLive
		}
	}
	if err == nil {
		stopCh := make(chan bool, 1)
		stoppedCh := make(chan bool, 1)
		rn = &RegNode{
			Acc:          acc.(*xt.TcpAcceptor),
			StopCh:       stopCh,
			StoppedCh:    stoppedCh,
			privCommsKey: cKey, // redundant, but provide visibility
			privSigKey:   sKey, //    in this package
			Node:         *n,
		}
	}
	return
}
