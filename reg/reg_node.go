package reg

// xlattice_go/reg/reg_node.go

// We collect functions and structures relating to the operation
// of the registry as a communicating server here.

import (
	"crypto/rsa"
	"fmt"
	xm "github.com/jddixon/xlattice_go/msg"
	xn "github.com/jddixon/xlattice_go/node"
	xi "github.com/jddixon/xlattice_go/nodeID"
	xo "github.com/jddixon/xlattice_go/overlay"
	// xr "github.com/jddixon/xlattice_go/rnglib"
	xt "github.com/jddixon/xlattice_go/transport"
)

var _ = fmt.Print

// options normally set from the command line or derived from those

type RegOptions struct {
	Name     string
	ID       *xi.NodeID
	Lfs      string
	Address  string
	Port     int
	EndPoint xt.EndPointI // derived from Address, Port
	Testing  bool
	Verbose  bool
}

type RegNode struct {
	Acc       xt.AcceptorI
	StopCh    chan bool       // volatile, so not serialized
	StoppedCh chan bool       // -ditto-
	ckPriv    *rsa.PrivateKey // duplicate to allow simple
	skPriv    *rsa.PrivateKey // access in this package
	xn.Node                   // name, id, ck, sk, etc, etc
}

func NewRegNode(name string, id *xi.NodeID, lfs string,
	commsKey, sigKey *rsa.PrivateKey,
	overlay xo.OverlayI, endPoint xt.EndPointI) (
	q *RegNode, err error) {

	var myNode *xn.Node
	var o []xo.OverlayI
	var e []xt.EndPointI
	var acc xt.AcceptorI

	if name == "" {
		name = "xlReg"
	}
	if id == nil {
		id, _ = xi.New(nil) // uses expensive SystemRNG to create a random ID
	}
	if lfs == "" {
		lfs = "/var/app/xlReg"
	}
	if overlay != nil {
		o = []xo.OverlayI{overlay}
	}
	if commsKey == nil || sigKey == nil {
		err = NilPrivateKey
	} else if endPoint == nil {
		endPoint, err = xt.NewTcpEndPoint("127.0.0.1:44444")
	}
	if err == nil {
		e = []xt.EndPointI{endPoint}
		myNode, err = xn.New(name, id, lfs, commsKey, sigKey, o, e, nil)
	}
	if err == nil {
		acc = myNode.GetAcceptor(0)
		if acc == nil {
			err = xm.AcceptorNotLive
		}
	}
	if err == nil {
		stopCh := make(chan bool, 1)
		stoppedCh := make(chan bool, 1)

		q = &RegNode{
			Acc:       acc,
			StopCh:    stopCh,
			StoppedCh: stoppedCh,
			ckPriv:    commsKey,
			skPriv:    sigKey,
			Node:      *myNode,
		}
	}
	return
}
