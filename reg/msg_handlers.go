package reg

// xlattice_go/reg/msg_handlers.go

import (
	//"crypto/aes"
	//"crypto/cipher"
	"fmt"
	xc "github.com/jddixon/xlattice_go/crypto"
	//"github.com/jddixon/xlattice_go/msg"
	xi "github.com/jddixon/xlattice_go/nodeID"
	//xt "github.com/jddixon/xlattice_go/transport"
)

var _ = fmt.Print

// THIS IS OLD CODE, but should convert easily to new pattern:

// Send the text of the error message to the peer and close the connection.
func (h *InHandler) errorReply(e error) (err error) {
	var reply XLRegMsg
	cmd := XLRegMsg_Error
	s := e.Error()
	reply.Op = &cmd
	reply.ErrDesc = &s
	h.writeMsg(&reply) // ignore any write error
	h.State = IN_CLOSED

	// XXX This would be a very strong action, given that we may have multiple
	// connections open to this peer.
	// h.Peer.MarkDown()

	return
}

// func (h *InHandler) simpleAck(msgN uint64) (err error) {
//
// 	var reply XLRegMsg
// 	cmd := XLRegMsg_Ack
// 	reply.Op = &cmd
// 	err = h.writeMsg(&reply) // this may yield an error ...
// 	h.State = HELLO_RCVD
// 	return err
// } //
//
//} GEEP

/////////////////////////////////////////////////////////////////////
// AES-BASED MESSAGE PAIRS
// All of these functions have the same signature, so that they can
// be invoked through a table.
/////////////////////////////////////////////////////////////////////

func (h *InHandler) badCombo() {
	h.errOut = InvalidMsgInForState
}
func (h *InHandler) doClientMsg() {

	var err error
	defer func() {
		h.errOut = err
	}()

	clientMsg := h.msgIn
	name := clientMsg.GetClientName()
	clientSpecs := clientMsg.GetClientSpecs()
	attrs := clientSpecs.GetAttrs()
	id, err := xi.New(clientSpecs.GetID())
	if err != nil {
		return
	}
	ck, err := xc.RSAPubKeyFromWire(clientSpecs.GetCommsKey())
	if err != nil {
		return
	}
	sk, err := xc.RSAPubKeyFromWire(clientSpecs.GetSigKey())
	if err != nil {
		return
	}
	myEnds := clientSpecs.GetMyEnds() // a string array
	cm, err := NewClusterMember(name, id, ck, sk, attrs, myEnds)
	if err != nil {
		return
	}
	h.thisMember = cm

	// Answer with ClientOK or error ----------------------

	// XXX STUB XXX

	// END HANDLE CLIENT

	return
}

func (h *InHandler) handleJoin() (err error) {
	//	go func() {
	//		<-h.StopCh
	//		h.Cnx.Close()
	//		h.StoppedCh <- true
	//	}()
	//	for err == nil {
	//		var m *XLRegMsg
	//		m, err = h.readMsg()
	//		if err == nil {
	//			cmd := m.GetOp()
	//			switch cmd {
	//			case XLRegMsg_Bye:
	//				err = h.checkMsgNbrAndAck(m)
	//				h.State = IN_CLOSED
	//			case XLRegMsg_KeepAlive:
	//				// XXX Update last-time-spoken-to for peer
	//				h.Peer.LastContact()
	//				err = h.checkMsgNbrAndAck(m)
	//			default:
	//				// XXX should log
	//				//fmt.Printf("handleJoin: UNEXPECTED MESSAGE TYPE %v\n", m.GetOp())
	//				err = msg.UnexpectedMsgType
	//				h.errorReply(err) // ignore any errors from the call itself
	//			}
	//		} else {
	//			break
	//		}
	//	}
	//	if err != nil {
	//		text := err.Error()
	//		if strings.HasSuffix(text, "use of closed network connection") {
	//			err = nil
	//		} else {
	//			fmt.Printf("    handleJoin gets %v\n", err) // DEBUG
	//		}
	//	}
	return // GEEP
}
