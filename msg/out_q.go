package msg

// xlattice_go/msg/out_q.go

import (
	"fmt"
	xc "github.com/jddixon/xlattice_go/crypto"
	xn "github.com/jddixon/xlattice_go/node"
	//xt "github.com/jddixon/xlattice_go/transport"
)

var _ = fmt.Print

type OutHandler struct {
	CnxHandler
}

func MakeHelloMsg(n *xn.Node) (m *XLatticeMsg, err error) {
	var ck, sk []byte
	cmd := XLatticeMsg_Hello
	ck, err = xc.RSAPubKeyToWire(n.GetCommsPublicKey())
	if err != nil {
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
