package reg

// xlattice_go/reg/reg.go

import (
	"fmt"
	xm "github.com/jddixon/xlattice_go/msg"
	xn "github.com/jddixon/xlattice_go/node"
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
	Acc       *xt.TcpAcceptor
	StopCh    chan bool
	StoppedCh chan bool
	xn.Node
}

// options normally set from the command line
type RegOptions struct {
	Lfs     string
	Port    int
	Testing bool
	Verbose bool
}

// Create a registry node around a live node which has an open acceptor

func New(n *xn.Node, stopCh, stoppedCh chan bool) (rn *RegNode, err error) {

	if n == nil {
		err = xm.NilNode
	}
	tcpAcc := n.GetAcceptor(0).(*xt.TcpAcceptor)
	if err == nil && tcpAcc == nil {
		err = xm.AcceptorNotLive
	}
	if err == nil && stopCh == nil {
		err = xm.NilControlCh
	}
	if err == nil {
		rn = &RegNode{
			Acc:       tcpAcc,
			StopCh:    stopCh,
			StoppedCh: stoppedCh,
			Node:      *n,
		}
	}
	return
}
