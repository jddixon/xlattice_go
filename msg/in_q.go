package msg

// xlattice_go/msg/in_q.go

import (
	// "code.google.com/p/goprotobuf/proto"
	"fmt"
	xc "github.com/jddixon/xlattice_go/crypto"
	xn "github.com/jddixon/xlattice_go/node"
	xt "github.com/jddixon/xlattice_go/transport"
	xu "github.com/jddixon/xlattice_go/util"
)

var _ = fmt.Print

// States through which the input msg queue may pass
const (
	IN_START = iota
	HELLO_RCVD

	IN_CLOSED
)

type InHandler struct {
	CnxHandler
}

// Given an open new connection, process a hello message for this node,
// returning a handler for the connection if that succeeds.
func NewInHandler(n *xn.Node, conn xt.ConnectionI) (h *InHandler, err error) {
	if n == nil {
		return nil, NilNode
	}
	if conn == nil {
		return nil, NilConnection
	}
	cnx := conn.(*xt.TcpConnection)
	h = &InHandler{CnxHandler{Cnx: cnx, State: IN_START}}
	err = h.handleHello(n)
	if err == nil {
		err = h.handleInMsg()
	}
	return
}

// Send the text of the error message to the peer and close the connection.
func (h *InHandler) errorReply(e error) (err error) {
	var reply *xn.XLatticeMsg
	cmd := xn.XLatticeMsg_Error
	s := e.Error()
	reply.Op = &cmd
	reply.MsgN = &ONE
	reply.ErrDesc = &s
	h.writeMsg(reply) // ignore any write error
	h.State = IN_CLOSED
	return
}
func (h *InHandler) simpleAck(msgN uint64) (err error) {
	var reply *xn.XLatticeMsg
	cmd := xn.XLatticeMsg_Ack
	h.MsgN = msgN + 1
	reply.Op = &cmd
	reply.MsgN = &h.MsgN
	reply.YourMsgN = &msgN  // XXX this field is pointless !
	err = h.writeMsg(reply) // this may yield an error ...
	h.State = HELLO_RCVD
	return err
}
func (h *InHandler) checkMsgNbrAndAck(m *xn.XLatticeMsg) (err error) {
	msgN := m.GetMsgN()
	if msgN == h.MsgN+1 {
		err = h.simpleAck(msgN)
	} else {
		err = WrongMsgNbr
		s := err.Error() // its serialization
		var reply *xn.XLatticeMsg
		h.MsgN += 2 // from my point of view
		cmd := xn.XLatticeMsg_Error
		reply.Op = &cmd
		reply.MsgN = &h.MsgN
		reply.ErrDesc = &s
		reply.YourMsgN = &msgN
		h.writeMsg(reply)
		h.State = IN_CLOSED
	}
	return
}
func (h *InHandler) handleInMsg() (err error) {
	for err == nil {
		var m *xn.XLatticeMsg
		m, err = h.readMsg()
		if err == nil {
			cmd := m.GetOp()
			switch cmd {
			case xn.XLatticeMsg_Bye:
				err = h.checkMsgNbrAndAck(m)
				h.State = IN_CLOSED
			case xn.XLatticeMsg_KeepAlive:
				// XXX Update last-time-spoken-to for peer
				err = h.checkMsgNbrAndAck(m)
			default:
				// DEBUG
				fmt.Printf("handleInMsg: UNEXPECTED MESSAGE TYPE %v\n",
					m)
				// END
				err = UnexpectedMsgType
				h.errorReply(err) // ignore any errors from the call itself
			}
		}
	}
	return
}
func (h *InHandler) handleHello(n *xn.Node) (err error) {
	var (
		m    *xn.XLatticeMsg
		id   []byte
		peer *xn.Peer
	)
	m, err = h.readMsg()

	// message must be a Hello
	if err == nil {
		if m.GetOp() != xn.XLatticeMsg_Hello {
			err = MissingHello
		}
	}
	if err == nil {
		//  the message is a hello; is its NodeID that of a known peer?
		id = m.GetID()
		if peer = n.FindPeer(id); peer == nil {
			err = xn.NotAKnownPeer
		} else {
			h.Peer = peer
		}
	}

	// On any error up to here close the connection and delete the handler.
	if err != nil {
		h.Cnx.Close()
		h = nil
		return
	}
	// message is a hello from a known peer -------------------------

	// MsgN must be 1
	msgN := m.GetMsgN()
	h.MsgN = msgN
	ck := m.GetCommsKey()
	sk := m.GetSigKey()
	var serCK, serSK []byte
	if h.MsgN != 1 {
		err = xn.ExpectedMsgOne
	}
	if err == nil {
		serCK, err = xc.RSAPubKeyToWire(peer.GetCommsPublicKey())
		if err == nil {
			if !xu.SameBytes(serCK, ck) {
				err = xn.NotExpectedCommsKey
			}
		}
	}
	if err == nil {
		serSK, err = xc.RSAPubKeyToWire(peer.GetSigPublicKey())
		if err == nil {
			if !xu.SameBytes(serSK, sk) {
				err = xn.NotExpectedSigKey
			}
		}
	}

	if err == nil {
		// Everything is good; so Ack, leaving cnx open.
		err = h.simpleAck(msgN)
	} else {
		// Send the text of the error to the peer; the send itself
		// may of course cause an error, but we will ignore that.
		h.errorReply(err)
	}
	return
}
