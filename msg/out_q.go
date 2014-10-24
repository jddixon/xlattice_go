package msg

// xlattice_go/msg/out_q.go

import (
	"fmt"
	xr "github.com/jddixon/rnglib_go"
	xc "github.com/jddixon/xlCrypto_go"
	xn "github.com/jddixon/xlNode_go"
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

type OutHandler struct {
	Node    *xn.Node
	PeerNdx int    // which one of the node's peers this handles
	MsgN    uint64 // message number, always odd
	CnxHandler
}

func MakeHelloMsg(n *xn.Node) (m *XLatticeMsg, err error) {
	var ck, sk, salt, sig []byte
	cmd := XLatticeMsg_Hello
	id := n.GetNodeID().Value()
	ck, err = xc.RSAPubKeyToWire(n.GetCommsPublicKey())
	if err == nil {
		sk, err = xc.RSAPubKeyToWire(n.GetSigPublicKey())
	}
	if err == nil {
		sysRNG := xr.MakeSystemRNG()
		salt = make([]byte, 8)
		sysRNG.NextBytes(salt)
		chunks := [][]byte{id, ck, sk, salt}
		sig, err = n.Sign(chunks)
	}
	if err == nil {
		m = &XLatticeMsg{
			Op:       &cmd,
			MsgN:     &ONE,
			ID:       id,
			CommsKey: ck,
			SigKey:   sk,
			Salt:     salt,
			Sig:      sig,
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
func (oh *OutHandler) SendHello() (err error) {
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
		hello, err = MakeHelloMsg(oh.Node)
	}
	if err == nil {
		oh.MsgN = uint64(1)
		err = oh.writeMsg(hello)
	}
	return
}

func NewOutHandler(n *xn.Node, k int, stopCh, stoppedCh chan (bool)) (err error) {

	// XXX STUB XXX
	return
}
