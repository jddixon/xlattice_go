package msg

// xlattice_go/msg/in_q.go

import (
	"bytes"
	cr "crypto"
	"crypto/rsa"
	"crypto/sha1"
	"fmt"
	xc "github.com/jddixon/xlCrypto_go"
	xn "github.com/jddixon/xlNode_go"
	xt "github.com/jddixon/xlTransport_go"
	"strings"
)

var _ = fmt.Print

// States through which the input msg queue may pass
const (
	IN_START = iota
	HELLO_RCVD

	IN_CLOSED
)

type InHandler struct {
	StopCh, StoppedCh chan bool
	CnxHandler
}

// Given an open new connection, process a hello message for this node,
// returning a handler for the connection if that succeeds.
func NewInHandler(n *xn.Node, conn xt.ConnectionI,
	stopCh, stoppedCh chan bool) (h *InHandler, err error) {

	if n == nil {
		return nil, NilNode
	}
	if conn == nil {
		return nil, NilConnection
	}
	cnx := conn.(*xt.TcpConnection)
	h = &InHandler{
		StopCh:     stopCh,
		StoppedCh:  stoppedCh,
		CnxHandler: CnxHandler{Cnx: cnx, State: IN_START}}
	err = h.handleHello(n)
	if err == nil {
		err = h.handleInMsg()
	}
	return
}

// Send the text of the error message to the peer and close the connection.
func (h *InHandler) errorReply(e error) (err error) {
	var reply XLatticeMsg
	cmd := XLatticeMsg_Error
	s := e.Error()
	reply.Op = &cmd
	h.MsgN += 2
	reply.MsgN = &h.MsgN
	reply.ErrDesc = &s
	h.writeMsg(&reply) // ignore any write error
	h.State = IN_CLOSED

	// XXX This would be a very strong action, given that we may have multiple
	// connections open to this peer.
	// h.Peer.MarkDown()

	return
}
func (h *InHandler) simpleAck(msgN uint64) (err error) {
	h.Peer.StillAlive() // update time of last contact

	var reply XLatticeMsg
	cmd := XLatticeMsg_Ack
	h.MsgN = msgN + 1
	reply.Op = &cmd
	reply.MsgN = &h.MsgN
	reply.YourMsgN = &msgN   // XXX this field is pointless !
	err = h.writeMsg(&reply) // this may yield an error ...
	h.State = HELLO_RCVD
	return err
}
func (h *InHandler) checkMsgNbrAndAck(m *XLatticeMsg) (err error) {
	msgN := m.GetMsgN()
	if msgN == h.MsgN+1 {
		err = h.simpleAck(msgN)
	} else {
		err = WrongMsgNbr
		s := err.Error() // its serialization
		var reply XLatticeMsg
		h.MsgN += 2 // from my point of view
		cmd := XLatticeMsg_Error
		reply.Op = &cmd
		reply.MsgN = &h.MsgN
		reply.ErrDesc = &s
		reply.YourMsgN = &msgN
		h.writeMsg(&reply)
		h.State = IN_CLOSED
	}
	return
}
func (h *InHandler) handleInMsg() (err error) {
	defer h.Cnx.Close()
	go func() {
		<-h.StopCh
		h.Cnx.Close()
		h.StoppedCh <- true
	}()
	for err == nil {
		var m *XLatticeMsg
		m, err = h.readMsg()
		if err == nil {
			cmd := m.GetOp()
			switch cmd {
			case XLatticeMsg_Bye:
				err = h.checkMsgNbrAndAck(m)
				h.State = IN_CLOSED
			case XLatticeMsg_KeepAlive:
				// XXX Update last-time-spoken-to for peer
				h.Peer.LastContact()
				err = h.checkMsgNbrAndAck(m)
			default:
				// XXX should log
				fmt.Printf("handleInMsg: UNEXPECTED MESSAGE TYPE %v\n", m.GetOp())
				err = UnexpectedMsgType
				h.errorReply(err) // ignore any errors from the call itself
			}
		} else {
			break
		}
	}
	if err != nil {
		text := err.Error()
		if strings.HasSuffix(text, "use of closed network connection") {
			err = nil
		} else {
			fmt.Printf("    handleInMsg gets %v\n", err) // DEBUG
		}
	}
	return
}
func (h *InHandler) handleHello(n *xn.Node) (err error) {
	var (
		m                     *XLatticeMsg
		msgN                  uint64
		id, ck, sk, sig, salt []byte
		peer                  *xn.Peer
	)
	m, err = h.readMsg()

	// message must be a Hello
	if err == nil {
		if m.GetOp() != XLatticeMsg_Hello {
			err = MissingHello
		}
	}
	if err == nil {
		//  the message is a hello; is its NodeID that of a known peer?
		id = m.GetID()
		if id == nil {
			// DEBUG
			fmt.Printf("handleHello: message has no ID field\n")
			// END
			err = NilPeer
		} else {
			peer, err = n.FindPeer(id)
			if err == nil {
				if peer == nil {
					err = xn.NotAKnownPeer
				} else {
					h.Peer = peer
				}
			}
		}
	}

	// On any error up to here silently close the connection and delete
	// the handler.
	if err != nil {
		h.Cnx.Close()
		h = nil
		return
	}
	// message is a hello from a known peer -------------------------

	// MsgN must be 1
	msgN = m.GetMsgN()
	h.MsgN = msgN
	ck = m.GetCommsKey() // comms key as byte slice
	sk = m.GetSigKey()   // sig key as byte slice
	salt = m.GetSalt()
	sig = m.GetSig() // digital signature

	var serCK, serSK []byte

	if h.MsgN != 1 {
		err = ExpectedMsgOne
	}
	if err == nil {
		peerID := peer.GetNodeID().Value()
		if !bytes.Equal(id, peerID) {
			fmt.Println("NOT SAME NODE ID") // XXX should log
			err = NotExpectedNodeID
		}
	}
	if err == nil {
		serCK, err = xc.RSAPubKeyToWire(peer.GetCommsPublicKey())
		if err == nil {
			if !bytes.Equal(serCK, ck) {
				fmt.Println("NOT SAME COMMS KEY") // XXX should log
				err = NotExpectedCommsKey
			}
		}
	}
	if err == nil {
		serSK, err = xc.RSAPubKeyToWire(peer.GetSigPublicKey())
		if err == nil {
			if !bytes.Equal(serSK, sk) {
				fmt.Println("NOT SAME SIG KEY") // XXX should log
				err = NotExpectedSigKey
			}
		}
	}
	if err == nil {
		sigPubKey := peer.GetSigPublicKey()
		d := sha1.New()
		d.Write(id)
		d.Write(ck)
		d.Write(sk)
		d.Write(salt)
		hash := d.Sum(nil)
		err = rsa.VerifyPKCS1v15(sigPubKey, cr.SHA1, hash, sig)
	}
	if err == nil {
		// Everything is good; so Ack, leaving cnx open.
		h.Peer.MarkUp() // we consider the peer live
		h.Peer.LastContact()
		err = h.simpleAck(msgN)
	} else {
		// Send the text of the error to the peer; the send itself
		// may of course cause an error, but we will ignore that.
		// The peer is NOT marked as up.
		h.errorReply(err)
	}
	return
}
