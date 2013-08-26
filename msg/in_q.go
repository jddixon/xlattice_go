package msg

// xlattice_go/msg/in_q.go

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	xc "github.com/jddixon/xlattice_go/crypto"
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


// States through which the input msg queue may pass
const (
	IN_START = iota
	HELLO_RCVD

	IN_CLOSED
)

type InHandler struct {
	Cnx     *xt.TcpConnection
	Peer	*xn.Peer
	MsgN    uint64
	State   int
}

// Given an open new connection, process a hello message for this node,
// returning a handler for the connection if that succeeds.
func NewInHandler(n *xn.Node, conn xt.ConnectionI)(h *InHandler, err error){
	if n == nil {
		return nil, NilNode
	}
	if conn == nil {
		return nil, NilConnection
	}
	cnx := conn.(*xt.TcpConnection)
	h = &InHandler{Cnx:cnx, State:IN_START}
	err = h.handleHello(n)
	if err == nil { 
		return
	} else { 
		return nil, err 
	}
}
func (h *InHandler) readMsg() (m *xn.XLatticeMsg, err error) {
	// read the next message
	inBuf := make([]byte, MSG_BUF_LEN)
	count, err := h.Cnx.Read(inBuf)
	if err == nil && count > 0 {
		inBuf = inBuf[:count]
		// parse = deserialize, unmarshal the message
		m, err = DecodePacket(inBuf)
	}
	return
}

func (h *InHandler) handleHello(n *xn.Node)(err error){
	var m *xn.XLatticeMsg

	m, err = h.readMsg()

	// message must be a Hello
	if err == nil {
		if m.GetOp() != xn.XLatticeMsg_Hello {
			err = MissingHello
		}
	}
	//  the message is a hello; is its NodeID that of a known peer?
	id   := m.GetID()
	peer := n.FindPeer(id)
	if peer == nil {
		err = xn.NotAKnownPeer
	}
	h.Peer = peer

	// on any error up to here close the connection and delete the handler
	if err != nil {
		h.Cnx.Close()
		h = nil
		return
	}
	// message is a hello from a known peer -------------------------
	h.State = HELLO_RCVD
	
	// MsgN must be 1
	h.MsgN = m.GetMsgN()
	ck := m.GetCommsKey()
	sk := m.GetSigKey()
	var serCK, serSK []byte
	if h.MsgN == 1 {
		err = xn.ExpectedMsgOne
	} 
	if err == nil {
		serCK, err = xc.RSAPubKeyToWire(peer.GetCommsPublicKey()) 
		if err == nil {
			if ! SameBytes(serCK, ck) {
				err = xn.NotExpectedCommsKey
			}
		}
	}
	if err == nil {
		serSK, err = xc.RSAPubKeyToWire(peer.GetSigPublicKey()) 
		if err == nil {
			if ! SameBytes(serSK, sk) {
				err = xn.NotExpectedSigKey
			}
		}
	}


	// if msgN, crypto pubKey, or sig pubKey not as expected, reply
	// with ErrMsg on a timeout.
	// XXX STUB XXX

		// whether or not there is a timeout, set State to IN_CLOSED,
		// and return suitable err
		// XXX STUB XXX

	// otherwise everything is good; so Ack, leaving cnx open
	// XXX STUB XXX

	return
}
