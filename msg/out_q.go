package msg

// xlattice_go/msg/out_q.go

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	xn "github.com/jddixon/xlattice_go/node"
	//xt "github.com/jddixon/xlattice_go/transport"
)

var _ = fmt.Print

func EncodePacket(msg *xn.XLatticeMsg) (data []byte, err error) {
	return proto.Marshal(msg)
}
