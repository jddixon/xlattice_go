package reg

// xlattice_go/reg/in_handler.go

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	// xc "github.com/jddixon/xlattice_go/crypto"
	"github.com/jddixon/xlattice_go/msg"
	// xi "github.com/jddixon/xlattice_go/nodeID"
	xt "github.com/jddixon/xlattice_go/transport"
)

var _ = fmt.Print

// States through which the input cnx may pass
const (
	HELLO_RCVD = iota
	CLIENT_DETAILS_RCVD
	CLUSTER_REQUEST_RCVD
	JOIN_RCVD
	BYE_RCVD
	IN_CLOSED
)

const (
	IN_STATE_COUNT    = 5
	MSG_HANDLER_COUNT = 5
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
	entryState                         int
	exitState                          int
	msgIn                              *XLRegMsg
	msgOut                             *XLRegMsg
	msgHandlers                        [][]interface{}
	errOut                             error
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
				Cnx: cnx,
			},
		}
	}
	h.msgHandlers = make([][]interface{}, IN_STATE_COUNT, MSG_HANDLER_COUNT)
	// XXX INCOMPLETE INITIALIZATION
	h.msgHandlers[0] = []interface{}{
		h.doClientMsg, h.badCombo, h.badCombo, h.badCombo, h.badCombo}
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
	if err != nil {
		return
	}
	err = h.SetUpSessionKey()
	if err != nil {
		return
	}
	for {
		// convert raw data off the wire into an XLRegMsg object
		var ciphertext []byte
		ciphertext, err = h.readData()
		if err != nil {
			return
		}
		// receive, decode the client request
		h.msgIn, err = DecryptUnpadDecode(ciphertext, h.decrypterS)
		if err != nil {
			return
		}
		// SUPERFLUOUS: this is now handled through the table
		if h.msgIn.GetOp() != XLRegMsg_Client {
			err = UnexpectedMsgType
			return
		}

		h.msgHandlers[h.entryState][0].(func())()
		if h.errOut != nil {
			return
		}

		// encode, pad, and encrypt the XLRegMsg object, then put it on the wire

		// XXX STUB
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

/////////////////////////////////////////////////////////////////////
// RSA-BASED MESSAGE PAIR
/////////////////////////////////////////////////////////////////////

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
		iv1, key1, salt1, version1,
			err = msg.ServerDecodeHello(ciphertext, rn.ckPriv)
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
