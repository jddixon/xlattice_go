package msg

// xlattice_go/msg/out_q.go

import (
	"errors"
	"fmt"
	xc "github.com/jddixon/xlattice_go/crypto"
	xn "github.com/jddixon/xlattice_go/node"
	//xt "github.com/jddixon/xlattice_go/transport"
)

var _ = fmt.Print

// XXX This is specific to a given node ! and so needs to either be
// a map or part of the Node data structure
var (
	helloMsg *XLatticeMsg // should be cached
)

const (
	OUT_START = iota
	HELLO_SENT

	OUT_CLOSED
)

var (
	CannotSendSecondHello = errors.New("can't send second hello")
)

type OutHandler struct {
	CnxHandler
}

func MakeHelloMsg(n *xn.Node) (m *XLatticeMsg, err error) {
	var ck, sk []byte
	cmd := XLatticeMsg_Hello
	ck, err = xc.RSAPubKeyToWire(n.GetCommsPublicKey())
	if err == nil {
		sk, err = xc.RSAPubKeyToWire(n.GetSigPublicKey())
	}
	if err == nil {
		m = &XLatticeMsg{
			Op:       &cmd,
			MsgN:     &ONE,
			ID:       n.GetNodeID().Value(),
			CommsKey: ck,
			SigKey:   sk,
		}
	}
	return
}

func (oh *OutHandler) SendBye() (err error) {
	// XXX should verify state
	cmd := XLatticeMsg_Bye
	oh.MsgN += 1
	bye := &XLatticeMsg{
		Op:   &cmd,
		MsgN: &oh.MsgN,
	}
	err = oh.writeMsg(bye)
	return
}
func (oh *OutHandler) SendHello(n *xn.Node) (err error) {
	var hello *XLatticeMsg
	if oh.Cnx == nil || oh.Peer == nil {
		err = MissingHandlerField
	}
	if err == nil {
		if oh.MsgN > 0 {
			err = CannotSendSecondHello
		}
	}
	if err == nil {
		hello, err = MakeHelloMsg(n)
	}
	if err == nil {
		oh.MsgN = uint64(1)
		err = oh.writeMsg(hello)
	}
	return
}
