package reg

// xlattice_go/reg/client_node.go

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	xc "github.com/jddixon/xlattice_go/crypto"
	xm "github.com/jddixon/xlattice_go/msg"
	xn "github.com/jddixon/xlattice_go/node"
	xi "github.com/jddixon/xlattice_go/nodeID"
	xt "github.com/jddixon/xlattice_go/transport"
	xu "github.com/jddixon/xlattice_go/util"
	// "io"
	"time"
)

var _ = fmt.Print // DEBUG

// client states
const (
	CLIENT_START = iota
	HELLO_SENT
	CLIENT_SENT
	CLUSTER_SENT
	JOIN_SENT
	GET_SENT
	BYE_SENT
	CLIENT_CLOSED
)

type ClientNode struct {

	// EPHEMERAL INFORMATION ========================================
	doneCh          chan bool
	err             error
	h               *CnxHandler
	proposedAttrs   uint64
	proposedVersion uint32 // proposed by client
	iv1, key1       []byte // one-shot
	iv2, key2       []byte // session
	salt1, salt2    []byte // not currently used
	engineC         cipher.Block
	encrypterC      cipher.BlockMode
	decrypterC      cipher.BlockMode

	// SOMETIMES (?) PERSISTED ======================================
	// The significance of these fields is different in different
	// subclasses.  To the user, this is the cluster being joined.
	// To an admin client, it is a cluster being created.

	clusterName  string
	clusterAttrs uint64
	clusterID    *xi.NodeID
	clusterSize  uint32 // this is a FIXED size, aka MaxSize

	// PERSISTED ====================================================
	// If the ClientNode is not ephemeral, this information is saved
	// as LFS/.xlattice/client.node

	serverName     string
	serverID       *xi.NodeID
	serverEnd      xt.EndPointI
	serverCK       *rsa.PublicKey
	serverSK       *rsa.PublicKey
	decidedAttrs   uint64
	decidedVersion uint32           // decreed by server
	members        []*ClusterMember // information on cluster members

	// By convention endPoints[0] is used for member-member communications
	// and [1] for comms with cluster clients, should they exist. Some or
	// all of these (the first epCount) are passed to other cluster
	// members via the registry.
	epCount uint32

	// INFORMATION PERSISTED AS Node CONFIGURATION ------------------
	// This is used to build the node and so is persisted as part of
	// the node when that is saved.
	endPoints      []xt.EndPointI
	lfs            string
	name           string
	clientID       *xi.NodeID
	ckPriv, skPriv *rsa.PrivateKey

	// XLattice Node ------------------------------------------------
	// This is created during the first session and serialized to
	// LFS/.xlattice/node.config if LFS != ""
	xn.Node
}

// Given contact information for a registry and the name of a cluster,
// the client joins the cluster, collects information on the other members,
// and terminates when it has info on the entire membership.

