package reg

// xlattice_go/reg/in_handler.go

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	xc "github.com/jddixon/xlattice_go/crypto"
	"github.com/jddixon/xlattice_go/msg"
	xi "github.com/jddixon/xlattice_go/nodeID"
	xt "github.com/jddixon/xlattice_go/transport"
)

var _ = fmt.Print

// States through which the input cnx may pass
const (
	IN_START = iota
	HELLO_RCVD
	CLIENT_DETAILS_RCVD
	CLUSTER_REQUEST_RCVD
	JOIN_RCVD
	GETTING
	BYE_RCVD
	IN_CLOSED
)

type InHandler struct {
	iv1, key1, iv2, key2, salt1, salt2 []byte
	engineS                            cipher.Block
	encrypterS                         cipher.BlockMode
	decrypterS                         cipher.BlockMode
	reg                                *Registry
	thisMember                         *ClusterMember
	clusterName                        string
	clusterID                          []byte
	clusterSize                        int
	version                            uint32 // protocol version used in session
	known                              uint64 // a bit vector:
	state                              int
	CnxHandler
}

// Given an open new connection, process a hello message for this node,
// returning a handler for the connection if that succeeds.  The hello
// consists of an AES Key+IV, a salt, and a requested protocol version.
// The salt must be at least eight bytes long.

func NewInHandler(reg *Registry, conn xt.ConnectionI) (
	h *InHandler, err error) {

	if reg == nil {
		return nil, NilRegistry
	}
	rn := &reg.RegNode
	if rn == nil {
		err = msg.NilNode
	} else if conn == nil {
		err = msg.NilConnection
	} else {
		cnx := conn.(*xt.TcpConnection)
		h = &InHandler{
			CnxHandler: CnxHandler{
				Cnx:   cnx,
				State: IN_START,
			},
		}
	}
	return
}

func (h *InHandler) SetUpSessionKey() (err error) {
	h.engineS, err = aes.NewCipher(h.key2)
	if err == nil {
		h.encrypterS = cipher.NewCBCEncrypter(h.engineS, h.iv2)
		h.decrypterS = cipher.NewCBCDecrypter(h.engineS, h.iv2)
	}
	return
}
func (h *InHandler) Run() (err error) {

	defer func() {
		if h.Cnx != nil {
			h.Cnx.Close()
		}
	}()

	err = h.handleHello() // this adds iv2, key2 to handler
	if err == nil {
		err = h.SetUpSessionKey()
	}
	if err == nil {
		err = h.handleClientMsg()
	}

	// Expect CreateMsg ---------------------------------------------

	// Answer with CreateReply ----------------------------

	// Expect JoinMsg -----------------------------------------------

	// Answer with JoinReply ------------------------------

	// Expect Get ---------------------------------------------------

	// Answer with Members --------------------------------

	// Repeat Get/Members or Expect Bye -----------------------------

	// Send Ack -------------------------------------------

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

// The client has sent the server a one-time AES key+iv encrypted with
// the server's RSA comms public key.  The server creates the real
// session iv+key and returns them to the client encrypted with the
// one-time key+iv.

func (h *InHandler) handleHello() (err error) {
	var (
		ciphertext, iv1, key1, salt1 []byte
		version1                     uint32
	)
	rn := &h.reg.RegNode
	ciphertext, err = h.readData()
	if err == nil {
		iv1, key1, salt1, version1, err = msg.ServerDecodeHello(ciphertext, rn.ckPriv)
	}
	if err == nil {
		version2 := version1 // accept whatever version they propose
		iv2, key2, salt2, ciphertextOut, err := msg.ServerEncodeHelloReply(
			iv1, key1, salt1, version2)
		if err == nil {
			err = h.writeData(ciphertextOut)
		}
		if err == nil {
			h.iv1 = iv1
			h.key1 = key1
			h.iv2 = iv2
			h.key2 = key2
			h.salt1 = salt1
			h.salt2 = salt2
			h.version = version2
			h.State = HELLO_RCVD
		}
	}
	// On any error silently close the connection and delete the handler,
	// an exciting thing to do.
	if err != nil {
		h.Cnx.Close()
		h = nil
		return
	}
	return
}
func (h *InHandler) handleClientMsg() (err error) {
	// BEGIN HANDLE CLIENT
	var (
		clientMsg *XLRegMsg
	)

	var ciphertext []byte
	ciphertext, err = h.readData()
	if err != nil {
		return
	}
	clientMsg, err = DecryptUnpadDecode(ciphertext, h.decrypterS)
	if err != nil {
		return
	}
	if clientMsg.GetOp() != XLRegMsg_Client {
		err = UnexpectedMsgType
		return
	}

	name := clientMsg.GetClientName()
	clientSpecs := clientMsg.GetClientSpecs()
	attrs := clientSpecs.GetAttrs()
	id, err := xi.New(clientSpecs.GetID())
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

	// XXX STUB

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
