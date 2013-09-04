package reg

// xlattice_go/reg/inHandler.go

/////////////////////////////////////////////////////////
// YANKED FROM ../msg, BEING HACKED INTO SOMETHING USEFUL
/////////////////////////////////////////////////////////

import (
	//cr "crypto"
	//"crypto/rsa"
	// "crypto/sha1"
	"fmt"
	//xc "github.com/jddixon/xlattice_go/crypto"
	xm "github.com/jddixon/xlattice_go/msg"
	xn "github.com/jddixon/xlattice_go/node"
	xt "github.com/jddixon/xlattice_go/transport"
	// xu "github.com/jddixon/xlattice_go/util"
//	"strings"
)

var _ = fmt.Print

// States through which the input cnx may pass
const (
	IN_START = iota
	HELLO_RCVD
	REPLY_SENT
	DETAILS_RCVD
	PEER_INFO_SENT
	IN_CLOSED
)

type InHandler struct {
	iv1, aes1, iv2, aes2 []byte
	CnxHandler
}

// Given an open new connection, process a hello message for this node,
// returning a handler for the connection if that succeeds.  The hello
// consists of an AES Key+IV and a salt which we require to be eight
// bytes long.

func NewInHandler(n *xn.Node, conn xt.ConnectionI) (h *InHandler, err error) {

	if n == nil {
		return nil, xm.NilNode
	}
	if conn == nil {
		return nil, xm.NilConnection
	}
	cnx := conn.(*xt.TcpConnection)
	h = &InHandler{
		CnxHandler: CnxHandler{Cnx: cnx, State: IN_START}}
	err = h.handleHello(n) // sends hello reply unless error
	if err == nil {
		err = h.handleJoin() // sends peer tokens unless error
	}
	return
}

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
//}

func (h *InHandler) handleHello(n *xn.Node) (err error) {
	var (
		m *XLRegMsg
		//id, ck, sk, sig, salt []byte
		//peer                  *xn.Peer
	)
	m, err = h.readMsg()

	// message must be a Hello
	if err == nil {
		if m.GetOp() != XLRegMsg_Hello {
			err = xm.MissingHello
		}
	}

	// On any error up to here silently close the connection and delete
	// the handler.
	if err != nil {
		h.Cnx.Close()
		h = nil
		return
	}
	// message is a hello -------------------------------------------

	//	// XXX THIS CODE DECIPHERS MyDesc
	//	ck = m.GetCommsKey() // comms key as byte slice
	//	sk = m.GetSigKey()   // sig key as byte slice
	//	salt = m.GetSalt1()
	//	// sig = m.GetSig() // digital signature
	//	_ = sig	// ???
	//
	//	var serCK, serSK []byte
	//
	//	if err == nil {
	//		peerID := peer.GetNodeID().Value()
	//		if !xu.SameBytes(id, peerID) {
	//			fmt.Println("NOT SAME NODE ID") // XXX should log
	//			err = xm.NotExpectedNodeID
	//		}
	//	}
	//	if err == nil {
	//		serCK, err = xc.RSAPubKeyToWire(peer.GetCommsPublicKey())
	//		if err == nil {
	//			if !xu.SameBytes(serCK, ck) {
	//				fmt.Println("NOT SAME COMMS KEY") // XXX should log
	//				err = xm.NotExpectedCommsKey
	//			}
	//		}
	//	}
	//	if err == nil {
	//		serSK, err = xc.RSAPubKeyToWire(peer.GetSigPublicKey())
	//		if err == nil {
	//			if !xu.SameBytes(serSK, sk) {
	//				fmt.Println("NOT SAME SIG KEY") // XXX should log
	//				err = xm.NotExpectedSigKey
	//			}
	//		}
	//	}
	//	if err == nil {
	//		sigPubKey := peer.GetSigPublicKey()
	//		d := sha1.New()
	//		d.Write(id)
	//		d.Write(ck)
	//		d.Write(sk)
	//		d.Write(salt)
	//		hash := d.Sum(nil)
	//		err = rsa.VerifyPKCS1v15(sigPubKey, cr.SHA1, hash, sig)
	//	} // FOO
	//	if err == nil {
	//		// Everything is good; so Ack, leaving cnx open.
	//		h.Peer.MarkUp() // we consider the peer live
	//		h.Peer.LastContact()
	//		err = h.simpleAck(msgN)
	//	} else {
	//		// Send the text of the error to the peer; the send itself
	//		// may of course cause an error, but we will ignore that.
	//		// The peer is NOT marked as up.
	//		h.errorReply(err)
	//	}
	return
}

func (h *InHandler) handleJoin() (err error) {
	defer h.Cnx.Close()
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
	//				err = xm.UnexpectedMsgType
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
