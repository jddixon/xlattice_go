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

// XXX Possibly a problem, possibly not: the message number / sequence
// number has disappeared.

/////////////////////////////////////////////////////////////////////
// AES-BASED MESSAGE PAIRS
// All of these functions have the same signature, so that they can
// be invoked through a table.
/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// PENDING CHANGE: ALL OF THESE SHOULD TAKE *InHandler AS AN ARGUMENT
// so that the dispatch table doesn't have to be rebuilt over and over
// again.
/////////////////////////////////////////////////////////////////////

func (h *InHandler) badCombo() {
	h.errOut = InvalidMsgInForState
}

// Handle the client message which opens the session by identifying
// the caller.

func (h *InHandler) doClientMsg() {
	var err error
	defer func() {
		h.errOut = err
	}()
	// Examine incoming message -------------------------------------

	// XXX We should accept EITHER clientName + token OR clientID
	// This implementation does neither!

	clientMsg := h.msgIn
	name := clientMsg.GetClientName()
	clientSpecs := clientMsg.GetClientSpecs()
	attrs := clientSpecs.GetAttrs()
	nodeID, err := xi.New(clientSpecs.GetID())
	if err != nil {
		// XXX In this approach to error handling, any error returned
		// should cause an error message to ge sent to the client and
		// h.exitState to be set to IN_CLOSED after the error message
		// has been sent.
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

	// Take appropriate action --------------------------------------
	cm, err := NewClusterMember(name, nodeID, ck, sk, attrs, myEnds)
	if err != nil {
		return
	}
	h.thisMember = cm

	// Prepare reply to client --------------------------------------
	// We simply accept the client's proposed attrs and ID.
	op := XLRegMsg_ClientOK
	h.msgOut = &XLRegMsg{
		Op:       &op,
		ClientID: nodeID.Value(),
		Attrs:    &attrs, // in production, review and limit
	}
	// Set exit state -----------------------------------------------
	h.exitState = CLIENT_DETAILS_RCVD
}

// CREATE AND CREATE_REPLY ==========================================

// Handle the Create message which associates a unique name with a
// cluster and specifies its proposed size.  The server replies with the
// cluster ID and its server-assigned size.

func (h *InHandler) doCreateMsg() {
	var err error
	defer func() {
		h.errOut = err
	}()
	// Examine incoming message -------------------------------------
	createMsg := h.msgIn
	clusterName := createMsg.GetClusterName()
	clusterSize := createMsg.GetClusterSize()

	// Take appropriate action --------------------------------------

	// XXX STUB ? error if name already in use

	if clusterSize < 2 {
		clusterSize = 2
	} else if clusterSize > 64 {
		clusterSize = 64
	}

	_, _, _ = createMsg, clusterName, clusterSize

	// Assign a quasi-random cluster ID
	nID, _ := xi.New(nil)

	// XXX STUB: add cluster to registry; both name and ID must be
	// unique

	// Prepare reply to client --------------------------------------
	op := XLRegMsg_CreateReply
	h.msgOut = &XLRegMsg{
		Op:          &op,
		ClusterID:   nID.Value(),
		ClusterSize: &clusterSize,
	}
	// Set exit state -----------------------------------------------
	h.exitState = CREATE_REQUEST_RCVD
}

// JOIN AND JOIN_REPLY ==============================================

// Tie this session to a specific cluster, either by supplying its
// name or using the clusterID.  Return the cluster ID and its size.
//

func (h *InHandler) doJoinMsg() {
	var err error
	defer func() {
		h.errOut = err
	}()
	// Examine incoming message -------------------------------------
	var (
		clusterName string
		clusterID   []byte
		clusterSize uint32
	)
	joinMsg := h.msgIn

	// Take appropriate action --------------------------------------

	// XXX Accept either cluster name or id.  If it's just the name,
	// attempt to retrieve the ID; it's an error if it does not exist
	// in the registry.  . In either case use the ID to retrieve the size.

	// XXX STUB

	_, _, _ = joinMsg, clusterName, clusterID

	// Prepare reply to client --------------------------------------
	// XXX If the cluster cannot be found, we will return an error
	// instead.
	op := XLRegMsg_JoinReply
	h.msgOut = &XLRegMsg{
		Op:          &op,
		ClusterID:   clusterID,
		ClusterSize: &clusterSize,
	}
	// Set exit state -----------------------------------------------
	h.exitState = JOIN_RCVD
}

// GET AND MEMBERS ==================================================

// Fetch from the registry details for the specified members for the
// cluster.  The cluster is identified by its ID.  Members requested
// are specified using a bit vector: we assume that members are stored
// in the order in which they joined, so if the Nth bit is set, we
// want a copy of the details for that member.  It is an error if the
// clusterID does not correspond to an existing cluster.  It is not
// an error if a member cannot be found for one of the bits set: the
// server returns a bit vector specifying which member tokens are being
// returned.

func (h *InHandler) doGetMsg() {
	var err error
	defer func() {
		h.errOut = err
	}()
	// Examine incoming message -------------------------------------
	getMsg := h.msgIn
	clusterID := getMsg.GetClusterID()
	whichIn := getMsg.GetWhich()

	// Take appropriate action --------------------------------------

	var tokens []*XLRegMsg
	var whichOut uint64

	_, _, _ = clusterID, whichIn, tokens

	// XXX STUB: collect these, setting bits in whichOut

	// Prepare reply to client --------------------------------------
	// XXX If the cluster cannot be found, we will return an error
	// instead.
	op := XLRegMsg_Members
	h.msgOut = &XLRegMsg{
		Op:        &op,
		ClusterID: clusterID,
		Which:     &whichOut,
		//Tokens:
	}
	// Set exit state -----------------------------------------------
	h.exitState = JOIN_RCVD // this is intentional !
}

// BYE AND ACK ======================================================

func (h *InHandler) doByeMsg() {
	var err error
	defer func() {
		h.errOut = err
	}()
	// Examine incoming message -------------------------------------
	//ByeMsg := h.msgIn

	// Take appropriate action --------------------------------------

	// Prepare reply to client --------------------------------------
	op := XLRegMsg_Ack
	h.msgOut = &XLRegMsg{
		Op: &op,
	}
	// Set exit state -----------------------------------------------
	h.exitState = BYE_RCVD
}

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

// func (h *InHandler) doJoinMsg() (err error) {
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
//				//fmt.Printf("doJoinMsg: UNEXPECTED MESSAGE TYPE %v\n", m.GetOp())
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
//			fmt.Printf("    doJoinMsg gets %v\n", err) // DEBUG
//		}
//	} // GEEP
