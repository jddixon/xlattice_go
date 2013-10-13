package reg

// xlattice_go/reg/msg_handlers.go

import (
	"code.google.com/p/go.crypto/sha3"
	"crypto/rsa"
	"fmt"
	xc "github.com/jddixon/xlattice_go/crypto"
	xi "github.com/jddixon/xlattice_go/nodeID"
	xu "github.com/jddixon/xlattice_go/util"
)

var _ = fmt.Print

// XXX Possibly a problem, possibly not: the message number / sequence
// number has disappeared.

/////////////////////////////////////////////////////////////////////
// AES-BASED MESSAGE PAIRS
// All of these functions have the same signature, so that they can
// be invoked through a dispatch table.
/////////////////////////////////////////////////////////////////////

// Dispatch table entry where a client message received is inappropriate
// the the state of the connection.  For example, if we haven't yet
// received information about the client, we should not be receiving a
// Join or Get message.
func badCombo(h *InHandler) {
	h.errOut = RcvdInvalidMsgForState
}

// Handle the message which gives us information about the client and
// so associates this connection with a specific user.

func doClientMsg(h *InHandler) {
	var err error
	defer func() {
		h.errOut = err
	}()
	// Examine incoming message -------------------------------------
	var (
		name   string
		attrs  uint64
		nodeID *xi.NodeID
		ck, sk *rsa.PublicKey
		myEnds []string
		cm     *ClusterMember
	)
	// XXX We should accept EITHER clientName + token OR clientID
	// This implementation only accepts a token.

	clientMsg := h.msgIn
	name = clientMsg.GetClientName()
	clientSpecs := clientMsg.GetClientSpecs()
	attrs = clientSpecs.GetAttrs()
	ckBytes := clientSpecs.GetCommsKey()
	skBytes := clientSpecs.GetSigKey()

	if err == nil {
		ck, err = xc.RSAPubKeyFromWire(ckBytes)
		if err == nil {
			sk, err = xc.RSAPubKeyFromWire(skBytes)
			if err == nil {
				myEnds = clientSpecs.GetMyEnds() // a string array
			}
		}
	}
	if err == nil {
		id := clientSpecs.GetID()
		if id == nil {
			sha := sha3.NewKeccak256()
			sha.Write(ckBytes)
			sha.Write(skBytes)
			// XXX WE NEED SOME RANDOMNESS HERE!
			id = sha.Sum(nil)
		}
		nodeID, err = xi.New(id)
	}
	// Take appropriate action --------------------------------------
	if err == nil {
		// The appropriate action is to hang a token for this client off
		// the InHandler.
		cm, err = NewClusterMember(name, nodeID, ck, sk, attrs, myEnds)
		if err == nil {
			fmt.Printf("'thisMember' for connection is %s\n", cm.GetName())
			h.thisMember = cm
		}
	}
	if err == nil {
		// Prepare reply to client --------------------------------------
		// In this implementation We simply accept the client's proposed
		// attrs and ID.
		op := XLRegMsg_ClientOK
		h.msgOut = &XLRegMsg{
			Op:          &op,
			ClientID:    nodeID.Value(),
			ClientAttrs: &attrs, // in production, review and limit
		}
		// Set exit state -----------------------------------------------
		h.exitState = CLIENT_DETAILS_RCVD
		// DEBUG
		fmt.Printf("server has received client details and sent OK\n")
	}
}

// CREATE AND CREATE_REPLY ==========================================

// Handle the Create message which associates a unique name with a
// cluster and specifies its proposed size.  The server replies with the
// cluster ID and its server-assigned size.
//
// XXX This implementation does not handle cluster attrs.

func doCreateMsg(h *InHandler) {
	var err error
	defer func() {
		h.errOut = err
	}()
	// Examine incoming message -------------------------------------
	var clusterID *xi.NodeID
	var index int

	createMsg := h.msgIn
	clusterName := createMsg.GetClusterName()
	endPointCount := createMsg.GetEndPointCount()
	clusterSize := createMsg.GetClusterSize()

	// Take appropriate action --------------------------------------

	// Determine whether the cluster exists.  If it does, we will just
	// use its existing properties.

	h.reg.mu.RLock()
	cluster, exists := h.reg.ClustersByName[clusterName]
	h.reg.mu.RUnlock()

	if exists {
		// XXX THIS NO LONGER MAKES ANY SENSE
		h.cluster = cluster

		clusterSize = uint32(cluster.maxSize)
		clusterID, _ = xi.New(cluster.ID)
	} else {
		attrs := uint64(0)
		if clusterSize < 2 {
			clusterSize = 2
		} else if clusterSize > 64 {
			clusterSize = 64
		}
		// Assign a quasi-random cluster ID
		clusterID, _ = xi.New(nil)
		cluster, err = NewRegCluster(attrs, clusterName, clusterID,
			uint(endPointCount), uint(clusterSize))
		if err == nil {
			h.cluster = cluster
			index, err = h.reg.AddCluster(cluster)
			// XXX index not used
		}
		// DEBUG
		fmt.Printf("cluster %s has epCount %d\n", clusterName, endPointCount)
		// END
	}
	_ = index // INDEX IS NOT BEING USED

	if err == nil {
		// Prepare reply to client --------------------------------------
		op := XLRegMsg_CreateReply
		id := clusterID.Value() // XXX blows up
		h.msgOut = &XLRegMsg{
			Op:            &op,
			ClusterID:     id,
			EndPointCount: &endPointCount,
			ClusterSize:   &clusterSize,
		}
		// Set exit state -----------------------------------------------
		h.exitState = CREATE_REQUEST_RCVD
	}
}

