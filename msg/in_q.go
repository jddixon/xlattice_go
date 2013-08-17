package msg

// xlattice_go/msg/in_q.go

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	xn "github.com/jddixon/xlattice_go/node"
	xt "github.com/jddixon/xlattice_go/transport"
)

var _ = fmt.Print

func DecodePacket(data []byte) (*xn.XLatticeMsg, error) {
	var m xn.XLatticeMsg
	err := proto.Unmarshal(data, &m)
	// XXX do some filtering, eg for nil op
	return &m, err
}

// states through which the input msg queue may pass
const (
	IN_START = iota
	HELLO_RCVD

	IN_CLOSED
)

type InHandler struct {
	cnx     *xt.TcpConnection
	peerNdx int
	seqN    int
	state   int
}

func NewInHandler(conn xt.ConnectionI) (h *InHandler, err error) {
	// if conn is nil, return err
	if conn == nil {
		return nil, NilConnection
	}
	cnx := conn.(*xt.TcpConnection)
	h = &InHandler{cnx, -1, -1, IN_START}
	return
}
func (h *InHandler) readMsg() (m *xn.XLatticeMsg, err error) {
	// read the next message
	inBuf := make([]byte, MSG_BUF_LEN)
	count, err := h.cnx.Read(inBuf)
	if err == nil && count > 0 {
		inBuf = inBuf[:count]
		// parse = deserialize, unmarshal the message
		m, err = DecodePacket(inBuf)
	}
	return
}
func HandleHello(node *xn.Node, conn xt.ConnectionI) (h *InHandler, err error) {
	var m *xn.XLatticeMsg

	h, err = NewInHandler(conn)
	if err == nil {
		m, err = h.readMsg()
	}
	// message must be a Hello
	if err == nil {
		if m.GetOp() != xn.XLatticeMsg_Hello {
			err = MissingHello
		}
	}
	//  the message is a hello; if the nodeID is unknown, close cnx, return err
	// XXX STUB XXX

	// on any error up to here close the connection and delete the handler
	if err != nil {
		h.cnx.Close()
		h = nil
		return
	}
	// message is a hello from a known peer -------------------------

	// create InHandler{cnx, peerNdx, seqN, HELLO_RCVD}, expecting seqN = 1
	// XXX STUB XXX

	// if msgN, crypto pubKey, or sig pubKey not as expected, reply
	// with ErrMsg on a timeout.
	// XXX STUB XXX

	// whether or not there is a timeout, set state to IN_CLOSED,
	// and return suitable err
	// XXX STUB XXX

	// otherwise everything is good; so try to send an Ack; this is
	// on a timeout, repeat up to three times
	// XXX STUB XXX

	// if Ack fails, set state to IN_CLOSED, return timeout err
	// XXX STUB XXX

	// else Ack was successful, set state to HELLO_RCVD and return
	// XXX STUB XXX

	return
}