func NewClientNode(
	name, lfs string, attrs uint64,
	serverName string, serverID *xi.NodeID, serverEnd xt.EndPointI,
	serverCK, serverSK *rsa.PublicKey,
	clusterName string, clusterAttrs uint64, clusterID *xi.NodeID, size int,
	epCount int, e []xt.EndPointI) (
	cn *ClientNode, err error) {

	var (
		isAdmin        = (attrs & ATTR_ADMIN) != 0
		ckPriv, skPriv *rsa.PrivateKey
		node           *xn.Node
	)

	// sanity checks on parameter list
	if serverName == "" || serverID == nil || serverEnd == nil ||
		serverCK == nil {

		err = MissingServerInfo
	}
	if err == nil && clusterName == "" {
		err = MissingClusterNameOrID
	}
	if err == nil && size < 2 {
		err = ClusterMustHaveTwo
	}
	if err == nil {
		// if the client is an admin client epCount applies to the cluster
		if epCount < 1 {
			epCount = 1
		}
		if !isAdmin {
			// XXX There is some confusion here: we don't require that
			// all members have the same number of endpoints
			actualEPCount := len(e)
			if actualEPCount == 0 {
				err = ClientMustHaveEndPoint
			} else if epCount > actualEPCount {
				epCount = actualEPCount
			}
		}
		// XXX This is a gross simplification.  If lfs is not specified,
		// this is an ephemeral node.  If lfs IS specified and
		// configuration files are present, we should deserialize the
		// configuration files, which creates the node.
		if lfs == "" {
			ckPriv, err = rsa.GenerateKey(rand.Reader, 2048)
			if err == nil {
				skPriv, err = rsa.GenerateKey(rand.Reader, 2048)
			}
		}
	}

	if err == nil && node == nil && (ckPriv == nil || skPriv == nil) {
		err = NoNodeNoKeys
	}
	if err == nil {
		cnxHandler := &CnxHandler{State: CLIENT_START}
		cn = &ClientNode{
			name:          name,
			lfs:           lfs, // if blank, node is ephemeral
			proposedAttrs: attrs,
			doneCh:        make(chan bool, 1),
			serverName:    serverName,
			serverID:      serverID,
			serverEnd:     serverEnd,
			serverCK:      serverCK,
			serverSK:      serverSK,
			clusterName:   clusterName,
			clusterAttrs:  clusterAttrs,
			clusterID:     clusterID,
			clusterSize:   uint32(size),
			h:             cnxHandler,
			epCount:       uint32(epCount),
			endPoints:     e,
			ckPriv:        ckPriv,
			skPriv:        skPriv,

			// Node NOT YET INITIALIZED
		}
	}
	return
}

// Read the next message over the connection
func (cn *ClientNode) readMsg() (m *XLRegMsg, err error) {
	inBuf, err := cn.h.readData()
	if err == nil && inBuf != nil {
		m, err = DecryptUnpadDecode(inBuf, cn.decrypterC)
	}
	return
}

// Write a message out over the connection
func (cn *ClientNode) writeMsg(m *XLRegMsg) (err error) {
	var data []byte
	// serialize, marshal the message
	data, err = EncodePadEncrypt(m, cn.encrypterC)
	if err == nil {
		err = cn.h.writeData(data)
	}
	return
}

// RUN CODE =========================================================

// Subclasses (UserClient, AdminClient, etc) use sequences of calls to
// these these functions to accomplish their purposes.

func (cn *ClientNode) SessionSetup(proposedVersion uint32) (
	cnx *xt.TcpConnection, decidedVersion uint32, err error) {
	var (
		ciphertext1, iv1, key1, salt1, salt1c []byte
		ciphertext2, iv2, key2, salt2         []byte
	)
	// Set up connection to server. -----------------------------
	ctor, err := xt.NewTcpConnector(cn.serverEnd)
	if err == nil {
		var conn xt.ConnectionI
		conn, err = ctor.Connect(nil)
		if err == nil {
			cnx = conn.(*xt.TcpConnection)
		}
	}
	// Send HELLO -----------------------------------------------
	if err == nil {
		cn.h.Cnx = cnx
		ciphertext1, iv1, key1, salt1,
			err = xm.ClientEncodeHello(proposedVersion, cn.serverCK)
	}
	if err == nil {
		err = cn.h.writeData(ciphertext1)
		// DEBUG
		if err != nil {
			fmt.Printf("SessionSetup, sending Hello: err after write is %v\n",
				err)
		}
		// END
	}
	// Process HELLO REPLY --------------------------------------
	if err == nil {
		ciphertext2, err = cn.h.readData()
		// DEBUG
		if err != nil {
			fmt.Printf("SessionSetup, sending Hello: err after read is %v\n",
				err)
		}
		// END
	}
	if err == nil {
		iv2, key2, salt2, salt1c, decidedVersion,
			err = xm.ClientDecodeHelloReply(ciphertext2, iv1, key1)
		_ = salt1c // XXX
	}
	// Set up AES engines ---------------------------------------
	if err == nil {
		cn.salt1 = salt1
		cn.iv2 = iv2
		cn.key2 = key2
		cn.salt2 = salt2
		cn.decidedVersion = decidedVersion
		cn.engineC, err = aes.NewCipher(key2)
		if err == nil {
			cn.encrypterC = cipher.NewCBCEncrypter(cn.engineC, iv2)
			cn.decrypterC = cipher.NewCBCDecrypter(cn.engineC, iv2)
		}
		// DEBUG
		fmt.Printf("client %s AES engines set up\n", cn.name)
		// END
	}
	return
}