// JOIN AND JOIN_REPLY ==============================================

// Tie this session to a specific cluster, either by supplying its
// name or using the clusterID.  Return the cluster ID and its size.
//

func doJoinMsg(h *InHandler) {
	var err error
	defer func() {
		h.errOut = err
	}()
	// Examine incoming message -------------------------------------
	var (
		cluster       *RegCluster
		clusterName   string
		clusterID     []byte
		endPointCount uint32
		clusterSize   uint32
	)
	joinMsg := h.msgIn

	// Take appropriate action --------------------------------------

	// Accept either cluster name or id.  If it's just the name,
	// attempt to retrieve the ID; it's an error if it does not exist
	// in the registry.  . In either case use the ID to retrieve the size.

	clusterName = joinMsg.GetClusterName() // will be "" if absent
	clusterID = joinMsg.GetClusterID()     // will be nil if absent

	if clusterID == nil && clusterName == "" {
		// if neither is present, we will use any cluster already
		// associated with this connection
		if h.cluster != nil {
			cluster = h.cluster
		} else {
			err = MissingClusterNameOrID
		}
	} else if clusterID != nil {
		// if an ID has Leen defined, we will try to use that
		h.reg.mu.RLock()
		cluster = h.reg.ClustersByID.FindBNI(clusterID).(*RegCluster)
		h.reg.mu.RUnlock()
		if cluster == nil {
			err = CantFindClusterByID
		}
	} else {
		// we have no ID and clusterName is not nil, so we will try to use that
		var ok bool
		h.reg.mu.RLock()
		if cluster, ok = h.reg.ClustersByName[clusterName]; !ok {
			err = CantFindClusterByName
		}
		h.reg.mu.RUnlock()
	}
	if err == nil {
		// if we get here, cluster is not nil
		h.cluster = cluster
		clusterID = cluster.ID
		endPointCount = uint32(h.cluster.epCount)
		clusterSize = uint32(h.cluster.maxSize)
		err = cluster.AddMember(h.thisMember)
	}
	if err == nil {
		// Prepare reply to client ----------------------------------
		// XXX If the cluster cannot be found, we must return an error
		// instead.
		op := XLRegMsg_JoinReply
		h.msgOut = &XLRegMsg{
			Op:            &op,
			ClusterID:     clusterID,
			EndPointCount: &endPointCount,
			ClusterSize:   &clusterSize,
		}
		// Set exit state -------------------------------------------
		h.exitState = JOIN_RCVD
	}
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

func doGetMsg(h *InHandler) {
	var err error
	defer func() {
		h.errOut = err
	}()
	// Examine incoming message -------------------------------------
	getMsg := h.msgIn
	clusterID := getMsg.GetClusterID()
	whichRequested := xu.NewBitMap64(getMsg.GetWhich())

	// Take appropriate action --------------------------------------
	var tokens []*XLRegMsg_Token
	whichReturned := xu.NewBitMap64(0)

	h.reg.mu.RLock()
	cluster := h.reg.ClustersByID.FindBNI(clusterID).(*RegCluster)
	h.reg.mu.RUnlock()

	if cluster == nil {
		err = CantFindClusterByID
	} else {
		size := uint(cluster.Size()) // actual size, not MaxSize
		if size > 64 {               // yes, should be impossible
			size = 64
		}
		weHave := xu.LowNMap(size)
		whichToSend := whichRequested.Intersection(weHave)
		// DEBUG
		fmt.Printf("doGetMsg: have 0x%x, client requests 0x%x, will send 0x%x\n",
			weHave.Bits, whichRequested.Bits, whichToSend.Bits)
		// END
		for i := uint(0); i < size; i++ {
			if whichToSend.Test(i) { // they want this one
				member := cluster.members[i]
				token, err := member.Token()
				if err == nil {
					tokens = append(tokens, token)
					whichReturned = whichReturned.Set(i)
				} else {
					// DEBUG
					fmt.Printf("ERROR seen while tokenizing member %d, %s\n",
						i, member.GetName())
					// END
					break
				}
			}
		}
	}
	if err == nil {
		// Prepare reply to client --------------------------------------
		op := XLRegMsg_ClusterMembers
		h.msgOut = &XLRegMsg{
			Op:        &op,
			ClusterID: clusterID,
			Which:     &whichReturned.Bits,
			Tokens:    tokens,
		}
		// Set exit state -----------------------------------------------
		h.exitState = JOIN_RCVD // the JOIN is intentional !

		// DEBUG
		fmt.Printf("server returning %d tokens\n", whichReturned.Count())
		// END
	}
}

// BYE AND ACK ======================================================

func doByeMsg(h *InHandler) {
	var err error
	defer func() {
		h.errOut = err
	}()

	// Examine incoming message -------------------------------------
	//ByeMsg := h.msgIn

	// Take appropriate action --------------------------------------
	// nothing to do

	// Prepare reply to client --------------------------------------
	op := XLRegMsg_Ack
	h.msgOut = &XLRegMsg{
		Op: &op,
	}
	// Set exit state -----------------------------------------------
	h.exitState = BYE_RCVD
}