func (cn *ClientNode) ClientAndOK() (err error) {

	var (
		ckBytes, skBytes []byte
		myEnds           []string
	)
	// XXX attrs not dealt with

	// Send CLIENT MSG ==========================================
	//ckBytes, err = xc.RSAPubKeyToWire(cn.GetCommsPublicKey())
	ckBytes, err = xc.RSAPubKeyToWire(&cn.ckPriv.PublicKey)
	if err == nil {
		//skBytes, err = xc.RSAPubKeyToWire(cn.GetSigPublicKey())
		skBytes, err = xc.RSAPubKeyToWire(&cn.skPriv.PublicKey)
		if err == nil {
			for i := 0; i < len(cn.endPoints); i++ {
				myEnds = append(myEnds, cn.endPoints[i].String())
			}
			token := &XLRegMsg_Token{
				Name:     &cn.name,
				Attrs:    &cn.proposedAttrs,
				CommsKey: ckBytes,
				SigKey:   skBytes,
				MyEnds:   myEnds,
			}

			op := XLRegMsg_Client
			request := &XLRegMsg{
				Op:          &op,
				ClientName:  &cn.name, // XXX redundant
				ClientSpecs: token,
			}
			// SHOULD CHECK FOR TIMEOUT
			err = cn.writeMsg(request)
			// DEBUG
			fmt.Printf("CLIENT_MSG for %s sent\n", cn.name)
			// END
		}
	}
	// Process CLIENT_OK --------------------------------------------
	// SHOULD CHECK FOR TIMEOUT
	response, err := cn.readMsg()
	if err == nil {
		id := response.GetClientID()
		cn.clientID, err = xi.New(id)

		// XXX err ignored

		cn.decidedAttrs = response.GetClientAttrs()
		// DEBUG
		fmt.Printf("    client %s has received ClientOK\n",
			cn.name)
		// END
	}
	return
} // GEEP2

func (cn *ClientNode) CreateAndReply() (err error) {

	var response *XLRegMsg

	// Send CREATE MSG ==========================================
	op := XLRegMsg_Create
	request := &XLRegMsg{
		Op:            &op,
		ClusterName:   &cn.clusterName,
		ClusterAttrs:  &cn.clusterAttrs,
		ClusterSize:   &cn.clusterSize,
		EndPointCount: &cn.epCount,
	}
	// SHOULD CHECK FOR TIMEOUT
	err = cn.writeMsg(request)
	// DEBUG
	fmt.Printf("client %s sends CREATE for cluster %s, epCount %d, size %d\n",
		cn.name, cn.clusterName, cn.epCount, cn.clusterSize)
	// END

	if err == nil {
		// Process CREATE REPLY -------------------------------------
		// SHOULD CHECK FOR TIMEOUT AND VERIFY THAT IT'S A CREATE REPLY
		response, err = cn.readMsg()
		op = response.GetOp()
		_ = op
		// DEBUG
		fmt.Printf("    client has received CreateReply; err is %v\n", err)
		// END
		if err == nil {
			id := response.GetClusterID()
			cn.clusterID, err = xi.New(id)
			cn.clusterAttrs = response.GetClusterAttrs()
			cn.clusterSize = response.GetClusterSize()
			// XXX no check on err
		}
	}
	return
}

func (cn *ClientNode) JoinAndReply() (err error) {

	// Send JOIN MSG ============================================
	op := XLRegMsg_Join
	request := &XLRegMsg{
		Op:          &op,
		ClusterName: &cn.clusterName,
	}
	// SHOULD CHECK FOR TIMEOUT
	err = cn.writeMsg(request)
	// DEBUG
	fmt.Printf("Client %s sends JOIN by name cluster %s\n",
		cn.name, cn.clusterName)
	// END

	// Process JOIN REPLY ---------------------------------------
	if err == nil {
		var response *XLRegMsg

		// SHOULD CHECK FOR TIMEOUT AND VERIFY THAT IT'S A JOIN REPLY
		response, err = cn.readMsg()
		op := response.GetOp()
		_ = op

		epCount := uint32(response.GetEndPointCount())
		// DEBUG
		fmt.Printf("    client has received JoinReply; epCount %d, err is %v\n",
			epCount, err)
		// END
		if err == nil {
			clusterSizeNow := response.GetClusterSize()
			if cn.clusterSize != clusterSizeNow {
				cn.clusterSize = clusterSizeNow
				cn.members = make([]*ClusterMember, cn.clusterSize)
			}
			cn.epCount = epCount
			// XXX This is just wrong: we already know the cluster ID
			id := response.GetClusterID()
			// cn.clusterID, err = xi.New(id)
			_ = id // DO SOMETHING WITH IT  XXX
		}
	} // GEEP3
	return
}

// Collect information on all cluster members
func (cn *ClientNode) GetAndMembers() (err error) {

	if cn.clusterID == nil {
		fmt.Printf("** ENTERING GetAndMembers for %s with nil clusterID! **\n",
			cn.name)
	}
	MAX_GET := 16
	if cn.members == nil {
		cn.members = make([]*ClusterMember, cn.clusterSize)
	}
	stillToGet := xu.LowNMap(uint(cn.clusterSize))
	for count := 0; count < MAX_GET && stillToGet.Any(); count++ {
		var response *XLRegMsg

		for i := uint(0); i < uint(cn.clusterSize); i++ {
			if cn.members[i] != nil {
				stillToGet = stillToGet.Clear(i)
			}
		}
		// DEBUG
		fmt.Printf("ClientNode %s sends GET for %d members (bits 0x%x)\n",
			cn.name, stillToGet.Count(), stillToGet.Bits)
		// END

		// Send GET MSG =========================================
		op := XLRegMsg_GetCluster
		request := &XLRegMsg{
			Op:        &op,
			ClusterID: cn.clusterID.Value(),
			Which:     &stillToGet.Bits,
		}
		// SHOULD CHECK FOR TIMEOUT
		err = cn.writeMsg(request)

		// Process MEMBERS = GET REPLY --------------------------
		if err != nil {
			break
		}
		response, err = cn.readMsg()
		// XXX HANDLE ANY ERROR
		op = response.GetOp()
		// XXX op MUST BE XLRegMsg_Members
		_ = op

		if err == nil {
			id := response.GetClusterID()
			_ = id // XXX ignore for now
			which := xu.NewBitMap64(response.GetWhich())
			// DEBUG
			fmt.Printf("    client has received %d MEMBERS\n",
				which.Count())
			// END
			tokens := response.GetTokens() // a slice
			if which.Any() {
				offset := 0
				for i := uint(0); i < uint(cn.clusterSize); i++ {
					if which.Test(i) {
						token := tokens[offset]
						offset++
						cn.members[i], err = NewClusterMemberFromToken(
							token)
						stillToGet = stillToGet.Clear(i)
					}
				}
			}
			if stillToGet.None() {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
	return
} // GEEP4

// Send Bye, wait for and process Ack.

func (cn *ClientNode) ByeAndAck() (err error) {

	op := XLRegMsg_Bye
	request := &XLRegMsg{
		Op: &op,
	}
	// SHOULD CHECK FOR TIMEOUT
	err = cn.writeMsg(request)
	// DEBUG
	fmt.Printf("client %s BYE sent\n", cn.name)
	// END

	// Process ACK = BYE REPLY ----------------------------------
	if err == nil {
		var response *XLRegMsg

		// SHOULD CHECK FOR TIMEOUT AND VERIFY THAT IT'S AN ACK
		response, err = cn.readMsg()
		op := response.GetOp()
		_ = op
		// DEBUG
		fmt.Printf("    client %s has received ACK; err is %v\n",
			cn.name, err)
		// END
	}
	return
} // GEEP6
